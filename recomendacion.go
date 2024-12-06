package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/books/v1"
)

func (b *Bot) Recomendar(id int64, enGoogleBooks bool) {
	if enGoogleBooks {
		token, _ := b.obtenerTokenAlmacenado(id)
		b.recomendarParaTi(id, token)
		return
	}

	b.API.Send(crearMenu(RECOMENDACION, id, enGoogleBooks))
}

/* Recomendados de la estanteria 'Para ti' de GoogleBooks */
func (b *Bot) recomendarParaTi(id int64, token *oauth2.Token) {
	client := b.OAuthConfig.Client(context.Background(), token)
	service, err := books.New(client)
	if err != nil {
		b.sendText(id, "Error al crear el cliente de Google Books: "+err.Error())
		return
	}

	bookShelves, _ := service.Mylibrary.Bookshelves.Volumes.List(COD_PARA_TI).Do()

	if len(bookShelves.Items) == 0 {
		b.sendText(id, "Ups, no hay recomendaciones para vos 😧😧!\n Te recomiendo interactuar un poco mas en con /googlebooks asi podemos descubrir un poco mas tus gustos ✨")
		return
	}

	b.sendText(id, "Veamos que te puedo recomendar de la estanteria Para Ti 📚")

	for i, item := range bookShelves.Items {
		var recomendacion string
		recomendacion += strconv.Itoa(i + 1)
		recomendacion += ". "
		recomendacion += item.VolumeInfo.Title
		recomendacion += "\n"
		recomendacion += "Género: "
		recomendacion += strings.Join(item.VolumeInfo.Categories, ", ")

		b.sendText(id, recomendacion)
	}
}

/* Recomendacion en base a lo que el cliente pide mediante filtros*/
func (b *Bot) recomendarLibros(msg *tgbotapi.Message, filtro string) {
	client := &http.Client{}
	service, err := books.New(client)
	if err != nil {
		b.sendText(msg.Chat.ID, "Error al crear el cliente de Google Books: "+err.Error())
		return
	}
	call := service.Volumes.List(filtro).MaxResults(3)
	resp, err := call.Do()

	if err != nil {
		b.sendText(msg.Chat.ID, "Error al buscar libros: "+err.Error())
		return
	}

	if len(resp.Items) == 0 {
		b.sendText(msg.Chat.ID, "No se encontraron libros.")
		return
	}

	b.sendText(msg.Chat.ID, "Te recomendamos los siguientes libros: ")

	for i, libro := range resp.Items {
		var recomendacion string
		recomendacion += strconv.Itoa(i + 1)
		recomendacion += ". "
		recomendacion += libro.VolumeInfo.Title
		recomendacion += "\n"
		recomendacion += libro.VolumeInfo.Description
		recomendacion += "\n"
		recomendacion += libro.VolumeInfo.InfoLink

		b.sendText(msg.Chat.ID, recomendacion)
		downloadLink := conseguirLink(libro)
		BookBD := BookBD{
			Title:   libro.VolumeInfo.Title,
			Link:    downloadLink,
			Periodo: time.Now(),
		}
		b.guardarRecomendaciones(BookBD, msg.Chat.ID)
	}

}
