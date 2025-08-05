package commands

import (
	"fmt"
	"redhat-bot/limiter"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CrsCommand struct {
	rateLimiter *limiter.RateLimiter
}

func NewCrsCommand(rateLimiter *limiter.RateLimiter) *CrsCommand {
	return &CrsCommand{
		rateLimiter: rateLimiter,
	}
}

func (r *CrsCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID

	// فرمت‌بندی پاسخ
	response := fmt.Sprintf("📊 *وضعیت بات کوو*\n\n" +
		"🎯 *وضعیت:* درخواست‌های نامحدود در دسترس\n" +
		"🔄 *بدون محدودیت روزانه*\n" +
		"⚡ *بدون تأخیر*\n\n" +
		"💡 از `/covo <سوال>` برای پرسش استفاده کنید!\n" +
		"😄 از `/covoJoke <موضوع>` برای جوک‌های خنده‌دار استفاده کنید!")

	msg := tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown

	return msg
}
