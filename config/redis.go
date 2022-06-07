package config

import "github.com/go-redis/redis/v8"

var Redis *redis.Client

// Create Redis client
func CreateRedisClient() {
	opt, err := redis.ParseURL("redis://localhost:6364/0")

	// Handle error -> Crash server
	if err != nil {
		panic(err)
	}

	redis := redis.NewClient(opt)
	Redis = redis
}
