package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/books/v1"
	"net/http"
)

const (
	FTITULO    = "intitle:"
	FAUTOR     = "inauthor:"
	FEDITORIAL = "inpublisher:"
	FGENERO    = "subject:"
)

func (b *Bot) buscarSinAuth(msg *tgbotapi.Message, filtro string) {

	// Create a new HTTP client without OAuth authentication
	client := &http.Client{}
	service, err := books.New(client)
	if err != nil {
		b.sendText(msg.Chat.ID, "Error al crear el cliente de Google Books: "+err.Error())
		return
	}

	//book := armarQuery(filtro, busqueda)
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

	//campo := getCampo(filtros, firstBook)

	// pasar a minuscula
	//busqueda = strings.ToLower(busqueda)
	//campo = strings.ToLower(campo)

	// transformar el download link a un archivo

	// hay veces que el campo esta mal ingresado

	b.sendText(msg.Chat.ID, fmt.Sprintf("El libro encontrado es %s.Descargalo en %s", titulo, downloadLink))

	/*if strings.Contains(campo, busqueda) {
		b.sendText(msg.Chat.ID, fmt.Sprintf("El libro encontrado es %s.Descargalo en %s", titulo, downloadLink))
		// ver de conseguir con download

	} else {
		b.sendText(msg.Chat.ID, "No se encontro el libro en el campo especificado. Verifica por errores ortograficos o de tipeo")
		b.sendText(msg.Chat.ID, fmt.Sprintf("el primer libro encontrado fue %s. Descargalo en %s", titulo, downloadLink))
		// mandar libro en formato epub
	}*/

}

/*
func getCampo(filtro []string, firstbook *books.Volume) string {



	switch filtro {
	case FTITULO:
		return firstbook.VolumeInfo.Title
	case FAUTOR:
		return firstbook.VolumeInfo.Authors[0] // ver de no siempre agarrrar el primero
	case FEDITORIAL:
		// escribir por consola
		fmt.Printf("Editorial: %s\n", firstbook.VolumeInfo.Publisher)

		return firstbook.VolumeInfo.Publisher
	case FGENERO:
		return firstbook.VolumeInfo.Categories[0] // ver de no siempre agarrrar el primero
	default:
		return "filtro no valido"
	}
}*/

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

func (b *Bot) realizarbusqueda(msg *tgbotapi.Message) {

	if b.Recomendacion {
	} else {
		b.buscarSinAuth(msg, b.filtro)
	}
	b.filtro = ""
	return

}
