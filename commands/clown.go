package commands

import (
	"fmt"
	"log"
	"redhat-bot/ai"
	"redhat-bot/limiter"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ClownCommand struct {
	aiClient    *ai.DeepSeekClient
	rateLimiter *limiter.RateLimiter
	bot         *tgbotapi.BotAPI
}

func NewClownCommand(aiClient *ai.DeepSeekClient, rateLimiter *limiter.RateLimiter, bot *tgbotapi.BotAPI) *ClownCommand {
	return &ClownCommand{
		aiClient:    aiClient,
		rateLimiter: rateLimiter,
		bot:         bot,
	}
}

func (r *ClownCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// بررسی محدودیت درخواست
	if allowed, message := r.rateLimiter.CheckRateLimit(userID); !allowed {
		return tgbotapi.NewMessage(chatID, message)
	}

	// استخراج نام مخاطب از دستور
	text := update.Message.Text
	targetName := strings.TrimSpace(strings.TrimPrefix(text, "/clown"))

	if targetName == "" {
		msg := tgbotapi.NewMessage(chatID, "🤡 *دستور دلقک*\n\nنحوه استفاده: `/clown <نام مخاطب>`\n\nمثال: `/clown علی` یا `/clown @username`")
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	}

	// ارسال پیام "در حال پردازش"
	processingMsg := tgbotapi.NewMessage(chatID, "🤡 در حال آماده‌سازی توهین - کمی صبر کنید...")
	sentMsg, err := r.bot.Send(processingMsg)
	if err != nil {
		log.Printf("خطا در ارسال پیام پردازش: %v", err)
	}

	// ساخت پرامپت برای توهین
	prompt := fmt.Sprintf("به %s چند تا فحش بده", targetName)

	// دریافت پاسخ از هوش مصنوعی
	response, err := r.aiClient.AskQuestion(prompt)
	if err != nil {
		log.Printf("خطا در دریافت پاسخ هوش مصنوعی: %v", err)
		// حذف پیام پردازش و ارسال پیام خطا
		if sentMsg.MessageID != 0 {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID)
			r.bot.Send(deleteMsg)
		}
		return tgbotapi.NewMessage(chatID, "❌ متأسفانه در پردازش درخواست شما مشکلی پیش آمد. لطفاً دوباره تلاش کنید.")
	}

	// فرمت‌بندی پاسخ
	formattedResponse := fmt.Sprintf("🤡 *دلقک به %s:*\n\n%s", targetName, response)

	// ویرایش پیام پردازش با پاسخ نهایی
	if sentMsg.MessageID != 0 {
		editMsg := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID, formattedResponse)
		editMsg.ParseMode = tgbotapi.ModeMarkdown
		_, err = r.bot.Send(editMsg)
		if err != nil {
			log.Printf("خطا در ویرایش پیام: %v", err)
			// اگر ویرایش ناموفق بود، پیام جدید ارسال کن
			msg := tgbotapi.NewMessage(chatID, formattedResponse)
			msg.ParseMode = tgbotapi.ModeMarkdown
			return msg
		}
		// پیام خالی برگردان چون پیام ویرایش شده
		return tgbotapi.MessageConfig{}
	}

	// اگر پیام پردازش ارسال نشده بود، پیام عادی ارسال کن
	msg := tgbotapi.NewMessage(chatID, formattedResponse)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}
