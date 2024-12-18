package pgs

import (
	"encoding/json"
	"fmt"
	"lab07/models"
	"os"
)

func ReadJSON(filePath string) ([]models.User, error) {
    var users []models.User

    jsonData, err := os.ReadFile("./data/1.json")
    if err != nil {
        return nil, fmt.Errorf("ошибка чтения JSON-файла: %v", err)
    }

    err = json.Unmarshal(jsonData, &users)
    if err != nil {
        return nil, fmt.Errorf("ошибка десериализации JSON: %v", err)
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

func UpdateJSON(filePath string, updateFunc func(user *models.User)) error {
    users, err := ReadJSON(filePath)
    if err != nil {
        return err
    }

    for i := range users {
        updateFunc(&users[i])
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


func AddToJSON(filePath string, newUsers []models.User) error {
    users, err := ReadJSON(filePath)
    if err != nil {
        return err
    }

    users = append(users, newUsers...)

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

