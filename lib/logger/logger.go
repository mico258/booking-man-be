package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Fields logrus.Fields

type Logger interface {
	WithFields(fields Fields) *logrus.Entry
	WithField(key string, val interface{}) *logrus.Entry
	WithError(err error) *logrus.Entry
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
	SetOut(out io.Writer)
}

type logger struct {
	base *logrus.Logger
}

func (l *logger) WithFields(fields Fields) *logrus.Entry {
	return l.base.WithFields(logrus.Fields(fields))
}

func (l *logger) WithField(key string, val interface{}) *logrus.Entry {
	return l.base.WithField(key, val)
}

func (l *logger) WithError(err error) *logrus.Entry {
	return l.base.WithError(err)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.base.Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.base.Infof(format, args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.base.Printf(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.base.Warnf(format, args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.base.Warningf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.base.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.base.Fatalf(format, args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.base.Panicf(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.base.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.base.Info(args...)
}

func (l *logger) Print(args ...interface{}) {
	l.base.Print(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.base.Warn(args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.base.Warning(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.base.Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.base.Fatal(args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.base.Panic(args...)
}

func (l *logger) Debugln(args ...interface{}) {
	l.base.Debugln(args...)
}

func (l *logger) Infoln(args ...interface{}) {
	l.base.Infoln(args...)
}

func (l *logger) Println(args ...interface{}) {
	l.base.Println(args...)
}

func (l *logger) Warnln(args ...interface{}) {
	l.base.Warnln(args...)
}

func (l *logger) Warningln(args ...interface{}) {
	l.base.Warningln(args...)
}

func (l *logger) Errorln(args ...interface{}) {
	l.base.Errorln(args...)
}

func (l *logger) Fatalln(args ...interface{}) {
	l.base.Fatalln(args...)
}

func (l *logger) Panicln(args ...interface{}) {
	l.base.Panicln(args...)
}

func (l *logger) SetOut(out io.Writer) {
	l.base.Out = out
}

func New() Logger {
	textFormatter := &logrus.TextFormatter{}
	base := &logrus.Logger{
		Formatter: &PrefixTextFormatter{textFormatter, prefix},
		Out:       os.Stdout,
		Level:     GetLogLevel(),
	}
	return &logger{base}
}
