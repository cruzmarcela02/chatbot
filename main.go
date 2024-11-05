package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/books/v1"
)

const (
	RECOMENDACION   = "/recomendacion"
	BUSQUEDA        = "/busqueda"
	HISTORIAL       = "/historial"
	GOOGLEBOOKS     = "/googlebooks"
	PERSONALIZACION = "/personalización"
	TITULO          = "Titulo"
	AUTOR           = "Autor"
	EDITORIAL       = "Editorial"
	GENERO          = "Genero"
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

func (b *Bot) manejarComando(id int64, msg string) {

	switch msg {
	case RECOMENDACION:
		b.Recomendacion = true
		menuOpciones := crearMenu(RECOMENDACION, id)
		b.API.Send(menuOpciones)
		b.recomendar(msg, id)
	case BUSQUEDA:
		b.Busqueda = true
		//b.sendText(id, "Estas son las busquedas:")

		menuOpciones := crearMenu(BUSQUEDA, id)

		b.API.Send(menuOpciones)

		b.buscar(msg, id)
	case HISTORIAL:
		b.Historial = true
		b.verHistorial(msg, id)
	case GOOGLEBOOKS:

		// sus botones
		//b.interactuarGoogleBooks(msg, boton seleccionado )
	case PERSONALIZACION:
		// code if no case matches
	default:
		menu_opciones := crearMenu("/start", id)
		b.API.Send(menu_opciones)
		// informar que boton se toco
	}

}

// comandos
func (b *Bot) buscar(msg string, id int64) {
	//remover menu

	if msg == AUTOR {
		// llamar a la api buscando por autior y mostrar los resultados
	}
	if msg == EDITORIAL {
		// llamar a la api buscando por editorial y mostrar los resultados
	}
	if msg == TITULO {
		b.sendText(id, "Por favor ingrese el titulo del libro")
		// agarrar el titulo del libro y buscarlo

		// llamar a la api buscando por titulo y mostrar los resultados
		// b.sendText(id, resultados)
	}
	if msg == GENERO {
		// llamar a la api buscando por genero y mostrar los resultados
	}

}

func (b *Bot) recomendar(msg string, id int64) {
	b.sendText(id, "Estas son las recomendaciones: bu bu bu ")
	if msg == EDITORIAL {
		// llamar a la api buscando por editorial y mostrar los resultados
	}
	if msg == TITULO {
		b.sendText(id, "Por favor ingrese el titulo del libro")
		// agarrar el titulo del libro y buscarlo

		// llamar a la api buscando por titulo y mostrar los resultados
		// b.sendText(id, resultados)
	}
	if msg == GENERO {
		// llamar a la api buscando por genero y mostrar los resultados
	}
}
func (b *Bot) verHistorial(msg string, id int64) {
	b.sendText(id, "Este es tu historial")
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
		b.manejarComando(msg.Chat.ID, msg.Text)
	} else if msg.Text == AUTOR || msg.Text == EDITORIAL || msg.Text == GENERO {
		b.sendText(msg.Chat.ID, "presionste el boton de "+msg.Text)
		// si se toca un boton de autor o editorial -> puede ser recomendacion o busqueda -> ver como diferenciar

	} else if msg.Text == TITULO {
		b.buscar(msg.Text, msg.Chat.ID)
	} else {
		b.sendText(msg.Chat.ID, "No se reconoce el comando")
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

func (b *Bot) onCallbackQuery(update tgbotapi.Update) {
	callback := update.CallbackQuery

	if callback != nil {

		data := callback.Data
		// es un camnando -> estamos en el /start
		if data == RECOMENDACION || data == BUSQUEDA || data == HISTORIAL || data == GOOGLEBOOKS || data == PERSONALIZACION {

			b.manejarComando(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

		}

	}
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
		if update.Message != nil {
			b.onUpdateReceived(update)
		}
		if update.CallbackQuery != nil {
			b.onCallbackQuery(update)
		}

	}
}
