package main

import "github.com/desertthunder/quotesky/log"

var logger = log.DefaultLogger()

func main() {
	logger.Debug("Debug message")
	logger.Debug("Debug message")
	logger.Debug("Debug message")
	logger.Debug("Debug message")
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")
	logger.Fatal("Fatal message")
}
