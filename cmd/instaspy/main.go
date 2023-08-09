package main

import (
	"fmt"
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
		logrus.Fatalf("DB connection failed at %s: %w", op, err)
	}
	//defer db.Close() - if uncomment - huge error stacktrace in terminal

	// Application core
	// Used for connecting to selenium image remotely
	conn, err := core.EstablishRemote()
	if err != nil {
		logrus.Fatalf("Failed connect to selenium at %s: %w", op, err)
	}
	defer conn.Quit()

	// Parse pictures and save them
	for _, username := range cfg.Usernames {
		pic, _, err := conn.Job(username)
		if err != nil {
			logger.HandleOpError(op, err)
			logrus.Fatal(err)
		}

		for _, image := range pic {
			fileInfo, err := save.Image(username, image, db)
			if err != nil {
				logger.HandleOpError(op, err)
				logrus.Fatal(err)
			}

			if fileInfo.Hash == "dont" {
				continue
			} else {
				res, _ := db.AddInfo(fileInfo)
				telegram.SendMessage(username, cfg.TelegramBotToken, cfg.ChatID, fileInfo.Picture_name)
				if res != 200 {
					fmt.Printf("%s: %s", op, err)
					os.Exit(1)
				}
			}
		}
	}
}
