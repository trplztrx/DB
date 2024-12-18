package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// 1
func GetCountOrder(ctx context.Context, dbPool *pgxpool.Pool) error {
	var result string
	query := `select count(*) from orders`
	err := dbPool.QueryRow(ctx, query).Scan(&result)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса GetCountOrder: %v", err)
	}
	fmt.Println(result)
	return nil
}

// 2
func GetSendersByCity(ctx context.Context, dbPool *pgxpool.Pool, city string) error {
	query := `
	select u.name, u.user_type, a.city
	from users u
	join orders o on o.sender_user_id = u.id
	join addresses a on o.pickup_address_id = a.id
	where a.city = $1`
	rows, err := dbPool.Query(ctx, query, city)
    if err != nil {
        return fmt.Errorf("ошибка выполнения запроса GetSendersByCity: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var name, userType, city string
        err := rows.Scan(&name, &userType, &city)
        if err != nil {
            return fmt.Errorf("ошибка чтения строк: %v", err)
        }
        fmt.Printf("%s, %s, %s\n", name, userType, city)
    }

    if rows.Err() != nil {
        return fmt.Errorf("ошибка обработки строк: %v", rows.Err())
    }

    return nil
}

// 3
func GetOrderDeliveryCostStat(ctx context.Context, dbPool *pgxpool.Pool) error {
    query := `
    with order_aggregates as (
        select 
            id as order_id,
            status,
            delivery_cost,
            min(delivery_cost) over (partition by status) as min_delivery_cost,
            max(delivery_cost) over (partition by status) as max_delivery_cost,
            avg(delivery_cost) over (partition by status) as avg_delivery_cost
        from orders
    )
    select 
        order_id,
        status,
        delivery_cost,
        min_delivery_cost,
        max_delivery_cost,
        avg_delivery_cost
    from order_aggregates
    order by status, delivery_cost;`

    rows, err := dbPool.Query(ctx, query)
    if err != nil {
        return fmt.Errorf("ошибка выполнения SQL-запроса: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var (
            orderID         int
            status          string
            deliveryCost    float64
            minDeliveryCost float64
            maxDeliveryCost float64
            avgDeliveryCost float64
        )

        err := rows.Scan(&orderID, &status, &deliveryCost, &minDeliveryCost, &maxDeliveryCost, &avgDeliveryCost)
        if err != nil {
            return fmt.Errorf("ошибка сканирования строки результата: %v", err)
        }

        fmt.Printf("%d %s %f %f %f %f\n", orderID, status, deliveryCost, minDeliveryCost, maxDeliveryCost, avgDeliveryCost)
    }

    if rows.Err() != nil {
        return fmt.Errorf("ошибка обработки строк результата: %v", rows.Err())
    }

    return nil
}

// 4
func GetTablesMetadata(ctx context.Context, dbPool *pgxpool.Pool) error {
    query := `
        select schemaname || '.' || tablename
        from pg_catalog.pg_tables
        where schemaname not in ('pg_catalog', 'information_schema')`
    rows, err := dbPool.Query(ctx, query)
    if err != nil {
        return fmt.Errorf("ошибка выполнения запроса к метаданным: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var tableName string
        if err := rows.Scan(&tableName); err != nil {
            return fmt.Errorf("ошибка чтения строки результата: %v", err)
        }
		fmt.Println(tableName)
    }

    if rows.Err() != nil {
        return fmt.Errorf("ошибка обработки строк результата: %v", rows.Err())
    }

    return nil
}


// 5
func GetTotalDeliveryCost(ctx context.Context, dbPool *pgxpool.Pool, senderID int) error {
	var result string

	query := `select get_total_delivery_cost($1)`

	err := dbPool.QueryRow(ctx, query, senderID).Scan(&result)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса GetTotalDeliveryCost: %v", err)
	}
	fmt.Println(result)
	return nil
}

// 6
func GetSenderOrderDetails(ctx context.Context, dbPool *pgxpool.Pool, senderID int) error { 
	query := `select * from get_sender_order_details($1)`

	rows, err := dbPool.Query(ctx, query, senderID)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса GetSenderOrderID: %v", err)
	}
	type OrderDetail struct {
        OrderID           int
        DeliveryCost      float64
        Status            string
        UpdatedAt         time.Time
        TotalOrders       int
        TotalDeliveryCost float64
    }
	for rows.Next() {
		var result OrderDetail
        err := rows.Scan(&result.OrderID, &result.DeliveryCost, &result.Status, &result.UpdatedAt, &result.TotalOrders, &result.TotalDeliveryCost)
        if err != nil {
            return fmt.Errorf("ошибка чтения строк: %v", err)
        }
        fmt.Printf("%d, %f, %s, %s, %d, %f\n", result.OrderID, result.DeliveryCost, result.Status, result.UpdatedAt.String(), result.TotalOrders, result.TotalDeliveryCost)
    }
	if rows.Err() != nil {
        return fmt.Errorf("ошибка обработки строк: %v", rows.Err())
    }
	return nil
}

