package bot

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"work_kg_backend/internal/database"
	"work_kg_backend/internal/models"
)

func sendWelcome(chatID int64) {
	text := `⚠️ Инструкция по использованию ⚠️

1️⃣ Заполнить анкету - Самый важный пункт. Для того чтобы с вами связались работники/работодатели

2️⃣ Поиск сотрудника - В этом разделе вы можете быстро найти временного и постоянного работника/сотрудника

3️⃣ Поиск работы - В этом разделе вы можете быстро найти временную и постоянную работу

4️⃣ Развлечение - В этом разделе вы можете разгрузить себя от суеты шутками и способами

5️⃣ Зарабатывать вместе - в этом разделе вы можете зарабатывать с нами, выполняя разные задачи.`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👍 Ознакомился", "menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Заполнить анкету", "fill_form"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendMainMenu(chatID int64) {
	text := "💵 💵 Главное меню 💵 💵"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Личный кабинет 📁", "profile"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Поиск сотрудника 👷", "search_employee"),
			tgbotapi.NewInlineKeyboardButtonData("Поиск работы 😌", "search_job"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Развлечение 😊", "entertainment"),
			tgbotapi.NewInlineKeyboardButtonData("Зарабатывать вместе 💸", "earn_together"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Приобрести подписку ✅", "subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("назад ⬅️", "back"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendHelp(chatID int64) {
	text := `❓ Помощь

Команды бота:
/start - Начать работу с ботом
/menu - Главное меню
/help - Помощь

По всем вопросам обращайтесь к администратору.`

	msg := tgbotapi.NewMessage(chatID, text)
	Bot.Send(msg)
}

func sendProfile(chatID int64, userID int64) {
	user, err := database.GetUserByTelegramID(userID)

	if err != nil {
		text := "📁 Личный кабинет\n\nПрофиль не найден. Заполните анкету для создания профиля."
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📝 Заполнить анкету", "fill_form"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu"),
			),
		)
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		Bot.Send(msg)
		return
	}

	text := "📁 Личный кабинет\n\n"
	text += fmt.Sprintf("👤 Имя: %s %s\n", user.FirstName, user.LastName)
	if user.Username != "" {
		text += fmt.Sprintf("📱 Username: @%s\n", user.Username)
	}
	if user.Phone != "" {
		text += fmt.Sprintf("📞 Телефон: %s\n", user.Phone)
	}
	if user.City != "" {
		text += fmt.Sprintf("📍 Город: %s\n", user.City)
	}
	if user.Specialty != "" {
		text += fmt.Sprintf("💼 Специальность: %s\n", user.Specialty)
	}
	if user.Experience != "" {
		text += fmt.Sprintf("📝 Опыт: %s\n", user.Experience)
	}
	text += fmt.Sprintf("📅 Дата регистрации: %s", user.CreatedAt.Format("02.01.2006"))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Редактировать анкету", "fill_form"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendFormInstructions(chatID int64, userID int64) {
	state := &models.UserState{State: "form_name", FormMessageIDs: []int{}}
	userStates[userID] = state

	text := "📝 Заполнение анкеты\n\nВведите ваше имя:"

	msg := tgbotapi.NewMessage(chatID, text)
	sentMsg, err := Bot.Send(msg)
	if err == nil {
		state.FormMessageIDs = append(state.FormMessageIDs, sentMsg.MessageID)
	}
}

