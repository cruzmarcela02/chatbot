package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	// Menu texts
	MenuStart         = "Bienvenido a tu bot literario."
	MenuBusqueda      = "Bajo que parametro deseas buscar?"
	MenuRecomendacion = "Recomendacion"
	MenuGoogleBooks   = "Google Books"
	//secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	recomendacion   = "Recomendacion"
	busqueda        = "Busqueda"
	historial       = "Historial"
	googlebooks     = "Googlebooks"
	personalización = "Personalización"

	/*	// Keyboard layout for the first menu. One button, one row
		MenuStartMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(recomendacion, RECOMENDACION),
				tgbotapi.NewInlineKeyboardButtonData(busqueda, BUSQUEDA),
				tgbotapi.NewInlineKeyboardButtonData(historial, HISTORIAL),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(googlebooks, GOOGLEBOOKS),
				tgbotapi.NewInlineKeyboardButtonData(personalización, PERSONALIZACION),
			),
		)*/
)

func crearMenu(comando string, id int64) tgbotapi.MessageConfig {

	var menu_opciones tgbotapi.MessageConfig
	switch comando {
	case "/start":
		menu_opciones = CrearMenuStart(id)
	case RECOMENDACION:
		menu_opciones = CrearMenuRecomendar(id)
	case BUSQUEDA:
		menu_opciones = CrearMenuBusqueda(id)
	}
	return menu_opciones
}

func CrearMenuStart(id int64) tgbotapi.MessageConfig {
	start := tgbotapi.NewMessage(id, MenuStart)
	start.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(recomendacion, RECOMENDACION),
			tgbotapi.NewInlineKeyboardButtonData(busqueda, BUSQUEDA),
			tgbotapi.NewInlineKeyboardButtonData(historial, HISTORIAL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(googlebooks, GOOGLEBOOKS),
			tgbotapi.NewInlineKeyboardButtonData(personalización, PERSONALIZACION),
		),
	)
	return start
}

func CrearMenuBusqueda(id int64) tgbotapi.MessageConfig {
	busqueda := tgbotapi.NewMessage(id, MenuBusqueda)

	busqueda.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(TITULO),
			tgbotapi.NewKeyboardButton(AUTOR),
			tgbotapi.NewKeyboardButton(EDITORIAL),
			tgbotapi.NewKeyboardButton(GENERO),
		),
	)

	return busqueda
}

func CrearMenuRecomendar(id int64) tgbotapi.MessageConfig {
	busqueda := tgbotapi.NewMessage(id, MenuBusqueda)
	busqueda.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(AUTOR, AUTOR),
			tgbotapi.NewInlineKeyboardButtonData(EDITORIAL, EDITORIAL),
			tgbotapi.NewInlineKeyboardButtonData(GENERO, GENERO),
		),
	)
	return busqueda
}
