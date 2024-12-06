package main

import (
	"context"
	"fmt"
	"time"
)

const (
	MENSUAL = "Ultimo mes"
	SEMANAL = "Ultima semana"
	DIARIO  = "Ultimo dia"
)

func consguirLibros(conjuntolibros map[string]BookBD, periodo time.Time) []BookBD {
	var libros []BookBD
	for _, query := range conjuntolibros {
		if query.Periodo.After(periodo) {
			libros = append(libros, query)
		}
	}
	return libros
}

func (b *Bot) generarAnalisis(id int64, periodo string) {
	client, err := initializeFirebase()
	if err != nil {
		return
	}

	refSearch := client.NewRef(fmt.Sprintf("searchResults/%d", id))
	var booksBusqueda map[string]BookBD
	if err := refSearch.Get(context.Background(), &booksBusqueda); err != nil {
		return
	}
	var busquedas []BookBD

	refRecomendaciones := client.NewRef(fmt.Sprintf("recomendaciones/%d", id))
	var booksRecomendacion map[string]BookBD
	if err := refRecomendaciones.Get(context.Background(), &booksRecomendacion); err != nil {
		return
	}
	var recomendaciones []BookBD

	switch periodo {
	case MENSUAL:
		busquedas = consguirLibros(booksBusqueda, time.Now().AddDate(0, -1, 0))
		recomendaciones = consguirLibros(booksRecomendacion, time.Now().AddDate(0, -1, 0))
	case SEMANAL:
		// Obtener las recomendaciones y busquedas de la ultima semana
		busquedas = consguirLibros(booksBusqueda, time.Now().AddDate(0, 0, -7))
		recomendaciones = consguirLibros(booksRecomendacion, time.Now().AddDate(0, 0, -7))
	case DIARIO:
		// Obtener las recomendaciones y busquedas del ultimo dia
		busquedas = consguirLibros(booksBusqueda, time.Now().AddDate(0, 0, -1))
		recomendaciones = consguirLibros(booksRecomendacion, time.Now().AddDate(0, 0, -1))
	}

	// Mostrar los libros
	var mensaje string
	for i, libro := range busquedas {
		mensaje += fmt.Sprintf("%d. %s\n", i+1, libro.Title)
	}
	b.sendText(id, "Tus busquedas en el periodo indicado son las siguientes\n"+mensaje)
	mensaje = ""
	for i, libro := range recomendaciones {
		mensaje += fmt.Sprintf("%d. %s\n", i+1, libro.Title)
	}
	b.sendText(id, "Tus recomendaciones en el periodo indicado son las siguientes\n"+mensaje)
	b.API.Send(RemoverMenu(id, "Espero que te haya sido de ayuda!"))
}
