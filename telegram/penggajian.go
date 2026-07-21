package telegram

import (
	"fmt"
	"main/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handlePenggajianCallback(chatID int64, messageID int, telegramID int64, action string) {
	switch action {
	case "lihat":
		penggajianLihat(chatID, messageID)
	case "tambah":
		penggajianMulaiTambah(chatID, messageID, telegramID)
	case "ubah":
		penggajianMulaiUbah(chatID, messageID, telegramID)
	case "hapus":
		penggajianMulaiHapus(chatID, messageID, telegramID)
	}
}

func penggajianLihat(chatID int64, messageID int) {
	var data []models.Penggajian
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "📋 <b>Data Penggajian</b>\n\n<i>Tidak ada data.</i>", BackButton("nav:penggajian"))
		return
	}

	text := "📋 <b>Data Penggajian</b>\n\n"
	for _, d := range data {
		text += fmt.Sprintf("<blockquote>🆔 <b>ID:</b> %d\n👤 <b>Nama:</b> %s\n💰 <b>Gaji Pokok:</b> Rp %.0f\n⏰ <b>Jam Lembur:</b> %d\n📊 <b>Gaji Kotor:</b> Rp %.0f\n🧾 <b>Pajak:</b> Rp %.0f\n💵 <b>Gaji Bersih:</b> Rp %.0f</blockquote>\n",
			d.ID, esc(d.NamaPegawai), d.GajiPokok, d.JamLembur, d.GajiKotor, d.Pajak, d.GajiBersih)
	}
	sendOrEdit(chatID, messageID, text, BackButton("nav:penggajian"))
}

func penggajianMulaiTambah(chatID int64, messageID int, telegramID int64) {
	state := &UserState{
		Module:         "penggajian",
		Action:         "tambah",
		Step:           0,
		TempData:       make(map[string]string),
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, "➕ <b>Tambah Penggajian</b>\n\n<i>Masukkan nama pegawai:</i>")
}

func penggajianMulaiUbah(chatID int64, messageID int, telegramID int64) {
	var data []models.Penggajian
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk diubah.</i>", BackButton("nav:penggajian"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🆔 %d  👤 %s  💵 Rp %.0f", d.ID, d.NamaPegawai, d.GajiBersih)
		callback := fmt.Sprintf("pick:penggajian:ubah:%d", d.ID)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:penggajian")),
	)
	sendOrEdit(chatID, messageID, "✏️ <b>Pilih data yang akan diubah:</b>", keyboard)
}

func penggajianMulaiHapus(chatID int64, messageID int, telegramID int64) {
	var data []models.Penggajian
	result := dbConn.Find(&data)
	if result.Error != nil || len(data) == 0 {
		sendOrEdit(chatID, messageID, "<i>Tidak ada data untuk dihapus.</i>", BackButton("nav:penggajian"))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, d := range data {
		label := fmt.Sprintf("🆔 %d  👤 %s  💵 Rp %.0f", d.ID, d.NamaPegawai, d.GajiBersih)
		callback := fmt.Sprintf("pick:penggajian:hapus:%d", d.ID)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(label, callback)),
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Batal", "nav:penggajian")),
	)
	sendOrEdit(chatID, messageID, "🗑️ <b>Pilih data yang akan dihapus:</b>", keyboard)
}

