package models

import (
	"gorm.io/gorm"
)

type Informasi struct {
	gorm.Model
	Judul      string
	Konten     string
	UrlDokumen string
}
