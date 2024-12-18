package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	UserType string `gorm:"type:varchar(50);not null;check:user_type IN ('физ','юр')"`
	Name     string `gorm:"type:varchar(255);not null"`
	Phone    string `gorm:"type:varchar(50);unique;not null"`
	Email    string `gorm:"type:varchar(100);unique;not null"`
}
