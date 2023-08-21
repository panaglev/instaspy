package logger

import (
	"fmt"
	"instaspy/src/config"
	"instaspy/telegram"

	"github.com/sirupsen/logrus"
)

type OperationError struct {
	Op  string
	Err error
}

func HandleOpError(op string, err error) {
	logrus.Errorf("Error at %s: %s", op, err)
}

func HandleOpErrorWithComment(op string, err error, message string) {
	logrus.Errorf("Error at %s: %s\n %s", op, err, message)
}

func HandleOpErrorTelegramMessage(op string, err error) {
	cfg := config.MustLoad()
	preparedMessage := fmt.Sprintf("Hey Boss, I'm down at %s\n\n stacktrace:\n\n %s", op, err)
	telegram.SendMessage(preparedMessage, cfg.TelegramBotToken, cfg.ChatID)
}

func HandleOpTelegramMessage(op, message string) {
	cfg := config.MustLoad()
	preparedMessage := fmt.Sprintf("%s at %s", message, op)
	telegram.SendMessage(preparedMessage, cfg.TelegramBotToken, cfg.ChatID)
}
