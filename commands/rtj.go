package commands

import (
	"fmt"
	"log"
	"redhat-bot/ai"
	"redhat-bot/limiter"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CovoJokeCommand struct {
	aiClient    *ai.DeepSeekClient
	rateLimiter *limiter.RateLimiter
	bot         *tgbotapi.BotAPI
}

func NewCovoJokeCommand(aiClient *ai.DeepSeekClient, rateLimiter *limiter.RateLimiter, bot *tgbotapi.BotAPI) *CovoJokeCommand {
	return &CovoJokeCommand{
		aiClient:    aiClient,
		rateLimiter: rateLimiter,
		bot:         bot,
	}
}

func (r *CovoJokeCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// بررسی محدودیت درخواست
	if allowed, message := r.rateLimiter.CheckRateLimit(userID); !allowed {
		return tgbotapi.NewMessage(chatID, message)
	}

	text := update.Message.Text
	topic := strings.TrimSpace(strings.TrimPrefix(text, "/covoJoke"))

	if topic == "" {
		msg := tgbotapi.NewMessage(chatID, "😄 *تولیدکننده جوک کوو*\n\nنحوه استفاده: `/covoJoke <موضوع>`\n\nمن یک جوک خنده‌دار درباره موضوع انتخابی شما تولید می‌کنم! 🎭")
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	}

	// ارسال پیام "در حال پردازش"
	processingMsg := tgbotapi.NewMessage(chatID, "درحال پردازش - کمی شکیبا باشید ✨")
	sentMsg, err := r.bot.Send(processingMsg)
	if err != nil {
		log.Printf("خطا در ارسال پیام پردازش: %v", err)
	}

	// ساخت درخواست برای هوش مصنوعی
	prompt := fmt.Sprintf("هی، یک جوک خنده‌دار و مناسب خانواده درباره '%s' تولید کن و ارسال کن.", topic)

	// استفاده از AskQuestion برای ارسال درخواست
	joke, err := r.aiClient.AskQuestion(prompt)
	if err != nil {
		log.Printf("خطا در تولید جوک: %v", err)
		// حذف پیام پردازش و ارسال پیام خطا
		if sentMsg.MessageID != 0 {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID)
			r.bot.Send(deleteMsg)
		}
		return tgbotapi.NewMessage(chatID, "❌ متأسفانه نتوانستم جوک تولید کنم. لطفاً دوباره تلاش کنید.")
	}

	// فرمت‌بندی پاسخ
	formattedResponse := fmt.Sprintf("😄 *تولیدکننده جوک کوو*\n\n*موضوع:* %s\n\n%s", topic, joke)

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
