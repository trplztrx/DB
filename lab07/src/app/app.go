package app

import (
	"fmt"
	"lab07/config"
	"lab07/models"
	pgs "lab07/postgres"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RunApp(cfg *config.Config) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password ,cfg.Host, cfg.Port, cfg.DBConfig.DatabaseName)
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	for {
		showMenu()
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Некорректный ввод")
			continue
		}
		switch choice {
		case 1:
			err = pgs.GetFirstUsers(db, "./data/1.json")
			if err != nil {
				log.Fatalf("Ошибка запроса %v", err)
			}
		case 2:
			var userID int
			fmt.Println("Введите id пользователя: ")
			fmt.Scanln(&userID)
			err = pgs.GetOrdersByUser(db, "./data/2.json", userID)
			if err != nil {
				log.Fatalf("Ошибка запроса %v", err)
			}
		case 3:
			var n int
			fmt.Println("Введите кол-во заказов N: ")
			fmt.Scanln(&n)
			err = pgs.GetCouriersWithOrdersN(db, "./data/3.json", n)
			if err != nil {
				log.Fatalf("Ошибка запроса %v", err)
			}
		case 4:
			var status string
			fmt.Println("Введите статус ('создан', 'в пути', 'доставлен', 'отменен'): ")
			fmt.Scanln(&status)
			err = pgs.GetUsersByOrderStatus(db, "./data/4.json", status)
			if err != nil {
				log.Fatalf("Ошибка запроса %v", err)
			}
		case 5:
			var prefix string
			fmt.Println("Введите Начало фамилии: ")
			fmt.Scanln(&prefix)
			err = pgs.GetCouriersBySurnamePrefix(db, "./data/5.json", prefix)
			if err != nil {
				log.Fatalf("Ошибка запроса %v", err)
			}
		case 10:
			_, err = pgs.ReadJSON("./data/10.json")
			if err != nil {
				log.Fatalf("Ошибка запроса %v", err)
			}
		case 11:
			err := pgs.UpdateJSON("./data/11.json", func(user *models.User) {
				if strings.HasPrefix(user.Phone, "+7") {
					user.Name = "Обновлённый пользователь"
				}
			})
			if err != nil {
				log.Fatalf("Ошибка запроса %v", err)
			}
		case 12:
			newUsers := []models.User{
				{ID: 5000, UserType: "физ", Name: "Иван Иванов", Phone: "+79161234567", Email: "ivan@example.com"},
			}
			err = pgs.AddToJSON("./data/12.json", newUsers)
			if err != nil {
				log.Fatalf("Ошибка запроса %v", err)
			}
		case 20:
			fmt.Println("Все физ лица")
			_, err := pgs.GetPhysicalUsers(db, "./data/20.json")
			if err != nil {
    			log.Fatalf("Ошибка: %v", err)
			}
		case 21:
			fmt.Println("Все заказы с полной ифнормацией")
			_, err := pgs.GetOrdersWithDetails(db, "./data/21.json")
			if err != nil {
    			log.Fatalf("Ошибка: %v", err)
			}
		case 22:
			newUser := &models.User{ID: 5000, UserType: "физ", Name: "Иван Иванов", Phone: "+79161234567", Email: "ivan@example.com"}
			fmt.Println("Добавить пользователя в БД")
			err := pgs.AddNewUser(db, newUser)
			if err != nil {
    			log.Fatalf("Ошибка: %v", err)
			}
		case 23:
			err := pgs.UpdateUser(db, 5000)
			if err != nil {
    			log.Fatalf("Ошибка: %v", err)
			}
		case 24:
			err := pgs.DeleteUser(db, 5000)
			if err != nil {
    			log.Fatalf("Ошибка: %v", err)
			}
		case 25:
			err = pgs.CallAddOrUpdateOrderItem(db, 3, "товар_А", 10)
			if err != nil {
				log.Fatalf("Ошибка выполнения процедуры: %v", err)
			}
		case 0:
			return
		default:
			fmt.Println("Некорректный выбор")
		}
	}
}

func showMenu() {
	fmt.Println(`
Задание 1
	1. Первые 10 пользователей, отсортированных по имени в алфавитном порядке
	2. Получение заказов выбранного пользователя
	3. Курьеры с более чем N заказами
	4. Отправители, с выбранным статсусом заказа
	5. Поиск курьеров по фамилии

Задание 2
	10. Чтение из 1.json. Вывод в 10.json
	11. Обновлние пользователей с номером телефона +7
	12. Добавление пользователя

Задание 3
	20. Однотабличный запрос на выборку
	21. Многотабличный запрос на выборку
	22. Запрос на добавление
	23. Запрос на изменение
	24. Запрос на удаение
	25. Хранимая процедура

0. Выход

Введите номер пункта:`)
}