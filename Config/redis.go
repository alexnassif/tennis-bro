package Config

import "github.com/go-redis/redis/v8"

var Redis *redis.Client

func CreateRedisClient() {
    opt, err := redis.ParseURL("redis://localhost:6379")
    if err != nil {
        panic(err)
    }

    redis := redis.NewClient(opt)
    Redis = redis
}