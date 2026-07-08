package controllers

import (
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Binding dari POST JSON
type StrukturPesan struct {
	Kode    string
	Balasan string
}

// Tambahkan ini
func PesanTampil(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//buat variabel array dari model suhu
	var modelPesan []models.Pesan
	hasil := db.Find(&modelPesan)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil Tampil data",
			"kesalahan": nil,
			"data":      modelPesan,
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
func PesanTambah(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data Suhu dengan struktur suhu
	// dan menangkap data dari request
	var dataPesan StrukturPesan
	if err := c.ShouldBindJSON(&dataPesan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat data baru dengan model suhu
	modelPesan := models.Pesan{
		Kode:      dataPesan.Kode,
		Balasan:   dataPesan.Balasan,
		CreatedAt: time.Now(),
	}
	hasil := db.Create(&modelPesan)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data",
			"kesalahan": nil,
			"data":      modelPesan,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tambah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelPesan,
		})
	}
}

func PesanUbah(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data Suhu dengan struktur suhu
	//dan menangkap data dari request
	var dataPesan StrukturPesan
	if err := c.ShouldBindJSON(&dataPesan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat variabel model suhu
	var modelPesan models.Pesan
	//mencari data suhu dan merubah datanya
	db.First(&modelPesan, dataPesan.Kode)
	modelPesan.Balasan = dataPesan.Balasan
	hasil := db.Save(&modelPesan)

	// Khusus Id tidak ada binding: required
	// karena hanya dibutuhkan saat mode
	// ubah dan hapus saja

	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil ubah data",
			"kesalahan": nil,
			"data":      modelPesan,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelPesan,
		})
	}
}

func PesanHapus(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data Suhu dengan struktur suhu
	//dan menangkap data dari request
	var dataPesan StrukturPesan
	if err := c.ShouldBindJSON(&dataPesan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat variabel model suhu
	var modelPesan models.Pesan
	//menghapus data suhu berdasarkan Id yang dikirim
	hasil := db.Delete(&modelPesan, dataPesan.Kode)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data",
			"kesalahan": nil,
			"data":      dataPesan,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus Data",
			"kesalahan": kesalahan.Error(),
			"data":      dataPesan,
		})
	}
}
