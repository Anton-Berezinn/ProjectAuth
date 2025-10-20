package main

import (
	"Project/internal/cache/redis"
	"Project/internal/configs"
	"Project/internal/handlers"
	"Project/internal/logger"
	"Project/internal/manager"
	"Project/internal/services"
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	config, err := configs.New()
	if err != nil {
		if errors.Is(err, configs.ErrData) {
			fmt.Println("Empty Data")
			return
		}
		panic(err)
	}

	logg := logger.New()
	db, err := services.NewDb(config)
	if err != nil {
		fmt.Println("Error to open database", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cacheUser, err := redis.NewClientUser(ctx)
	if err != nil {
		fmt.Println("Error to open cache", err)
		return
	}
	cacheCatalog, err := redis.NewClientCatalog(ctx)
	if err != nil {
		fmt.Println("Error to open cache", err)
		return
	}
	Model, err := manager.NewProject(db, logg, cacheUser, cacheCatalog)
	if err != nil {
		fmt.Println("error to NewProject", err)
		return
	}
	go Model.SomeWork()

	handlers.StartEcho(ctx, config.Project.Port, Model)
}
