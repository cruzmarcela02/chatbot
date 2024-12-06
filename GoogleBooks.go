package main

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/books/v1"
)

const (
	COD_FAVORITOS        = "0"
	COD_LEER             = "2"
	COD_LEYENDO          = "3"
	COD_LEIDOS           = "4"
	COD_VISTOS_RECIENTES = "6"
	COD_PARA_TI          = "8"

	FAVORITOS        = "Favoritos"
	POR_LEER         = "Por Leer"
	LEYENDO_AHORA    = "Leyendo Ahora"
	LEIDOS           = "Leidos"
	LEIDOSB          = "leidos"
	VISTOS_RECIENTES = "Vistos Recientes"
	NO_AGREGAR       = "No agregar"
)

func (b *Bot) interactuarGoogleBooks(id int64) {
	// Verificar si el usuario está autenticado
	token, err := b.obtenerTokenAlmacenado(id)
	if err != nil || token.AccessToken == "" {
		b.GoogleBooksAuth(id)
	} else {
		b.autenticado = true
		b.API.Send(crearMenu(GOOGLEBOOKS, id, false))
	}

}

func (b *Bot) buscarlibro(filtro string, id int64, token *oauth2.Token) {
	client := b.OAuthConfig.Client(context.Background(), token)
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

	b.sendText(id, fmt.Sprintf("El libro encontrado es %s.\nDescargalo en %s", titulo, downloadLink))

	b.bufferLibro = book
	b.API.Send(CrearMenuAgregar(id))
}

func (b *Bot) agregarLibro(id int64, estanteria string) {
	token, err := b.obtenerTokenAlmacenado(id)
	if err != nil {
		b.sendText(id, "Error al obtener el token almacenado: "+err.Error())
		return
	}

	client := b.OAuthConfig.Client(context.Background(), token)
	service, _ := books.New(client)
	libro := b.bufferLibro

	var mensaje string
	switch estanteria {
	case FAVORITOS:
		service.Mylibrary.Bookshelves.AddVolume(COD_FAVORITOS, libro.Id).Do()
		mensaje = fmt.Sprintf("%s esta en tus Favoritos, fijate si luego le pegas una releida 👀.", libro.VolumeInfo.Title)

	case POR_LEER:
		service.Mylibrary.Bookshelves.AddVolume(COD_LEER, libro.Id).Do()
		mensaje = fmt.Sprintf("Uh ahora %s esta en Por Leer! no cuelgues y leelo 🤨😒.", libro.VolumeInfo.Title)

	case LEYENDO_AHORA:
		service.Mylibrary.Bookshelves.AddVolume(COD_LEYENDO, libro.Id).Do()
		mensaje = fmt.Sprintf("Asi que estas Leyendo %s 😦??\n Si te llega a gustar muchos fijate de agregarlo a favoritos mas tarde 😌", libro.VolumeInfo.Title)

	case LEIDOSB:
		service.Mylibrary.Bookshelves.AddVolume(COD_LEIDOS, libro.Id).Do()
		mensaje = fmt.Sprintf("Genial! Agregaste %s a tu coleccion de Leidos a la coleccion 🤓", libro.VolumeInfo.Title)

	default:
		mensaje = "Decidiste no agregar el libro a ninguna de las estanterias"
	}

	removerMenu := RemoverMenu(id, mensaje)
	b.API.Send(removerMenu)
	GuardarVistosRecientesGB(libro.VolumeInfo.Title, id)
}
