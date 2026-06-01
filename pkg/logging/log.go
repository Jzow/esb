package logging

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/file"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Level int

var scrollDateMutex sync.Mutex

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	scrollDate = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// Setup initialize the log instance
func Setup() {
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F, err = file.MustOpen(fileName, filePath)
	if err != nil {
		log.Fatalf("logging.Setup err: %v", err)
	}

	logger = log.New(io.MultiWriter(os.Stdout, F), DefaultPrefix, log.LstdFlags|log.Lmicroseconds)
	scrollDate = time.Now().Format(setting.AppSetting.TimeFormat)
}

// Debug output logs at debug level
func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

// Info output logs at info level
func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}
func Infof(format string, v ...interface{}) {
	setPrefix(INFO)
	logger.Printf(format, v...)
}

// Warn output logs at warn level
func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}
func Warnf(format string, v ...interface{}) {
	setPrefix(WARNING)
	logger.Printf(format, v...)
}

// Error output logs at error level
func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}
func Errorf(format string, v ...interface{}) {
	setPrefix(ERROR)
	logger.Printf(format, v...)
}

// Fatal output logs at fatal level
func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}
func Fatalf(format string, v ...interface{}) {
	setPrefix(FATAL)
	logger.Printf(format, v...)
}

// setPrefix set the prefix of the log output
func setPrefix(level Level) {
	if time.Now().Format(setting.AppSetting.TimeFormat) != scrollDate {
		scrollDateMutex.Lock()
		defer scrollDateMutex.Unlock()
		if time.Now().Format(setting.AppSetting.TimeFormat) != scrollDate {
			Setup()
		}
	}
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
