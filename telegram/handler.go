package telegram

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func esc(s string) string {
	return html.EscapeString(s)
}

type UserState struct {
	Module          string
	Action          string
	Step            int
	TempData        map[string]string
	EditID          uint
	PanelMessageID  int
}

var userStates = make(map[int64]*UserState)

func MainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌡️  Suhu", "nav:suhu"),
			tgbotapi.NewInlineKeyboardButtonData("📰  Informasi", "nav:informasi"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰  Penggajian", "nav:penggajian"),
			tgbotapi.NewInlineKeyboardButtonData("💬  Pesan", "nav:pesan"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤  User", "nav:user"),
			tgbotapi.NewInlineKeyboardButtonData("📄  Dokumen", "nav:dokumen"),
		),
	)
}

func SubmenuKeyboard(module string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋  Lihat Semua", "act:"+module+":lihat"),
			tgbotapi.NewInlineKeyboardButtonData("➕  Tambah", "act:"+module+":tambah"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✏️  Ubah", "act:"+module+":ubah"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️  Hapus", "act:"+module+":hapus"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙  Kembali", "nav:main"),
		),
	)
}

func SubmenuKeyboardReadOnly(module string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋  Lihat Semua", "act:"+module+":lihat"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙  Kembali", "nav:main"),
		),
	)
}

func BackButton(data string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙  Kembali", data),
		),
	)
}

func ConfirmKeyboard(module, action string, id uint) tgbotapi.InlineKeyboardMarkup {
	confirmData := fmt.Sprintf("confirm:%s:%s:%d", module, action, id)
	cancelData := fmt.Sprintf("nav:%s", module)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅  Simpan", confirmData),
			tgbotapi.NewInlineKeyboardButtonData("❌  Batal", cancelData),
		),
	)
}

func sendOrEdit(chatID int64, messageID int, text string, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.Message {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ReplyMarkup = &keyboard
	edit.ParseMode = "HTML"
	msg, err := botClient.Send(edit)
	if err != nil {
		newMsg := tgbotapi.NewMessage(chatID, text)
		newMsg.ReplyMarkup = keyboard
		newMsg.ParseMode = "HTML"
		msg, _ = botClient.Send(newMsg)
	}
	return msg
}

func sendMessage(chatID int64, text string, keyboard interface{}) tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	sent, _ := botClient.Send(msg)
	return sent
}

// panelPrompt mengedit panel aktif menjadi prompt teks input dengan tombol Batal.
// Jika edit gagal (fallback pesan baru), PanelMessageID di-update ke ID baru
// sehingga langkah berikutnya edit ke pesan yang pasti bisa di-edit (self-healing).
func panelPrompt(state *UserState, chatID int64, text string) {
	cancel := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌  Batal", "nav:"+state.Module),
		),
	)
	msg := sendOrEdit(chatID, state.PanelMessageID, text, cancel)
	if msg.MessageID != 0 {
		state.PanelMessageID = msg.MessageID
	}
}

// panelConfirm mengedit panel aktif menjadi pesan konfirmasi dengan keyboard Simpan/Batal.
func panelConfirm(state *UserState, chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := sendOrEdit(chatID, state.PanelMessageID, text, keyboard)
	if msg.MessageID != 0 {
		state.PanelMessageID = msg.MessageID
	}
}

func moduleTitle(module string) string {
	titles := map[string]string{
		"suhu":        "🌡️ Menu Suhu",
		"informasi":   "📰 Menu Informasi",
		"penggajian":  "💰 Menu Penggajian",
		"pesan":       "💬 Menu Pesan",
		"user":        "👤 Menu User",
		"dokumen":     "📄 Menu Dokumen",
	}
	if t, ok := titles[module]; ok {
		return t
	}
	return module
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cb := update.CallbackQuery
	callback := tgbotapi.NewCallback(cb.ID, "")
	bot.Request(callback)

	data := cb.Data
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID
	telegramID := cb.From.ID

	if data == "cancel" {
		delete(userStates, telegramID)
		sendOrEdit(chatID, messageID, "❌ <i>Dibatalkan.</i>", MainMenuKeyboard())
		return
	}

	parts := strings.Split(data, ":")

	switch parts[0] {
	case "nav":
		delete(userStates, telegramID)
		if parts[1] == "main" {
			sendOrEdit(chatID, messageID, "🤖 <b>PCC Bot</b>\n<i>Silakan pilih fitur:</i>", MainMenuKeyboard())
		} else {
			module := parts[1]
			if module == "dokumen" {
				sendOrEdit(chatID, messageID, "<b>"+moduleTitle(module)+"</b>", SubmenuKeyboardReadOnly(module))
			} else {
				sendOrEdit(chatID, messageID, "<b>"+moduleTitle(module)+"</b>", SubmenuKeyboard(module))
			}
		}

	case "act":
		module := parts[1]
		action := parts[2]
		switch module {
		case "suhu":
			handleSuhuCallback(chatID, messageID, telegramID, action)
		case "informasi":
			handleInformasiCallback(chatID, messageID, telegramID, action)
		case "penggajian":
			handlePenggajianCallback(chatID, messageID, telegramID, action)
		case "pesan":
			handlePesanCallback(chatID, messageID, telegramID, action)
		case "user":
			handleUserCallback(chatID, messageID, telegramID, action)
		case "dokumen":
			handleDokumenCallback(chatID, messageID, action)
		}

	case "pick":
		module := parts[1]
		action := parts[2]
		id, _ := strconv.ParseUint(parts[3], 10, 64)
		switch module {
		case "suhu":
			suhuPickItem(chatID, messageID, telegramID, action, uint(id))
		case "informasi":
			informasiPickItem(chatID, messageID, telegramID, action, uint(id))
		case "penggajian":
			penggajianPickItem(chatID, messageID, telegramID, action, uint(id))
		case "pesan":
			pesanPickItem(chatID, messageID, telegramID, action, parts[3])
		case "user":
			userPickItem(chatID, messageID, telegramID, action, uint(id))
		}

	case "confirm":
		module := parts[1]
		action := parts[2]
		var id uint
		if len(parts) > 3 {
			parsedID, _ := strconv.ParseUint(parts[3], 10, 64)
			id = uint(parsedID)
		}
		switch module {
		case "suhu":
			suhuConfirm(chatID, messageID, telegramID, action, id)
		case "informasi":
			informasiConfirm(chatID, messageID, telegramID, action, id)
		case "penggajian":
			penggajianConfirm(chatID, messageID, telegramID, action, id)
		case "pesan":
			pesanConfirm(chatID, messageID, telegramID, action, parts[3])
		case "user":
			userConfirm(chatID, messageID, telegramID, action, id)
		}
	}
}

func handleTextMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	messageID := update.Message.MessageID

	state, exists := userStates[telegramID]
	if !exists {
		return
	}

	// hapus pesan input user agar panel tetap jadi satu-satunya pesan yang terlihat
	deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Request(deleteConfig)

	switch state.Module {
	case "suhu":
		suhuHandleInput(chatID, telegramID, text, state)
	case "informasi":
		informasiHandleInput(chatID, telegramID, text, state)
	case "penggajian":
		penggajianHandleInput(chatID, telegramID, text, state)
	case "pesan":
		pesanHandleInput(chatID, telegramID, text, state)
	case "user":
		userHandleInput(chatID, telegramID, text, state)
	}
}
