package main

import (
	"fmt"
	"instaspy/core"
	"instaspy/src/config"
	"instaspy/src/save"
	sqlite "instaspy/src/storage"
	"os"
)

func main() {
	const op = "cmd.instaparser.main"

	cfg := config.MustLoad()

	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		fmt.Printf("error connection to db %s", err)
	}

	conn, err := core.EstablishRemote()
	if err != nil {
		fmt.Printf("Error connecting to WebDriver at %s: %s", op, err)
		os.Exit(1)
	}
	defer conn.Quit()

	for _, username := range cfg.Usernames {
		pic, _, err := conn.Job(username)
		if err != nil {
			fmt.Printf("Problem during parse job at %s: %s", op, err)
			os.Exit(1)
		}

		for _, image := range pic {
			fileInfo, err := save.Image(username, image, db)
			if err != nil {
				fmt.Println(err)
			}

			if fileInfo.Hash == "dont" {
				continue
			} else {
				res, _ := db.AddInfo(fileInfo)
				if res != 200 {
					fmt.Printf("%s: %s", op, err)
				}
			}

		}
	}
}
