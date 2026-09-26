package main

import (
	"database/sql"
	"log"

	"otomedia/task-managment/internal/cache"
	"otomedia/task-managment/internal/config"
	"otomedia/task-managment/internal/handler"
	"otomedia/task-managment/internal/repository"
	"otomedia/task-managment/internal/router"
	"otomedia/task-managment/internal/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("failed to open mysql connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to mysql: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	taskRepo := repository.NewMySQLTaskRepository(db)
	redisCache := cache.NewRedisCache(redisClient)
	taskService := service.NewTaskService(taskRepo, redisCache)
	taskHandler := handler.NewTaskHandler(taskService)

	r := router.New(taskHandler)

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}