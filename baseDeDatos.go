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

func (b *Bot) saveSearchResult(book BookBD) error {
	client, err := initializeFirebase()
	if err != nil {
		return err
	}

	ref := client.NewRef("searchResults")
	if _, err := ref.Push(context.Background(), book); err != nil {
		return fmt.Errorf("error saving search result: %v", err)
	}

	return nil
}
func (b *Bot) getSavedSearchResults() ([]BookBD, error) {
	client, err := initializeFirebase()
	if err != nil {
		return nil, err
	}

	ref := client.NewRef("searchResults")
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
