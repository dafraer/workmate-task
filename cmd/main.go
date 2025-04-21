package main

import (
	"context"
	"fmt"
	"github.com/dafraer/workmate-task/api"
	"github.com/dafraer/workmate-task/tasks"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {
	//Check that port has been specified
	if len(os.Args) < 2 {
		panic("Port needs to be passed as an argument")
	}
	port := os.Args[1]

	// Create a task manager
	tm := tasks.NewTaskManager()

	//Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("error while creating new Logger, %v ", err))
	}

	//Create a new service
	service := api.New(logger.Sugar(), tm)

	//Context that cancels on os.Interrupt
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	//Run the service
	if err := service.Run(ctx, port); err != nil {
		panic(err)
	}
}
