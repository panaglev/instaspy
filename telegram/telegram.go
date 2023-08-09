package telegram

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(username, telegram_token, chat_id string, filename int) {
	const op = "telegram.SendMessage"

	bot, err := tgbotapi.NewBotAPI(telegram_token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	chatID, err := strconv.Atoi(chat_id)
	if err != nil {
		fmt.Println("Не удалось преобразовать чатайди к инту")
		panic(1)
	}

	filePath := fmt.Sprintf("/Users/panaglev/Desktop/golang/instaspy/images/%s/%d.jpg", username, filename)
	data, _ := ioutil.ReadFile(filePath)
	b := tgbotapi.FileBytes{Name: "picture", Bytes: data}
	message := tgbotapi.NewPhoto(int64(chatID), b)
	message.Caption = fmt.Sprintf("inst: %s", username)

	_, err = bot.Send(message)
	if err != nil {
		log.Panic(err)
	}
}
