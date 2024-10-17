package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	API           *tgbotapi.BotAPI
	Recomendacion bool // si se toca el comando /recomendacion o si el usuario ingresa recomendaciono alguna variante
	Busqueda      bool // si se toca el comando /busqueda o si el usuario ingresa busqueda o alguna variante
	Historial     bool // si se toca el comando /historial o si el usuario ingresa historial o alguna variante
}

func (b *Bot) getBotUsername() string {
	return b.API.Self.UserName
}

func (b *Bot) getBotToken() string {
	return b.API.Token
}

func (b *Bot) manejarComando(msg *tgbotapi.Message) {
	switch msg.Text {
	case "/start":
		_ = b.sendText(msg.Chat.ID, "mostrar un menu con todos los demas comandos")
	case "/recomendacion":
		b.Recomendacion = true
		b.recomendar(msg)
	case "/busqueda":
		b.Busqueda = true
		b.buscar(msg)
	case "/historial":
		b.Historial = true
		b.verHistorial(msg)
	default:
		// code if no case matches
	}

}
func (b *Bot) buscar(msg *tgbotapi.Message) {
	b.sendText(msg.Chat.ID, "Estas son las busquedas:")
}
func (b *Bot) recomendar(msg *tgbotapi.Message) {
	b.sendText(msg.Chat.ID, "Estas son las recomendaciones: bu bu bu ")
}
func (b *Bot) verHistorial(msg *tgbotapi.Message) {
	b.sendText(msg.Chat.ID, "Este es tu historial")
}

func (b *Bot) onUpdateReceived(update tgbotapi.Update) { // lee los mensajes
	msg := update.Message
	user := msg.From
	// si se toca algun comando -> llamar a manejarComando
	if msg.IsCommand() {
		b.manejarComando(msg)
	} else {
		b.sendText(user.ID, "No entiendo tu mensaje. por favor usar los comandos")
		// llamar a "/start"
	}

}

func (b *Bot) sendText(who int64, what string) error {

	msg := tgbotapi.NewMessage(who, what)
	send, err := b.API.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %s", err)
	}

	log.Printf("Sent message to %s", send.Chat.FirstName)
	return err
}

// Every time someone sends a private message to your bot,
//your onUpdateReceived method will be called automatically and you'll be able to handle the update parameter,
//which contains the message, along with a great deal of other info which you can see detailed here.
// The user - Who sent the message. Access it via update.getMessage().getFrom().
//The message - What was sent. Access it via update.getMessage().

func main() {
	bot, err := tgbotapi.NewBotAPI("8040461009:AAGk-uZFfkIR5-mX5OI7XmNVIlwJseS6iPE")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	b := &Bot{API: bot}
	id := int64(0)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		b.onUpdateReceived(update)
		usr := update.Message.From
		id = usr.ID
		log.Printf("id is  %d", id)
		err = b.sendText(id, "Hello, World!")
	}

}

// paraa mandar mensajes necesitamos el userID
