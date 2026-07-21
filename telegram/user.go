package telegram

import (
	"crypto/sha1"
	"fmt"
	"main/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleUserCallback(chatID int64, messageID int, telegramID int64, action string) {
	switch action {
	case "lihat":
		userLihat(chatID, messageID)
	case "tambah":
		userMulaiTambah(chatID, messageID, telegramID)
	case "ubah":
		userMulaiUbah(chatID, messageID, telegramID)
	case "hapus":
		userMulaiHapus(chatID, messageID, telegramID)
	}
}

func userLihat(chatID int64, messageID int) {
	var data []models.User
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "📋 <b>Data User</b>\n\n<i>Tidak ada data.</i>", BackButton("nav:user"))
		return
	}

	text := "📋 <b>Data User</b>\n\n"
	for _, d := range data {
		text += fmt.Sprintf("<blockquote>🆔 <b>ID:</b> %d | 👤 <b>Nama:</b> %s\n🔑 <b>Username:</b> %s</blockquote>\n", d.ID, esc(d.Nama), esc(d.Username))
	}
	sendOrEdit(chatID, messageID, text, BackButton("nav:user"))
}

func userMulaiTambah(chatID int64, messageID int, telegramID int64) {
	state := &UserState{
		Module:         "user",
		Action:         "tambah",
		Step:           0,
		TempData:       make(map[string]string),
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, "➕ <b>Tambah User</b>\n\n<i>Masukkan nama:</i>")
}

func userMulaiUbah(chatID int64, messageID int, telegramID int64) {
	var data []models.User
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk diubah.</i>", BackButton("nav:user"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🆔 %d  👤 %s  🔑 %s", d.ID, d.Nama, d.Username)
		callback := fmt.Sprintf("pick:user:ubah:%d", d.ID)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:user")),
	)
	sendOrEdit(chatID, messageID, "✏️ <b>Pilih data yang akan diubah:</b>", keyboard)
}

func userMulaiHapus(chatID int64, messageID int, telegramID int64) {
	var data []models.User
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk dihapus.</i>", BackButton("nav:user"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🆔 %d  👤 %s  🔑 %s", d.ID, d.Nama, d.Username)
		callback := fmt.Sprintf("pick:user:hapus:%d", d.ID)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:user")),
	)
	sendOrEdit(chatID, messageID, "🗑️ <b>Pilih data yang akan dihapus:</b>", keyboard)
}

func userPickItem(chatID int64, messageID int, telegramID int64, action string, id uint) {
	var data models.User
	result := dbConn.First(&data, id)
	if result.Error != nil {
		sendOrEdit(chatID, messageID, "<i>Data tidak ditemukan.</i>", BackButton("nav:user"))
		return
	}

	if action == "hapus" {
		text := fmt.Sprintf("🗑️ <b>Yakin hapus data?</b>\n\n<blockquote>🆔 <b>ID:</b> %d | 👤 <b>Nama:</b> %s\n🔑 <b>Username:</b> %s</blockquote>", data.ID, esc(data.Nama), esc(data.Username))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅  Ya, Hapus", fmt.Sprintf("confirm:user:hapus:%d", id)),
				tgbotapi.NewInlineKeyboardButtonData("❌  Batal", "nav:user"),
			),
		)
		sendOrEdit(chatID, messageID, text, keyboard)
		return
	}

	state := &UserState{
		Module:         "user",
		Action:         "ubah",
		Step:           0,
		EditID:         id,
		TempData:       make(map[string]string),
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, fmt.Sprintf("✏️ <b>Ubah User</b> (ID: %d)\n\n<i>Masukkan nama baru (ketik </i><code>-</code><i> untuk skip):</i>", id))
}

func hashSHA1(text string) string {
	sha := sha1.New()
	sha.Write([]byte(text))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func userHandleInput(chatID int64, telegramID int64, text string, state *UserState) {
	if state.Action == "tambah" {
		switch state.Step {
		case 0:
			state.TempData["nama"] = text
			state.Step = 1
			panelPrompt(state, chatID, "➕ <b>Tambah User</b>\n\n<i>Masukkan username:</i>")
		case 1:
			state.TempData["username"] = text
			state.Step = 2
			panelPrompt(state, chatID, "➕ <b>Tambah User</b>\n\n<i>Masukkan password:</i>")
		case 2:
			state.TempData["password"] = text
			text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Simpan</b>\n\n<blockquote>👤 <b>Nama:</b> %s\n🔑 <b>Username:</b> %s\n🔒 <b>Password:</b> ********</blockquote>", esc(state.TempData["nama"]), esc(state.TempData["username"]))
			keyboard := ConfirmKeyboard("user", "tambah", 0)
			panelConfirm(state, chatID, text_confirm, keyboard)
		}
	} else if state.Action == "ubah" {
		switch state.Step {
		case 0:
			if text != "-" {
				state.TempData["nama"] = text
			}
			state.Step = 1
			panelPrompt(state, chatID, "✏️ <b>Ubah User</b>\n\n<i>Masukkan username baru (ketik </i><code>-</code><i> untuk skip):</i>")
		case 1:
			if text != "-" {
				state.TempData["username"] = text
			}
			state.Step = 2
			panelPrompt(state, chatID, "✏️ <b>Ubah User</b>\n\n<i>Masukkan password baru (ketik </i><code>-</code><i> untuk skip):</i>")
		case 2:
			if text != "-" {
				state.TempData["password"] = text
			}
			text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Ubah</b>\n\n<blockquote>👤 <b>Nama:</b> %s\n🔑 <b>Username:</b> %s</blockquote>",
				displayOrDefault(esc(state.TempData["nama"]), "-"), displayOrDefault(esc(state.TempData["username"]), "-"))
			keyboard := ConfirmKeyboard("user", "ubah", state.EditID)
			panelConfirm(state, chatID, text_confirm, keyboard)
		}
	}
}

func userConfirm(chatID int64, messageID int, telegramID int64, action string, id uint) {
	state, exists := userStates[telegramID]
	if action != "hapus" && !exists {
		sendOrEdit(chatID, messageID, "❌ <i>Sesi expired. Silakan ulangi.</i>", BackButton("nav:user"))
		return
	}

	switch action {
	case "tambah":
		newData := models.User{
			Nama:     state.TempData["nama"],
			Username: state.TempData["username"],
			Password: hashSHA1(state.TempData["password"]),
		}
		result := dbConn.Create(&newData)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menyimpan data.</i>", BackButton("nav:user"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data user berhasil ditambahkan!</b>", BackButton("nav:user"))
		}
	case "ubah":
		var data models.User
		dbConn.First(&data, state.EditID)
		if v, ok := state.TempData["nama"]; ok {
			data.Nama = v
		}
		if v, ok := state.TempData["username"]; ok {
			data.Username = v
		}
		if v, ok := state.TempData["password"]; ok {
			data.Password = hashSHA1(v)
		}
		result := dbConn.Save(&data)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal mengubah data.</i>", BackButton("nav:user"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data user berhasil diubah!</b>", BackButton("nav:user"))
		}
	case "hapus":
		result := dbConn.Delete(&models.User{}, id)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menghapus data.</i>", BackButton("nav:user"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data user berhasil dihapus!</b>", BackButton("nav:user"))
		}
	}
	delete(userStates, telegramID)
}
