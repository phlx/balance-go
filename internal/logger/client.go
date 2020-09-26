package logger

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func New(debug bool) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	if debug {
		logger.SetFormatter(&logrus.TextFormatter{})
	}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	}
	return logger
}

func NewForTest() *logrus.Logger {
	output := TestOutput{Written: &[]byte{}}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(output)
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

type TestOutput struct {
	io.Writer
	Written *[]byte
}

func (o TestOutput) Write(p []byte) (n int, err error) {
	*o.Written = append(*o.Written, p...)
	return len(p), nil
}

func (o TestOutput) Lines() []string {
	return strings.Split(string(*o.Written), "\n")
}
