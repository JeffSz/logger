package main

import (
	"logger"
	"fmt"
	"time"
)

func main(){
	//logger.BasicConfig("test.log", logger.INFO, "[{datetime}] {message}", "2006-01-02")
	handler, _ := logger.NewTimeRotateHandler(logger.Minute, logger.INFO, "test1.log")
	logger.Logger.AddHandler(handler)
	logger.Logger.SetFormat("[{level}] {message}")

	fmt.Println("Debug")
	logger.Debug("DDDDDDDDDD")
	fmt.Println("Info")
	logger.Info("IIIIIIIIII")
	fmt.Println("Error")
	logger.Error("EEEEEEEEEE")
	fmt.Println("Fatal")
	logger.Fatal("FFFFFFFFFF")
	fmt.Println("Warn")
	logger.Warn("WWWWWWWWWW")
	time.Sleep(time.Second)
	fmt.Println("Done!")
}
