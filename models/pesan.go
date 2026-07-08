package models

import (
	"time"
)

type Pesan struct {
	Kode      string
	Balasan   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
