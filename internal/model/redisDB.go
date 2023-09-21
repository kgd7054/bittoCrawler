package model

import (
	conf "bittoCralwer/config"
	"crypto/tls"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)

type RedisDB struct {
	client *redis.Client
}

// NewRedisDB : RedisDB 객체 할당 및 반환
// func NewRedisDB(config map[string]interface{}, root *Repositories) (IRepository, error) {
func NewRedisDB(config *conf.Config, root *Repositories) (IRepository, error) {
	redisConfig, ok := config.Repositories["redis-db"]["db"].(string)
	if !ok || len(redisConfig) < 0 {
		redisConfig = "0"
	}
	redisDB, err := strconv.Atoi(redisConfig)
	if err != nil {
		redisDB = 0
	}

	redisOption := redis.Options{
		Addr:      config.Repositories["redis-db"]["datasource"].(string),
		Password:  config.Repositories["redis-db"]["pass"].(string), // no password set
		DB:        redisDB,                                          // use default DB
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if strings.EqualFold(config.Common.ServiceId, "alpha") {
		redisOption.TLSConfig = nil
	}

	client := redis.NewClient(&redisOption)

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	r := &RedisDB{
		client: client,
	}

	log.Trace("load repository : RedisDB")
	return r, nil
}

func (p *RedisDB) Start() error {
	return nil
}

func (p *RedisDB) SetCache(key, value string) error {
	err := p.client.Set(key, value, time.Duration(60)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (p *RedisDB) GetCache(key string) (string, error) {
	result, err := p.client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
