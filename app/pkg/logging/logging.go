package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
)

type Logger struct {
	*logrus.Entry
}

var instanse Logger
var once sync.Once

func Init(level string) Logger {
	once.Do(func() {
		logrusLevel, err := logrus.ParseLevel(level)
		if err != nil {
			log.Fatalln(err)
		}

		l := logrus.New()
		l.ReportCaller = true
		l.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				filename := path.Base(frame.File)
				return fmt.Sprintf("%s:%d", filename, frame.Line), filename
			},
			DisableColors: false,
			FullTimestamp: true,
		}

		l.SetOutput(os.Stdout)

		l.SetLevel(logrusLevel)

		instanse = Logger{logrus.NewEntry(l)}
	})

	return instanse
}
