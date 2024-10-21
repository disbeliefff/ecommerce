package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
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
	return err
}

func (hook *WriteHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{e}
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

	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s:%d", filename, frame.Line), filename
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	l.SetOutput(io.Discard)

	l.AddHook(&WriteHook{
		Writer:    []io.Writer{l.Writer()},
		LogLevels: []logrus.Level{logrusLevel},
	})

	l.SetLevel(logrusLevel)

	e = logrus.NewEntry(l)
}
