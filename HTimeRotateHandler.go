package logger

import (
	"os"
	"time"
	"sync"
)

type HTimeRotateHandler struct{
	TimeRotateHandler
	buffer chan string
}

func NewHTimeRotateHandler(when When, level LevelType, filePath string, backup int) (*HTimeRotateHandler, error) {
	var fd *os.File
	suffixFormat := fileSuffixFormat(when)
	if state, err := os.Stat(filePath); os.IsNotExist(err) {
		// If file not exists, then create one.
		if fd, err = os.Create(filePath); err != nil {
			return nil, err
		}
	} else {
		// If file modified on the same rotate, then append it, or create new one.
		if state.ModTime().Format(suffixFormat) == time.Now().Format(fileSuffixFormat(when)){
			if fd, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600); err != nil {
				return nil, err
			}
		}else{
			os.Rename(filePath, filePath + "." + state.ModTime().Format(suffixFormat))
			if fd, err = os.Create(filePath); err != nil {
				return nil, err
			}
		}
	}

	handler := &HTimeRotateHandler{
		TimeRotateHandler: TimeRotateHandler{
			level:         level,
			fileName:      filePath,
			file:          fd,
			mu:            new(sync.RWMutex),
			rotateFormat:  suffixFormat,
			currentRotate: time.Now().Format(fileSuffixFormat(when)),
			backup: backup},
		buffer: make(chan string),
	}

	go func(){
		var bytes []byte
		for{
			handler.Rotate()
			select {
			case message := <- handler.buffer:
				bytes = append(bytes, []byte(message)...)
			}
			for ; len(bytes) > 0;  {
				count, _ := handler.TimeRotateHandler.file.Write(bytes)
				if count < len(bytes){
					handler.file.Sync()
				}
				bytes = bytes[count:]
			}
		}
	}()
	return handler, nil
}

func (handler *HTimeRotateHandler) Log(message string, level LevelType) error {
	if level < handler.level {
		return nil
	}
	handler.Rotate()
	message += "\r\n"
	handler.buffer <- message
	return nil
}