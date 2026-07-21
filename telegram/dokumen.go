package telegram

import (
	"fmt"
	"main/models"
)

func handleDokumenCallback(chatID int64, messageID int, action string) {
	switch action {
	case "lihat":
		dokumenLihat(chatID, messageID)
	}
}

func dokumenLihat(chatID int64, messageID int) {
	var data []models.Dokumen
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "📋 <b>Data Dokumen</b>\n\n<i>Tidak ada data.</i>", BackButton("nav:dokumen"))
		return
	}

	text := "📋 <b>Data Dokumen</b>\n\n"
	for _, d := range data {
		text += fmt.Sprintf("<blockquote>🆔 <b>ID:</b> %d\n📄 <b>Nama:</b> %s\n🔗 <b>URL:</b> %s</blockquote>\n", d.ID, esc(d.NamaDokumen), esc(d.FileUrl))
	}
	sendOrEdit(chatID, messageID, text, BackButton("nav:dokumen"))
}
