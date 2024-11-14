package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	// Menu texts
	MenuStart         = "Bienvenido a tu bot literario."
	MenuBusqueda      = "Bajo que parametro deseas buscar?"
	MenuRecomendacion = "Recomendacion"

	MenuGoogleBooks = "Google Books"

	//secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	recomendacion   = "Recomendacion"
	busqueda        = "Busqueda"
	historial       = "Historial"
	googlebooks     = "Googlebooks"
	informe         = "Informe"
	personalización = "Personalización"
)

func crearMenu(comando string, id int64) tgbotapi.MessageConfig {

	switch comando {
	case START:
		return CrearMenuStart(id)
	case RECOMENDACION:
		return CrearMenuRecomendar(id)
	default:
		return CrearMenuBusqueda(id)
	}
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
			tgbotapi.NewInlineKeyboardButtonData(informe, INFORME),
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
	recomendar := tgbotapi.NewMessage(id, MenuRecomendacion)

	recomendar.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(AUTOR),
			tgbotapi.NewKeyboardButton(EDITORIAL),
			tgbotapi.NewKeyboardButton(GENERO),
		),
	)
	return recomendar
}

func RemoverMenu(id int64, mensaje string) tgbotapi.MessageConfig {
	menu := tgbotapi.NewMessage(id, mensaje)
	menu.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return menu
}