func penggajianPickItem(chatID int64, messageID int, telegramID int64, action string, id uint) {
	var data models.Penggajian
	result := dbConn.First(&data, id)
	if result.Error != nil {
		sendOrEdit(chatID, messageID, "<i>Data tidak ditemukan.</i>", BackButton("nav:penggajian"))
		return
	}

	if action == "hapus" {
		text := fmt.Sprintf("🗑️ <b>Yakin hapus data?</b>\n\n<blockquote>🆔 <b>ID:</b> %d | 👤 <b>Nama:</b> %s\n💵 <b>Gaji Bersih:</b> Rp %.0f</blockquote>", data.ID, esc(data.NamaPegawai), data.GajiBersih)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅  Ya, Hapus", fmt.Sprintf("confirm:penggajian:hapus:%d", id)),
				tgbotapi.NewInlineKeyboardButtonData("❌  Batal", "nav:penggajian"),
			),
		)
		sendOrEdit(chatID, messageID, text, keyboard)
		return
	}

	state := &UserState{
		Module:         "penggajian",
		Action:         "ubah",
		Step:           0,
		EditID:         id,
		TempData:       make(map[string]string),
		PanelMessageID: messageID,
	}
	userStates[telegramID] = state
	panelPrompt(state, chatID, fmt.Sprintf("✏️ <b>Ubah Penggajian</b> (ID: %d)\n\n<i>Masukkan nama pegawai baru (ketik </i><code>-</code><i> untuk skip):</i>", id))
}

func hitungPenggajian(gajiPokok float64, jamLembur int) (gajiKotor, pajak, gajiBersih float64) {
	uangLembur := float64(jamLembur) * 50000
	gajiKotor = gajiPokok + uangLembur
	pajak = 0.0
	gajiBersih = gajiKotor - pajak
	if gajiKotor > 5000000 {
		pajak = gajiKotor * 0.05
		gajiBersih = gajiKotor - pajak
	}
	return
}

func penggajianHandleInput(chatID int64, telegramID int64, text string, state *UserState) {
	if state.Action == "tambah" {
		switch state.Step {
		case 0:
			state.TempData["nama"] = text
			state.Step = 1
			panelPrompt(state, chatID, "➕ <b>Tambah Penggajian</b>\n\n<i>Masukkan gaji pokok:</i>")
		case 1:
			_, err := strconv.ParseFloat(text, 64)
			if err != nil {
				panelPrompt(state, chatID, "❌ <i>Gaji pokok harus berupa angka. Coba lagi:</i>")
				return
			}
			state.TempData["gaji_pokok"] = text
			state.Step = 2
			panelPrompt(state, chatID, "➕ <b>Tambah Penggajian</b>\n\n<i>Masukkan jam lembur:</i>")
		case 2:
			_, err := strconv.Atoi(text)
			if err != nil {
				panelPrompt(state, chatID, "❌ <i>Jam lembur harus berupa angka bulat. Coba lagi:</i>")
				return
			}
			state.TempData["jam_lembur"] = text

			gajiPokok, _ := strconv.ParseFloat(state.TempData["gaji_pokok"], 64)
			jamLembur, _ := strconv.Atoi(state.TempData["jam_lembur"])
			gajiKotor, pajak, gajiBersih := hitungPenggajian(gajiPokok, jamLembur)

			text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Simpan</b>\n\n<blockquote>👤 <b>Nama:</b> %s\n💰 <b>Gaji Pokok:</b> %s\n⏰ <b>Jam Lembur:</b> %s\n📊 <b>Gaji Kotor:</b> %.0f\n🧾 <b>Pajak:</b> %.0f\n💵 <b>Gaji Bersih:</b> %.0f</blockquote>",
				esc(state.TempData["nama"]), esc(state.TempData["gaji_pokok"]), esc(state.TempData["jam_lembur"]), gajiKotor, pajak, gajiBersih)
			keyboard := ConfirmKeyboard("penggajian", "tambah", 0)
			panelConfirm(state, chatID, text_confirm, keyboard)
		}
	} else if state.Action == "ubah" {
		switch state.Step {
		case 0:
			if text != "-" {
				state.TempData["nama"] = text
			}
			state.Step = 1
			panelPrompt(state, chatID, "✏️ <b>Ubah Penggajian</b>\n\n<i>Masukkan gaji pokok baru (ketik </i><code>-</code><i> untuk skip):</i>")
		case 1:
			if text != "-" {
				_, err := strconv.ParseFloat(text, 64)
				if err != nil {
					panelPrompt(state, chatID, "❌ <i>Gaji pokok harus berupa angka. Coba lagi:</i>")
					return
				}
				state.TempData["gaji_pokok"] = text
			}
			state.Step = 2
			panelPrompt(state, chatID, "✏️ <b>Ubah Penggajian</b>\n\n<i>Masukkan jam lembur baru (ketik </i><code>-</code><i> untuk skip):</i>")
		case 2:
			if text != "-" {
				_, err := strconv.Atoi(text)
				if err != nil {
					panelPrompt(state, chatID, "❌ <i>Jam lembur harus berupa angka bulat. Coba lagi:</i>")
					return
				}
				state.TempData["jam_lembur"] = text
			}
			text_confirm := fmt.Sprintf("📝 <b>Konfirmasi Ubah</b>\n\n<blockquote>👤 <b>Nama:</b> %s\n💰 <b>Gaji Pokok:</b> %s\n⏰ <b>Jam Lembur:</b> %s</blockquote>",
				displayOrDefault(esc(state.TempData["nama"]), "-"), displayOrDefault(esc(state.TempData["gaji_pokok"]), "-"), displayOrDefault(esc(state.TempData["jam_lembur"]), "-"))
			keyboard := ConfirmKeyboard("penggajian", "ubah", state.EditID)
			panelConfirm(state, chatID, text_confirm, keyboard)
		}
	}
}

