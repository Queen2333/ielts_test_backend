package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	// 连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "172.25.138.133:6379",
		Password: "Yx180236",
		DB:       0,
	})

	ctx := context.Background()

	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)

	// 获取所有 keys
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		log.Fatalf("Failed to get keys: %v", err)
	}

	if len(keys) == 0 {
		fmt.Println("No keys found in Redis")
		return
	}

	fmt.Printf("Found %d keys:\n", len(keys))
	for _, key := range keys {
		fmt.Printf("  - %s\n", key)
	}

	// 删除所有 keys
	fmt.Println("\nDeleting all keys...")
	for _, key := range keys {
		err := rdb.Del(ctx, key).Err()
		if err != nil {
			fmt.Printf("Failed to delete key %s: %v\n", key, err)
		} else {
			fmt.Printf("Deleted: %s\n", key)
		}
	}

	fmt.Println("\n✅ All Redis data cleared successfully!")
}
