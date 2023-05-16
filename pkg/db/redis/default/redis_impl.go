package redisImpl

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	redisrepo "github.com/thteam47/go-agent-pc/pkg/db/redis"
)

type RedisRepositoryImpl struct {
	redis   *cache.Cache
	timeOut time.Duration
}

func NewRedisRepo(redis *cache.Cache, timeOut time.Duration) redisrepo.RedisRepository {
	return &RedisRepositoryImpl{
		redis:   redis,
		timeOut: timeOut,
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func (inst *RedisRepositoryImpl) SetKeyToListKeyCache(key string, listCache string) {
	ctx, cancel := context.WithTimeout(context.Background(), inst.timeOut)
	defer cancel()
	keyCache := listCache
	var keyList []string
	keyCacheStatus := inst.redis.Get(ctx, keyCache, &keyList)
	if keyCacheStatus != nil {
		keyList = append(keyList, key)
	} else {
		if !stringInSlice(key, keyList) {
			keyList = append(keyList, key)
		}
	}

	if err := inst.redis.Set(&cache.Item{
		Ctx:   ctx,
		Key:   keyCache,
		Value: keyList,
	}); err != nil {
		panic(err)
	}
}
func (inst *RedisRepositoryImpl) RemoveListKeyCache(keyListCache string) {
	ctx, cancel := context.WithTimeout(context.Background(), inst.timeOut)
	defer cancel()
	var keyIndex []string
	keyCacheIndex := inst.redis.Get(ctx, keyListCache, &keyIndex)
	if keyCacheIndex == nil {
		for _, key := range keyIndex {
			inst.RemoveValueCache(key)
		}
	}
}
func (inst *RedisRepositoryImpl) GetValueCache(key string, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), inst.timeOut)
	defer cancel()
	err := inst.redis.Get(ctx, key, result)
	if err != nil {
		return err
	}
	return nil
}
func (inst *RedisRepositoryImpl) SetValueCache(key string, data interface{}, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), inst.timeOut)
	defer cancel()
	if err := inst.redis.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: data,
		TTL:   exp,
	}); err != nil {
		return err
	}
	return nil
}
func (inst *RedisRepositoryImpl) RemoveValueCache(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), inst.timeOut)
	defer cancel()
	err := inst.redis.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
