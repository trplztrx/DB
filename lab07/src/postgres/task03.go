package pgs

import (
	"encoding/json"
	"fmt"
	"lab07/models"
	"os"


	"gorm.io/gorm"
)

func GetPhysicalUsers(db *gorm.DB, filePath string) ([]models.User, error) {
    var users []models.User

    result := db.Where("user_type = ?", "физ").Find(&users)
    if result.Error != nil {
        return nil, fmt.Errorf("ошибка при выборке физических пользователей: %v", result.Error)
    }
	formattedJSON, err := json.MarshalIndent(users, "", "    ")
    if err != nil {
        return nil, fmt.Errorf("ошибка сериализации JSON: %v", err)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return nil, fmt.Errorf("ошибка создания файла %s: %v",filePath, err)
    }
    defer file.Close()
    
    _, err = file.Write(formattedJSON)
    if err != nil {
        return nil, fmt.Errorf("ошибка записи в файл %s: %v", filePath, err)
    }

    return users, nil
}

func GetOrdersWithDetails(db *gorm.DB, filePath string) ([]models.Order, error) {
    var orders []models.Order

    result := db.
        Preload("SenderUser").
        Preload("ReceiverUser").
        Preload("PickupAddress").
        Preload("DeliveryAddress").
        Find(&orders)
    if result.Error != nil {
        return nil, fmt.Errorf("ошибка при многотабличной выборке заказов: %v", result.Error)
    }

	formattedJSON, err := json.MarshalIndent(orders, "", "    ")
    if err != nil {
        return nil, fmt.Errorf("ошибка сериализации JSON: %v", err)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return nil, fmt.Errorf("ошибка создания файла %s: %v",filePath, err)
    }
    defer file.Close()
    
    _, err = file.Write(formattedJSON)
    if err != nil {
        return nil, fmt.Errorf("ошибка записи в файл %s: %v", filePath, err)
    }

    return orders, nil
}

func AddNewUser(db *gorm.DB, newUser *models.User) error {
	result := db.Create(newUser)
	if result.Error != nil {
		return fmt.Errorf("ошибка выполнения функции Create: %v", result.Error)
	}

	fmt.Println("Успешное добавление элемента")
	return nil
}

func UpdateUser(db *gorm.DB, userID uint) error {
	result := db.Model(&models.User{}).Where("id = ?", userID).Update("Name", "SHKA")
	if result.Error != nil {
		return fmt.Errorf("ошибка выполнения функции Update: %v", result.Error)
	}

	fmt.Println("Успешное обновление элемента")
	return nil
}

func DeleteUser(db *gorm.DB, userID uint) error {
	result := db.Delete(&models.User{}, userID)
	if result.Error != nil {
		return fmt.Errorf("ошибка выполнения функции Delete: %v", result.Error)
	}

	fmt.Println("Успешное удаление элемента")
	return nil
}

func CallAddOrUpdateOrderItem(db *gorm.DB, orderID int, itemName string, quantity int) error {
    result := db.Exec(
        "call add_or_update_order_item(?, ?, ?)",
        orderID, itemName, quantity,
    )

    if result.Error != nil {
        return fmt.Errorf("ошибка вызова хранимой процедуры add_or_update_order_item: %v", result.Error)
    }

    return nil
}

