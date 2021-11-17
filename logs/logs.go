package logs

import (
	"github.com/sirupsen/logrus"
)

type ILogger interface {
	Error(msg map[string]interface{}, tips string)
	Debug(msg map[string]interface{}, tips string)
	Info(msg map[string]interface{}, tips string)
	Warn(msg map[string]interface{}, tips string)
	Fatal(msg map[string]interface{}, tips string)
}

type Logger struct {
	Module string
	log    *logrus.Logger
}

func NewLogger(module string) *Logger {
	return &Logger{
		Module: module,
		log:    logrus.New(),
	}
}

func (n *Logger) Error(msg map[string]interface{}, tips string) {
	n.log.WithFields(newFields(msg)).Error(tips)
}
func (n *Logger) Debug(msg map[string]interface{}, tips string) {
	n.log.WithFields(newFields(msg)).Debug(tips)
}

func (n *Logger) Info(msg map[string]interface{}, tips string) {
	n.log.WithFields(newFields(msg)).Info(tips)
}

func (n *Logger) Warn(msg map[string]interface{}, tips string) {
	n.log.WithFields(newFields(msg)).Warn(tips)
}

func (n *Logger) Fatal(msg map[string]interface{}, tips string) {
	n.log.WithFields(newFields(msg)).Fatal(tips)
}

func newFields(msg map[string]interface{}) logrus.Fields {
	result := make(logrus.Fields)
	for k, v := range msg {
		result[k] = v
	}
	return result
}
