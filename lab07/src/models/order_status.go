package models

import "time"

type OrderStatus struct {
	ID        uint      `gorm:"primaryKey"`
	OrderID   uint      `gorm:"not null"`
	Status    string    `gorm:"type:varchar(50);not null;check:status IN ('создан','отправлен','доставлен','отменен')"`
	UpdatedAt time.Time `gorm:"not null;default:current_timestamp"`
	Order     Order     `gorm:"foreignKey:OrderID"`
}
