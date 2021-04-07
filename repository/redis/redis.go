package redis

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var addr = os.Getenv("redis__addr")
var password = os.Getenv("redis__password")

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     addr,
	Password: password,
})

func Set(key, value string) {
	err := rdb.Set(ctx, key, value, 0).Err()

	if err != nil {
		panic(err)
	}
}

func Get(key string) string {
	value, err := rdb.Get(ctx, key).Result()

	if err != nil {
		panic(err)
	}

	return value
}
