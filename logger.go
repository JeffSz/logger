package logger

import (
	"sync"
	"strings"
	"time"
	"regexp"
	"fmt"
)

type WriteFunc func(string, LevelType) error

// Predefined message format
type FormatType string

const (
	MessageDF  = "{message}"
	LevelDF    = "{level}"
	DatetimeDf = "{datetime}"
)

var (
	_LOGGERS = make(map[string]*LogType)
	_Df_REG  = regexp.MustCompile("\\{\\w\\}")
)

type LogType struct {
	handlers   []Handler
	format     WriteFunc
	dateFormat FormatType
}

var (
	_LocalMu = new(sync.RWMutex)
)

func NewLogger(logger string) *LogType {

	if _, ok := _LOGGERS[logger]; !ok {
		_LocalMu.Lock()
		defer _LocalMu.Unlock()
		if _, ok = _LOGGERS[logger]; !ok {
			_LOGGERS[logger] = &LogType{dateFormat: "2006-01-02 15:04:05"}
		}
	}
	return _LOGGERS[logger]
}

func (logger *LogType) SetFormat(format FormatType) error {

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
			handler.Log(msg, levelType)
		}
		return nil
	}

	return nil
}

func (logger *LogType) AddHandler(handler Handler)  {
	logger.handlers = append(logger.handlers, handler)
}

func (logger *LogType) CleanHandler() {
	logger.handlers = make([]Handler, 0)
}

func (logger *LogType) SetDatetimeFormat(format FormatType) {
	logger.dateFormat = format
}

func (logger *LogType) Info(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg, INFO)
}

func (logger *LogType) Debug(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg, DEBUG)
}

func (logger *LogType) Warn(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg, WARN)
}

func (logger *LogType) Error(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg, ERROR)
}

func (logger *LogType) Fatal(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.format(msg, FATAL)
}

var(
	Logger = &LogType{dateFormat: "2006-01-02 15:04:05"}
)

func basicConfig(logFile string, level LevelType, format FormatType, datetimeFormat FormatType){
	handler, err := NewBasicHandler(level, logFile)
	if err != nil{
		panic(err)
	}
	Logger.SetDatetimeFormat(datetimeFormat)
	Logger.SetFormat(format)
	Logger.AddHandler(handler)
}

func BasicConfig(logFile string, level LevelType, format FormatType, datetimeFormat FormatType){
	_LocalMu.Lock()
	defer _LocalMu.Unlock()
	basicConfig(logFile, level, format, datetimeFormat)
}

func initBasic(){
	if len(Logger.handlers) == 0{
		_LocalMu.Lock()
		defer _LocalMu.Unlock()
		if len(Logger.handlers) == 0{
			basicConfig("log.log", INFO, "[{datetime}] [{level}] {message}", "2006-01-02 15:04:05")
		}
	}
}

func Info(v ...interface{}){
	initBasic()
	Logger.Info(v...)
}

func Debug(v ...interface{}) {
	initBasic()
	Logger.Debug(v...)
}

func Warn(v ...interface{}) {
	initBasic()
	Logger.Warn(v...)
}

func Error(v ...interface{}) {
	initBasic()
	Logger.Error(v...)
}

func Fatal(v ...interface{}) {
	initBasic()
	Logger.Fatal(v...)
}