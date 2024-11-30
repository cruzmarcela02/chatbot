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
	COD_NAVEGACION       = "9"

	FAVORITOS        = "Favoritos"
	POR_LEER         = "Por Leer"
	LEYENDO_AHORA    = "Leyendo Ahora"
	LEIDOS           = "Leidos"
	VISTOS_RECIENTES = "Vistos Recientes"
	LIBROS_PARA_TI   = "Libros para ti"
	NAVEGACION       = "De navegación"
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

	b.sendText(id, fmt.Sprintf("El libro encontrado es %s.Descargalo en %s", titulo, downloadLink))

	service.Mylibrary.Bookshelves.Volumes.List(COD_LEIDOS)
	buffer := service.Mylibrary.Bookshelves.AddVolume(COD_LEIDOS, book.Id)
	buffer.Do()

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

	recuperarLibro := service.Mylibrary.Bookshelves.Volumes.List(COD_LEIDOS).MaxResults(1)
	llamado, _ := recuperarLibro.Do()
	libro := llamado.Items[0]

	if estanteria == FAVORITOS {
		favoritos := service.Mylibrary.Bookshelves.AddVolume(COD_FAVORITOS, libro.Id)
		favoritos.Do()
		b.sendText(id, fmt.Sprintf("El libro '%s' ha sido agregado a tus favoritos.", libro.VolumeInfo.Title))

	} else if estanteria == POR_LEER {
		porLeer := service.Mylibrary.Bookshelves.AddVolume(COD_LEER, libro.Id)
		porLeer.Do()
		b.sendText(id, fmt.Sprintf("El libro '%s' ha sido agregado a tus libros por leer.", libro.VolumeInfo.Title))

	} else if estanteria == LEYENDO_AHORA {
		leyendo := service.Mylibrary.Bookshelves.AddVolume(COD_LEYENDO, libro.Id)
		leyendo.Do()
		b.sendText(id, fmt.Sprintf("El libro '%s' ha sido agregado a tus libros que estas leyendo.", libro.VolumeInfo.Title))

	} else {
		b.sendText(id, "No se agregara el libro a ninguna de esas estanterias")
	}
	buffer := service.Mylibrary.Bookshelves.RemoveVolume("4", libro.Id)
	buffer.Do()

}

/*call := service.Volumes.List(busqueda).MaxResults(5)
	resp, err := call.Do()
	if err != nil {
		b.sendText(msg.Chat.ID, "Error al buscar libros: "+err.Error())
		return
	}

	if len(resp.Items) == 0 {
		b.sendText(msg.Chat.ID, "No se encontraron libros.")
		return
	}

	// Mostrar los resultados y pedir al usuario que elija un libro
	for i, item := range resp.Items {
		b.sendText(msg.Chat.ID, fmt.Sprintf("%d. %s por %s", i+1, item.VolumeInfo.Title, item.VolumeInfo.Authors[0]))
	}
	b.sendText(msg.Chat.ID, "Por favor, elige un número de libro para agregar a tu bookshelf:")

	// Aquí deberías implementar una lógica para esperar la respuesta del usuario
	// Por ahora, elegiremos el primer libro para el ejemplo
	libroElegido := resp.Items[0]

	// Agregar el libro al bookshelf del usuario
	_, err = service.Mylibrary.Bookshelves.AddVolume("0", libroElegido.Id).Do()
	if err != nil {
		b.sendText(msg.Chat.ID, "Error al agregar el libro a tu bookshelf: "+err.Error())
		return
	}

	b.sendText(msg.Chat.ID, fmt.Sprintf("El libro '%s' ha sido agregado a tu bookshelf.", libroElegido.VolumeInfo.Title))
}*/

// dentro de la busqueda o la recomendacion de google books -> ver de si queremos agregar  a favoritos o a por leer
