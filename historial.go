package main

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/books/v1"
)

/* Envia las opciones de historiales posibles */
func (b *Bot) MostrarMenuHistorial(id int64, enGoogleBooks bool) {
	if enGoogleBooks {
		b.sendText(id, "Historial para gbooks")
		b.API.Send(crearMenu(HISTORIAL, id, enGoogleBooks))
		return
	}

	b.sendText(id, "Historial comun")
	b.API.Send(crearMenu(HISTORIAL, id, enGoogleBooks))
}

/* Muestra el historial del usuario, gb o chat */
func (b *Bot) verHistorial(msg *tgbotapi.Message, filtro string, enGoogleBooks bool) {
	if enGoogleBooks {
		b.mostrarHistorialGoogleBooks(msg.Chat.ID, filtro)
		return
	}

	b.armarHistorial(msg, filtro)
}

/* Historiales - GoogleBooks: Vistos Recientes o Leidos */
func (b *Bot) mostrarHistorialGoogleBooks(id int64, estanteria string) {
	token, _ := b.obtenerTokenAlmacenado(id)
	service := autenticarCliente(b, id, token)

	if estanteria == VISTOS_RECIENTES {
		historial := armarHistorialGoogleBooks(COD_VISTOS_RECIENTES, service)
		b.sendText(id, "Historial de tus vistos recientes")
		b.sendText(id, historial)

	} else {
		historial := armarHistorialGoogleBooks(COD_LEIDOS, service)
		b.sendText(id, "Historial de tus leidos")
		b.sendText(id, historial)

	}
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

/* Retorna un string con el historial armado */
func armarHistorialGoogleBooks(cod_estanteria string, service *books.Service) string {
	var historial string

	bookshelf, err := service.Mylibrary.Bookshelves.Volumes.List(cod_estanteria).Do()
	if err != nil {
		historial = "SURGIO UN ERROR: " + err.Error()
		return historial
	}

	if len(bookshelf.Items) == 0 {
		historial = "Estanteria vacia"
		return historial
	}

	for i, libro := range bookshelf.Items {
		historial += strconv.Itoa(i + 1)
		historial += ". "
		historial += libro.VolumeInfo.Title
		historial += "\n"
	}

	return historial
}

// Historial del chat. Fuera de google books
func (b *Bot) armarHistorial(msg *tgbotapi.Message, filtro string) {
	if filtro == RECOMENDACIONES {
		b.sendText(msg.Chat.ID, "Historial de tus recomendaciones")
		recomendaciones, err := b.obtenerRecomendaciones(msg.Chat.ID)
		if err != nil {

		}
		for i, recomendacion := range recomendaciones {
			b.sendText(msg.Chat.ID, fmt.Sprintf("%d. Titulo:%s\n Link:%s ", i+1, recomendacion.Title, recomendacion.Link))
		}
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
