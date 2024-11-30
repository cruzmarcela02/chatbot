package main

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"fmt"
	"google.golang.org/api/option"
)

type BookBD struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func initializeFirebase() (*db.Client, error) {
	opt := option.WithCredentialsFile("chatbottdl-firebase-adminsdk-fpk9k-22b8516b13.json")
	config := &firebase.Config{
		DatabaseURL: "https://chatbottdl-default-rtdb.firebaseio.com/",
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Database(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing database client: %v", err)
	}

	return client, nil
}

func (b *Bot) saveSearchResult(book BookBD, id int64) error {
	client, err := initializeFirebase()
	if err != nil {
		return err
	}

	ref := client.NewRef(fmt.Sprintf("searchResults/%d", id))
	if _, err := ref.Push(context.Background(), book); err != nil {
		return fmt.Errorf("error saving search result: %v", err)
	}

	return nil
}
func (b *Bot) getSavedSearchResults(id int64) ([]BookBD, error) {
	client, err := initializeFirebase()
	if err != nil {
		return nil, err
	}

	ref := client.NewRef(fmt.Sprintf("searchResults/%d", id))
	var booksMap map[string]BookBD
	if err := ref.Get(context.Background(), &booksMap); err != nil {
		return nil, fmt.Errorf("error retrieving search results: %v", err)
	}

	var books []BookBD
	for _, book := range booksMap {
		books = append(books, book)
	}

	return books, nil
}

func guardarFiltroGlobal(id int64, filtro string) error {
	client, err := initializeFirebase()
	if err != nil {
		return err
	}

	ref := client.NewRef(fmt.Sprintf("filtros/%d", id))
	err = ref.Set(context.Background(), filtro)
	return nil
}

func obtenerFiltroGlobal(id int64) (string, error) {
	client, err := initializeFirebase()
	if err != nil {
		return "", err
	}

	ref := client.NewRef(fmt.Sprintf("filtros/%d", id))
	var filtro string
	err = ref.Get(context.Background(), &filtro)
	if err != nil {
		return "", err
	}
	if filtro == "" {
		return "", fmt.Errorf("la tabla no existe o el filtro está vacío")
	}
	return filtro, nil
}

func eliminarFiltrosBD(id int64) error {
	client, err := initializeFirebase()
	if err != nil {
		return err
	}

	ref := client.NewRef(fmt.Sprintf("filtros/%d", id))
	err = ref.Delete(context.Background())
	return err
}

func (b *Bot) obtenerRecomendaciones(id int64) ([]BookBD, error) {
	client, err := initializeFirebase()
	if err != nil {
	}

	ref := client.NewRef(fmt.Sprintf("recomendaciones/%d", id))
	var recomendacionesMap map[string]BookBD
	if err := ref.Get(context.Background(), &recomendacionesMap); err != nil {
		return nil, fmt.Errorf("error retrieving recomendaciones: %v", err)
	}

	var recomendaciones []BookBD
	for _, recomendacion := range recomendacionesMap {
		recomendaciones = append(recomendaciones, recomendacion)
	}
	return recomendaciones, nil
}

func (b *Bot) guardarRecomendaciones(book BookBD, id int64) error {
	client, err := initializeFirebase()
	if err != nil {
		return err
	}

	ref := client.NewRef(fmt.Sprintf("recomendaciones/%d", id))
	if _, err := ref.Push(context.Background(), book); err != nil {
		return fmt.Errorf("error saving recomendaciones: %v", err)
	}

	return nil
}
