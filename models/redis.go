package models

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client
var ctx = context.Background()

func RedisInit() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully!")
}

func SaveMessageToRedis(sender, message string) error {
	// 创建带时间戳的消息结构
	msg := map[string]interface{}{
		"sender":  sender,
		"message": message,
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	}

	// 序列化为JSON
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 使用RPush而不是LPush，这样新消息会追加到列表末尾
	err = RedisClient.RPush(ctx, "chat:messages", msgJSON).Err()
	if err != nil {
		return err
	}
	// 修剪列表，只保留最近的100条消息
	return RedisClient.LTrim(ctx, "chat:messages", -100, -1).Err()
}

func GetHistoryMessages() ([]map[string]string, error) {
	// 获取最近的100条消息（从旧到新）
	messages, err := RedisClient.LRange(ctx, "chat:messages", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var result []map[string]string
	for _, msg := range messages {
		var data map[string]string
		if err := json.Unmarshal([]byte(msg), &data); err == nil {
			result = append(result, data)
		}
	}
	return result, nil
}
