package main

import (
	"Project/internal/configs"
	"Project/internal/handlers"
	"Project/internal/logger"
	"Project/internal/manager"
	"Project/internal/services/users"
	"errors"
	"fmt"
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
	db, err := users.NewDb(config)
	if err != nil {
		fmt.Println("Error to open database", err)
		return
	}
	Model, err := manager.NewProject(db, logg)
	if err != nil {
		fmt.Println("error to NewProject", err)
		return
	}
	go Model.SomeWork()

	handlers.StartEcho(config.Project.Port, Model)
}
