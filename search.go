package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/books/v1"
)

const (
	FTITULO    = "intitle:"
	FAUTOR     = "inauthor:"
	FEDITORIAL = "inpublisher:"
	FGENERO    = "subject:"
)

func (b *Bot) obtenerLibroRandom(id int64, filtro string) {
	client := &http.Client{}
	service, err := books.New(client)
	if err != nil {
		return
	}

	// longitud random
	randomNumber := rand.Intn(10)
	call := service.Volumes.List(filtro).MaxResults(10)
	resp, err := call.Do()
	if err != nil {
		b.sendText(id, "Error al buscar libros: "+err.Error())
		return
	}
	if len(resp.Items) == 0 {
		b.sendText(id, "No se encontraron libros.")
		return
	}
	book := resp.Items[randomNumber]
	downloadLink := conseguirLink(book)
	titulo := book.VolumeInfo.Title
	b.sendText(id, fmt.Sprintf("El libro es %s.\nDescargalo en %s", titulo, downloadLink))
}

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
		Title:   titulo,
		Link:    downloadLink,
		Periodo: time.Now(),
	}
	b.saveSearchResult(BookBD, id)

}

func conseguirLink(firstBook *books.Volume) string {
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
			b.ultimoComando = RECOMENDACION
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
