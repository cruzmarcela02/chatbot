package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	MenuStart         = "Bienvenido a tu bot literario."
	MenuBusqueda      = "Bajo que parametro deseas buscar?"
	MenuRecomendacion = "Recomendacion"
	MenuHistorial     = "Historial"
	MenuGoogleBooks   = "Bienvenido a tu Google Books"
	MenuAgregar       = "¿Deseas agregar el libro a alguna estanteria?"
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
	case HISTORIAL:
		return CrearMenuHistorial(id)
	case GOOGLEBOOKS:
		return CrearMenuGoogleBooks(id)
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
			tgbotapi.NewKeyboardButton(AUTOR)),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(EDITORIAL),
			tgbotapi.NewKeyboardButton(GENERO),
			tgbotapi.NewKeyboardButton(TERMINAR),
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
			tgbotapi.NewKeyboardButton(TERMINAR),
		),
	)

	return recomendar
}

func CrearMenuHistorial(id int64) tgbotapi.MessageConfig {
	verHistorial := tgbotapi.NewMessage(id, MenuHistorial)
	verHistorial.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BUSQUEDAS),
			tgbotapi.NewKeyboardButton(RECOMENDACIONES),
		),
	)
	return verHistorial
}

func RemoverMenu(id int64, mensaje string) tgbotapi.MessageConfig {
	menu := tgbotapi.NewMessage(id, mensaje)
	menu.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return menu
}

func CrearMenuGoogleBooks(id int64) tgbotapi.MessageConfig {
	menuGoogleBooks := tgbotapi.NewMessage(id, MenuGoogleBooks)
	menuGoogleBooks.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(recomendacion, RECOMENDACION),
			tgbotapi.NewInlineKeyboardButtonData(busqueda, BUSQUEDA),
			tgbotapi.NewInlineKeyboardButtonData(historial, HISTORIAL),
		),
	)
	return menuGoogleBooks
}

func CrearMenuAgregar(id int64) tgbotapi.MessageConfig {
	menuAgregar := tgbotapi.NewMessage(id, MenuAgregar)
	menuAgregar.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(FAVORITOS),
			tgbotapi.NewKeyboardButton(POR_LEER),
			tgbotapi.NewKeyboardButton(LEYENDO_AHORA),
			tgbotapi.NewKeyboardButton(NO_AGREGAR),
		),
	)
	return menuAgregar
}
