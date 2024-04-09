package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	var (
		rhost  = flag.String("rhost", "localhost", "redis server host")
		rport  = flag.Int("rport", 6379, "redis server port")
		prefix = flag.String("pre", "seed:", "redis key prefix to watch")
	)
	flag.Parse()
	redisAddr := fmt.Sprintf("%s:%d", *rhost, *rport)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})

	ctx := context.Background()
	pubsub := rdb.PSubscribe(ctx, fmt.Sprintf("__keyspace@0__:%s", *prefix))

	ch := pubsub.Channel()
	fmt.Println("hello world")
	for rmsg := range ch {
		key := rmsg.Payload
		val, err := rdb.Get(ctx, key).Result()
		handleChange(key, val, err)
	}

}

func handleChange(redisKey string, redisVal string, err error) {
	msg := "Key change"

	if err != nil && err != redis.Nil {
		slog.Error(
			msg+": couldn't retrieve",
			slog.String("key", redisKey),
			slog.Any("err", err),
		)
		return
	}

	if err == redis.Nil {
		msg += ":key deleted"
	}
	slog.Info(
		msg,
		slog.String("key", redisKey),
		slog.String("val", redisVal),
	)
}
