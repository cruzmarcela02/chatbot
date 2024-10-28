package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/books/v1"
	"net/http"
	"strconv"
	"strings"
)

var userTokens = make(map[int64]*oauth2.Token)
var stateTokens = make(map[int64]string)

func obtenerTokenAlmacenado(userID int64) (*oauth2.Token, error) {
	token, ok := userTokens[userID]
	if !ok {
		return nil, fmt.Errorf("token no encontrado para el usuario")
	}
	return token, nil
}

func (b *Bot) GoogleBooksAuth(msg *tgbotapi.Message) {

	// Generar un token de estado único para este usuario
	state := "state-" + strconv.FormatInt(msg.From.ID, 10)
	stateTokens[msg.From.ID] = state // Almacenar el token asociado al usuario

	authURL := b.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	b.sendText(msg.Chat.ID, "Por favor, autorizanos a acceder a google books visitando esta URL: "+authURL)

}

func obtenerUserIDDesdeState(state string) int64 {
	parts := strings.Split(state, "-")
	if len(parts) < 2 {
		return 0
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0
	}
	return userID
}

func (b *Bot) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")

	// Verificar que el estado recibido coincida con el almacenado
	userID := obtenerUserIDDesdeState(state)
	if storedState, ok := stateTokens[userID]; !ok || state != storedState {
		http.Error(w, "Invalid state token", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := b.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = b.almacenarToken(userID, token)
	if err != nil {
		http.Error(w, "Failed to store token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Autenticacion realizada ocn exito.")

}

func (b *Bot) almacenarToken(userID int64, token *oauth2.Token) error {
	// Implementar la lógica para almacenar el token
	// Por ejemplo, puedes usar un mapa en memoria (no recomendado para producción)
	// o almacenarlo en una base de datos
	userTokens[userID] = token
	return nil
}

var waitingForSearchTerm = make(map[int64]bool)

func (b *Bot) interactuarGoogleBooks(msg *tgbotapi.Message, update tgbotapi.Update) {
	// Verificar si el usuario está autenticado
	token, err := obtenerTokenAlmacenado(msg.From.ID)
	if err != nil {
		b.GoogleBooksAuth(msg)
		return
	}

	// Si el usuario está autenticado, pedir el término de búsqueda
	b.sendText(msg.Chat.ID, "Por favor, ingresa el término de búsqueda para Google Books:")
	// quedarse esparando u nnuevo mensaje con el termino de busqueda y realizar la busqueda
	waitingForSearchTerm[msg.From.ID] = true

	if waitingForSearchTerm[msg.From.ID] {
		if update.Message == nil {
			b.sendText(msg.Chat.ID, "Por favor, ingresa el término de búsqueda para Google Books:")
		}
		waitingForSearchTerm[msg.From.ID] = false
		// esperar a una nueva update

		b.buscarYAgregarLibro(msg, token)

	}

}

func (b *Bot) buscarYAgregarLibro(msg *tgbotapi.Message, token *oauth2.Token) {
	client := b.OAuthConfig.Client(context.Background(), token)
	_, err := books.New(client)
	if err != nil {
		b.sendText(msg.Chat.ID, "Error al crear el cliente de Google Books: "+err.Error())
		return
	}

	// Realizar la búsqueda
	b.sendText(msg.Chat.ID, msg.Text)

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

	b.sendText(msg.Chat.ID, fmt.Sprintf("El libro '%s' ha sido agregado a tu bookshelf.", libroElegido.VolumeInfo.Title))*/
}

// main

// menu -> comandos y botones

// googleOuth -> manejo de autenticacion -> requiere base de datos o json

// googleBooks -> manejo de google books -> usa googleOuth