func penggajianConfirm(chatID int64, messageID int, telegramID int64, action string, id uint) {
	state, exists := userStates[telegramID]
	if action != "hapus" && !exists {
		sendOrEdit(chatID, messageID, "❌ <i>Sesi expired. Silakan ulangi.</i>", BackButton("nav:penggajian"))
		return
	}

	switch action {
	case "tambah":
		gajiPokok, _ := strconv.ParseFloat(state.TempData["gaji_pokok"], 64)
		jamLembur, _ := strconv.Atoi(state.TempData["jam_lembur"])
		gajiKotor, pajak, gajiBersih := hitungPenggajian(gajiPokok, jamLembur)

		newData := models.Penggajian{
			NamaPegawai: state.TempData["nama"],
			GajiPokok:   gajiPokok,
			JamLembur:   jamLembur,
			GajiKotor:   gajiKotor,
			Pajak:       pajak,
			GajiBersih:  gajiBersih,
		}
		result := dbConn.Create(&newData)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menyimpan data.</i>", BackButton("nav:penggajian"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data penggajian berhasil ditambahkan!</b>", BackButton("nav:penggajian"))
		}
	case "ubah":
		var data models.Penggajian
		dbConn.First(&data, state.EditID)
		if v, ok := state.TempData["nama"]; ok {
			data.NamaPegawai = v
		}
		if v, ok := state.TempData["gaji_pokok"]; ok {
			data.GajiPokok, _ = strconv.ParseFloat(v, 64)
		}
		if v, ok := state.TempData["jam_lembur"]; ok {
			data.JamLembur, _ = strconv.Atoi(v)
		}
		gajiKotor, pajak, gajiBersih := hitungPenggajian(data.GajiPokok, data.JamLembur)
		data.GajiKotor = gajiKotor
		data.Pajak = pajak
		data.GajiBersih = gajiBersih

		result := dbConn.Save(&data)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal mengubah data.</i>", BackButton("nav:penggajian"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data penggajian berhasil diubah!</b>", BackButton("nav:penggajian"))
		}
	case "hapus":
		result := dbConn.Delete(&models.Penggajian{}, id)
		if result.Error != nil {
			sendOrEdit(chatID, messageID, "❌ <i>Gagal menghapus data.</i>", BackButton("nav:penggajian"))
		} else {
			sendOrEdit(chatID, messageID, "✅ <b>Data penggajian berhasil dihapus!</b>", BackButton("nav:penggajian"))
		}
	}
	delete(userStates, telegramID)
}
