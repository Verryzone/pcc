package telegram

import (
	"fmt"
	"main/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleSuhuCallback(chatID int64, messageID int, telegramID int64, action string) {
	switch action {
	case "lihat":
		suhuLihat(chatID, messageID)
	case "tambah":
		suhuMulaiTambah(chatID, messageID, telegramID)
	case "ubah":
		suhuMulaiUbah(chatID, messageID, telegramID)
	case "hapus":
		suhuMulaiHapus(chatID, messageID, telegramID)
	}
}

func suhuLihat(chatID int64, messageID int) {
	var data []models.Suhu
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "📋 <b>Data Suhu</b>\n\n<i>Tidak ada data.</i>", BackButton("nav:suhu"))
		return
	}

	text := "📋 <b>Data Suhu</b>\n\n"
	for _, d := range data {
		text += fmt.Sprintf("<blockquote>🆔 <b>ID:</b> %d\n📍 <b>Lokasi:</b> %s\n🌡️ <b>Suhu:</b> %.2f°C</blockquote>\n", d.Id, esc(d.Lokasi), d.Suhu)
	}
	sendOrEdit(chatID, messageID, text, BackButton("nav:suhu"))
}

func suhuMulaiTambah(chatID int64, messageID int, telegramID int64) {
	state := &UserState{
		Module:         "suhu",
		Action:         "tambah",
		Step:           0,
		TempData:       make(map[string]string),
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, "➕ <b>Tambah Suhu</b>\n\n<i>Masukkan lokasi:</i>")
}

func suhuMulaiUbah(chatID int64, messageID int, telegramID int64) {
	var data []models.Suhu
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk diubah.</i>", BackButton("nav:suhu"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🆔 %d  📍 %s  🌡️ %.2f°C", d.Id, d.Lokasi, d.Suhu)
		callback := fmt.Sprintf("pick:suhu:ubah:%d", d.Id)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:suhu")),
	)
	sendOrEdit(chatID, messageID, "✏️ Pilih data yang akan diubah:", keyboard)
}

func suhuMulaiHapus(chatID int64, messageID int, telegramID int64) {
	var data []models.Suhu
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk dihapus.</i>", BackButton("nav:suhu"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🆔 %d  📍 %s  🌡️ %.2f°C", d.Id, d.Lokasi, d.Suhu)
		callback := fmt.Sprintf("pick:suhu:hapus:%d", d.Id)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:suhu")),
	)
	sendOrEdit(chatID, messageID, "🗑️ Pilih data yang akan dihapus:", keyboard)
}

func suhuPickItem(chatID int64, messageID int, telegramID int64, action string, id uint) {
	var data models.Suhu
	result := dbConn.First(&data, id)
	if result.Error != nil {
		sendOrEdit(chatID, messageID, "<i>Data tidak ditemukan.</i>", BackButton("nav:suhu"))
		return
	}

	if action == "hapus" {
		text := fmt.Sprintf("🗑️ <b>Yakin hapus data?</b>\n\n<blockquote>🆔 <b>ID:</b> %d\n📍 <b>Lokasi:</b> %s\n🌡️ <b>Suhu:</b> %.2f°C</blockquote>", data.Id, esc(data.Lokasi), data.Suhu)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅  Ya, Hapus", fmt.Sprintf("confirm:suhu:hapus:%d", id)),
				tgbotapi.NewInlineKeyboardButtonData("❌  Batal", "nav:suhu"),
			),
		)
		sendOrEdit(chatID, messageID, text, keyboard)
		return
	}

	state := &UserState{
		Module:         "suhu",
		Action:         "ubah",
		Step:           0,
		EditID:         id,
		TempData:       make(map[string]string),
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, fmt.Sprintf("✏️ <b>Ubah Suhu</b> (ID: %d)\n\n<i>Masukkan lokasi baru (ketik </i><code>-</code><i> untuk skip):</i>", id))
}

