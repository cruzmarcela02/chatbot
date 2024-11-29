package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (b *Bot) obtenerTokenAlmacenado(userID int64) (*oauth2.Token, error) {
	/*	token, ok := userTokens[userID]
		if !ok {
			return nil, fmt.Errorf("token no encontrado para el usuario")
		}
		return token, nil*/

	//USAR FIREBASE
	client, err := initializeFirebase()
	if err != nil {
		return nil, err
	}

	ref := client.NewRef(fmt.Sprintf("GBooksTokens/%d", userID))
	var token oauth2.Token
	err = ref.Get(context.Background(), &token)

	if err != nil {
		return nil, fmt.Errorf("error retrieving token: %v", err)
	}

	if token.Expiry.Before(time.Now()) {
		refresco := b.OAuthConfig.TokenSource(context.Background(), &token)
		nuevo, a := refresco.Token()
		if a != nil {
			return nil, fmt.Errorf("error retrieving token: %v", err)
		}
		err = b.almacenarToken(userID, nuevo)
		if err != nil {
			return nil, err
		}
		return nuevo, nil
	}

	return &token, err
}

func alamacenarestado(id int64, estado string) error {
	client, err := initializeFirebase()
	if err != nil {
		return err
	}
	ref := client.NewRef(fmt.Sprintf("GBooksEstados/%d", id))
	err = ref.Set(context.Background(), estado)
	return err
}

func obtenerEstadoAlmacenado(id int64) (string, error) {
	client, err := initializeFirebase()
	if err != nil {

		return "", err
	}

	ref := client.NewRef(fmt.Sprintf("GBooksEstados/%d", id))
	var estado string
	err = ref.Get(context.Background(), &estado)

	return estado, err
}

func (b *Bot) GoogleBooksAuth(id int64) {
	// Generar un token de estado único para este usuario
	estado := "state-" + strconv.FormatInt(id, 10)
	err := alamacenarestado(id, estado)
	if err != nil {
		return
	}

	authURL := b.OAuthConfig.AuthCodeURL(estado, oauth2.AccessTypeOffline)
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
	storedState, err := obtenerEstadoAlmacenado(userID)
	if err != nil || state != storedState {
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
	//userTokens[userID] = token

	//USAR FIREBASE
	client, err := initializeFirebase()
	if err != nil {
		b.sendText(userID, "Error al inicializar Firebase: "+err.Error())
		return err
	}
	ref := client.NewRef(fmt.Sprintf("GBooksTokens/%d", userID))
	if err := ref.Set(context.Background(), token); err != nil {
		return fmt.Errorf("error saving token: %v", err)
	}

	return nil
}
