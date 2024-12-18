package app

import (
	"context"
	"fmt"
	"lab06/config"
	postgres "lab06/requests"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunApp(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password ,cfg.Host, cfg.Port, cfg.DBConfig.DatabaseName)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Ошибка подключения к постгрес: %v", err.Error())
	}
	defer pool.Close()

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
			fmt.Println("Количество заказов в таблице orders")
			err = postgres.GetCountOrder(ctx, pool)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 2:
			var city string
			fmt.Println("Список отправителей из выбранного города")
			fmt.Println("Введите город: ")
			fmt.Scanln(&city)
			err = postgres.GetSendersByCity(ctx, pool, city)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 3:
			fmt.Println("min/max/avg cатистика по стоимости доставки заказов")
			err = postgres.GetOrderDeliveryCostStat(ctx, pool)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 4:
			fmt.Println("Список всех отношений в БД")
			err = postgres.GetTablesMetadata(ctx, pool)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 5:
			var senderID int
			fmt.Println("Общая стоимость доставки для заказов, сделанных выбранным отправителем")
			fmt.Println("Введите id отправителя: ")
			fmt.Scanln(&senderID)
			err = postgres.GetTotalDeliveryCost(ctx, pool, senderID)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 6:
			var senderID int
			fmt.Println("Список заказо, сделанных выбранным отправителем")
			fmt.Println("Введите id отправителя: ")
			fmt.Scanln(&senderID)
			err = postgres.GetSenderOrderDetails(ctx, pool, senderID)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 7:
			var orderID, quantity int
			var itemName string
			fmt.Println("Список заказо, сделанных выбранным отправителем")
			fmt.Println("Введите id заказа: ")
			fmt.Scanln(&orderID)
			fmt.Println("Введите нименование товара: ")
			fmt.Scanln(&itemName)
			fmt.Println("Введите количество товара: ")
			fmt.Scanln(&quantity)
			err = postgres.AddOrderItem(ctx, pool, orderID, itemName, quantity)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 8:
			fmt.Println("Текущее имя БД")
			err = postgres.GetCurrentDBName(ctx, pool)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 9:
			fmt.Println("Создание таблицы Items")
			err = postgres.CreateTableItems(ctx, pool)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 10:
			fmt.Println("Выполнение вставки данных в таблицу Items")
			err = postgres.UpdateTableItems(ctx, pool)
			if err != nil {
				log.Fatalf("Ошибка запроса: %v", err)
			}
		case 11:
			_ = postgres.SelectAll(ctx, pool)
		case 0:
			return
		default:
			fmt.Println("Некорректный выбор")
		}
	}
}

func showMenu() {
		fmt.Println(`
Меню:
	1. Выполнить скалярный запрос
	2. Выполнить запрос с несколькими соединениями (JOIN)
	3. Выполнить запрос с ОТВ (CTE) и оконными функциями
	4. Выполнить запрос к метаданным
	5. Вызвать скалярную функцию (третья лабораторная работа)
	6. Вызвать многооператорную или табличную функцию (третья лабораторная работа)
	7. Вызвать хранимую процедуру (третья лабораторная работа)
	8. Вызвать системную функцию или процедуру
	9. Создать таблицу в базе данных, соответствующую тематике БД
	10. Выполнить вставку данных в созданную таблицу с использованием INSERT или COPY
	0. Выход

Введите номер пункта:`)
}