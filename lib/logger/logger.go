package logger

import (
	"bytes"
	"os"

	"github.com/sirupsen/logrus"
)

// InitializeLog setup configuration and return logrus object
func init() {
	textFormatter := &logrus.TextFormatter{}
	logrus.SetFormatter(&PrefixTextFormatter{textFormatter, prefix})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(GetLogLevel())
}

func GetLogLevel() logrus.Level {
	switch logLevel := os.Getenv("LOG_LEVEL"); logLevel {
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

type PrefixTextFormatter struct {
	base   *logrus.TextFormatter
	prefix string
}

func (f *PrefixTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg, err := f.base.Format(entry)
	if err != nil {
		return msg, err
	}

	msg = bytes.Join([][]byte{[]byte(f.prefix), msg}, nil)

	return msg, err
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logrus.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Print(args ...interface{}) {
	logrus.Print(args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warning(args ...interface{}) {
	logrus.Warning(args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

func Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

func Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

func Println(args ...interface{}) {
	logrus.Println(args...)
}

func Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}

func Warningln(args ...interface{}) {
	logrus.Warningln(args...)
}

func Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	logrus.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	logrus.Panicln(args...)
}
