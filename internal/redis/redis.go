package redis

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func InitRedis() (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	res, err := client.Get(context.Background(), "vote:ip:a").Result()
	if err != nil && err != redis.Nil {
		log.Printf("Failed to get key: %v", err)
		return nil, err
	}

	log.Printf("Got key: %v", res)

	return &Redis{Client: client}, nil
}

func (r *Redis) SetVote(ipAddress string, cookie string) error {
	timeIn24Hours := time.Now().Add(time.Hour * 24)

	pipe := r.Client.TxPipeline()
	pipe.HSet(context.TODO(), "vote:cookie:"+cookie, "t", timeIn24Hours.Unix(), "i", ipAddress)
	pipe.HSet(context.TODO(), "vote:ip:"+ipAddress, "t", timeIn24Hours.Unix(), "c", cookie)
	pipe.Expire(context.TODO(), "vote:cookie:"+cookie, time.Until(timeIn24Hours))
	pipe.Expire(context.TODO(), "vote:ip:"+ipAddress, time.Until(timeIn24Hours))
	_, err := pipe.Exec(context.TODO())
	return err

}

func (r *Redis) SetIP(ipAddress string, timestamp time.Time, cookie string) error {
	pipe := r.Client.TxPipeline()
	pipe.HSet(context.TODO(), "vote:ip:"+ipAddress, "t", timestamp.Unix(), "c", cookie)
	pipe.Expire(context.TODO(), "vote:ip:"+ipAddress, time.Until(timestamp))
	_, err := pipe.Exec(context.TODO())
	if err != redis.Nil {
		return err
	}

	return nil
}

func (r *Redis) SetCookie(cookie string, timestamp time.Time, ipHash string) error {
	pipe := r.Client.TxPipeline()
	pipe.HSet(context.TODO(), "vote:cookie:"+cookie, "t", timestamp.Unix(), "i", ipHash)
	pipe.Expire(context.TODO(), "vote:cookie:"+cookie, time.Until(timestamp))
	_, err := pipe.Exec(context.TODO())

	if err != redis.Nil {
		return err
	}

	return nil
}

type RedisCookie struct {
	Timestamp int64
	Ip        string
	Pending   int64
}

func (r *Redis) CheckCookie(cookie string) (*RedisCookie, error) {
	result, err := r.Client.HGetAll(context.TODO(), "vote:cookie:"+cookie).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	redisCookie := RedisCookie{}
	for key, value := range result {
		if key == "t" {
			redisCookie.Timestamp, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
		} else if key == "i" {
			redisCookie.Ip = value
		}
	}

	redisCookie.Pending = redisCookie.Timestamp - time.Now().Unix()
	if redisCookie.Pending < 0 {
		return nil, nil
	}

	return &redisCookie, nil
}

type RedisIp struct {
	Timestamp int64
	Cookie    string
	Pending   int64
}

func (r *Redis) CheckIP(ipAddress string) (*RedisIp, error) {
	result, err := r.Client.HGetAll(context.TODO(), "vote:ip:"+ipAddress).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	redisIp := RedisIp{}
	for key, value := range result {
		if key == "t" {
			redisIp.Timestamp, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
		} else if key == "c" {
			redisIp.Cookie = value
		}
	}

	redisIp.Pending = redisIp.Timestamp - time.Now().Unix()

	if redisIp.Pending < 0 {
		return nil, nil
	}

	return &redisIp, nil
}
