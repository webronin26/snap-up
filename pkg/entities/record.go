package entities

import "time"

// 這邊紀錄產品被購買的紀錄
type Record struct {
	ID         int32     `gorm:"primary_key"` // 每筆訂單的代號
	ProductID  int32     // 產品的 ID
	CustomerID int32     // 購買人的 ID
	Number     int32     // 此次購買的數量
	Date       time.Time // 購買的時間
}
