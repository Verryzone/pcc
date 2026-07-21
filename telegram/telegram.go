package telegram

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var botClient *tgbotapi.BotAPI
var dbConn *gorm.DB

func StartBot(db *gorm.DB) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		fmt.Println("[Telegram] TELEGRAM_BOT_TOKEN tidak ditemukan, bot tidak dijalankan.")
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("[Telegram] Gagal membuat bot: %v\n", err)
		return
	}

	botClient = bot
	dbConn = db

	fmt.Printf("[Telegram] Bot berhasil diinisialisasi: @%s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update)
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommand(bot, update)
			} else {
				handleTextMessage(bot, update)
			}
		}
	}
}

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	switch msg.Command() {
	case "start":
		text := "🤖 <b>PCC Bot</b>\n<i>Silakan pilih fitur:</i>"
		keyboard := MainMenuKeyboard()
		reply := tgbotapi.NewMessage(msg.Chat.ID, text)
		reply.ReplyMarkup = keyboard
		reply.ParseMode = "HTML"
		bot.Send(reply)
	default:
		reply := tgbotapi.NewMessage(msg.Chat.ID, "⚠️ Perintah tidak dikenali. Ketik /start untuk memulai.")
		bot.Send(reply)
	}
}
