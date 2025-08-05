package commands

import (
	"fmt"
	"log"
	"redhat-bot/ai"
	"redhat-bot/limiter"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CovoCommand struct {
	aiClient    *ai.DeepSeekClient
	rateLimiter *limiter.RateLimiter
	bot         *tgbotapi.BotAPI
}

func NewCovoCommand(aiClient *ai.DeepSeekClient, rateLimiter *limiter.RateLimiter, bot *tgbotapi.BotAPI) *CovoCommand {
	return &CovoCommand{
		aiClient:    aiClient,
		rateLimiter: rateLimiter,
		bot:         bot,
	}
}

func (r *CovoCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// بررسی محدودیت درخواست
	if allowed, message := r.rateLimiter.CheckRateLimit(userID); !allowed {
		return tgbotapi.NewMessage(chatID, message)
	}

	// استخراج سوال از دستور
	text := update.Message.Text
	question := strings.TrimSpace(strings.TrimPrefix(text, "/covo"))

	if question == "" {
		msg := tgbotapi.NewMessage(chatID, "🤖 *دستیار هوشمند کوو*\n\nنحوه استفاده: `/covo <سوال شما>`\n\nهر سوالی دارید بپرسید! من اینجا هستم تا کمک کنم. 💡")
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	}

	// ارسال پیام "در حال پردازش"
	processingMsg := tgbotapi.NewMessage(chatID, "درحال پردازش - کمی شکیبا باشید ✨")
	sentMsg, err := r.bot.Send(processingMsg)
	if err != nil {
		log.Printf("خطا در ارسال پیام پردازش: %v", err)
	}

	// دریافت پاسخ از هوش مصنوعی
	response, err := r.aiClient.AskQuestion(question)
	if err != nil {
		log.Printf("خطا در دریافت پاسخ هوش مصنوعی: %v", err)
		// حذف پیام پردازش و ارسال پیام خطا
		if sentMsg.MessageID != 0 {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID)
			r.bot.Send(deleteMsg)
		}
		return tgbotapi.NewMessage(chatID, "❌ متأسفانه در پردازش سوال شما مشکلی پیش آمد. لطفاً دوباره تلاش کنید.")
	}

	// فرمت‌بندی پاسخ
	formattedResponse := fmt.Sprintf("🤖 *هوش مصنوعی کوو*\n\n%s", response)

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
