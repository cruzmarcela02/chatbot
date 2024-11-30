package main

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/books/v1"
)

/* Muestra el historial del usuario, gb o chat */
func (b *Bot) verHistorial(msg *tgbotapi.Message, filtro string, enGoogleBooks bool) {
	if enGoogleBooks {
		b.armarHistorialGoogleBooks(msg.Chat.ID, filtro)
		return
	}

	b.armarHistorial(msg, filtro)
}

// Historiales - GoogleBooks: Vistos Recientes o Leidos
func (b *Bot) armarHistorialGoogleBooks(id int64, estanteria string) {
	token, _ := b.obtenerTokenAlmacenado(id)
	service := autenticarCliente(b, id, token)

	if estanteria == VISTOS_RECIENTES {
		b.sendText(id, "Historial de tus vistos recientes")
		b.historialNavegacion(id, COD_VISTOS_RECIENTES, service)
		return
	}
	b.sendText(id, "Historial de tus leidos")
	b.historialNavegacion(id, COD_LEIDOS, service)
}

func autenticarCliente(b *Bot, id int64, token *oauth2.Token) *books.Service {
	client := b.OAuthConfig.Client(context.Background(), token)
	service, err := books.New(client)
	if err != nil {
		b.sendText(id, "Error al crear el cliente de Google Books: "+err.Error())
		return nil
	}
	return service
}

// Recently Viewed: 6
func (b *Bot) historialNavegacion(id int64, cod_estanteria string, service *books.Service) {
	var historial string

	// Chequeamos los Recientemente vistos
	bookshelf, err := service.Mylibrary.Bookshelves.Volumes.List(cod_estanteria).Do()
	if err != nil {
		b.sendText(id, "SURGIO UN ERROR: "+err.Error())
	}

	if len(bookshelf.Items) == 0 {
		b.sendText(id, "Estanteria vacia")
		return
	}

	for i, libro := range bookshelf.Items {
		historial += strconv.Itoa(i + 1)
		historial += ". "
		historial += libro.VolumeInfo.Title
		historial += "\n"
	}
	b.sendText(id, historial)
}

// Historial del chat. Fuera de google books
func (b *Bot) armarHistorial(msg *tgbotapi.Message, filtro string) {
	if filtro == RECOMENDACIONES {
		b.sendText(msg.Chat.ID, "Historial de tus recomendaciones")
		return
	}
	b.sendText(msg.Chat.ID, "Historial de tus busquedas")

	books, err := b.getSavedSearchResults(msg.Chat.ID)
	if err != nil {
		b.sendText(msg.Chat.ID, "Error al obtener el historial de búsquedas: "+err.Error())
		return
	}

	if len(books) == 0 {
		b.sendText(msg.Chat.ID, "No se encontraron resultados de búsqueda guardados.")
		return
	}

	b.sendText(msg.Chat.ID, "Historial de búsquedas:")
	for i, book := range books {
		b.sendText(msg.Chat.ID, fmt.Sprintf("%d. Titulo:%s , Link:%s ", i+1, book.Title, book.Link))
	}
}
