package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	var (
		rhost = flag.String("h", "localhost", "redis server host")
		rport = flag.Int("p", 6379, "redis server port")
		key   = flag.String("k", "seed:yell_e2e_three", "redis key to watch")
	)
	flag.Parse()

	redisAddr := fmt.Sprintf("%s:%d", *rhost, *rport)
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()

	var value string
	slog.Info("watching key: " + *key)
	for {
		res, err := rdb.Get(ctx, *key).Result()
		if value != res {
			handleChange(*key, res, err)
		}
		value = res
		time.Sleep(time.Millisecond * 200)
	}
}

func handleChange(redisKey string, redisVal string, err error) {
	if err != nil && err != redis.Nil {
		slog.Error(
			"couldn't retrieve",
			slog.String("key", redisKey),
			slog.Any("err", err),
		)
		return
	}
	slog.Info(fmt.Sprintf("%s::%s", redisKey, redisVal))
}
