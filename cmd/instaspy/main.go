package main

import (
	"fmt"
	"instaspy/core"
	"instaspy/pkg/config"
	"instaspy/pkg/save"
	"os"
)

func main() {
	const op = "cmd.instaparser.main"

	cfg := config.MustLoad()

	conn, err := core.EstablishRemote()
	if err != nil {
		fmt.Printf("Error connecting to WebDriver at %s: %s", op, err)
		os.Exit(1)
	}
	defer conn.Quit()

	for _, username := range cfg.Usernames {
		res, err := conn.Job(username)
		if err != nil {
			fmt.Printf("Problem during parse job at %s: %s", op, err)
			os.Exit(1)
		}

		for _, image := range res {
			save.Image(username, image)
		}
	}
}
