package logger

import (
	"os"
	"sync"
	"time"
)

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

func (handler *BasicHandler) log(message string, level LevelType) error {
	if level < handler.level {
		return nil
	}
	handler.file.WriteString(message + "\r\n")
	return nil
}

type When string

const (
	Minute = "m"
	Hour  = "H"
	Day   = "D"
	Month = "M"
	Year  = "Y"
	Mon   = "W1"
	Tue   = "W2"
	Wed   = "W3"
	Thu   = "W4"
	Fri   = "W5"
	Sat   = "W6"
	Sun   = "W7"
)

type TimeRotateHandler struct {
	level         LevelType
	fileName      string
	file          *os.File
	mu            *sync.RWMutex
	rotateFormat  string
	currentRotate string
}

func fileSuffixFormat(when When) string {
	timeFormat := "2006-01-02 15:04:05"
	switch when {
	case Minute:
		timeFormat = "2006-01-02@15:04"
	case Hour:
		timeFormat = "2006-01-02@15"
	case Day:
		timeFormat = "2006-01-02"
	case Month:
		timeFormat = "2006-01"
	case Year:
		timeFormat = "2006"
	case Mon, Tue, Wed, Thu, Fri, Sat, Sun:
		timeFormat = "2006-01-02"
	}
	return timeFormat
}

func NewTimeRotateHandler(when When, level LevelType, filePath string) (*TimeRotateHandler, error) {
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

	return &TimeRotateHandler{
		level:         level,
		fileName:      filePath,
		file:          fd,
		mu:            new(sync.RWMutex),
		rotateFormat:  suffixFormat,
		currentRotate: time.Now().Format(fileSuffixFormat(when)),
	}, nil
}

func (handler *TimeRotateHandler) shouldRotate() bool {
	return time.Now().Format(handler.rotateFormat) != handler.currentRotate
}

func (handler *TimeRotateHandler) Rotate() error {
	if handler.shouldRotate(){
		handler.mu.Lock()
		defer handler.mu.Unlock()
		if handler.shouldRotate(){
			handler.file.Close()
			os.Rename(handler.fileName, handler.fileName + "." + handler.currentRotate)
			if fd, err := os.Create(handler.fileName); err != nil {
				return err
			}else{
				handler.file = fd
				handler.currentRotate = time.Now().Format(handler.rotateFormat)
			}
		}
	}
	return nil
}

func (handler *TimeRotateHandler) log(message string, level LevelType) error {
	if level < handler.level {
		return nil
	}
	handler.Rotate()
	handler.file.WriteString(message + "\r\n")
	return nil
}
