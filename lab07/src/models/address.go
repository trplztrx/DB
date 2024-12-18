package models

type Address struct {
	ID            uint   `gorm:"primaryKey"`
	Region        string `gorm:"type:varchar(100);not null"`
	City          string `gorm:"type:varchar(100);not null"`
	AddressStreet string `gorm:"type:varchar(255);not null"`
	AddressType   string `gorm:"type:varchar(50);not null;check:address_type IN ('отправитель','получатель')"`
}
