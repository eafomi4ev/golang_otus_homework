package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func New(ctx context.Context, level string, path string) *Logger {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(fmt.Errorf("could not open log file: %w", err))
	}

	logg := logrus.New()
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatal(fmt.Errorf("error on parsing logger level: %w", err))
	}
	logg.SetLevel(lvl)
	logg.SetOutput(io.MultiWriter(file, os.Stdout))
	// logg.SetFormatter(&logrus.JSONFormatter{})

	go func() {
		<-ctx.Done()
		if err := file.Close(); err != nil {
			log.Println("Cannot close the log file")
			return
		}
		fmt.Println("Log file has been closed")
	}()

	return &Logger{logger: logg}
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}
