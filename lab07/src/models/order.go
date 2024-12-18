package models

import "time"

type Order struct {
	ID                uint      `gorm:"primaryKey"`
	CreatedAt         time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt         time.Time `gorm:"not null;default:current_timestamp"`
	Status            string    `gorm:"type:varchar(50);not null;check:status IN ('создан','отправлен','доставлен','отменен')"`
	DeliveryCost      float64   `gorm:"type:decimal(10,2);not null;check:delivery_cost > 0"`
	SenderUserID      uint      `gorm:"not null"`
	ReceiverUserID    uint      `gorm:"not null"`
	PickupAddressID   uint      `gorm:"not null"`
	DeliveryAddressID uint      `gorm:"not null"`
	SenderUser        User      `gorm:"foreignKey:SenderUserID"`
	ReceiverUser      User      `gorm:"foreignKey:ReceiverUserID"`
	PickupAddress     Address   `gorm:"foreignKey:PickupAddressID"`
	DeliveryAddress   Address   `gorm:"foreignKey:DeliveryAddressID"`
}
