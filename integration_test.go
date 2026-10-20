package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithRedis(t *testing.T) {
	ctx := context.Background()
	req := tc.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop redis: %s", err)
		}
	}()

	ctx = context.Background()
	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("couldn't get endpoint %s", err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: endpoint})

	rdb.Set(ctx, "testkey", "helloworld", 10*time.Second)
	res, err := rdb.Get(ctx, "testkey").Result()
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	log.Printf("res: %v", res)
}
