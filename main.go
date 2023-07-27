package main

import (
	"fmt"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("something")
	logger.Error("something else")
	fmt.Println("Hello world. This is the ground.")
}
