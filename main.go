package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("8040461009:AAGk-uZFfkIR5-mx50i7XmNVILwJseS6iPE")
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
} //TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