func suhuHandleInput(chatID int64, telegramID int64, text string, state *UserState) {
	if state.Action == "tambah" {
		switch state.Step {
		case 0:
			state.TempData["lokasi"] = text
			state.Step = 1
			panelPrompt(state, chatID, "➕ <b>Tambah Suhu</b>\n\n<i>Masukkan suhu (°C):</i>")
		case 1:
			_, err := strconv.ParseFloat(text, 32)
			if err != nil {
				panelPrompt(state, chatID, "❌ <i>Suhu harus berupa angka. Coba lagi:</i>")
				return
			}
			state.TempData["suhu"] = text
			text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Simpan</b>\n\n<blockquote>📍 <b>Lokasi:</b> %s\n🌡️ <b>Suhu:</b> %s°C</blockquote>", esc(state.TempData["lokasi"]), esc(state.TempData["suhu"]))
			keyboard := ConfirmKeyboard("suhu", "tambah", 0)
			panelConfirm(state, chatID, text_confirm, keyboard)
		}
	} else if state.Action == "ubah" {
		switch state.Step {
		case 0:
			if text != "-" {
				state.TempData["lokasi"] = text
			}
			state.Step = 1
			panelPrompt(state, chatID, "✏️ <b>Ubah Suhu</b>\n\n<i>Masukkan suhu baru (ketik </i><code>-</code><i> untuk skip):</i>")
		case 1:
			if text != "-" {
				_, err := strconv.ParseFloat(text, 32)
				if err != nil {
					panelPrompt(state, chatID, "❌ <i>Suhu harus berupa angka. Coba lagi:</i>")
					return
				}
				state.TempData["suhu"] = text
			}
			text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Ubah</b>\n\n<blockquote>📍 <b>Lokasi:</b> %s\n🌡️ <b>Suhu:</b> %s</blockquote>", displayOrDefault(esc(state.TempData["lokasi"]), "-"), displayOrDefault(esc(state.TempData["suhu"]), "-"))
			keyboard := ConfirmKeyboard("suhu", "ubah", state.EditID)
			panelConfirm(state, chatID, text_confirm, keyboard)
		}
	}
}

func suhuConfirm(chatID int64, messageID int, telegramID int64, action string, id uint) {
	state, exists := userStates[telegramID]
	if action != "hapus" && !exists {
		sendOrEdit(chatID, messageID, "❌ <i>Sesi expired. Silakan ulangi.</i>", BackButton("nav:suhu"))
		return
	}

	switch action {
	case "tambah":
		suhuVal, _ := strconv.ParseFloat(state.TempData["suhu"], 32)
		newData := models.Suhu{
			Lokasi: state.TempData["lokasi"],
			Suhu:   float32(suhuVal),
		}
		result := dbConn.Create(&newData)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menyimpan data.</i>", BackButton("nav:suhu"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data suhu berhasil ditambahkan!</b>", BackButton("nav:suhu"))
		}
	case "ubah":
		var data models.Suhu
		dbConn.First(&data, state.EditID)
		if v, ok := state.TempData["lokasi"]; ok {
			data.Lokasi = v
		}
		if v, ok := state.TempData["suhu"]; ok {
			suhuVal, _ := strconv.ParseFloat(v, 32)
			data.Suhu = float32(suhuVal)
		}
		result := dbConn.Save(&data)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal mengubah data.</i>", BackButton("nav:suhu"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data suhu berhasil diubah!</b>", BackButton("nav:suhu"))
		}
	case "hapus":
		result := dbConn.Delete(&models.Suhu{}, id)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menghapus data.</i>", BackButton("nav:suhu"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data suhu berhasil dihapus!</b>", BackButton("nav:suhu"))
		}
	}
	delete(userStates, telegramID)
}

func displayOrDefault(val, defaultVal string) string {
	if val == "" {
		return defaultVal
	}
	return val
}
