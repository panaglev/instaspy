package telegram

import (
	"fmt"
	"io/ioutil"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func SendPicture(username, filename, telegram_token, chat_id string) error {
	const op = "telegram.SendPicture"

	bot, err := tgbotapi.NewBotAPI(telegram_token)
	if err != nil {
		logrus.Errorf("Error at %s: %s", op, err)
		return err
	}

	chatID, err := strconv.Atoi(chat_id)
	if err != nil {
		logrus.Errorf("Error at %s: %s", op, err)
		return err
	}

	filePath := fmt.Sprintf("/exbestfriend/images/%s/%s.jpg", username, filename)
	data, _ := ioutil.ReadFile(filePath)
	b := tgbotapi.FileBytes{Name: "picture", Bytes: data}
	message := tgbotapi.NewPhoto(int64(chatID), b)
	message.Caption = fmt.Sprintf("inst: %s", username)

	_, err = bot.Send(message)
	if err != nil {
		logrus.Errorf("Error at %s: %s", op, err)
		return err
	}
	return nil
}

func SendMessage(message, telegram_token, chat_id string) error {
	const op = "telegram.SendMessage"

	bot, err := tgbotapi.NewBotAPI(telegram_token)
	if err != nil {
		logrus.Errorf("Error at %s: %s", op, err)
		return err
	}

	chatID, err := strconv.Atoi(chat_id)
	if err != nil {
		logrus.Errorf("Error at %s: %s", op, err)
		return err
	}

	msg := tgbotapi.NewMessage(int64(chatID), message)

	_, err = bot.Send(msg)
	return err
}
