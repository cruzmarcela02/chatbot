package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/books/v1"
	"log"
	"net/http"
	"os"
)

const (
	START           = "/start"
	RECOMENDACION   = "/recomendacion"
	BUSQUEDA        = "/busqueda"
	HISTORIAL       = "/historial"
	GOOGLEBOOKS     = "/gbooks"
	INFORME         = "/analisisbusqueda"
	PERSONALIZACION = "/personalización"
	TITULO          = "Titulo"
	AUTOR           = "Autor"
	EDITORIAL       = "Editorial"
	GENERO          = "Genero"
	RECOMENDACIONES = "Mis recomendaciones"
	BUSQUEDAS       = "Mis Busquedas"
	TERMINAR        = "Terminar"
	FAVORITOS       = "Favoritos"
	POR_LEER        = "Por Leer"
	LEYENDO_AHORA   = "Leyendo Ahora"
	NO_AGREGAR      = "No agregar"
)

type Bot struct {
	API           *tgbotapi.BotAPI
	Recomendacion bool // si se toca el comando /recomendacion o si el usuario ingresa recomendaciono alguna variante
	filtro        string
	filwait       bool
	OAuthConfig   *oauth2.Config
	autenticado   bool
}

/*func (b *Bot) getBotUsername() string {
	return b.API.Self.UserName
}

func (b *Bot) getBotToken() string {
	return b.API.Token
}*/

func (b *Bot) manejarComando(id int64, msg string) { // maneja los comandos historial, personalizacion e informe
	b.Recomendacion = false

	switch msg {
	case RECOMENDACION:
		b.Recomendacion = true
		b.API.Send(crearMenu(RECOMENDACION, id))
		// escuchar una nueva actualizacion
		// nueva actualizacion -> set filtro
		// escuchar nueva actualizacion
		// hacer recomendacion con el filtro
	case BUSQUEDA:
		b.API.Send(crearMenu(BUSQUEDA, id))

		// escuchar una nueva actualizacion
		// nueva actualizacion -> set filtro
		// escuchar nueva actualizacion
		// hacer busqueda con el filtro
	case HISTORIAL:
		b.API.Send(crearMenu(HISTORIAL, id))

	case GOOGLEBOOKS:
		b.interactuarGoogleBooks(id)

		// sus botones
		//b.interactuarGoogleBooks(msg, boton seleccionado )

	case INFORME:
		// realizar informe con todas las busquedas y las recomendaciones del ultimo mes

	case PERSONALIZACION:
		// editable
		// filtro global → si esta activo mostrar un menu
		//
		//desea usar su filtro global ?
		//
		//1) Usar filtro global → hace la busqueda agarrando el filtro de la BD
		//
		//2) busqueda personalizada → ingresar el campo a buscar (actual)
		//
		//si no existe muestra ‘
		//
		//Por favor ingrese el %s a buscar, campo

	default:
		b.API.Send(crearMenu(START, id))
		// informar que boton se toco
	}

}

// comandos
func (b *Bot) verificarFiltro(msg *tgbotapi.Message, filtro string) {
	switch filtro {
	case AUTOR:
		b.filtro += FAUTOR
	case EDITORIAL:
		b.filtro += FEDITORIAL
	case TITULO:
		b.filtro += FTITULO
	case GENERO:
		b.filtro += FGENERO
	}
	b.sendText(msg.Chat.ID, fmt.Sprintf("Por favor ingrese el %s a buscar", msg.Text))
}

// func (b *Bot) verHistorial(msg string, id int64) {
// 	b.sendText(id, "Este es tu historial")
// }

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
		return
	}

	if msg.Text == RECOMENDACIONES || msg.Text == BUSQUEDAS {
		removerMenu := RemoverMenu(msg.Chat.ID, "Queres ver el historial: "+msg.Text)
		b.API.Send(removerMenu)
		b.verHistorial(msg, msg.Text)
		return
	}

	// Busqueda de libro o recomendacion
	if msg.Text == TERMINAR {
		// Caso de marcar TERMINAR sin agregar ningun filtro
		if !b.filwait {
			removerMenu := RemoverMenu(msg.Chat.ID, "Se cancelo el proceso")
			b.API.Send(removerMenu)
			return
		}
		removerMenu := RemoverMenu(msg.Chat.ID, "Filtros ingresados con exito")
		b.API.Send(removerMenu)
		b.realizarbusqueda(msg)
		return
	}

	if msg.Text == AUTOR || msg.Text == EDITORIAL || msg.Text == GENERO || msg.Text == TITULO { //  -> hacer comando

		if !b.filwait {
			b.sendText(msg.Chat.ID, "Si no desea agregar otro filtro toque el boton cancelar, en caso contrario toque el boton de filtro")
			b.filwait = true
		}
		b.verificarFiltro(msg, msg.Text)
		return
	}

	if (msg.Text == FAVORITOS || msg.Text == POR_LEER || msg.Text == LEYENDO_AHORA || msg.Text == NO_AGREGAR) && b.autenticado {
		removerMenu := RemoverMenu(msg.Chat.ID, "Su opereacion se realizo con exito")
		b.API.Send(removerMenu)
		b.agregarLibro(msg.Chat.ID, msg.Text)
		return
	}

	if b.filwait {
		b.filtro += "\"" + msg.Text + "\" "
		return
	} else {
		b.sendText(msg.Chat.ID, "No se reconoce el comando, usar alguno de los comandos del menu")
		return
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
		if data == RECOMENDACION || data == BUSQUEDA || data == HISTORIAL || data == GOOGLEBOOKS || data == PERSONALIZACION { // lo dejamos o lo hacemos menu adentro del teclado
			b.manejarComando(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

		}

	}
}

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
		RedirectURL:  "https://relaxed-stunning-stag.ngrok-free.app/oauth2callback", // ngrok
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
