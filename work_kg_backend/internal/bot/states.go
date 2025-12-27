package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"work_kg_backend/internal/database"
	"work_kg_backend/internal/models"
)

func handleStateInput(chatID int64, userID int64, message *tgbotapi.Message, state *models.UserState) {
	text := message.Text

	switch state.State {
	case "awaiting_job_title":
		if state.TempJob == nil {
			state.TempJob = &models.Job{}
		}
		state.TempJob.Title = text
		state.State = "awaiting_job_description"
		msg := tgbotapi.NewMessage(chatID, "Введите описание вакансии:")
		Bot.Send(msg)

	case "awaiting_job_description":
		state.TempJob.Description = text
		state.State = "awaiting_job_salary"
		msg := tgbotapi.NewMessage(chatID, "Введите зарплату (например: 30000-50000 сом):")
		Bot.Send(msg)

	case "awaiting_job_salary":
		state.TempJob.Salary = text
		state.State = "awaiting_job_phone"
		msg := tgbotapi.NewMessage(chatID, "Введите контактный телефон:")
		Bot.Send(msg)

	case "awaiting_job_phone":
		state.TempJob.Phone = text
		state.State = "awaiting_job_company"
		msg := tgbotapi.NewMessage(chatID, "Введите название компании (или '-' если нет):")
		Bot.Send(msg)

	case "awaiting_job_company":
		if text != "-" {
			state.TempJob.Company = text
		}
		// Save job
		state.TempJob.CreatedBy = userID
		state.TempJob.Source = "telegram"
		database.SaveJob(state.TempJob)
		delete(userStates, userID)

		msg := tgbotapi.NewMessage(chatID, "✅ Вакансия успешно добавлена!")
		Bot.Send(msg)
		sendMainMenu(chatID)

	// Step-by-step form handling
	case "form_name":
		state.FormMessageIDs = append(state.FormMessageIDs, message.MessageID)
		state.FormName = text
		state.State = "form_phone"
		askFormQuestion(chatID, state, "Введите ваш номер телефона (+996 XXX XXX XXX):")

	case "form_phone":
		state.FormMessageIDs = append(state.FormMessageIDs, message.MessageID)
		state.FormPhone = text
		state.State = "form_city"
		askFormCityQuestion(chatID, state)

	case "form_city":
		state.FormMessageIDs = append(state.FormMessageIDs, message.MessageID)
		state.FormCity = text
		state.State = "form_specialty"
		askFormQuestion(chatID, state, "Введите вашу специальность:")

	case "form_specialty":
		state.FormMessageIDs = append(state.FormMessageIDs, message.MessageID)
		state.FormSpecialty = text
		state.State = "form_experience"
		askFormQuestion(chatID, state, "Опишите ваш опыт работы:")

	case "form_experience":
		state.FormMessageIDs = append(state.FormMessageIDs, message.MessageID)
		state.FormExperience = text
		// Delete ALL collected messages at the end
		deleteAllFormMessages(chatID, state)
		// Save form data and show confirmation
		saveFormData(userID, state)
		showFormSummary(chatID, state)
		delete(userStates, userID)
	}
}

func deleteAllFormMessages(chatID int64, state *models.UserState) {
	for _, msgID := range state.FormMessageIDs {
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
		Bot.Request(deleteMsg)
	}
}

func askFormQuestion(chatID int64, state *models.UserState, question string) {
	msg := tgbotapi.NewMessage(chatID, question)
	sentMsg, err := Bot.Send(msg)
	if err == nil {
		state.FormMessageIDs = append(state.FormMessageIDs, sentMsg.MessageID)
	}
}

func askFormCityQuestion(chatID int64, state *models.UserState) {
	text := "Выберите ваш город или введите свой:"

	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Бишкек", "form_city:Бишкек"),
		tgbotapi.NewInlineKeyboardButtonData("Ош", "form_city:Ош"),
		tgbotapi.NewInlineKeyboardButtonData("Талас", "form_city:Талас"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Нарын", "form_city:Нарын"),
		tgbotapi.NewInlineKeyboardButtonData("Каракол", "form_city:Каракол"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Жалал-Абад", "form_city:Жалал-Абад"),
		tgbotapi.NewInlineKeyboardButtonData("Чолпон-Ата", "form_city:Чолпон-Ата"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	sentMsg, err := Bot.Send(msg)
	if err == nil {
		state.FormMessageIDs = append(state.FormMessageIDs, sentMsg.MessageID)
	}
}

func saveFormData(userID int64, state *models.UserState) {
	database.UpdateUserFormData(userID, state.FormName, state.FormPhone, state.FormCity, state.FormSpecialty, state.FormExperience)

	username := database.GetUsernameByTelegramID(userID)

	database.SaveResume(userID, username, state.FormName, state.FormPhone, state.FormCity, state.FormSpecialty, state.FormExperience)
}
