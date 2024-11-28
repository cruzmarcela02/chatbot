package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
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

func (b *Bot) GoogleBooksAuth(id int64) {
	// Generar un token de estado único para este usuario
	state := "state-" + strconv.FormatInt(id, 10)
	stateTokens[id] = state // Almacenar el token asociado al usuario

	authURL := b.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	b.sendText(id, "Por favor, autorizanos a acceder a google books visitando esta URL: "+authURL)
	b.autenticado = false
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
	fmt.Fprint(w, "Ya estas autenticado. Puedes cerrar esta ventana.")
	b.autenticado = true
	b.manejarComando(userID, GOOGLEBOOKS)
}

func (b *Bot) almacenarToken(userID int64, token *oauth2.Token) error {
	// Implementar la lógica para almacenar el token
	// Por ejemplo, puedes usar un mapa en memoria (no recomendado para producción)
	// o almacenarlo en una base de datos
	userTokens[userID] = token
	return nil
}
