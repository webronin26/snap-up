package entities

import "time"

type Product struct {
	ID            int32     `gorm:"primary_key"` // 產品的貨號
	MaxBuy        int32     `gorm:"not null;"`   // 一次最多可以買的購買量
	ItemNumber    int32     `gorm:"not null;"`   // 目前產品的剩餘數量
	AvailableTime time.Time // 什麼時候開放購買
}
