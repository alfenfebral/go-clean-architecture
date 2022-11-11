package logger

import "github.com/sirupsen/logrus"

func Error(err error) {
	logrus.Error(err)
}
