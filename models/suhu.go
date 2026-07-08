package models

import (
	"time"
)

type Suhu struct {
	Id        uint      `gorm:"primaryKey"`
	Lokasi    string    `gorm:"type:varchar(255)"`
	Suhu      float32   `gorm:"type:decimal(10,2)"`
	CreatedAt time.Time `gorm:"datetime"`
}
