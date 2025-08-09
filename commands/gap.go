package commands

import (
	"fmt"
	"redhat-bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GapCommand struct {
	bot          *tgbotapi.BotAPI
	storage      *storage.MySQLStorage
	hafezCommand *HafezCommand
}

func NewGapCommand(bot *tgbotapi.BotAPI, storage *storage.MySQLStorage, hafezCommand *HafezCommand) *GapCommand {
	return &GapCommand{
		bot:          bot,
		storage:      storage,
		hafezCommand: hafezCommand,
	}
}

func (r *GapCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID

	// پیام راهنمای دستورات
	response := `📱 *دستورات ربات کوو*\n\nبرای استفاده از دستورات، روی دکمه‌های زیر کلیک کنید:`

	// ساخت کیبورد اینلاین
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		// ردیف اول - دستورات اصلی
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 وضعیت ربات", "status"),
			tgbotapi.NewInlineKeyboardButtonData("📕 فال حافظ", "hafez"),
		),
		// ردیف دوم - دستورات کراش
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💘 فعال کردن کراش", "enable_crush"),
			tgbotapi.NewInlineKeyboardButtonData("💔 غیرفعال کردن کراش", "disable_crush"),
		),
		// ردیف سوم - دستورات کراش و دلقک
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👀 وضعیت کراش", "crush_status"),
			tgbotapi.NewInlineKeyboardButtonData("🤡 دلقک", "clown_help"),
		),
		// ردیف چهارم - راهنما
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 راهنمای کامل", "full_help"),
			tgbotapi.NewInlineKeyboardButtonData("❓ راهنمای گروه", "group_help"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard
	return msg
}

// HandleCallback handles the callback queries from inline keyboard
func (r *GapCommand) HandleCallback(update tgbotapi.Update) tgbotapi.CallbackConfig {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID

	// ساخت یک پیام جدید برای نمایش نتیجه عملیات
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = tgbotapi.ModeMarkdown

	switch data {
	case "hafez":
		// ارسال فال حافظ
		msg.Text = "در حال دریافت فال..."
		r.bot.Send(msg)
		response := r.hafezCommand.Handle(update) // پاس دادن کل update
		r.bot.Send(response)
		return tgbotapi.NewCallback(update.CallbackQuery.ID, "✅")

	case "status":
		// نمایش وضعیت ربات
		msg.Text = `📊 *وضعیت ربات:*

✅ ربات فعال و آماده به کار است
⚡️ سرعت پاسخگویی: عالی
🔋 وضعیت سرور: آنلاین
🤖 نسخه ربات: 1.0.0`

	case "enable_crush":
		// فعال کردن کراش
		if err := r.storage.SetCrushEnabled(chatID, true); err != nil {
			msg.Text = "❌ خطا در فعال‌سازی قابلیت کراش"
		} else {
			msg.Text = "💘 *قابلیت کراش با موفقیت فعال شد!* ✅\n\n🔥 از این لحظه هر 10 ساعت یک بار، دو نفر از اعضای گروه به صورت تصادفی به عنوان کراش انتخاب می‌شوند!\n\n👀 منتظر اعلام اولین جفت کراش باشید..."
		}

	case "disable_crush":
		// غیرفعال کردن کراش
		if err := r.storage.SetCrushEnabled(chatID, false); err != nil {
			msg.Text = "❌ خطا در غیرفعال‌سازی قابلیت کراش"
		} else {
			msg.Text = "💔 *قابلیت کراش غیرفعال شد!* ❌\n\n🚫 دیگر اعلام خودکار کراش در این گروه انجام نخواهد شد."
		}

	case "crush_status":
		// نمایش وضعیت کراش
		enabled, err := r.storage.IsCrushEnabled(chatID)
		if err != nil {
			msg.Text = "❌ خطا در بررسی وضعیت کراش"
		} else {
			status := "فعال ✅"
			if !enabled {
				status = "غیرفعال ❌"
			}
			msg.Text = fmt.Sprintf(`💘 *وضعیت قابلیت کراش:*

🎯 وضعیت: %s
⏰ زمان اعلام: هر 10 ساعت
👥 نحوه انتخاب: تصادفی از بین اعضای گروه`, status)
		}

	case "clown_help":
		msg.Text = `🤡 *راهنمای قابلیت دلقک*

با این قابلیت می‌توانید به صورت هوشمند به افراد توهین کنید!

👉 برای استفاده:
1. نام فرد یا @username او را تایپ کنید
2. روی دکمه 🤡 دلقک کلیک کنید
3. منتظر پاسخ هوشمندانه ربات باشید!`

	case "full_help":
		msg.Text = `📚 *راهنمای کامل ربات کوو*

🤖 *قابلیت‌های اصلی:*
• پرسش و پاسخ هوشمند
• ساخت جوک
• پیشنهاد موزیک
• قابلیت کراش
• قابلیت دلقک

💡 *نکات مهم:*
• ربات در گروه‌ها و چت خصوصی کار می‌کند
• پاسخ‌ها با هوش مصنوعی تولید می‌شوند
• قابلیت کراش هر 10 ساعت یکبار اجرا می‌شود

برای اطلاعات بیشتر روی دکمه‌های مختلف کلیک کنید.`

	case "group_help":
		msg.Text = `❓ *راهنمای گروه*

👥 *قابلیت‌های گروه:*
• اعلام خودکار کراش
• توهین هوشمند به اعضا
• ثبت پیام‌های گروه
• خلاصه روزانه

⚙️ *تنظیمات:*
• فعال/غیرفعال کردن کراش
• تنظیم زمان اعلام کراش
• مدیریت اعضای گروه

برای استفاده از هر قابلیت، روی دکمه مربوطه کلیک کنید.`

	}

	// ارسال پیام نتیجه
	r.bot.Send(msg)

	// تایید دریافت callback
	return tgbotapi.NewCallback(update.CallbackQuery.ID, "✅")
}
