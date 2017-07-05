package logger

import (
"os"
)

type LevelType int32

const (
	ALL   LevelType = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

type Handler interface {
	Log(string, LevelType) error
}

type BasicHandler struct {
	level         LevelType
	file          *os.File
}

func NewBasicHandler(level LevelType, filePath string) (*BasicHandler, error) {
	var fd *os.File
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if fd, err = os.Create(filePath); err != nil {
			return nil, err
		}
	} else {
		if fd, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600); err != nil {
			return nil, err
		}
	}

	return &BasicHandler{
		level:         level,
		file:          fd,
	}, nil
}

func (handler BasicHandler) Log(message string, level LevelType) error {
	if level < handler.level {
		return nil
	}
	handler.file.WriteString(message + "\r\n")
	return nil
}
