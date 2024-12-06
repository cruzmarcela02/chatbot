package main

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/books/v1"
)

/* Envia las opciones de historiales posibles */
func (b *Bot) DarHistorial(id int64, enGoogleBooks bool) {
	if enGoogleBooks {
		b.mostrarHistorialGoogleBooks(id)
		return
	}
	b.API.Send(crearMenu(HISTORIAL, id))
}

/* Historiales - GoogleBooks: Leidos */
func (b *Bot) mostrarHistorialGoogleBooks(id int64) {
	token, _ := b.obtenerTokenAlmacenado(id)
	service := autenticarCliente(b, id, token)

	removerMenu := RemoverMenu(id, "Dale, repasaremos los libros leidos hasta ahora 📘✅")
	b.API.Send(removerMenu)
	historial := armarHistorialGoogleBooks(service)
	b.sendText(id, historial)
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
func armarHistorialGoogleBooks(service *books.Service) string {
	var historial string

	bookshelf, err := service.Mylibrary.Bookshelves.Volumes.List(COD_LEIDOS).Do()
	if err != nil {
		historial = "SURGIO UN ERROR: " + err.Error()
		return historial
	}

	if len(bookshelf.Items) == 0 {
		historial = "Upss, la estanteria esta vacia. No leiste nada aun"
		return historial
	}
	emojis := []string{"1️⃣ ", "2️⃣ ", "3️⃣ ", "4️⃣ ", "5️⃣ ", "6️⃣ ", "7️⃣ ", "8️⃣ ", "9️⃣ ", "🔟 "}
	for i, libro := range bookshelf.Items {
		historial += emojis[i]
		historial += libro.VolumeInfo.Title
		historial += "\n"
	}

	return historial
}

// Historial del chat. Fuera de google books
func (b *Bot) armarHistorial(msg *tgbotapi.Message, filtro string) {
	if filtro == RECOMENDACIONES {
		removerMenu := RemoverMenu(msg.Chat.ID, "Okey! Veamos tus libros recomendados 📚\nTe mostraré tus últimas 10 recomendaciones 📋")
		b.API.Send(removerMenu)
		recomendaciones, err := b.obtenerRecomendaciones(msg.Chat.ID)
		if err != nil {

		}
		for i, recomendacion := range recomendaciones {
			if i > 9 {
				break
			}

			var mensaje string
			mensaje += strconv.Itoa(i + 1)
			mensaje += ". Titulo: "
			mensaje += recomendacion.Title
			mensaje += "\nLink: "
			mensaje += recomendacion.Link
			b.sendText(msg.Chat.ID, mensaje)
		}
		return
	}

	removerMenu := RemoverMenu(msg.Chat.ID, "¿Asi que quéres ver tus busquedas 🔎🤔? \nTe mostraré tus últimas 10 busquedas 😉")
	b.API.Send(removerMenu)

	books, err := b.getSavedSearchResults(msg.Chat.ID)
	if err != nil {
		b.sendText(msg.Chat.ID, "Error al obtener el historial de búsquedas: "+err.Error())
		return
	}

	if len(books) == 0 {
		b.sendText(msg.Chat.ID, "No se encontraron resultados de búsqueda guardados.")
		return
	}

	for i, book := range books {
		if i > 9 {
			break
		}

		var mensaje string
		mensaje += strconv.Itoa(i + 1)
		mensaje += ". Titulo: "
		mensaje += book.Title
		mensaje += "\nLink: "
		mensaje += book.Link

		b.sendText(msg.Chat.ID, mensaje)
	}
}
