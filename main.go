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

const (
	START                 = "/start"
	RECOMENDACION         = "/recomendacion"
	BUSQUEDA              = "/busqueda"
	HISTORIAL             = "/historial"
	GOOGLEBOOKS           = "/gbooks"
	INFORME               = "/analisisbusqueda"
	PERSONALIZACION       = "/personalizar"
	TITULO                = "Titulo"
	AUTOR                 = "Autor"
	EDITORIAL             = "Editorial"
	GENERO                = "Genero"
	RECOMENDACIONES       = "Mis recomendaciones"
	BUSQUEDAS             = "Mis Busquedas"
	TERMINAR              = "Terminar"
	FAVORITOS             = "Favoritos"
	PORLEER               = "Por Leer"
	LEYENDOAHORA          = "Leyendo Ahora"
	NOAGREGAR             = "No agregar"
	ELIMINARFILTROS       = "Eliminar Filtros"
	BUSQUEDAGLOBAL        = "Aplicando filtros globales"
	BUSQUEDAPERSONALIZADA = "Sin aplicar filtros globales"
)

type Bot struct {
	API           *tgbotapi.BotAPI
	Recomendacion bool // si se toca el comando /recomendacion o si el usuario ingresa recomendaciono alguna variante
	filtro        string
	filwait       bool
	OAuthConfig   *oauth2.Config
	autenticado   bool
	filtroGLobal  bool
	ultimoComando string
}

func (b *Bot) manejarComando(id int64, msg string) { // maneja los comandos historial, personalizacion e informe
	b.Recomendacion = false
	b.filtroGLobal = false
	switch msg {
	case RECOMENDACION:
		b.Recomendacion = true
		if b.autenticado {
			token, _ := b.obtenerTokenAlmacenado(id)
			b.recomendarParaTi(id, token)
			return
		}
		b.API.Send(crearMenu(RECOMENDACION, id))
	case BUSQUEDA:
		if b.autenticado && b.ultimoComando == GOOGLEBOOKS {
			b.API.Send(crearMenu(BUSQUEDA, id))
		} else {
			b.API.Send(CrearMenuFiltros(id))
		}
	case HISTORIAL:
		b.API.Send(crearMenu(HISTORIAL, id))

	case GOOGLEBOOKS:
		b.interactuarGoogleBooks(id)
	case INFORME:
		// realizar informe con todas las busquedas y las recomendaciones del ultimo mes

	case PERSONALIZACION:
		b.API.Send(crearMenu(PERSONALIZACION, id))
		b.filtroGLobal = true
	default:
		b.API.Send(crearMenu(START, id))
		// informar que boton se toco
	}

}

// comandos
func (b *Bot) onUpdateReceived(update tgbotapi.Update) { // lee los mensajes
	msg := update.Message
	id := msg.Chat.ID
	if msg.IsCommand() {
		b.manejarComando(id, msg.Text)
		b.ultimoComando = msg.Text
		return
	}
	if msg.Text == BUSQUEDAGLOBAL {
		b.BusquedaFiltroGlobal(msg)
		return
	}

	if msg.Text == BUSQUEDAPERSONALIZADA {
		b.sendText(id, "No se aplicaran sus filtros globales")
		b.API.Send(crearMenu(BUSQUEDA, id))
		return
	}

	if msg.Text == ELIMINARFILTROS {
		eliminarFiltrosBD(id)
		b.sendText(id, "Todos sus filtros globales han sido eliminados")
	}

	if msg.Text == RECOMENDACIONES || msg.Text == BUSQUEDAS {
		removerMenu := RemoverMenu(id, "Queres ver el historial: "+msg.Text)
		b.API.Send(removerMenu)
		b.verHistorial(msg, msg.Text)
		return
	}

	if msg.Text == TERMINAR {
		if !b.registroFiltros(id, msg.Text) {
			// Caso de marcar TERMINAR sin agregar ningun filtro
			return
		}
		if b.filtroGLobal {
			guardarFiltroGlobal(id, b.filtro)
		} else {
			b.realizarquery(msg)
		}
		return
	}

	if msg.Text == AUTOR || msg.Text == EDITORIAL || msg.Text == GENERO || msg.Text == TITULO {
		if !b.filwait {
			b.sendText(msg.Chat.ID, "Si no desea agregar otro filtro toque el boton cancelar, en caso contrario toque el boton de filtro")
			b.filwait = true
		}
		b.verificarFiltro(msg, msg.Text)
		return
	}

	if (msg.Text == FAVORITOS || msg.Text == PORLEER || msg.Text == LEYENDOAHORA || msg.Text == NOAGREGAR) && b.autenticado {
		removerMenu := RemoverMenu(msg.Chat.ID, "Su opereacion se realizo con exito")
		b.API.Send(removerMenu)
		b.agregarLibro(msg.Chat.ID, msg.Text)
		return
	}

	if b.filwait {
		b.filtro += "\"" + msg.Text + "\" "
		return
	} else {
		b.sendText(msg.Chat.ID, "No se reconoce el comando, usar alguno de los comandos del menu"+msg.Text)
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
