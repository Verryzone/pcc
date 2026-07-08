package fungsi

import (
	"github.com/gin-gonic/gin"
)

type Mahasiswa struct {
	Nim          string ` json:"nim" `
	Nama         string ` json:"nama" `
	TanggalLahir string ` json:"tanggal_lahir" `
}
type ProgramStudi struct {
	KodeProdi string      ` json:"kode_prodi" `
	NamaProdi string      ` json:"nama_prodi" `
	Mahasiswa []Mahasiswa ` json:"mahasiswa"`
}
type ResponProgramStudi struct {
	Status bool         ` json:"status" `
	Pesan  string       ` json:"pesan" `
	Data   ProgramStudi ` json:"data"`
}

func BacaDataProdi(c *gin.Context) {
	var vprodi ProgramStudi
	c.ShouldBind(&vprodi)
	var responnya ResponProgramStudi
	responnya.Status = true
	responnya.Pesan = "Berhasil tampil."
	responnya.Data = vprodi
	c.JSON(200, responnya)
}
