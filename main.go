package main

import (
	"log"

	"github.com/qreepex/voting-backend/internal/data"
	"github.com/qreepex/voting-backend/internal/redis"
	"github.com/qreepex/voting-backend/internal/web"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := data.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	redisClient, err := redis.InitRedis()
	if err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}

	web.Init(db, redisClient)
}
