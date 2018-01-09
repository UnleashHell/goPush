package log

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	debugLogger *loggerObj
	infoLogger  *loggerObj
	errorLogger *loggerObj
	fatalLogger *loggerObj
}

type loggerObj struct {
	file     *os.File
	obj      *log.Logger
	lastDate *time.Time
	mu       *sync.RWMutex
}

const (
	DEBUG = iota
	INFO
	ERROR
	FATAL
	DATEFORMATE = "2006-01-02"
)

var Instance *Logger
var isConsole = true
var level = 1
var logDir = ""

func (this *Logger) InitLogger(dir string) (e error) {
	logDir = dir
	err := createDir(logDir)
	if err != nil {
		fmt.Println("mkdir dir failed")
		e = err
	} else {
		this.debugLogger = new(loggerObj)
		this.infoLogger = new(loggerObj)
		this.errorLogger = new(loggerObj)
		this.fatalLogger = new(loggerObj)
		makeLoggerObj(this.debugLogger, "debug")
		makeLoggerObj(this.infoLogger, "info")
		makeLoggerObj(this.errorLogger, "error")
		makeLoggerObj(this.fatalLogger, "fatal")
	}
	return
}

func (this *Logger) Info(format string, args ...interface{}) {
	if level <= INFO {
		this.infoLogger.write("info", format, args...)
	}
}

func (this *Logger) Fatal(format string, args ...interface{}) {
	if level <= FATAL {
		this.fatalLogger.write("fatal", format, args...)
	}
}

func (this *Logger) Error(format string, args ...interface{}) {
	if level <= ERROR {
		this.errorLogger.write("error", format, args...)
	}
}

func (this *Logger) Debug(format string, args ...interface{}) {
	if level <= DEBUG {
		this.debugLogger.write("debug", format, args...)
	}
}

func (this *Logger) SetLevel(userLevel int) {
	level = userLevel
}

func (this *Logger) SetConsole(console bool) {
	isConsole = console
}

func (this *loggerObj) isNewDay() bool {
	now := time.Now().Format(DATEFORMATE)
	t, _ := time.Parse(DATEFORMATE, now)
	return t.After(*this.lastDate)

}

func (this *loggerObj) write(levelString string, format string, args ...interface{}) {
	isNewDay := this.isNewDay()
	if isNewDay {
		makeLoggerObj(this, levelString)
	}
	str := fmt.Sprintf(format, args...)
	if this.obj != nil {
		this.obj.Println(str)
	}
	if isConsole {
		fmt.Println(str)
	}

}

func createDir(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			if os.IsPermission(err) {
				fmt.Println("create dir error:", err.Error())
				e = err
			}
		}
	}
	return
}

func makeLoggerObj(l *loggerObj, name string) {
	now := time.Now().Format(DATEFORMATE)
	l.mu = new(sync.RWMutex)
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		l.file.Close()
	}
	t, _ := time.Parse(DATEFORMATE, now)
	l.lastDate = &t
	fileName := logDir + "/" + name + "-" + now + ".log"
	f, _err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if _err == nil {
		l.file = f
		l.obj = log.New(l.file, "["+name+"] ", log.Ldate|log.Ltime)
	}
}
