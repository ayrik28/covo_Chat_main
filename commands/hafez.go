package commands

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HafezCommand struct {
	bot *tgbotapi.BotAPI
}

type FalResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Interpreter string `json:"interpreter"`
}

func NewHafezCommand(bot *tgbotapi.BotAPI) *HafezCommand {
	return &HafezCommand{
		bot: bot,
	}
}

func (r *HafezCommand) getHafezFal() (string, error) {
	// خواندن فایل JSON
	jsonFile, err := os.Open(filepath.Join("jsonfile", "fal.json"))
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()

	// خواندن محتوای فایل
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	var fals []FalResponse
	if err := json.Unmarshal(byteValue, &fals); err != nil {
		return "", err
	}

	// اگر آرایه خالی باشد
	if len(fals) == 0 {
		return "", nil
	}

	// انتخاب یک فال رندوم
	randomIndex := rand.Intn(len(fals))
	fal := fals[randomIndex]

	// ساختن متن فال
	result := "🎭 *فال حافظ*\n\n" +
		"📜 *عنوان فال:* " + fal.Title + "\n" +
		"🔢 *شماره فال:* " + strconv.Itoa(fal.Id) + "\n\n" +
		"📝 *تفسیر فال:*\n" + fal.Interpreter

	return result, nil
}

func (r *HafezCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	}

	// دریافت فال
	text, err := r.getHafezFal()
	if err != nil {
		log.Printf("Error getting hafez: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ متأسفانه در دریافت فال خطایی رخ داد. لطفاً دوباره تلاش کنید.")
	}

	if text == "" {
		return tgbotapi.NewMessage(chatID, "❌ متأسفانه فالی یافت نشد. لطفاً دوباره تلاش کنید.")
	}

	// اضافه کردن دکمه فال جدید به پیام
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎲 فال جدید", "new_hafez"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard
	return msg
}

// HandleCallback handles the callback queries from inline keyboard
func (r *HafezCommand) HandleCallback(update tgbotapi.Update) tgbotapi.CallbackConfig {
	if update.CallbackQuery.Data == "new_hafez" {
		// ارسال فال جدید
		msg := r.Handle(update)
		r.bot.Send(msg)
	}

	// تایید دریافت callback
	return tgbotapi.NewCallback(update.CallbackQuery.ID, "✅")
}

// HandleStart returns the start message for hafez command
func (r *HafezCommand) HandleStart() string {
	return `🎭 *فال حافظ*

با استفاده از این قابلیت می‌توانید فال حافظ بگیرید!

📜 *دستورات:*
• /فال - گرفتن فال حافظ
• دکمه "فال جدید" - گرفتن فال جدید

💫 *ویژگی‌ها:*
• نمایش عنوان فال
• شماره فال
• تعبیر و تفسیر فال
• امکان گرفتن فال‌های متنوع

برای شروع، دستور /فال را بفرستید یا روی دکمه فال حافظ کلیک کنید.`
}
