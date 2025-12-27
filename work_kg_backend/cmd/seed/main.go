package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"work_kg_backend/internal/config"
	"work_kg_backend/internal/database"
	"work_kg_backend/internal/models"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize database schema
	database.InitSchema()

	// Seed jobs
	seedJobs()

	log.Println("Seeding completed successfully!")
}

func seedJobs() {
	// Check if jobs already exist
	stats := database.GetStats()
	if stats.TotalJobs >= 100 {
		log.Printf("Database already has %d jobs, skipping seed", stats.TotalJobs)
		return
	}

	rand.Seed(time.Now().UnixNano())

	jobData := []struct {
		Category      string
		Subcategories []string
		Titles        []string
		Salaries      []string
		Companies     []string
	}{
		{
			Category:      "Строительство",
			Subcategories: []string{"Каменщик", "Кладка", "Электрик", "Сантехник", "Сварщик", "Отделочник", "Плиточник", "Фасадчик", "Монолитчик", "Разнорабочий"},
			Titles: []string{
				"Требуется %s на постоянную работу",
				"Срочно нужен %s",
				"%s с опытом работы",
				"Ищем %s в бригаду",
				"%s на строительный объект",
			},
			Salaries:  []string{"30000-40000 сом", "40000-50000 сом", "50000-60000 сом", "60000-80000 сом", "от 1500 сом/день"},
			Companies: []string{"СтройМастер", "БишкекСтрой", "Нур-Строй", "Альфа Констракшн", "МегаСтрой", ""},
		},
		{
			Category:      "Общепит",
			Subcategories: []string{"Повар", "Официант", "Бармен", "Посудомойщик", "Администратор", "Кассир"},
			Titles: []string{
				"Требуется %s в ресторан",
				"%s в кафе",
				"Ищем %s на постоянную работу",
				"%s с опытом",
				"Срочно нужен %s",
			},
			Salaries:  []string{"20000-25000 сом", "25000-30000 сом", "30000-35000 сом", "35000-45000 сом", "от 800 сом/день"},
			Companies: []string{"Кофемания", "Вкусно и Точка", "Навруз", "Дастархан", "Супара", ""},
		},
		{
			Category:      "Швейный цех",
			Subcategories: []string{"Швея", "Закройщик", "Упаковщик", "Технолог", "Контролер качества"},
			Titles: []string{
				"Требуется %s в швейный цех",
				"%s на постоянную работу",
				"Ищем опытного %s",
				"%s (срочно)",
				"Нужен %s в производство",
			},
			Salaries:  []string{"25000-30000 сом", "30000-40000 сом", "от 500 сом/день", "35000-45000 сом", "сдельная оплата"},
			Companies: []string{"Текстиль Плюс", "Модный дом", "Швейпром", "Кыргыз Текстиль", "АзияТекс", ""},
		},
		{
			Category:      "IT",
			Subcategories: []string{"Программист", "Дизайнер", "Тестировщик", "Системный администратор"},
			Titles: []string{
				"Требуется %s",
				"%s в IT компанию",
				"Junior/Middle %s",
				"Ищем %s в команду",
				"%s (удаленно)",
			},
			Salaries:  []string{"50000-70000 сом", "70000-100000 сом", "100000-150000 сом", "от 1000$", "договорная"},
			Companies: []string{"Nambavan", "Mad Devs", "Zensoft", "Beeline", "O!", ""},
		},
		{
			Category:      "Продажи",
			Subcategories: []string{"Продавец", "Менеджер", "Консультант", "Кассир"},
			Titles: []string{
				"Требуется %s в магазин",
				"%s-консультант",
				"Ищем %s на постоянную работу",
				"%s (график 2/2)",
				"Нужен %s",
			},
			Salaries:  []string{"20000-25000 сом", "25000-35000 сом", "оклад + %", "30000-40000 сом", "от 700 сом/день"},
			Companies: []string{"Бета Сторес", "Народный", "Глобус", "ЦУМ", "Детский мир", ""},
		},
		{
			Category:      "Транспорт",
			Subcategories: []string{"Водитель", "Курьер", "Экспедитор", "Диспетчер"},
			Titles: []string{
				"Требуется %s",
				"%s категории B/C",
				"Ищем %s на личном авто",
				"%s (срочно)",
				"Нужен %s в службу доставки",
			},
			Salaries:  []string{"35000-45000 сом", "40000-55000 сом", "от 1000 сом/день", "50000-70000 сом", "сдельная"},
			Companies: []string{"Яндекс Доставка", "Глово", "Достар", "Карго Транс", "Логистик Сервис", ""},
		},
	}

	descriptions := []string{
		"Требования: опыт работы от 1 года, ответственность, пунктуальность. Мы предлагаем: стабильную зарплату, официальное трудоустройство.",
		"Обязанности: выполнение поставленных задач качественно и в срок. Условия: график 5/2, обед, соц. пакет.",
		"Ищем ответственного сотрудника. Опыт приветствуется. Обучение на месте. Дружный коллектив.",
		"Работа на постоянной основе. Своевременная оплата. Возможен карьерный рост.",
		"Требуется сотрудник с опытом. Высокая зарплата. Хорошие условия труда.",
	}

	phones := []string{
		"+996 555 123 456",
		"+996 700 987 654",
		"+996 770 111 222",
		"+996 550 333 444",
		"+996 705 555 666",
		"+996 777 888 999",
		"+996 500 000 111",
		"+996 707 222 333",
	}

	jobCount := 0
	for jobCount < 100 {
		for _, data := range jobData {
			if jobCount >= 100 {
				break
			}

			subcategory := data.Subcategories[rand.Intn(len(data.Subcategories))]
			titleTemplate := data.Titles[rand.Intn(len(data.Titles))]
			title := fmt.Sprintf(titleTemplate, subcategory)
			city := models.Cities[rand.Intn(len(models.Cities))]
			salary := data.Salaries[rand.Intn(len(data.Salaries))]
			company := data.Companies[rand.Intn(len(data.Companies))]
			description := descriptions[rand.Intn(len(descriptions))]
			phone := phones[rand.Intn(len(phones))]

			job := &models.Job{
				Title:       title,
				Description: description,
				Category:    data.Category,
				Subcategory: subcategory,
				City:        city,
				Salary:      salary,
				Phone:       phone,
				Company:     company,
				IsActive:    true,
				Source:      "seed",
			}

			if err := database.SaveJob(job); err != nil {
				log.Printf("Error inserting job: %v", err)
				continue
			}

			jobCount++
			log.Printf("Created job %d: %s in %s", jobCount, title, city)
		}
	}

	log.Printf("Successfully created %d jobs", jobCount)
}
