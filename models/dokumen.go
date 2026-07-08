package models

import (
	"gorm.io/gorm"
)

type Dokumen struct {
	gorm.Model
	NamaDokumen string
	FileId      string
	FileUrl     string
}
