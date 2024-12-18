package pgs

import (
	"encoding/json"
	"fmt"
	"lab07/models"
	"os"

	"gorm.io/gorm"
)

func GetFirstUsers(db *gorm.DB, filePath string) error {
    var users []models.User
    result := db.
        Select("*").
        Where("user_type = ?", "физ").
        Order("name asc").
        Limit(10).
        Find(&users)
    if result.Error != nil {
        return fmt.Errorf("ошибка выполнения запроса GetFirstUsers: %v", result.Error)
    }

    formattedJSON, err := json.MarshalIndent(users, "", "    ")
    if err != nil {
        return fmt.Errorf("ошибка сериализации JSON: %v", err)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("ошибка создания файла %s: %v",filePath, err)
    }
    defer file.Close()

    _, err = file.Write(formattedJSON)
    if err != nil {
        return fmt.Errorf("ошибка записи в файл %s: %v", filePath, err)
    }

    return nil
}

func GetOrdersByUser(db *gorm.DB, filePath string, userID int) error {
    var orders []models.Order
    result := db.
        Preload("PickupAddress").
        Preload("DeliveryAddress").
        Preload("SenderUser").
        Preload("ReceiverUser").
        Joins("join users as sender on orders.sender_user_id = sender.id").
        Joins("join users as receiver on orders.receiver_user_id = receiver.id").
        Where("sender.id = ? or receiver.id = ?", userID, userID).
        Find(&orders)
    if result.Error != nil {
        return fmt.Errorf("ошибка создания файла %s: %v",filePath, result.Error)
    }

    formattedJSON, err := json.MarshalIndent(orders, "", "    ")
    if err != nil {
        return fmt.Errorf("ошибка сериализации JSON: %v", err)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("ошибка создания файла %s: %v",filePath, err)
    }
    defer file.Close()
    
    _, err = file.Write(formattedJSON)
    if err != nil {
        return fmt.Errorf("ошибка записи в файл %s: %v", filePath, err)
    }
    return nil
}

func GetCouriersWithOrdersN(db *gorm.DB, filePath string, n int) error {
    var couriers []models.Courier

    subquery := db.
        Table("ordercourier").
        Select("courier_id").
        Group("courier_id").
        Having("count(order_id) > ?", n)

    result := db.
        Model(&models.Courier{}).
        Where("id in (?)", subquery).
        Find(&couriers)

    if result.Error != nil {
        return fmt.Errorf("ошибка выполнения запроса GetCouriersWithOrdersN: %v", result.Error)
    }

    formattedJSON, err := json.MarshalIndent(couriers, "", "    ")
    if err != nil {
        return fmt.Errorf("ошибка сериализации JSON: %v", err)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("ошибка создания файла %s: %v", filePath, err)
    }
    defer file.Close()

    _, err = file.Write(formattedJSON)
    if err != nil {
        return fmt.Errorf("ошибка записи в файл %s: %v", filePath, err)
    }

    return nil
}


func GetItemsFromOrder(db *gorm.DB, filePath string, orderStatus string) error {
    var items []models.OrderItem
    
    subQuery := db.
        Table("orders").
        Select("id").
        Where("status = ?", orderStatus)

    result := db.
        Table("orderitem").
        Where("order_id in (?)", subQuery).
        Find(&items)
    if result.Error != nil {
        return fmt.Errorf("ошибка выполнения запроса GetItemsFromOrders: %v", result.Error)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("ошибка создания файла %s: %v", filePath, err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    err = encoder.Encode(items)
    if err != nil {
        return fmt.Errorf("ошибка записи данных в файл %s: %v", filePath, err)
    }

    return nil
}



func GetCouriersByStatus(db *gorm.DB, filePath string, orderStatus string) error {
    var couriers []models.Courier
    result := db.Select("*").
        Where("status = ?", orderStatus).
        Order("surname asc").
        Find(&couriers)
    if result.Error != nil {
        return fmt.Errorf("ошибка выполнения запроса GetCouriersInTransit: %v", result.Error)
    }

    formattedJSON, err := json.MarshalIndent(couriers, "", "    ")
    if err != nil {
        return fmt.Errorf("ошибка сериализации JSON: %v", err)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("ошибка создания файла %s: %v", filePath, err)
    }
    defer file.Close()

    _, err = file.Write(formattedJSON)
    if err != nil {
        return fmt.Errorf("ошибка записи в файл %s: %v", filePath, err)
    }

    return nil
}

func GetUsersByOrderStatus(db *gorm.DB, filePath string, orderStatus string) error {
    var users []models.User

    subQuery := db.Model(&models.Order{}).
        Select("sender_user_id").
        Where("status = ?", orderStatus)

    result := db.Model(&models.User{}).
        Where("id in (?)", subQuery).
        Find(&users)

    if result.Error != nil {
        return fmt.Errorf("ошибка выполнения запроса GetUsersByOrderStatus: %v", result.Error)
    }

    formattedJSON, err := json.MarshalIndent(users, "", "    ")
    if err != nil {
        return fmt.Errorf("ошибка сериализации JSON: %v", err)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("ошибка создания файла %s: %v", filePath, err)
    }
    defer file.Close()

    _, err = file.Write(formattedJSON)
    if err != nil {
        return fmt.Errorf("ошибка записи в файл %s: %v", filePath, err)
    }

    return nil
}

func GetCouriersBySurnamePrefix(db *gorm.DB, filePath string, prefix string) error {
    var couriers []models.Courier

    likePattern := prefix + "%"

    result := db.
        Where("surname like ?", likePattern).
        Find(&couriers)

    if result.Error != nil {
        return fmt.Errorf("ошибка выполнения запроса GetCouriersBySurnamePrefix: %v", result.Error)
    }

    formattedJSON, err := json.MarshalIndent(couriers, "", "    ")
    if err != nil {
        return fmt.Errorf("ошибка сериализации JSON: %v", err)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("ошибка создания файла %s: %v", filePath, err)
    }
    defer file.Close()

    _, err = file.Write(formattedJSON)
    if err != nil {
        return fmt.Errorf("ошибка записи в файл %s: %v", filePath, err)
    }

    return nil
}
