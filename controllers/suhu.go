package controllers

import (
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Binding dari POST JSON
type StrukturSuhu struct {
	Id     uint
	Lokasi string  `binding:"required"`
	Suhu   float32 `binding:"required"`
}

type StrukturSuhuHapus struct {
	Id uint `binding:"required"`
}

// Tambahkan ini
func Tampil(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//buat variabel array dari model suhu
	var modelSuhu []models.Suhu
	hasil := db.Find(&modelSuhu)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil Tampil data",
			"kesalahan": nil,
			"data":      modelSuhu,
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
func Tambah(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data Suhu dengan struktur suhu
	// dan menangkap data dari request
	var dataSuhu StrukturSuhu
	if err := c.ShouldBindJSON(&dataSuhu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat data baru dengan model suhu
	modelSuhu := models.Suhu{
		Lokasi:    dataSuhu.Lokasi,
		Suhu:      dataSuhu.Suhu,
		CreatedAt: time.Now(),
	}
	hasil := db.Create(&modelSuhu)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data",
			"kesalahan": nil,
			"data":      modelSuhu,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tambah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelSuhu,
		})
	}
}

func Ubah(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data Suhu dengan struktur suhu
	//dan menangkap data dari request
	var dataSuhu StrukturSuhu
	if err := c.ShouldBindJSON(&dataSuhu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat variabel model suhu
	var modelSuhu models.Suhu
	//mencari data suhu dan merubah datanya
	db.First(&modelSuhu, dataSuhu.Id)
	modelSuhu.Lokasi = dataSuhu.Lokasi
	modelSuhu.Suhu = dataSuhu.Suhu
	hasil := db.Save(&modelSuhu)

	// Khusus Id tidak ada binding: required
	// karena hanya dibutuhkan saat mode
	// ubah dan hapus saja

	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil ubah data",
			"kesalahan": nil,
			"data":      modelSuhu,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelSuhu,
		})
	}
}

func Hapus(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data Suhu dengan struktur suhu
	//dan menangkap data dari request
	var dataSuhu StrukturSuhuHapus
	if err := c.ShouldBindJSON(&dataSuhu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat variabel model suhu
	var modelSuhu models.Suhu
	//menghapus data suhu berdasarkan Id yang dikirim
	hasil := db.Delete(&modelSuhu, dataSuhu.Id)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data",
			"kesalahan": nil,
			"data":      dataSuhu,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus Data",
			"kesalahan": kesalahan.Error(),
			"data":      dataSuhu,
		})
	}
}
