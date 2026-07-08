package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"main/models"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func DriveUpload(c *gin.Context) {
	// Ambil data dari form
	fileName := c.PostForm("fileName")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  "File tidak ditemukan",
		})
		return
	}

	mimeType := file.Header.Get("Content-Type")

	// Buka file
	fileOpen, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membaca file",
		})
		return
	}
	defer fileOpen.Close()

	// Baca isi file
	fileData, err := ioutil.ReadAll(fileOpen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membaca isi file",
		})
		return
	}

	// Encode file ke Base64
	data := base64.StdEncoding.EncodeToString(fileData)

	// Payload untuk Google Apps Script
	postBody, err := json.Marshal(map[string]string{
		"fileName": fileName,
		"mimeType": mimeType,
		"data":     data,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membuat payload",
		})
		return
	}

	requestBody := bytes.NewBuffer(postBody)

	// Kirim request ke Google Apps Script
	res, err := http.Post(
		"https://script.google.com/macros/s/AKfycbzLMAEU44QS0dYjUD149KEw5bmsIntwe6jGML6jq7eht38FsFpzaciPhqZxrzg1EVmB/exec",
		"application/json; charset=UTF-8",
		requestBody,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"kode_error": "ERR-DRIVE",
			"pesan":      "Gagal Upload",
		})
		return
	}
	defer res.Body.Close()

	// Baca response
	hasilBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membaca response",
		})
		return
	}

	// Debug: Log raw response
	println("Response Status:", res.StatusCode)
	println("Response Body:", string(hasilBody))

	// Cek status HTTP
	if res.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  "Google Apps Script error",
			"code":   res.StatusCode,
			"body":   string(hasilBody),
		})
		return
	}

	// Parse JSON response
	var hasilJSON map[string]interface{}
	if err := json.Unmarshal(hasilBody, &hasilJSON); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Response tidak valid",
			"debug":  string(hasilBody), // Tambahan untuk melihat response asli
		})
		return
	}

	//ambil koneksi ke variabel db
	db := c.MustGet("db").(*gorm.DB)
	//membuat objek dokumen baru
	dokumenBaru := models.Dokumen{
		NamaDokumen: hasilJSON["filename"].(string),
		FileId:      hasilJSON["fileId"].(string),
		FileUrl:     hasilJSON["fileUrl"].(string),
	}
	//membuat record baru di dokumen
	hasilDokumen := db.Create(&dokumenBaru)

	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"pesan":     "Berhasil Upload",
		"data":      hasilJSON,
		"tersimpan": hasilDokumen.RowsAffected,
	})
}

// Menampilkan seluruh data dokumen
func DriveTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dokumen []models.Dokumen

	if err := db.Find(&dokumen).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil Tampil",
		"data":   dokumen,
	})
}

// Download file dari Google Drive melalui Apps Script
func DriveUnduh(c *gin.Context) {
	id := c.Param("id")

	// Request ke Google Apps Script
	res, err := http.Get(
		"https://script.google.com/macros/s/AKfycbzLMAEU44QS0dYjUD149KEw5bmsIntwe6jGML6jq7eht38FsFpzaciPhqZxrzg1EVmB/exec?id=" + id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal Unduh",
		})
		return
	}
	defer res.Body.Close()

	// Baca response
	hasilBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membaca response",
		})
		return
	}

	// Parse JSON
	var hasilJSON map[string]interface{}
	if err := json.Unmarshal(hasilBody, &hasilJSON); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Response tidak valid",
		})
		return
	}

	// Ambil data file
	fileBase64, ok := hasilJSON["file"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Data file tidak ditemukan",
		})
		return
	}

	// Ambil mime type
	mimeType, ok := hasilJSON["mimeType"].(string)
	if !ok {
		mimeType = "application/octet-stream"
	}

	// Decode Base64 menjadi file
	fileData, err := base64.StdEncoding.DecodeString(fileBase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal decode file",
		})
		return
	}

	// Kirim file ke browser
	c.Header("Content-Type", mimeType)
	c.Writer.Write(fileData)
}
