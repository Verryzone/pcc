package telegram

import (
	"fmt"
	"main/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handlePesanCallback(chatID int64, messageID int, telegramID int64, action string) {
	switch action {
	case "lihat":
		pesanLihat(chatID, messageID)
	case "tambah":
		pesanMulaiTambah(chatID, messageID, telegramID)
	case "ubah":
		pesanMulaiUbah(chatID, messageID, telegramID)
	case "hapus":
		pesanMulaiHapus(chatID, messageID, telegramID)
	}
}

func pesanLihat(chatID int64, messageID int) {
	var data []models.Pesan
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "📋 <b>Data Pesan</b>\n\n<i>Tidak ada data.</i>", BackButton("nav:pesan"))
		return
	}

	text := "📋 <b>Data Pesan</b>\n\n"
	for _, d := range data {
		text += fmt.Sprintf("<blockquote>🏷️ <b>Kode:</b> %s\n💬 <b>Balasan:</b> %s</blockquote>\n", esc(d.Kode), esc(d.Balasan))
	}
	sendOrEdit(chatID, messageID, text, BackButton("nav:pesan"))
}

func pesanMulaiTambah(chatID int64, messageID int, telegramID int64) {
	state := &UserState{
		Module:         "pesan",
		Action:         "tambah",
		Step:           0,
		TempData:       make(map[string]string),
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, "➕ <b>Tambah Pesan</b>\n\n<i>Masukkan kode:</i>")
}

func pesanMulaiUbah(chatID int64, messageID int, telegramID int64) {
	var data []models.Pesan
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk diubah.</i>", BackButton("nav:pesan"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🏷️ %s  💬 %s", d.Kode, truncate(d.Balasan, 30))
		callback := fmt.Sprintf("pick:pesan:ubah:%s", d.Kode)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:pesan")),
	)
	sendOrEdit(chatID, messageID, "✏️ <b>Pilih data yang akan diubah:</b>", keyboard)
}

func pesanMulaiHapus(chatID int64, messageID int, telegramID int64) {
	var data []models.Pesan
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk dihapus.</i>", BackButton("nav:pesan"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🏷️ %s  💬 %s", d.Kode, truncate(d.Balasan, 30))
		callback := fmt.Sprintf("pick:pesan:hapus:%s", d.Kode)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:pesan")),
	)
	sendOrEdit(chatID, messageID, "🗑️ <b>Pilih data yang akan dihapus:</b>", keyboard)
}

func pesanPickItem(chatID int64, messageID int, telegramID int64, action string, kode string) {
	var data models.Pesan
	result := dbConn.Where("kode = ?", kode).First(&data)
	if result.Error != nil {
		sendOrEdit(chatID, messageID, "<i>Data tidak ditemukan.</i>", BackButton("nav:pesan"))
		return
	}

	if action == "hapus" {
		text := fmt.Sprintf("🗑️ <b>Yakin hapus data?</b>\n\n<blockquote>🏷️ <b>Kode:</b> %s\n💬 <b>Balasan:</b> %s</blockquote>", esc(data.Kode), esc(data.Balasan))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅  Ya, Hapus", fmt.Sprintf("confirm:pesan:hapus:%s", kode)),
				tgbotapi.NewInlineKeyboardButtonData("❌  Batal", "nav:pesan"),
			),
		)
		sendOrEdit(chatID, messageID, text, keyboard)
		return
	}

	state := &UserState{
		Module:         "pesan",
		Action:         "ubah",
		Step:           0,
		TempData:       map[string]string{"kode": kode},
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, fmt.Sprintf("✏️ <b>Ubah Pesan</b> (Kode: %s)\n\n<i>Masukkan balasan baru:</i>", esc(kode)))
}

func pesanHandleInput(chatID int64, telegramID int64, text string, state *UserState) {
	if state.Action == "tambah" {
		switch state.Step {
		case 0:
			state.TempData["kode"] = text
			state.Step = 1
			panelPrompt(state, chatID, "➕ <b>Tambah Pesan</b>\n\n<i>Masukkan balasan:</i>")
		case 1:
			state.TempData["balasan"] = text
			text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Simpan</b>\n\n<blockquote>🏷️ <b>Kode:</b> %s\n💬 <b>Balasan:</b> %s</blockquote>", esc(state.TempData["kode"]), esc(state.TempData["balasan"]))
			keyboard := ConfirmKeyboard("pesan", "tambah", 0)
			panelConfirm(state, chatID, text_confirm, keyboard)
		}
	} else if state.Action == "ubah" {
		state.TempData["balasan"] = text
		text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Ubah</b>\n\n<blockquote>🏷️ <b>Kode:</b> %s\n💬 <b>Balasan:</b> %s</blockquote>", esc(state.TempData["kode"]), esc(state.TempData["balasan"]))
		keyboard := ConfirmKeyboard("pesan", "ubah", 0)
		panelConfirm(state, chatID, text_confirm, keyboard)
	}
}

func pesanConfirm(chatID int64, messageID int, telegramID int64, action string, kode string) {
	state, exists := userStates[telegramID]
	if action != "hapus" && !exists {
		sendOrEdit(chatID, messageID, "❌ <i>Sesi expired. Silakan ulangi.</i>", BackButton("nav:pesan"))
		return
	}

	switch action {
	case "tambah":
		newData := models.Pesan{
			Kode:    state.TempData["kode"],
			Balasan: state.TempData["balasan"],
		}
		result := dbConn.Create(&newData)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menyimpan data.</i>", BackButton("nav:pesan"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data pesan berhasil ditambahkan!</b>", BackButton("nav:pesan"))
		}
	case "ubah":
		var data models.Pesan
		dbConn.Where("kode = ?", state.TempData["kode"]).First(&data)
		data.Balasan = state.TempData["balasan"]
		result := dbConn.Save(&data)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal mengubah data.</i>", BackButton("nav:pesan"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data pesan berhasil diubah!</b>", BackButton("nav:pesan"))
		}
	case "hapus":
		result := dbConn.Where("kode = ?", kode).Delete(&models.Pesan{})
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menghapus data.</i>", BackButton("nav:pesan"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data pesan berhasil dihapus!</b>", BackButton("nav:pesan"))
		}
	}
	delete(userStates, telegramID)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

