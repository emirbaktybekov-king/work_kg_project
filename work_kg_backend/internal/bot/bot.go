package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"work_kg_backend/internal/database"
	"work_kg_backend/internal/models"
)

var Bot *tgbotapi.BotAPI
var userStates = make(map[int64]*models.UserState)

func Start(token string) {
	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("Failed to create bot: %v", err)
		return
	}

	Bot.Debug = false
	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallback(update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		handleMessage(update.Message)
	}
}

func handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := message.From.ID

	// Save user
	saveUser(message.From)

	state := userStates[userID]

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			sendWelcome(chatID)
		case "menu":
			sendMainMenu(chatID)
		case "help":
			sendHelp(chatID)
		default:
			sendMainMenu(chatID)
		}
		return
	}

	// Handle state-based input
	if state != nil {
		handleStateInput(chatID, userID, message, state)
		return
	}

	sendMainMenu(chatID)
}

func saveUser(from *tgbotapi.User) {
	username := ""
	if from.UserName != "" {
		username = from.UserName
	}
	database.SaveUser(from.ID, username, from.FirstName, from.LastName, "")
}
