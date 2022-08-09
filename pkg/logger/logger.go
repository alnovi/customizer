package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(os.Stdout)
}

func SetFormatter(format logrus.Formatter) {
	log.SetFormatter(format)
}

func SetLevel(laval logrus.Level) {
	log.SetLevel(laval)
}

func SetOutput(output io.Writer) {
	log.SetOutput(output)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return log.WithError(err)
}

func Logf(level logrus.Level, format string, args ...interface{}) {
	log.Logf(level, format, args)
}

func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args)
}

func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args)
}

func Log(level logrus.Level, args ...interface{}) {
	log.Log(level, args)
}

func Trace(args ...interface{}) {
	log.Trace(args)
}

func Debug(args ...interface{}) {
	log.Debug(args)
}

func Info(args ...interface{}) {
	log.Info(args)
}

func Print(args ...interface{}) {
	log.Print(args)
}

func Warn(args ...interface{}) {
	log.Warn(args)
}

func Warning(args ...interface{}) {
	log.Warning(args)
}

func Error(args ...interface{}) {
	log.Error(args)
}

func Fatal(args ...interface{}) {
	log.Fatal(args)
}

func Panic(args ...interface{}) {
	log.Panic(args)
}

func Exit(code int) {
	log.Exit(code)
}
