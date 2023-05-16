package db

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/thteam47/go-agent-pc/pkg/configs"
	redisrepo "github.com/thteam47/go-agent-pc/pkg/db/redis"
	redisImpl "github.com/thteam47/go-agent-pc/pkg/db/redis/default"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	MongoDB         *mongo.Collection
	RedisRepository redisrepo.RedisRepository
	GrpcConnConfig  configs.GrpcClientConnConfig
	JwtKey          string
	Exp             time.Duration
	TotpSecret      string
	TimeRequestId   time.Duration
	TimeEmailOtp    time.Duration
}

func NewHandlerWithConfig(config *configs.Config) (*Handler, error) {
	mongodb, err := connectMongo(config.MongoDb)
	if err != nil {
		return nil, errors.WithMessage(err, "connectMongo")
	}
	redisCache, err := connectRedis(config.RedisCache)
	if err != nil {
		return nil, errors.WithMessage(err, "connectRedis")
	}
	redisRepo := redisImpl.NewRedisRepo(redisCache, config.TimeoutRedis)
	return &Handler{
		MongoDB:         mongodb,
		RedisRepository: redisRepo,
		JwtKey:          config.KeyJwt,
		GrpcConnConfig:  config.GrpcClientConn,
		Exp:             config.Exp,
		TotpSecret:      config.TotpSecret,
		TimeRequestId:   config.TimeRequestId,
		TimeEmailOtp:    config.TimeEmailOtp,
	}, nil
}

func connectMongo(cf configs.MongoDB) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cf.Url))
	if err != nil {
		return nil, errors.WithMessage(err, "mongo.Connect")
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, errors.WithMessage(err, "client.Ping")
	}
	// password := "admin"
	// passHash, _ := HashPassword(password)
	// user := User {
	// 	Username: "admin",
	// 	Password: passHash,
	// }
	// DB.Collection("user").InsertOne(ctx, user)
	// if err != nil {
	// 	panic(err)
	// }
	return client.Database(cf.DbName).Collection(cf.Collection), nil
}
func connectRedis(cf configs.Redis) (*cache.Cache, error) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			cf.Address: cf.Url,
		},
	})

	err := ring.Ping(context.Background()).Err()
	if err != nil {
		return nil, errors.WithMessage(err, "ring.Ping")
	}
	memCache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	return memCache, nil
}
