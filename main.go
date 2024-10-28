package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/books/v1"
	"log"
	"net/http"
	"os"
)

var (
	// Menu texts
	firstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	nextButton     = "Next"
	backButton     = "Back"
	tutorialButton = "Tutorial"

	bot Bot

	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nextButton, nextButton),
		),
	)
)

type Bot struct {
	API           *tgbotapi.BotAPI
	Recomendacion bool // si se toca el comando /recomendacion o si el usuario ingresa recomendaciono alguna variante
	Busqueda      bool // si se toca el comando /busqueda o si el usuario ingresa busqueda o alguna variante
	Historial     bool // si se toca el comando /historial o si el usuario ingresa historial o alguna variante
	OAuthConfig   *oauth2.Config
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
	case "/googlebooks":
		// sus botones
		//b.interactuarGoogleBooks(msg, boton seleccionado )
	case "/personalización":
		// code if no case matches
	}

}

// comandos
func (b *Bot) buscar(msg *tgbotapi.Message) {
	b.sendText(msg.Chat.ID, "Estas son las busquedas:")
	menu_opciones := tgbotapi.NewMessage(msg.Chat.ID, firstMenu)
	menu_opciones.ReplyMarkup = firstMenuMarkup
	b.API.Send(menu_opciones)
}
func (b *Bot) recomendar(msg *tgbotapi.Message) {
	b.sendText(msg.Chat.ID, "Estas son las recomendaciones: bu bu bu ")
}
func (b *Bot) verHistorial(msg *tgbotapi.Message) {
	b.sendText(msg.Chat.ID, "Este es tu historial")
}

/*
	func (b *Bot) interactuarGoogleBooks(msg *tgbotapi.Message) {
		// implementar nuevos subcomandos
		// agregar
		// eliminar
		// ver
		// buscar

		// Si ya tenemos el token, buscar un libro de ejemplo

}
*/
func (b *Bot) onUpdateReceived(update tgbotapi.Update) { // lee los mensajes
	msg := update.Message
	//user := msg.From
	// si se toca algun comando -> llamar a manejarComando
	if msg.IsCommand() {
		b.manejarComando(msg)
	} else {
		//b.sendText(user.ID, "No entiendo tu mensaje. por favor usar los comandos")
		// usar nlp
		b.interactuarGoogleBooks(msg, update)

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
// your onUpdateReceived method will be called automatically and you'll be able to handle the update parameter,
// which contains the message, along with a great deal of other info which you can see detailed here.
// The user - Who sent the message. Access it via update.getMessage().getFrom().
// The message - What was sent. Access it via update.getMessage().
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

	// Configurar OAuth
	b.OAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{books.BooksScope},
		RedirectURL:  "https://relaxed-stunning-stag.ngrok-free.app/oauth2callback", // Cambiar cada vez que se levanta
		Endpoint:     google.Endpoint,
	}

	// Configurar la ruta para el callback de Google OAuth
	http.HandleFunc("/oauth2callback", b.handleGoogleCallback)

	// Iniciar el servidor HTTP en una goroutine separada
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		b.onUpdateReceived(update)
	}
}
