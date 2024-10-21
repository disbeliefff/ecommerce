package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"runtime"
)

type WriteHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *WriteHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return nil
}

func (hook *WriteHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var loggerInstance *logrus.Logger
var entryInstance *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	if entryInstance == nil {
		log.Fatal("Logger not initialized. Please call Init() before using the logger.")
	}
	return &Logger{entryInstance}
}

func (l *Logger) LWithField(k string, v any) *Logger {
	return &Logger{l.WithField(k, v)}
}

func (l *Logger) LWithFields(fields map[string]any) *Logger {
	return &Logger{l.WithFields(fields)}
}

func Init(level string) {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatalln(err)
	}

	loggerInstance = logrus.New()
	loggerInstance.SetReportCaller(true)
	loggerInstance.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s:%d", filename, frame.Line), filename
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	// Измените на os.Stdout или файл, если необходимо
	loggerInstance.SetOutput(os.Stdout)

	loggerInstance.AddHook(&WriteHook{
		Writer:    []io.Writer{loggerInstance.Writer()},
		LogLevels: []logrus.Level{logrusLevel},
	})

	loggerInstance.SetLevel(logrusLevel)

	entryInstance = logrus.NewEntry(loggerInstance)
}
