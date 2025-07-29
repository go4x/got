package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	redismock "github.com/go-redis/redismock/v9"
	rds "github.com/redis/go-redis/v9"
)

func MockRedis() (*rds.Client, redismock.ClientMock) {
	return redismock.NewClientMock()
}

func NewMiniRedis() *rds.Client {
	// 测试用miniredis
	mr, err := miniredis.Run()
	if err != nil {
		panic(fmt.Errorf("new test redis error: %v", err))
	}
	// 使用miniredis创建client
	client := rds.NewClient(&rds.Options{
		Addr:         mr.Addr(),
		Password:     "",
		DB:           0,
		MaxRetries:   5,
		MinIdleConns: 2,
		TLSConfig:    nil,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis error: %s", err.Error()))
	}
	log.Printf("redis connected, url: %s\n", client.Conn().String())
	return client
}

func NewRedisCluster() *rds.Client {
	// TODO mock redis cluster
	return nil
}
