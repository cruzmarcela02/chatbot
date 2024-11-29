package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) verificarFiltro(msg *tgbotapi.Message, filtro string) {
	switch filtro {
	case AUTOR:
		b.filtro += FAUTOR
	case EDITORIAL:
		b.filtro += FEDITORIAL
	case TITULO:
		b.filtro += FTITULO
	case GENERO:
		b.filtro += FGENERO
	}
	if b.filtroGLobal {
		b.sendText(msg.Chat.ID, fmt.Sprintf("Ingrese el %s a guardar para sus busquedas y recomendaciones mientras no este autenticado", msg.Text))
	} else {
		b.sendText(msg.Chat.ID, fmt.Sprintf("Por favor ingrese el %s a buscar", msg.Text))
	}
}

func (b *Bot) registroFiltros(id int64, mensaje string) bool {
	if !b.filwait && b.filtro == "" {
		removerMenu := RemoverMenu(id, "Se cancelo el proceso")
		b.API.Send(removerMenu)
		b.ultimoComando = mensaje
		return false
	}

	removerMenu := RemoverMenu(id, "Filtros ingresados con exito.")
	b.API.Send(removerMenu)
	return true
}

func (b *Bot) BusquedaFiltroGlobal(msg *tgbotapi.Message) {
	filtrosGlobales, err := obtenerFiltroGlobal(msg.Chat.ID)
	if err != nil && err.Error() == "la tabla no existe o el filtro está vacío" {
		b.sendText(msg.Chat.ID, "No tiene filtros globales guardados, se lo redigira a la personalizacion")
		b.manejarComando(msg.Chat.ID, PERSONALIZACION)
		return
	}
	b.filtro += filtrosGlobales + " "
	b.sendText(msg.Chat.ID, "Sus filtros globales van a ser aplicados "+b.filtro)
	b.sendText(msg.Chat.ID, "El ultimo comando fue "+b.ultimoComando)
	b.API.Send(crearMenu(BUSQUEDA, msg.Chat.ID))
	return
}
