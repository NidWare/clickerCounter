package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var adminIDs = map[int64]bool{
	180247888:  true,
	5908745143: true,
	462450231:  true,
	89174232:   true,
}

func StartTelegramBot(ctx context.Context, db *sql.DB) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	for {
		select {
		case update := <-updates:
			handleUpdate(bot, db, update)
		case <-ctx.Done():
			return
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, db *sql.DB, update tgbotapi.Update) {
	if !adminIDs[update.Message.Chat.ID] {
		return
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter key and URL separated by a space.")
			bot.Send(msg)
		}
	} else {
		params := strings.Fields(update.Message.Text)
		if len(params) == 2 {
			key, url := params[0], params[1]
			_, err := db.Exec("INSERT INTO redirects VALUES (?, ?)", key, url)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Insertion failed"))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Insertion successful"))
			}
		}
	}
}
