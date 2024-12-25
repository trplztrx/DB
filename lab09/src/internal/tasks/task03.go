package task

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OperationType string

const (
	Insert OperationType = "insert"
	Update OperationType = "update"
	Delete OperationType = "delete"
	None   OperationType = "none"
)

func RunScenario(ctx context.Context, pool *pgxpool.Pool, opType OperationType) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Сценарий %s завершен", opType)
			return
		case <-ticker.C:
			executeOperation(ctx, pool, opType)
		}
	}
}

func executeOperation(ctx context.Context, pool *pgxpool.Pool, opType OperationType) {
	switch opType {
	case Insert:
		deliveryAddrID := rand.Intn(1000) + 1
		addNewRows(ctx, pool, deliveryAddrID)
	case Delete:
		orderID := rand.Intn(1000) + 1
		deleteRows(ctx, pool, orderID)
	case Update:
		deliveryAddrID := rand.Intn(1000) + 1
		orderID := rand.Intn(1000) + 1
		updateRows(ctx, pool, deliveryAddrID, orderID)
	case None:
		log.Println("Сценарий без изменений: пропускаем операции")
	default:
		log.Printf("Неизвестный тип операции: %s", opType)
	}
}

func addNewRows(ctx context.Context, pool *pgxpool.Pool, deliveryAddrID int) {
	query := `
		insert into orders (created_at, updated_at, status, delivery_cost, sender_user_id, receiver_user_id, pickup_address_id, delivery_address_id)
		values (now(), now(), 'доставлен', 500.00, 1, 2, 10, $1);
	`
	_, err := pool.Exec(ctx, query, deliveryAddrID)
	if err != nil {
		log.Printf("Ошибка добавления строк: %v", err)
	} else {
		log.Println("Добавлена новая строка")
	}
}

func deleteRows(ctx context.Context, pool *pgxpool.Pool, orderID int) {
	_, err := pool.Exec(ctx, `delete from orders where id = $1;`, orderID)
	if err != nil {
		log.Printf("Ошибка удаления строк: %v", err)
	} else {
		log.Println("Удалена строка")
	}
}

func updateRows(ctx context.Context, pool *pgxpool.Pool, deliveryAddrID int, orderID int) {
	query := `
		update orders
		set delivery_address_id = $1
		where id = $2;
	`
	_, err := pool.Exec(ctx, query, deliveryAddrID, orderID)
	if err != nil {
		log.Printf("Ошибка изменения строк: %v", err)
	} else {
		log.Println("Обновлена строка")
	}
}
