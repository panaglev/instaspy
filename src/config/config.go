package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Usernames        []string `yaml:"usernames" env-required:"true"`
	StoragePath      string   `yaml:"storage_path" env-required:"true"`
	DownloadPath     string   `yaml:"download_path" env-required:"true"`
	TelegramBotToken string   `yaml:"telegram_bot"`
	ChatID           string   `yaml:"chat_id"`
}

func MustLoad() *Config {
	const op = "pkg.config.MustLoad"

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH doesn't set at %s", op)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file doesn't exist at %s", op)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config at %s: %s", op, err)
	}

	cfg.TelegramBotToken = os.Getenv("TELEGRAM_BOT")
	cfg.ChatID = os.Getenv("CHAT_ID")

	return &cfg
}
