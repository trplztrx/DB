package models

import "time"

type OrderCourier struct {
	CourierID  uint      `gorm:"primaryKey"`
	OrderID    uint      `gorm:"primaryKey"`
	AssignedAt time.Time `gorm:"not null;default:current_timestamp"`
	Courier    Courier   `gorm:"foreignKey:CourierID"`
	Order      Order     `gorm:"foreignKey:OrderID"`
}
