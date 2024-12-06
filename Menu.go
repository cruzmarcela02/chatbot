package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	MenuStart         = "Bienvenido a tu bot literario."
	MenuBusqueda      = "Bajo que parametro deseas buscar?"
	MenuRecomendacion = "Veo que buscas una /recomendacion, elegi algunos filtros para ver que puedo recomendarte"
	MenuHistorial     = "Buenas! Seleccionaste /historial\n📚Por favor elegi que historial queres ver"
	MenuGoogleBooks   = "Bienvenido a tu Google Books"
	MenuAgregar       = "¿Deseas agregar el libro a alguna estanteria?"
	MenuPersonalizar  = "Seleccione alguno de los filtros. Estos filtros globales se aplicaran a sus busquedas y recomendaciones que no sean dentro de " + GOOGLEBOOKS + " si asi lo desea."
	menuFiltros       = "¿Que tipo de query desea realizar?"
	menuInforme       = "Buenas! Seleccionaste /informe\n 📋Pero necesito que me indiques de cuál periodo lo queres"
	// Button texts
	recomendacion   = "Recomendacion"
	busqueda        = "Busqueda"
	historial       = "Historial"
	googlebooks     = "Googlebooks"
	informe         = "Informe"
	personalización = "Personalización"
)

func crearMenu(comando string, id int64, enGoogleBooks bool) tgbotapi.MessageConfig {

	switch comando {
	case START:
		return CrearMenuStart(id)
	case RECOMENDACION:
		return CrearMenuRecomendar(id)
	case HISTORIAL:
		return CrearMenuHistorial(id, enGoogleBooks)
	case GOOGLEBOOKS:
		return CrearMenuGoogleBooks(id)
	case PERSONALIZACION:
		return CrearMenuPersonalizacion(id)
	case INFORME:
		return CrearMenuInforme(id)
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

func CrearMenuHistorial(id int64, enGoogleBooks bool) tgbotapi.MessageConfig {
	if enGoogleBooks {
		return CrearMenuHistorialGoogleBooks(id)
	}

	return CrearMenuHistorialComun(id)
}

func CrearMenuHistorialComun(id int64) tgbotapi.MessageConfig {
	menuHistorial := tgbotapi.NewMessage(id, MenuHistorial)
	menuHistorial.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BUSQUEDAS),
			tgbotapi.NewKeyboardButton(RECOMENDACIONES),
		),
	)

	return menuHistorial
}

func CrearMenuHistorialGoogleBooks(id int64) tgbotapi.MessageConfig {
	menuHistorial := tgbotapi.NewMessage(id, MenuHistorial)
	menuHistorial.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(VISTOS_RECIENTES),
			tgbotapi.NewKeyboardButton(LEIDOS),
		),
	)

	return menuHistorial
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
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(LEIDOSB),
			tgbotapi.NewKeyboardButton(NO_AGREGAR),
		),
	)
	return menuAgregar
}

func CrearMenuPersonalizacion(id int64) tgbotapi.MessageConfig {
	menuPersonalizacion := tgbotapi.NewMessage(id, MenuPersonalizar)
	menuPersonalizacion.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(AUTOR),
			tgbotapi.NewKeyboardButton(GENERO),
			tgbotapi.NewKeyboardButton(EDITORIAL),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(TERMINAR),
			tgbotapi.NewKeyboardButton(ELIMINAR_FILTROS),
		),
	)
	return menuPersonalizacion
}

func CrearMenuFiltros(id int64) tgbotapi.MessageConfig {
	tipoBusqueda := tgbotapi.NewMessage(id, menuFiltros)
	tipoBusqueda.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BUSQUEDA_GLOBAL),
			tgbotapi.NewKeyboardButton(BUSQUEDA_PERSONALIZADA),
		),
	)
	return tipoBusqueda
}

func CrearMenuInforme(id int64) tgbotapi.MessageConfig {
	analisis := tgbotapi.NewMessage(id, menuInforme)
	analisis.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(MENSUAL),
			tgbotapi.NewKeyboardButton(SEMANAL),
			tgbotapi.NewKeyboardButton(DIARIO),
		),
	)

	return analisis
}
