package models

import (
	"gorm.io/gorm"
)

type Penggajian struct {
	gorm.Model
	NamaPegawai string
	GajiPokok   float64
	JamLembur   int
	GajiKotor   float64
	Pajak       float64
	GajiBersih  float64
}
