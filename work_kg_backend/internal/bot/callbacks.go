package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"work_kg_backend/internal/models"
)

func handleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	userID := callback.From.ID
	data := callback.Data
	messageID := callback.Message.MessageID

	// Answer callback
	Bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	// Delete the message that contained the button (clean chat)
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	Bot.Request(deleteMsg)

	parts := strings.Split(data, ":")

	switch parts[0] {
	case "menu":
		sendMainMenu(chatID)

	case "profile":
		sendProfile(chatID, userID)

	case "search_employee":
		userStates[userID] = &models.UserState{State: "search_employee", SearchType: "employee"}
		sendCategorySelection(chatID, "employee")

	case "search_job":
		userStates[userID] = &models.UserState{State: "search_job", SearchType: "job"}
		sendCategorySelection(chatID, "job")

	case "entertainment":
		sendEntertainment(chatID)

	case "earn_together":
		sendEarnTogether(chatID)

	case "subscription":
		sendSubscription(chatID)

	case "category":
		if len(parts) > 2 {
			category := parts[1]
			searchType := parts[2]
			state := userStates[userID]
			if state == nil {
				state = &models.UserState{}
				userStates[userID] = state
			}
			state.Category = category
			state.SearchType = searchType
			sendSubcategorySelection(chatID, category, searchType)
		}

	case "subcategory":
		if len(parts) > 2 {
			subcategory := parts[1]
			searchType := parts[2]
			state := userStates[userID]
			if state == nil {
				state = &models.UserState{}
				userStates[userID] = state
			}
			state.Subcategory = subcategory
			sendCitySelection(chatID, searchType)
		}

	case "city":
		if len(parts) > 2 {
			city := parts[1]
			searchType := parts[2]
			state := userStates[userID]
			if state == nil {
				state = &models.UserState{SearchType: searchType}
				userStates[userID] = state
			}
			state.City = city

			if searchType == "job" {
				showJobs(chatID, state)
			} else if searchType == "employee" {
				sendAddVacancyPrompt(chatID, state)
			}
		}

	case "add_vacancy":
		state := userStates[userID]
		if state == nil {
			state = &models.UserState{}
			userStates[userID] = state
		}
		state.State = "awaiting_job_title"
		state.TempJob = &models.Job{
			Category:    state.Category,
			Subcategory: state.Subcategory,
			City:        state.City,
		}
		msg := tgbotapi.NewMessage(chatID, "Введите название вакансии:")
		Bot.Send(msg)

	case "fill_form":
		sendFormInstructions(chatID, userID)

	case "form_city":
		if len(parts) > 1 {
			city := parts[1]
			state := userStates[userID]
			if state != nil && state.State == "form_city" {
				state.FormCity = city
				state.State = "form_specialty"
				askFormQuestion(chatID, state, "Введите вашу специальность:")
			}
		}

	case "back":
		sendMainMenu(chatID)
	}
}
