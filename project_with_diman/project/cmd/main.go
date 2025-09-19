package main

import (
	"Project/internal/cache/redis"
	"Project/internal/configs"
	"Project/internal/handlers"
	"Project/internal/logger"
	"Project/internal/manager"
	"Project/internal/services/users"
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	//Todo: prometheus,grafana,kafka,hash256 password, проверить ручку, что по версии токена можно заходить
	config, err := configs.New()
	if err != nil {
		if errors.Is(err, configs.ErrData) {
			fmt.Println("Empty Data")
			return
		}
		panic(err)
	}

	logg := logger.New()
	db, err := users.NewDb(config)
	if err != nil {
		fmt.Println("Error to open database", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cache, err := redis.NewClient(ctx)
	if err != nil {
		fmt.Println("Error to open cache", err)
		return
	}
	Model, err := manager.NewProject(db, logg, cache)
	if err != nil {
		fmt.Println("error to NewProject", err)
		return
	}
	go Model.SomeWork()

	handlers.StartEcho(ctx, config.Project.Port, Model)
}
