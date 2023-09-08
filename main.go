package main

import (
	"fmt"
	"os"

	"shorty/shortener"

	"github.com/go-redis/redis/v8"
)

func main() {
	port := os.Getenv("PORT")
	redisAddr := os.Getenv("REDIS_ADDRESS")
	domain := os.Getenv("DOMAIN")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
	realRedisClient := shortener.NewRealRedisClient(redisClient)
	shorty := shortener.NewRedisURLShortener(realRedisClient)

	if port == "" {
		port = "8080"
	}

	r := setupRouter(shorty, domain)
	err := r.Run(":" + port)
	if err != nil {
		fmt.Println("Error starting Gin server: ", err)
	}
}
