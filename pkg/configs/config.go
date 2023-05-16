package configs

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GrpcPort       string               `mapstructure:"grpc_port"`
	HttpPort       string               `mapstructure:"http_port"`
	MongoDb        MongoDB              `mapstructure:"mongo_db"`
	RedisCache     Redis                `mapstructure:"redis_cache"`
	KeyJwt         string               `mapstructure:"key_jwt"`
	GrpcClientConn GrpcClientConnConfig `mapstructure:"grpc_conn"`
	Exp            time.Duration        `mapstructure:"exp"`
	TotpSecret     string               `mapstructure:"totp_secret"`
	TimeoutRedis   time.Duration        `mapstructure:"time_out_redis"`
	TimeRequestId  time.Duration        `mapstructure:"time_request_id"`
	TimeEmailOtp   time.Duration        `mapstructure:"time_email_otp"`
}

type MongoDB struct {
	Url        string `mapstructure:"url"`
	DbName     string `mapstructure:"db_name"`
	Collection string `mapstructure:"collection"`
}

type Redis struct {
	Address string `mapstructure:"address"`
	Url     string `mapstructure:"url"`
}

type GrpcClientConnConfig struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	AccessToken string        `mapstructure:"access_token"`
}

func Get(this interface{}, key string) interface{} {
	return this.(map[interface{}]interface{})[key]
}
func String(payload interface{}) string {
	var load string
	if pay, oh := payload.(string); oh {
		load = pay
	} else {
		load = ""
	}
	return load
}
func LoadConfig() (*Config, error) {
	filename, err := filepath.Abs("../pkg/configs/config.yaml")
	if err != nil {
		return nil, errors.WithMessage(err, "filepath.Abs")
	}
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, errors.WithMessage(err, "ioutil.ReadFile")
	}
	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, errors.WithMessage(err, "yaml.Unmarshal")
	}
	return &config, nil
}
