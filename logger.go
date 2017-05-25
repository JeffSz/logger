package logger

import (
	"sync"
	"strings"
	"time"
	"regexp"
	"fmt"
)

type WriteFunc func(string, LevelType) error

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

// Predefined message format
type FormatType string

const (
	MessageDF  = "{message}"
	LevelDF    = "{level}"
	DatetimeDf = "{datetime}"
)

var (
	_LOGGERS = make(map[string]*_Logger)
	_Df_REG  = regexp.MustCompile("\\{\\w\\}")
)

type Handler interface {
	log(string, LevelType) error
}

type _Logger struct {
	handlers   []Handler
	format     WriteFunc
	dateFormat FormatType
}

var (
	_LocalMu = new(sync.RWMutex)
)

func NewLogger(logger string) *_Logger {

	if _, ok := _LOGGERS[logger]; !ok {
		_LocalMu.Lock()
		defer _LocalMu.Unlock()
		if _, ok = _LOGGERS[logger]; !ok {
			_LOGGERS[logger] = &_Logger{dateFormat: "2006-01-02 15:04:05"}
		}
	}
	return _LOGGERS[logger]
}

func (logger *_Logger) SetFormat(format FormatType) error {

	params := _Df_REG.FindAllString(string(format), 0)
	for _, value := range params {
		if value != MessageDF && value != LevelDF && value != DatetimeDf {
			return NewError("Not valid predefined param. {message}, {level}, {datetime} is valid.")
		}
	}
	if Contains(params, MessageDF) {
		return NewError("No message part found!")
	}

	logger.format = func(message string, levelType LevelType) error {

		var levelString string
		switch levelType {
		case ALL:
			levelString = "ALL"
		case DEBUG:
			levelString = "DEBUG"
		case INFO:
			levelString = "INFO"
		case WARN:
			levelString = "WARN"
		case ERROR:
			levelString = "ERROR"
		case FATAL:
			levelString = "FATAL"
		default:
			levelString = "INFO"
		}

		msg := string(format)
		if strings.Contains(msg, LevelDF) {
			msg = strings.Replace(msg, LevelDF, levelString, -1)
		}
		if strings.Contains(msg, DatetimeDf) {
			msg = strings.Replace(msg, DatetimeDf, time.Now().Format(string(logger.dateFormat)), -1)
		}
		if strings.Contains(msg, MessageDF) {
			msg = strings.Replace(msg, MessageDF, message, -1)
		}

		for _, handler := range logger.handlers {
			handler.log(msg, levelType)
		}
		return nil
	}

	return nil
}

func (logger *_Logger) AddHandler(handler Handler)  {
	logger.handlers = append(logger.handlers, handler)
}

func (logger *_Logger) CleanHandler(handler Handler) {
	logger.handlers = make([]Handler, 0)
}

func (logger *_Logger) SetDatetimeFormat(format FormatType) {
	logger.dateFormat = format
}

func (logger *_Logger) Info(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg[1: len(msg)-1], INFO)
}

func (logger *_Logger) Debug(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg[1: len(msg)-1], DEBUG)
}

func (logger *_Logger) Warn(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg[1: len(msg)-1], WARN)
}

func (logger *_Logger) Error(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg[1: len(msg)-1], ERROR)
}

func (logger *_Logger) Fatal(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg[1: len(msg)-1], FATAL)
}

var(
	Logger = &_Logger{dateFormat: "2006-01-02 15:04:05"}
)

func BasicConfig(logFile string, level LevelType, format FormatType, datetimeFormat FormatType){
	_LocalMu.Lock()
	defer _LocalMu.Unlock()
	handler, err := NewBasicHandler(level, logFile)
	if err != nil{
		panic(err)
	}
	Logger.SetDatetimeFormat(datetimeFormat)
	Logger.SetFormat(format)
	Logger.AddHandler(handler)
}

func initBasic(){
	if len(Logger.handlers) == 0{
		_LocalMu.Lock()
		defer _LocalMu.Unlock()
		if len(Logger.handlers) == 0{
			BasicConfig("log.log", INFO, "[{datetime}] [{level}] {message}", "2006-01-02 15:04:05")
		}
	}
}

func Info(v ...interface{}){
	initBasic()
	Logger.Info(v)
}

func Debug(v ...interface{}) {
	initBasic()
	Logger.Debug(v)
}

func Warn(v ...interface{}) {
	initBasic()
	Logger.Warn(v)
}

func Error(v ...interface{}) {
	initBasic()
	Logger.Error(v)
}

func Fatal(v ...interface{}) {
	initBasic()
	Logger.Fatal(v)
}