package main

import (
	"instaspy/core"
	"instaspy/src/config"
	"instaspy/src/logger"
	"instaspy/src/save"
	sqlite "instaspy/src/storage"
	"instaspy/telegram"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	const op = "cmd.instaparser.main"

	cfg := config.MustLoad()

	// Setting logger up to write in file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Set output logger file
		logrus.SetOutput(logFile)
	}
	defer logFile.Close()

	// Label application start
	logrus.Info("Starting ExBestFriend at %s", time.Now())

	// Connection to database
	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logrus.Fatalf("DB connection failed at %s: %s", op, err)
	}
	//defer db.Close() - if uncomment - huge error stacktrace in terminal

	// Application core
	// Used for connecting to selenium image remotely
	conn, err := core.EstablishRemote()
	if err != nil {
		logrus.Fatalf("Failed connect to selenium at %s: %s", op, err)
	}
	defer conn.Quit()

	for _, username := range cfg.Usernames {
		pic, _, err := conn.Job(username)
		if err != nil {
			logger.HandleOpError(op, err)
			// If parse attempt not successfull -> continue or repeat?
			continue
		}

		for _, image := range pic {
			fileInfo, err := save.Image(username, image, db)
			if err != nil {
				logger.HandleOpError(op, err)
				continue
			}

			if fileInfo.Hash == "" {
				continue
			} else {
				err = db.AddInfo(fileInfo)
				if err != nil {
					logger.HandleOpError(op, err)
					// just realized that if I already downloaded image and not added info about it might have a copy
					// I should fix logic I guess...
					break
				}
				err = telegram.SendPicture(username, fileInfo.Picture_name, cfg.TelegramBotToken, cfg.ChatID)
				if err != nil {
					logger.HandleOpError(op, err)
				}
			}
		}
	}
}
