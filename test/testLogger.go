package main

import (
	"logger"
	"fmt"
	"time"
	"os"
)

func main(){
	//logger.BasicConfig("test.log", logger.INFO, "[{datetime}] {message}", "2006-01-02")
	handler, _ := logger.NewHTimeRotateHandler(logger.Minute, logger.INFO, "/Users/Jeff/go/src/logger/test1.log", 3)
	logger.Logger.AddHandler(handler)
	logger.Logger.SetFormat("[{level}] {message}")

	logger.Debug("This is debug test")
	logger.Info("This is info test")
	logger.Error("This is error test")
	logger.Fatal("This is fatal test")
	logger.Error("This is error test")
	time.Sleep(time.Minute * 2)
	logger.Error("This is rotate test")
	time.Sleep(time.Minute * 2)
	logger.Error("This is rotate test")
	time.Sleep(time.Minute * 2)
	logger.Error("This is backup test")
	time.Sleep(time.Minute * 2)
	logger.Error("This is backup test")
	fmt.Println("Done!")
	os.Exit(0)
}
