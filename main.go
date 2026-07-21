package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"main/ai"
	"main/fungsi"
	"main/telegram"

	// "main/wa"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"

	"main/models"

	"time"

	jwtV3 "github.com/appleboy/gin-jwt/v3"

	"github.com/golang-jwt/jwt/v5"

	"main/controllers"
)

func main() {
	// read .env file
	godotenv.Load()

	// panggil koneksi
	db := koneksi()

	db.AutoMigrate(&models.Suhu{})
	db.AutoMigrate(&models.Informasi{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Penggajian{})
	db.AutoMigrate(&models.Dokumen{})
	db.AutoMigrate(&models.Pesan{})

	r := gin.Default()

	// JWT Middleware
	key_jwt := os.Getenv("KEY_JWT")
	authMiddleware, err := jwtV3.New(&jwtV3.GinJWTMiddleware{
		Realm:       "fikom UDB",
		Key:         []byte(key_jwt),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: "id",

		PayloadFunc: func(data any) jwt.MapClaims {
			value, ok := data.(models.User)
			if ok {
				return jwt.MapClaims{
					"id":   value.ID,
					"nama": value.Nama,
				}
			}
			return jwt.MapClaims{}
		},

		Authenticator: controllers.UserLogin,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	// jika route tidak ada
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status": false,
			"pesan":  "Route tidak ditemukan!",
		})
	})

	// membuat variable db untuk membawa koneksi
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// ROUTE
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": true,
			"pesan":  "Berhasil tampil",
		})
	})

	// route tanpa middleware (PUBLIC)
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/programstudi", fungsi.BacaDataProdi)

	// route group dengan middleware (PROTECTED - Memerlukan JWT)
	auth := r.Group("/backend", authMiddleware.MiddlewareFunc())

	auth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": true,
			"pesan":  "Berhasil Tampil",
		})
	})

	// route user (PROTECTED)
	auth.GET("/user", controllers.UserTampil)
	auth.POST("/user", controllers.UserTambah)
	auth.PUT("/user", controllers.UserUbah)
	auth.DELETE("/user", controllers.UserHapus)

	// route pesan (PROTECTED)
	auth.GET("/pesan", controllers.PesanTampil)
	auth.POST("/pesan", controllers.PesanTambah)
	auth.PUT("/pesan", controllers.PesanUbah)
	auth.DELETE("/pesan", controllers.PesanHapus)

	// route suhu (PROTECTED)
	auth.GET("/suhu", controllers.Tampil)
	auth.POST("/suhu", controllers.Tambah)
	auth.PUT("/suhu", controllers.Ubah)
	auth.DELETE("/suhu", controllers.Hapus)

	// route informasi (PROTECTED)
	auth.GET("/informasi", controllers.InformasiTampil)
	auth.POST("/informasi", controllers.InformasiTambah)
	auth.PUT("/informasi", controllers.InformasiUbah)
	auth.DELETE("/informasi", controllers.InformasiHapus)

	// route penggajian (PROTECTED)
	auth.GET("/penggajian", controllers.TampilPenggajian)
	auth.POST("/penggajian", controllers.TambahPenggajian)
	auth.PUT("/penggajian", controllers.UbahPenggajian)
	auth.DELETE("/penggajian", controllers.HapusPenggajian)

	//METHOD DRIVE
	auth.POST("/drive", controllers.DriveUpload)
	auth.GET("/drive", controllers.DriveTampil)
	auth.GET("/drive/:id", controllers.DriveUnduh)

	// membaca nilai port dari .env
	port := os.Getenv("PORT")
	go r.Run(":" + port)

	go telegram.StartBot(db)
	ai.InitAi()
	// wa.KonekWa(db)

	// blokir main goroutine agar server & bot tetap berjalan sampai diberhentikan
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
