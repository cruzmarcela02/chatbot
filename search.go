package main

import (
	"fmt"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/books/v1"
)

const (
	FTITULO    = "intitle:"
	FAUTOR     = "inauthor:"
	FEDITORIAL = "inpublisher:"
	FGENERO    = "subject:"
)

func (b *Bot) buscarSinAuth(id int64, filtro string) {

	// Create a new HTTP client without OAuth authentication
	client := &http.Client{}
	service, err := books.New(client)
	if err != nil {
		b.sendText(id, "Error al crear el cliente de Google Books: "+err.Error())
		return
	}

	call := service.Volumes.List(filtro).MaxResults(3)
	resp, err := call.Do()
	if err != nil {
		b.sendText(id, "Error al buscar libros: "+err.Error())
		return
	}
	if len(resp.Items) == 0 {
		b.sendText(id, "No se encontraron libros.")
		return
	}
	book := resp.Items[0]
	completa := false
	for _, item := range resp.Items {
		if item.AccessInfo.AccessViewStatus != "NONE" && !completa {
			book = item
			completa = true
		}
	}
	downloadLink := conseguirLink(book)
	titulo := book.VolumeInfo.Title

	b.sendText(id, fmt.Sprintf("El libro encontrado es %s.Descargalo en %s", titulo, downloadLink))
	BookBD := BookBD{
		Title: titulo,
		Link:  downloadLink,
	}
	b.saveSearchResult(BookBD, id)

}

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

		b.sendText(msg.Chat.ID, recomendacion)
	}
}

func conseguirLink(firstBook *books.Volume) string {
	/*
		if firstBook.AccessInfo.Epub.IsAvailable != false {

			if firstBook.AccessInfo.Epub.DownloadLink != "" {
				// Get the file as epub
				return firstBook.AccessInfo.Epub.DownloadLink
			} else {
				return firstBook.AccessInfo.Epub.AcsTokenLink
			}

		} else if firstBook.AccessInfo.Pdf.IsAvailable != false {

			if firstBook.AccessInfo.Pdf.DownloadLink != "" {
				// Get the file as pdf
				return firstBook.AccessInfo.Pdf.DownloadLink
			} else {
				return firstBook.AccessInfo.Pdf.AcsTokenLink
			}

		} else {
			return firstBook.VolumeInfo.PreviewLink

		}*/

	if firstBook.AccessInfo.Epub.IsAvailable && firstBook.AccessInfo.Epub.DownloadLink != "" {
		return firstBook.AccessInfo.Epub.DownloadLink
	}
	if firstBook.AccessInfo.Pdf.IsAvailable && firstBook.AccessInfo.Pdf.DownloadLink != "" {
		return firstBook.AccessInfo.Pdf.DownloadLink
	}
	return firstBook.VolumeInfo.PreviewLink
}

func (b *Bot) realizarquery(msg *tgbotapi.Message) {

	if b.Recomendacion {
		if b.autenticado && b.ultimoComando == GOOGLEBOOKS {
			token, _ := b.obtenerTokenAlmacenado(msg.Chat.ID)
			b.recomendarParaTi(msg.Chat.ID, token)
		} else {
			b.recomendarLibros(msg, b.filtro)
		}

	} else {
		if b.autenticado && b.ultimoComando == GOOGLEBOOKS {
			token, _ := b.obtenerTokenAlmacenado(msg.Chat.ID)
			b.buscarlibro(b.filtro, msg.Chat.ID, token)
			b.ultimoComando = BUSQUEDA

		} else {
			b.buscarSinAuth(msg.Chat.ID, b.filtro)
		}
	}

	b.filwait = false
	b.filtro = ""
}
