package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	var (
		rhost = flag.String("host", "localhost", "redis server host")
		rport = flag.Int("port", 6379, "redis server port")
		key   = flag.String("key", "namespace:key", "redis key to watch")
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
			handleChange(*key, res, value, err)
		}
		value = res
		time.Sleep(time.Millisecond * 200)
	}
}

func timeSince(unix string) time.Duration {
	unixTimestamp, err := strconv.ParseInt(unix, 10, 64)
	if err != nil {
		slog.Error("couldn't convert to unix timestamp")
	}
	t := time.Unix(unixTimestamp/1000, 0)
	duration := time.Since(t)
	return duration
}

func handleChange(redisKey, redisVal, oldVal string, err error) {
	var duration time.Duration
	if err != nil && err != redis.Nil {
		slog.Error(
			"couldn't retrieve",
			slog.String("key", redisKey),
			slog.Any("err", err),
		)
		return
	}

	if oldVal != "" {
		duration = timeSince(oldVal)
	}

	slog.Info(fmt.Sprintf("%s::%s [%v]", redisKey, redisVal, duration))
}
