package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GapCommand struct {
	bot *tgbotapi.BotAPI
}

func NewGapCommand(bot *tgbotapi.BotAPI) *GapCommand {
	return &GapCommand{
		bot: bot,
	}
}

func (r *GapCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID

	// پیام راهنمای دستورات گروه
	response := `📱 *دستورات مخصوص گروه*

🔹 *دستورات عمومی:*
• /covo <سوال> - هر سوالی دارید بپرسید!
• /cj <موضوع> - جوک خنده‌دار درباره هر موضوعی
• /music - پیشنهاد موسیقی بر اساس سلیقه شما
• /crs - بررسی وضعیت بات

🔸 *دستورات مخصوص گروه:*
• /crushon - فعال‌سازی قابلیت کراش (هر 10 ساعت)
• /crushoff - غیرفعال‌سازی قابلیت کراش
• /کراشوضعیت - مشاهده وضعیت قابلیت کراش
• /clown <نام> - توهین هوشمند به شخص مورد نظر

💡 *راهنما:*
• /gap - نمایش این پیام راهنما
• /help - راهنمای کامل بات

برای استفاده از هر دستور، کافیست آن را در گروه تایپ کنید.`

	msg := tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}
