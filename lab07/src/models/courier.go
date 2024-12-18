package models

type Courier struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"type:varchar(255);not null"`
	Surname    string `gorm:"type:varchar(255);not null"`
	Patronymic string `gorm:"type:varchar(255);not null"`
	Phone      string `gorm:"type:varchar(50);unique;not null"`
	Status     string `gorm:"type:varchar(50);not null;check:status in ('свободен','в пути','занят')"`
}
