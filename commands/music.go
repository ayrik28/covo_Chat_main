package commands

import (
	"fmt"
	"log"
	"redhat-bot/ai"
	"redhat-bot/limiter"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MusicCommand struct {
	aiClient    *ai.DeepSeekClient
	rateLimiter *limiter.RateLimiter
	bot         *tgbotapi.BotAPI
}

func NewMusicCommand(aiClient *ai.DeepSeekClient, rateLimiter *limiter.RateLimiter, bot *tgbotapi.BotAPI) *MusicCommand {
	return &MusicCommand{
		aiClient:    aiClient,
		rateLimiter: rateLimiter,
		bot:         bot,
	}
}

func (r *MusicCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// بررسی محدودیت درخواست
	if allowed, message := r.rateLimiter.CheckRateLimit(userID); !allowed {
		return tgbotapi.NewMessage(chatID, message)
	}

	// اگر پیام ریپلای است، آن را پردازش کن
	if update.Message.ReplyToMessage != nil {
		return r.handleReply(update)
	}

	// ارسال پیام اولیه برای درخواست موسیقی
	response := `🎵 *پیشنهاد موسیقی*

چه نوع آهنگی دوست داری؟ حس و حال الانت چیه؟ (غمگین، شاد، آروم، انگیزشی...)

لطفاً جواب رو ریپلای کن به همین پیام.`

	msg := tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

func (r *MusicCommand) handleReply(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	userPreference := update.Message.Text

	// ارسال پیام "در حال پردازش"
	processingMsg := tgbotapi.NewMessage(chatID, "درحال پردازش - کمی شکیبا باشید ✨")
	sentMsg, err := r.bot.Send(processingMsg)
	if err != nil {
		log.Printf("خطا در ارسال پیام پردازش: %v", err)
	}

	// ساخت درخواست برای هوش مصنوعی
	prompt := fmt.Sprintf(`"%s" این آهنگ با این حس و حال لینک بفرس

توضیحات خیلی کوتاه باشه و لینک یوتیوب و اسپاتیفای درجا بده.`, userPreference)

	// دریافت پاسخ از هوش مصنوعی
	response, err := r.aiClient.AskQuestion(prompt)
	if err != nil {
		log.Printf("خطا در دریافت پیشنهاد موسیقی: %v", err)
		// حذف پیام پردازش و ارسال پیام خطا
		if sentMsg.MessageID != 0 {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID)
			r.bot.Send(deleteMsg)
		}
		return tgbotapi.NewMessage(chatID, "❌ متأسفانه نتوانستم پیشنهاد موسیقی ارائه دهم. لطفاً دوباره تلاش کنید.")
	}

	// فرمت‌بندی پاسخ
	formattedResponse := fmt.Sprintf("🎵 *پیشنهاد موسیقی*\n\n%s", response)

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
