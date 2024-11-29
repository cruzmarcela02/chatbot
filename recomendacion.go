package main

import (
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/api/books/v1"
	"strconv"
	"strings"
)

const (
	COD_PARA_TI    = "8"
	LIBROS_PARA_TI = "Libros para ti"
)

func (b *Bot) recomendarParaTi(id int64, token *oauth2.Token) {
	client := b.OAuthConfig.Client(context.Background(), token)
	service, err := books.New(client)
	if err != nil {
		b.sendText(id, "Error al crear el cliente de Google Books: "+err.Error())
		return
	}

	bookShelves, _ := service.Mylibrary.Bookshelves.Volumes.List(COD_PARA_TI).Do()

	if len(bookShelves.Items) == 0 {
		b.sendText(id, "No hay libros para recomendar")
		return
	}

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
