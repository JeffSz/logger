package logger

import (
	"os"
	"time"
	"sync"
	"path"
	"strings"
	"io/ioutil"
	"fmt"
	"sort"
)

type When string

const (
	Minute = "m"
	Hour  = "H"
	Day   = "D"
	Month = "M"
	Year  = "Y"
)

type TimeRotateHandler struct {
	level         LevelType
	fileName      string
	file          *os.File
	mu            *sync.RWMutex
	rotateFormat  string
	currentRotate string
	backup int
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
	}
	return timeFormat
}

func NewTimeRotateHandler(when When, level LevelType, filePath string, backup int) (*TimeRotateHandler, error) {
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
		backup: backup,
	}, nil
}

func (handler *TimeRotateHandler) shouldRotate() bool {
	return time.Now().Format(handler.rotateFormat) != handler.currentRotate
}

// Keep the given number of backups.
func (handler *TimeRotateHandler) keepBackUp(){
	dirPath, fileName := path.Split(handler.fileName)
	if dirPath == ""{
		dirPath = "."
	}else{
		dirPath = strings.TrimRight(dirPath, string(os.PathSeparator))
	}
	dir, _ := ioutil.ReadDir(dirPath)
	fileNames := make([]string, 0)
	for _, f := range dir{
		if !f.IsDir() && f.Name() != fileName && strings.HasPrefix(f.Name(), fileName){
			if _, err := time.Parse(handler.rotateFormat, strings.TrimLeft(f.Name(), fileName + ".")); err == nil{
				fileNames = append(fileNames, f.Name())
			}
		}
	}
	if len(fileNames) > handler.backup{
		fmt.Println(fileNames)
		fmt.Println(handler.backup)
		sort.Strings(fileNames)
		for _, f := range fileNames[: len(fileNames) - handler.backup]{
			fmt.Println(dirPath + string(os.PathSeparator) + f)
			os.Remove(dirPath + string(os.PathSeparator) + f)
		}
	}
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

			go handler.keepBackUp()
		}
	}
	return nil
}

func (handler *TimeRotateHandler) Log(message string, level LevelType) error {
	if level < handler.level {
		return nil
	}
	handler.Rotate()
	handler.file.WriteString(message + "\r\n")
	return nil
}