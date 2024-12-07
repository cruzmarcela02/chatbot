package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/books/v1"
)

const (
	START                  = "/start"
	RECOMENDACION          = "/recomendacion"
	BUSQUEDA               = "/busqueda"
	HISTORIAL              = "/historial"
	GOOGLEBOOKS            = "/gbooks"
	INFORME                = "/informe"
	PERSONALIZACION        = "/personalizar"
	TITULO                 = "Titulo"
	AUTOR                  = "Autor"
	EDITORIAL              = "Editorial"
	GENERO                 = "Genero"
	RECOMENDACIONES        = "Mis recomendaciones"
	BUSQUEDAS              = "Mis Busquedas"
	TERMINAR               = "Terminar"
	ELIMINAR_FILTROS       = "Eliminar Filtros"
	BUSQUEDA_GLOBAL        = "Aplicando filtros globales"
	BUSQUEDA_PERSONALIZADA = "Sin aplicar filtros globales"
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
	bufferLibro   *books.Volume
}

var chatsConcurrencia = struct {
	sync.RWMutex
	m map[int64]chan tgbotapi.Update
}{m: make(map[int64]chan tgbotapi.Update)}

func (b *Bot) manejarComando(id int64, msg string) { // maneja los comandos historial, personalizacion e informe
	b.Recomendacion = false
	b.filtroGLobal = false
	enGoogleBooks := b.autenticado && b.ultimoComando == GOOGLEBOOKS
	filtrosGlobales, _ := formatearFiltros(id)

	switch msg {
	case RECOMENDACION:
		b.Recomendacion = true
		b.API.Send(CrearMenuFiltros(id))
		b.sendText(id, "Tus filtros globales actuales son: "+filtrosGlobales)
	case BUSQUEDA:
		if b.autenticado && b.ultimoComando == GOOGLEBOOKS {
			b.API.Send(crearMenu(BUSQUEDA, id))
		} else {
			b.sendText(id, "Tus filtros globales actuales son: "+filtrosGlobales)
			b.API.Send(CrearMenuFiltros(id))

		}
	case HISTORIAL:
		b.DarHistorial(id, enGoogleBooks)

	case GOOGLEBOOKS:
		b.interactuarGoogleBooks(id)

	case INFORME:
		b.API.Send(crearMenu(INFORME, id))

	case PERSONALIZACION:
		b.API.Send(crearMenu(PERSONALIZACION, id))
		b.filtroGLobal = true
		b.sendText(id, "filtros globales aplicados: "+filtrosGlobales)
	default:
		b.API.Send(crearMenu(START, id))
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

	if msg.Text == MENSUAL || msg.Text == SEMANAL || msg.Text == DIARIO {
		b.generarAnalisis(id, msg.Text)
		return
	}

	if msg.Text == BUSQUEDA_GLOBAL {
		b.BusquedaFiltroGlobal(msg)
		return
	}

	if msg.Text == BUSQUEDA_PERSONALIZADA {
		if b.Recomendacion {
			b.API.Send(crearMenu(RECOMENDACION, id))
			return
		}
		b.sendText(id, "No se aplicaran sus filtros globales")
		b.API.Send(crearMenu(BUSQUEDA, id))
		return
	}

	if msg.Text == ELIMINAR_FILTROS {
		err := eliminarFiltrosBD(id)
		if err != nil {
			b.sendText(id, "Error al eliminar los filtros globales")
		}
		removerMenu := RemoverMenu(msg.Chat.ID, "🧹Has eliminado todos tus filtros globales")
		b.API.Send(removerMenu)
		return
	}

	if msg.Text == RECOMENDACIONES || msg.Text == BUSQUEDAS {
		b.armarHistorial(msg, msg.Text)
		return
	}

	if msg.Text == TERMINAR {
		if !b.registroFiltros(id, msg.Text) && !b.filtroGLobal {
			b.obtenerLibroRandom(id, "subject:Novel")
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
			b.filwait = true
		}
		b.verificarFiltro(msg, msg.Text)
		return
	}

	if (msg.Text == FAVORITOS || msg.Text == POR_LEER || msg.Text == LEYENDO_AHORA || msg.Text == NO_AGREGAR || msg.Text == LEIDOSB) && b.autenticado {
		b.agregarLibro(msg.Chat.ID, msg.Text)
		return
	}

	if b.filwait {
		b.filtro += "\"" + msg.Text + "\""
		b.sendText(msg.Chat.ID, "¿Querés agregar otro filtro? ➡️ Agregalo\nSi no es así, apreta Terminar ✅")
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
		if data == RECOMENDACION || data == BUSQUEDA || data == HISTORIAL || data == GOOGLEBOOKS || data == PERSONALIZACION { // lo dejamos o lo hacemos menu adentro del teclado
			b.manejarComando(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		}

		if data == GOOGLEBOOKS {
			b.ultimoComando = GOOGLEBOOKS
		}

		if data == PARA_TI {
			b.recomendarParaTi(update.CallbackQuery.Message.Chat.ID)
			return
		}

		if data == LEIDOSH {
			b.mostrarHistorialGoogleBooks(update.CallbackQuery.Message.Chat.ID)
			return
		}
		log.Printf("Unknown callback data: %s", data) //
	}
}

func (b *Bot) manejarActualizaciones(updates chan tgbotapi.Update) {
	for update := range updates {
		if update.Message != nil {
			b.onUpdateReceived(update)
		}
		if update.CallbackQuery != nil {
			b.onCallbackQuery(update)
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
		var chatID int64
		if update.Message != nil {
			chatID = update.Message.Chat.ID
			log.Printf("Received message from %d", update.Message.Chat.ID)
		}
		if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
			log.Printf("Received callback query from %d", update.CallbackQuery.Message.Chat.ID)
		}
		chatsConcurrencia.Lock()
		if _, ok := chatsConcurrencia.m[chatID]; !ok {
			chatsConcurrencia.m[chatID] = make(chan tgbotapi.Update)
			go b.manejarActualizaciones(chatsConcurrencia.m[chatID])
		}
		chatsConcurrencia.Unlock()
		chatsConcurrencia.RLock()
		chatsConcurrencia.m[chatID] <- update
		chatsConcurrencia.RUnlock()
	}
}
