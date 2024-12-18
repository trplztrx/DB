package models

type OrderItem struct {
	ID       uint   `gorm:"primaryKey"`
	OrderID  uint   `gorm:"not null"`
	ItemName string `gorm:"type:varchar(255);not null"`
	Quantity int    `gorm:"not null;check:quantity > 0"`
	Order    Order  `gorm:"foreignKey:OrderID"`
}
