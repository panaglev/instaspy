package logger

import "github.com/sirupsen/logrus"

type OperationError struct {
	Op  string
	Err error
}

func HandleOpError(op string, err error) {
	logrus.Errorf("Error at %s: %s", op, err)
}