// 7
func AddOrderItem(ctx context.Context, dbPool *pgxpool.Pool, orderID int, itemName string, quantity int) error {
	query := `call add_or_update_order_item($1, $2, $3)`
	_, err := dbPool.Exec(ctx, query, orderID, itemName, quantity)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса AddOrderItem: %v", err)
	}
	return nil
}

// 8
func GetCurrentDBName(ctx context.Context, dbPool *pgxpool.Pool) error {
    query := `select current_database()`
    var currentDB string

    err := dbPool.QueryRow(ctx, query).Scan(&currentDB)
    if err != nil {
        return fmt.Errorf("ошибка вызова системной функции current_database(): %v", err)
    }

	fmt.Println(currentDB)
    return nil
}

// 9
func CreateTableItems(ctx context.Context, dbPool *pgxpool.Pool) error {
	query := `
	create table if not exists items (
    id serial primary key,               -- Уникальный идентификатор товара
    name varchar(150) not null,          -- Название товара
    description text,                    -- Описание товара
    quantity integer not null default 0, -- Количество товара на складе
    price decimal(10, 2) not null,       -- Цена товара
    created_at timestamp default now(),  -- Дата добавления товара
    updated_at timestamp default now()   -- Дата последнего обновления информации о товаре
	);`
	_, err := dbPool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса CreateTableItems: %v", err)
	}
	fmt.Println("Таблица Itmes успешно создана")
	return nil
}

// 10
func UpdateTableItems(ctx context.Context, dbPool *pgxpool.Pool) error {
	query := `
	insert into items (name, description, quantity, price) values
    ('Ноутбук', 'Ноутбук с процессором Intel Core i7, 16GB RAM, 512GB SSD', 25, 75000.00),
    ('Смартфон', 'Смартфон с экраном 6.5 дюймов, 128GB памяти', 50, 40000.00),
    ('Планшет', 'Планшет с экраном 10 дюймов, 64GB памяти', 30, 25000.00),
    ('Монитор', 'Монитор 24 дюйма, Full HD', 10, 12000.00);`
    _, err := dbPool.Exec(ctx, query)
    if err != nil {
        return fmt.Errorf("ошибка выполнения запроса UpdateTableItems: %v", err)
    }
    fmt.Println("Таблица Itmes успешно обновлена")
	return nil
}

func SelectAll(ctx context.Context, dbPool *pgxpool.Pool) error {
	query := "select * from items"
	rows, err := dbPool.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса SelectAll: %v", err)
	}
	defer rows.Close()

	type Item struct {
		Id          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Quantity    int       `json:"quantity"`
		Price       float64   `json:"price"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	var items []Item

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Id, &item.Name, &item.Description, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return fmt.Errorf("ошибка чтения строк: %v", err)
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return fmt.Errorf("ошибка обработки строк: %v", rows.Err())
	}

	outputFile := "./data/items.json"

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("ошибка создания файла JSON: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(items)
	if err != nil {
		return fmt.Errorf("ошибка записи данных в JSON: %v", err)
	}

	fmt.Printf("Данные успешно записаны в файл: %s\n", outputFile)
	return nil
}