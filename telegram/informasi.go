package telegram

import (
	"fmt"
	"main/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleInformasiCallback(chatID int64, messageID int, telegramID int64, action string) {
	switch action {
	case "lihat":
		informasiLihat(chatID, messageID)
	case "tambah":
		informasiMulaiTambah(chatID, messageID, telegramID)
	case "ubah":
		informasiMulaiUbah(chatID, messageID, telegramID)
	case "hapus":
		informasiMulaiHapus(chatID, messageID, telegramID)
	}
}

func informasiLihat(chatID int64, messageID int) {
	var data []models.Informasi
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "📋 <b>Data Informasi</b>\n\n<i>Tidak ada data.</i>", BackButton("nav:informasi"))
		return
	}

	text := "📋 <b>Data Informasi</b>\n\n"
	for _, d := range data {
		text += fmt.Sprintf("<blockquote>🆔 <b>ID:</b> %d\n📝 <b>Judul:</b> %s\n📄 <b>Konten:</b> %s\n🔗 <b>URL:</b> %s</blockquote>\n", d.ID, esc(d.Judul), esc(d.Konten), esc(d.UrlDokumen))
	}
	sendOrEdit(chatID, messageID, text, BackButton("nav:informasi"))
}

func informasiMulaiTambah(chatID int64, messageID int, telegramID int64) {
	state := &UserState{
		Module:         "informasi",
		Action:         "tambah",
		Step:           0,
		TempData:       make(map[string]string),
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, "➕ <b>Tambah Informasi</b>\n\n<i>Masukkan judul:</i>")
}

func informasiMulaiUbah(chatID int64, messageID int, telegramID int64) {
	var data []models.Informasi
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk diubah.</i>", BackButton("nav:informasi"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🆔 %d  📝 %s", d.ID, d.Judul)
		callback := fmt.Sprintf("pick:informasi:ubah:%d", d.ID)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:informasi")),
	)
	sendOrEdit(chatID, messageID, "✏️ <b>Pilih data yang akan diubah:</b>", keyboard)
}

func informasiMulaiHapus(chatID int64, messageID int, telegramID int64) {
	var data []models.Informasi
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk dihapus.</i>", BackButton("nav:informasi"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🆔 %d  📝 %s", d.ID, d.Judul)
		callback := fmt.Sprintf("pick:informasi:hapus:%d", d.ID)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:informasi")),
	)
	sendOrEdit(chatID, messageID, "🗑️ <b>Pilih data yang akan dihapus:</b>", keyboard)
}

func informasiPickItem(chatID int64, messageID int, telegramID int64, action string, id uint) {
	var data models.Informasi
	result := dbConn.First(&data, id)
	if result.Error != nil {
		sendOrEdit(chatID, messageID, "<i>Data tidak ditemukan.</i>", BackButton("nav:informasi"))
		return
	}

	if action == "hapus" {
		text := fmt.Sprintf("🗑️ <b>Yakin hapus data?</b>\n\n<blockquote>🆔 <b>ID:</b> %d\n📝 <b>Judul:</b> %s</blockquote>", data.ID, esc(data.Judul))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅  Ya, Hapus", fmt.Sprintf("confirm:informasi:hapus:%d", id)),
				tgbotapi.NewInlineKeyboardButtonData("❌  Batal", "nav:informasi"),
			),
		)
		sendOrEdit(chatID, messageID, text, keyboard)
		return
	}

	state := &UserState{
		Module:         "informasi",
		Action:         "ubah",
		Step:           0,
		EditID:         id,
		TempData:       make(map[string]string),
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, fmt.Sprintf("✏️ <b>Ubah Informasi</b> (ID: %d)\n\n<i>Masukkan judul baru (ketik </i><code>-</code><i> untuk skip):</i>", id))
}

func informasiHandleInput(chatID int64, telegramID int64, text string, state *UserState) {
	if state.Action == "tambah" {
		switch state.Step {
		case 0:
			state.TempData["judul"] = text
			state.Step = 1
			panelPrompt(state, chatID, "➕ <b>Tambah Informasi</b>\n\n<i>Masukkan konten:</i>")
		case 1:
			state.TempData["konten"] = text
			state.Step = 2
			panelPrompt(state, chatID, "➕ <b>Tambah Informasi</b>\n\n<i>Masukkan URL dokumen:</i>")
		case 2:
			state.TempData["url_dokumen"] = text
			text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Simpan</b>\n\n<blockquote>📝 <b>Judul:</b> %s\n📄 <b>Konten:</b> %s\n🔗 <b>URL:</b> %s</blockquote>", esc(state.TempData["judul"]), esc(state.TempData["konten"]), esc(state.TempData["url_dokumen"]))
			keyboard := ConfirmKeyboard("informasi", "tambah", 0)
			panelConfirm(state, chatID, text_confirm, keyboard)
		}
	} else if state.Action == "ubah" {
		switch state.Step {
		case 0:
			if text != "-" {
				state.TempData["judul"] = text
			}
			state.Step = 1
			panelPrompt(state, chatID, "✏️ <b>Ubah Informasi</b>\n\n<i>Masukkan konten baru (ketik </i><code>-</code><i> untuk skip):</i>")
		case 1:
			if text != "-" {
				state.TempData["konten"] = text
			}
			state.Step = 2
			panelPrompt(state, chatID, "✏️ <b>Ubah Informasi</b>\n\n<i>Masukkan URL dokumen baru (ketik </i><code>-</code><i> untuk skip):</i>")
		case 2:
			if text != "-" {
				state.TempData["url_dokumen"] = text
			}
			text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Ubah</b>\n\n<blockquote>📝 <b>Judul:</b> %s\n📄 <b>Konten:</b> %s\n🔗 <b>URL:</b> %s</blockquote>", displayOrDefault(esc(state.TempData["judul"]), "-"), displayOrDefault(esc(state.TempData["konten"]), "-"), displayOrDefault(esc(state.TempData["url_dokumen"]), "-"))
			keyboard := ConfirmKeyboard("informasi", "ubah", state.EditID)
			panelConfirm(state, chatID, text_confirm, keyboard)
		}
	}
}

func informasiConfirm(chatID int64, messageID int, telegramID int64, action string, id uint) {
	state, exists := userStates[telegramID]
	if action != "hapus" && !exists {
		sendOrEdit(chatID, messageID, "❌ <i>Sesi expired. Silakan ulangi.</i>", BackButton("nav:informasi"))
		return
	}

	switch action {
	case "tambah":
		newData := models.Informasi{
			Judul:      state.TempData["judul"],
			Konten:     state.TempData["konten"],
			UrlDokumen: state.TempData["url_dokumen"],
		}
		result := dbConn.Create(&newData)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menyimpan data.</i>", BackButton("nav:informasi"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data informasi berhasil ditambahkan!</b>", BackButton("nav:informasi"))
		}
	case "ubah":
		var data models.Informasi
		dbConn.First(&data, state.EditID)
		if v, ok := state.TempData["judul"]; ok {
			data.Judul = v
		}
		if v, ok := state.TempData["konten"]; ok {
			data.Konten = v
		}
		if v, ok := state.TempData["url_dokumen"]; ok {
			data.UrlDokumen = v
		}
		result := dbConn.Save(&data)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal mengubah data.</i>", BackButton("nav:informasi"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data informasi berhasil diubah!</b>", BackButton("nav:informasi"))
		}
	case "hapus":
		result := dbConn.Delete(&models.Informasi{}, id)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menghapus data.</i>", BackButton("nav:informasi"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data informasi berhasil dihapus!</b>", BackButton("nav:informasi"))
		}
	}
	delete(userStates, telegramID)
}
