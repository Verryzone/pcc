package controllers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Binding dari POST JSON
type StrukturPenggajian struct {
	Id          uint
	NamaPegawai string  `binding:"required"`
	GajiPokok   float64 `binding:"required"`
	JamLembur   int     `binding:"required"`
	GajiKotor   float64
	Pajak       float64
	GajiBersih  float64
}

type StrukturPenggajianHapus struct {
	Id uint `binding:"required"`
}

// Tambahkan ini
func TampilPenggajian(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//buat variabel array dari model Penggajian
	var modelPenggajian []models.Penggajian
	hasil := db.Find(&modelPenggajian)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil Tampil data",
			"kesalahan": nil,
			"data":      modelPenggajian,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tampil Data",
			"kesalahan": kesalahan.Error(),
			"data":      nil,
		})
	}
}
func TambahPenggajian(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data Penggajian dengan struktur Penggajian
	// dan menangkap data dari request
	var dataPenggajian StrukturPenggajian
	if err := c.ShouldBindJSON(&dataPenggajian); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	// perhitungan gaji kotor, pajak, dan gaji bersih
	var jam_lembur = dataPenggajian.JamLembur
	var uang_lembur = jam_lembur * 50000

	var gaji_pokok = dataPenggajian.GajiPokok
	var gaji_kotor = float64(uang_lembur) + gaji_pokok

	var pajak = 0.0
	var gaji_bersih = gaji_kotor - pajak

	if gaji_kotor > 5000000 {
		pajak = gaji_kotor * 0.05
		gaji_bersih = gaji_kotor - pajak
	}

	//membuat data baru dengan model Penggajian
	modelPenggajian := models.Penggajian{
		NamaPegawai: dataPenggajian.NamaPegawai,
		GajiPokok:   dataPenggajian.GajiPokok,
		JamLembur:   dataPenggajian.JamLembur,
		GajiKotor:   gaji_kotor,
		Pajak:       pajak,
		GajiBersih:  gaji_bersih,
	}
	hasil := db.Create(&modelPenggajian)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data",
			"kesalahan": nil,
			"data":      modelPenggajian,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tambah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelPenggajian,
		})
	}
}

func UbahPenggajian(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data Penggajian dengan struktur Penggajian
	//dan menangkap data dari request
	var dataPenggajian StrukturPenggajian
	if err := c.ShouldBindJSON(&dataPenggajian); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat variabel model Penggajian
	var modelPenggajian models.Penggajian
	//mencari data Penggajian dan merubah datanya
	db.First(&modelPenggajian, dataPenggajian.Id)

	// perhitungan gaji kotor, pajak, dan gaji bersih
	var jam_lembur = dataPenggajian.JamLembur
	var uang_lembur = jam_lembur * 50000

	var gaji_pokok = dataPenggajian.GajiPokok
	var gaji_kotor = float64(uang_lembur) + gaji_pokok

	var pajak = 0.0
	var gaji_bersih = gaji_kotor - pajak

	if gaji_kotor > 5000000 {
		pajak = gaji_kotor * 0.05
		gaji_bersih = gaji_kotor - pajak
	}

	modelPenggajian.NamaPegawai = dataPenggajian.NamaPegawai
	modelPenggajian.GajiPokok = dataPenggajian.GajiPokok
	modelPenggajian.JamLembur = dataPenggajian.JamLembur
	modelPenggajian.GajiKotor = gaji_kotor
	modelPenggajian.Pajak = pajak
	modelPenggajian.GajiBersih = gaji_bersih

	hasil := db.Save(&modelPenggajian)

	// Khusus Id tidak ada binding: required
	// karena hanya dibutuhkan saat mode
	// ubah dan hapus saja

	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil ubah data",
			"kesalahan": nil,
			"data":      modelPenggajian,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelPenggajian,
		})
	}
}

func HapusPenggajian(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data Penggajian dengan struktur Penggajian
	//dan menangkap data dari request
	var dataPenggajian StrukturPenggajianHapus
	if err := c.ShouldBindJSON(&dataPenggajian); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat variabel model Penggajian
	var modelPenggajian models.Penggajian
	//menghapus data Penggajian berdasarkan Id yang dikirim
	hasil := db.Delete(&modelPenggajian, dataPenggajian.Id)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data",
			"kesalahan": nil,
			"data":      dataPenggajian,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus Data",
			"kesalahan": kesalahan.Error(),
			"data":      dataPenggajian,
		})
	}
}
