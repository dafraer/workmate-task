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
	if err != nil {
		panic(fmt.Errorf("error while creating new Logger, %v ", err))
	}

	//Create a new service
	service := api.New(logger.Sugar(), tm)

	//Run the service
	if err := service.Run(context.Background(), ":8080"); err != nil {
		panic(err)
	}
}
