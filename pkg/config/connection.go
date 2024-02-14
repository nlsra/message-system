package config

import (
	"github.com/go-redis/redis"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (a *App) DBConnection() {
	var error error

	postgresDSN := os.Getenv("DB_DRIVER") + "://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" +
		os.Getenv("DB_NAME")

	a.DB, error = gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if error != nil {
		log.Fatal(error)
	} else {
		log.Println("Database connection successful")
	}
}

func (a *App) RedisConnection() {
	var error error

	redisOptions := &redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	a.RedisClient = redis.NewClient(redisOptions)
	if error != nil {
		log.Fatal(error)
	} else {
		log.Println("Redis connection successful")
	}
}
