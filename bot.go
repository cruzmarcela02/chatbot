package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	API *tgbotapi.BotAPI
}

func (b *Bot) getBotUsername() string {
	return b.API.Self.UserName
}

func (b *Bot) getBotToken() string {
	return b.API.Token
}

func (b *Bot) onUpdateReceived(update tgbotapi.Update) {
	// Handle the update
}
