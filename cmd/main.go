package main

import (
	"context"
	"fmt"

	"github.com/dafraer/workmate-task/api"
	"github.com/dafraer/workmate-task/tasks"
	"go.uber.org/zap"
)

func main() {
	// Create a task manager
	tm := tasks.NewTaskManager()

	//Create logger
	logger, err := zap.NewDevelopment()
	var sugar *zap.SugaredLogger
	if err != nil {
		panic(fmt.Errorf("error while creating new Logger, %v ", err))
	}

	//Create sugared logger
	if logger != nil {
		sugar = logger.Sugar()
	}

	//Create a new service
	service := api.NewService(sugar, tm)

	//Run the service
	if err := service.Run(context.Background(), "localhost:8080"); err != nil {
		panic(err)
	}
}
