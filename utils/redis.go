package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// InitRedis 初始化Redis客户端
func InitRedis() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.100.213:6379", // Redis服务器地址
		Password: "Yx180236",                   // Redis密码，如果没有设置密码则为空
		DB:       0,                    // Redis数据库索引（默认为0）
	})
	// 可以在这里添加其他的初始化配置，例如设置连接池大小等

	// 测试连接是否成功
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return err
	}
	fmt.Println(pong) // 打印结果: PONG
	return nil
}

// Set 设置一个键值对 (如果过期时间小于等于0，则表示不设置过期时间)
func Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()

	if expiration <= 0 {
		// 如果过期时间小于等于0，则表示不设置过期时间
		err := rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			return err
		}
	} else {
		err := rdb.Set(ctx, key, value, expiration).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// Get 获取一个键的值
func Get(key string) (string, error) {
	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// 测试 Redis 写的一段代码

/****** 这一段写在main函数里 *****/
//ctx := context.Background()
//client := connectToRedis()
//
//// 测试连接是否成功
//pingRedis(ctx, client)
//
//// 设置一个键值对
//err := setKeyValue(ctx, client, "key100", "zheshiyiduanceshi")
//if err != nil {
//	fmt.Println("Error setting value:", err)
//	return
//}
//
//// 获取一个键的值
//value, err := getValueByKey(ctx, client, "key100")
//if err == redis.Nil {
//	fmt.Println("Key does not exist")
//} else if err != nil {
//	fmt.Println("Error getting value:", err)
//} else {
//	fmt.Println("mykey:", value)
//}
//
//// 关闭连接
//client.Close()

/****** 这一段是链接redis定义的函数 *****/

// func connectToRedis() *redis.Client {
// 	return redis.NewClient(&redis.Options{
// 		Addr: "10.244.74.249:6379", // Redis服务器地址和端口
// 	})
// }

// func pingRedis(ctx context.Context, client *redis.Client) {
// 	pong, err := client.Ping(ctx).Result()
// 	if err != nil {
// 		fmt.Println("Error connecting to Redis:", err)
// 		return
// 	}
// 	fmt.Println("Connected to Redis:", pong)
// }

// func setKeyValue(ctx context.Context, client *redis.Client, key, value string) error {
// 	return client.Set(ctx, key, value, 0).Err()
// }

// func getValueByKey(ctx context.Context, client *redis.Client, key string) (string, error) {
// 	return client.Get(ctx, key).Result()
// }
