package controllers

import (
	"crypto/sha1"
	"fmt"
	"main/models"
	"net/http"

	jwtV3 "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Binding dari POST JSON
type StrukturUserTambah struct {
	Nama     string `binding:"required"`
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type StrukturUserUbah struct {
	Id       uint
	Nama     string `binding:"required"`
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type StrukturUserHapus struct {
	Id uint `binding:"required"`
}

type StrukturLogin struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

// Tambahkan ini
func UserTampil(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//buat variabel array dari model suhu
	var modelUser []models.User
	hasil := db.Find(&modelUser)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil Tampil data",
			"kesalahan": nil,
			"data":      modelUser,
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
func UserTambah(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data User dengan struktur User
	// dan menangkap data dari request
	var dataUser StrukturUserTambah
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	// enkripsi password dengan sha1
	var sha = sha1.New()
	sha.Write([]byte(dataUser.Password))
	var encrypted = fmt.Sprintf("%x", sha.Sum(nil))
	var encrytedString = fmt.Sprintf("%s", encrypted)

	//membuat data baru dengan model User
	modelUser := models.User{
		Nama:     dataUser.Nama,
		Username: dataUser.Username,
		Password: encrytedString,
	}
	hasil := db.Create(&modelUser)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data",
			"kesalahan": nil,
			"data":      modelUser,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tambah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelUser,
		})
	}
}

func UserUbah(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data User dengan struktur User
	//dan menangkap data dari request
	var dataUser StrukturUserUbah
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat variabel model User
	var modelUser models.User

	// mencari data user dan merubah datanya
	cekUser := db.First(&modelUser, dataUser.Id)
	if cekUser.Error != nil {
		// encripsi password dengan sha1
		var sha = sha1.New()
		sha.Write([]byte(dataUser.Password))
		var encrypted = fmt.Sprintf("%x", sha.Sum(nil))
		var encrytedString = fmt.Sprintf("%s", encrypted)

		modelUser.Nama = dataUser.Nama
		modelUser.Username = dataUser.Username
		modelUser.Password = encrytedString

		db.First(&modelUser, dataUser.Id)
		modelUser.Nama = dataUser.Nama
		modelUser.Username = dataUser.Username
		modelUser.Password = dataUser.Password
		hasil := db.Save(&modelUser)

		kesalahan := hasil.Error
		if hasil.Error == nil {
			c.JSON(http.StatusOK, gin.H{
				"status":    true,
				"pesan":     "Berhasil ubah data",
				"kesalahan": nil,
				"data":      modelUser,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":    false,
				"pesan":     "Gagal ubah Data",
				"kesalahan": kesalahan.Error(),
				"data":      modelUser,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah data",
			"kesalahan": cekUser.Error.Error(),
			"data":      modelUser,
		})
	}

}

func UserHapus(c *gin.Context) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	//membuat variabel data User dengan struktur User
	//dan menangkap data dari request
	var dataUser StrukturUserHapus
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}
	//membuat variabel model User
	var modelUser models.User
	//menghapus data User berdasarkan Id yang dikirim
	hasil := db.Delete(&modelUser, dataUser.Id)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data",
			"kesalahan": nil,
			"data":      dataUser,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus Data",
			"kesalahan": kesalahan.Error(),
			"data":      dataUser,
		})
	}
}

func UserLogin(c *gin.Context) (any, error) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	// membuat variabel data User dengan struktur user dan menangkap data dari request
	var dataUser StrukturLogin
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		//kembalikan data kosong dan eror input login
		return nil, jwtV3.ErrMissingLoginValues
	}
	//enkripsi password dengan sha1
	var sha = sha1.New()
	sha.Write([]byte(dataUser.Password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)
	//membuat variabel model user
	var modelUser models.User
	//mencari data user berdasarkan username dan password
	cekUser := db.Where("username = ?", dataUser.Username).Where("password = ?", encryptedString).First(&modelUser)
	if cekUser.Error == nil {
		//kembalikan data user dan eror=nil
		return modelUser, nil
	} else {
		//kembalikan data kosong dan eror gagal login
		return nil, jwtV3.ErrFailedAuthentication
	}
}
