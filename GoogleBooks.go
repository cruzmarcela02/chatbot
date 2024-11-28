package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/books/v1"
)

func (b *Bot) interactuarGoogleBooks(id int64) {
	// Verificar si el usuario está autenticado
	_, err := obtenerTokenAlmacenado(id)
	if err != nil {
		b.GoogleBooksAuth(id)
	}
	b.sendText(id, fmt.Sprintf("el valor de autenticado es: %s", b.autenticado))
	if b.autenticado {
		b.API.Send(crearMenu(GOOGLEBOOKS, id))
	}
}

func (b *Bot) buscarlibro(filtro string, id int64, token *oauth2.Token) {
	client := b.OAuthConfig.Client(context.Background(), token)
	service, err := books.New(client)
	if err != nil {
		b.sendText(id, "Error al crear el cliente de Google Books: "+err.Error())
		return
	}

	// Realizar la búsqueda
	b.sendText(id, fmt.Sprintf("el cliente es %s", id))

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

	service.Mylibrary.Bookshelves.Volumes.List("0")
	fav := service.Mylibrary.Bookshelves.AddVolume("0", book.Id)
	fav.Do()

	favoritos := service.Mylibrary.Bookshelves.Volumes.List("0").MaxResults(5)

	a, err := favoritos.Do()
	favorito := a.Items[0]

	b.sendText(id, fmt.Sprintf("el valor de favoritos es %s", favorito.VolumeInfo.Title))
}

func (b *Bot) agregarLibro(id int64) {

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