func showFormSummary(chatID int64, state *models.UserState) {
	text := fmt.Sprintf(`✅ Ваша анкета сохранена!

📋 Ваши данные:

👤 Имя: %s
📞 Телефон: %s
📍 Город: %s
💼 Специальность: %s
📝 Опыт: %s

Работодатели смогут с вами связаться.`, state.FormName, state.FormPhone, state.FormCity, state.FormSpecialty, state.FormExperience)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Изменить анкету", "fill_form"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendEntertainment(chatID int64) {
	jokes := []string{
		"Почему программисты не любят природу? Слишком много багов! 🐛",
		"Как называется группа программистов? Git-ара! 🎸",
		"Почему Java-разработчик носит очки? Потому что он не видит C#! 👓",
	}

	text := "😊 Развлечение\n\n" + jokes[time.Now().Unix()%int64(len(jokes))]

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("😂 Ещё шутку", "entertainment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendEarnTogether(chatID int64) {
	text := `💸 Зарабатывать вместе

Приглашайте друзей и получайте бонусы!

За каждого приглашённого друга вы получите:
• 100 бонусных баллов
• Приоритетный показ вашей анкеты

Ваша реферальная ссылка: t.me/work_kg_bot?start=ref_` + fmt.Sprintf("%d", chatID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendSubscription(chatID int64) {
	text := `✅ Подписка

Преимущества подписки:
• Приоритетный показ вашей анкеты
• Доступ к премиум вакансиям
• Уведомления о новых вакансиях

Стоимость: 500 сом/месяц

Для оплаты свяжитесь с администратором.`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendCategorySelection(chatID int64, searchType string) {
	var text string
	if searchType == "employee" {
		text = "Мы в разделе поиска сотрудника!\nВыберите в какой сфере ищем. 👇"
	} else {
		text = "Мы в разделе поиска работы!\nВыберите в какой сфере ищем. 👇"
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for category := range models.Categories {
		emoji := models.CategoryEmojis[category]
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(category+" "+emoji, fmt.Sprintf("category:%s:%s", category, searchType)),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendSubcategorySelection(chatID int64, category string, searchType string) {
	text := fmt.Sprintf("Выберите узкую специальность %s! 👇", strings.ToLower(category))

	subcategories := models.Categories[category]
	var rows [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subcategories); i += 2 {
		if i+1 < len(subcategories) {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(subcategories[i], fmt.Sprintf("subcategory:%s:%s", subcategories[i], searchType)),
				tgbotapi.NewInlineKeyboardButtonData(subcategories[i+1], fmt.Sprintf("subcategory:%s:%s", subcategories[i+1], searchType)),
			))
		} else {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(subcategories[i], fmt.Sprintf("subcategory:%s:%s", subcategories[i], searchType)),
			))
		}
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", fmt.Sprintf("search_%s", searchType)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendCitySelection(chatID int64, searchType string) {
	text := "Выберите ваш город 👇"

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Бишкек 🇰🇬", fmt.Sprintf("city:Бишкек:%s", searchType)),
		tgbotapi.NewInlineKeyboardButtonData("Ош 🇰🇬", fmt.Sprintf("city:Ош:%s", searchType)),
		tgbotapi.NewInlineKeyboardButtonData("Талас 🇰🇬", fmt.Sprintf("city:Талас:%s", searchType)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Нарын 🇰🇬", fmt.Sprintf("city:Нарын:%s", searchType)),
		tgbotapi.NewInlineKeyboardButtonData("Каракол 🇰🇬", fmt.Sprintf("city:Каракол:%s", searchType)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Жалал-Абад 🇰🇬", fmt.Sprintf("city:Жалал-Абад:%s", searchType)),
		tgbotapi.NewInlineKeyboardButtonData("Чолпон-Ата 🇰🇬", fmt.Sprintf("city:Чолпон-Ата:%s", searchType)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func showJobs(chatID int64, state *models.UserState) {
	jobs, err := database.SearchJobs(state.Category, state.Subcategory, state.City)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Ошибка при поиске вакансий")
		Bot.Send(msg)
		return
	}

	if len(jobs) == 0 {
		text := "😔 К сожалению, вакансий по вашему запросу не найдено.\n\nПопробуйте изменить параметры поиска."
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔍 Искать снова", "search_job"),
				tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "menu"),
			),
		)
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		Bot.Send(msg)
		return
	}

	for _, job := range jobs {
		text := fmt.Sprintf("📋 *%s*\n\n", job.Title)
		text += fmt.Sprintf("📍 Город: %s\n", job.City)
		text += fmt.Sprintf("📂 Категория: %s / %s\n", job.Category, job.Subcategory)
		if job.Salary != "" {
			text += fmt.Sprintf("💰 Зарплата: %s\n", job.Salary)
		}
		if job.Company != "" {
			text += fmt.Sprintf("🏢 Компания: %s\n", job.Company)
		}
		if job.Description != "" {
			text += fmt.Sprintf("\n📝 %s\n", job.Description)
		}
		text += fmt.Sprintf("\n📞 Контакт: %s", job.Phone)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		Bot.Send(msg)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 Искать ещё", "search_job"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "menu"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Найдено %d вакансий", len(jobs)))
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}

func sendAddVacancyPrompt(chatID int64, state *models.UserState) {
	text := fmt.Sprintf("📋 Поиск сотрудника\n\n📂 Категория: %s / %s\n📍 Город: %s\n\nВы можете добавить вакансию, чтобы найти сотрудника.",
		state.Category, state.Subcategory, state.City)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить вакансию", "add_vacancy"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	Bot.Send(msg)
}
