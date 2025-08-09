package commands

import (
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
			tgbotapi.NewInlineKeyboardButtonData("🎛️ قابلیت‌ها", "features"),
		),
		// ردیف چهارم - قفل‌ها
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔒 قفل", "locks"),
		),
		// ردیف پنجم - راهنما
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
	case "features":
		// نمایش دکمه‌های قابلیت‌ها (کِراش و فال)
		crushEnabled, _ := r.storage.IsCrushEnabled(chatID)
		hafezEnabled, _ := r.storage.IsFeatureEnabled(chatID, "hafez")

		crushIcon := "❌"
		if crushEnabled {
			crushIcon = "✅"
		}
		hafezIcon := "❌"
		if hafezEnabled {
			hafezIcon = "✅"
		}

		featuresKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("💘 کراش "+crushIcon, "toggle_crush"),
				tgbotapi.NewInlineKeyboardButtonData("📕 فال "+hafezIcon, "toggle_hafez"),
			),
		)
		msg.Text = "🎛️ تنظیمات قابلیت‌ها:\n\nبا دکمه‌های زیر می‌توانید قابلیت‌ها را فعال/غیرفعال کنید."
		msg.ReplyMarkup = featuresKeyboard

	case "status":
		// نمایش وضعیت ربات
		msg.Text = `📊 *وضعیت ربات:*

✅ ربات فعال و آماده به کار است
⚡️ سرعت پاسخگویی: عالی
🔋 وضعیت سرور: آنلاین
🤖 نسخه ربات: 1.0.0`

	case "toggle_crush":
		// تغییر وضعیت کراش + ارسال پیام معادل دستور رسمی
		enabled, err := r.storage.IsCrushEnabled(chatID)
		if err != nil {
			msg.Text = "❌ خطا در بررسی وضعیت کراش"
			break
		}
		newEnabled := !enabled
		if err := r.storage.SetCrushEnabled(chatID, newEnabled); err != nil {
			msg.Text = "❌ خطا در تغییر وضعیت کراش"
			break
		}
		if newEnabled {
			msg.Text = "💘 *قابلیت کراش با موفقیت فعال شد!* ✅\n\n🔥 از این لحظه هر 10 ساعت یک بار، دو نفر از اعضای گروه به صورت تصادفی به عنوان کراش انتخاب می‌شوند!\n\n👀 منتظر اعلام اولین جفت کراش باشید..."
			msg.ParseMode = tgbotapi.ModeMarkdown
		} else {
			msg.Text = "💔 *قابلیت کراش غیرفعال شد!* ❌\n\n🚫 دیگر اعلام خودکار کراش در این گروه انجام نخواهد شد.\n\n✅ برای فعال‌سازی مجدد از دستور `/crushon` استفاده کنید."
			msg.ParseMode = tgbotapi.ModeMarkdown
		}
		// بازسازی کیبورد قابلیت‌ها
		crushEnabled, _ := r.storage.IsCrushEnabled(chatID)
		hafezEnabled, _ := r.storage.IsFeatureEnabled(chatID, "hafez")
		crushIcon := "❌"
		if crushEnabled {
			crushIcon = "✅"
		}
		hafezIcon := "❌"
		if hafezEnabled {
			hafezIcon = "✅"
		}
		featuresKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("💘 کراش "+crushIcon, "toggle_crush"),
				tgbotapi.NewInlineKeyboardButtonData("📕 فال "+hafezIcon, "toggle_hafez"),
			),
		)
		msg.ReplyMarkup = featuresKeyboard

	case "toggle_hafez":
		// تغییر وضعیت فال
		enabled, err := r.storage.IsFeatureEnabled(chatID, "hafez")
		if err != nil {
			msg.Text = "❌ خطا در بررسی وضعیت فال"
			break
		}
		if err := r.storage.SetFeatureEnabled(chatID, "hafez", !enabled); err != nil {
			msg.Text = "❌ خطا در تغییر وضعیت فال"
			break
		}
		// بازسازی کیبورد قابلیت‌ها
		crushEnabled, _ := r.storage.IsCrushEnabled(chatID)
		hafezEnabled, _ := r.storage.IsFeatureEnabled(chatID, "hafez")
		crushIcon := "❌"
		if crushEnabled {
			crushIcon = "✅"
		}
		hafezIcon := "❌"
		if hafezEnabled {
			hafezIcon = "✅"
		}
		featuresKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("💘 کراش "+crushIcon, "toggle_crush"),
				tgbotapi.NewInlineKeyboardButtonData("📕 فال "+hafezIcon, "toggle_hafez"),
			),
		)
		msg.Text = "وضعیت قابلیت‌ها بروز شد."
		msg.ReplyMarkup = featuresKeyboard

	case "clown_help":
		msg.Text = `🤡 *راهنمای قابلیت دلقک*

با این قابلیت می‌توانید به صورت هوشمند به افراد توهین کنید!

👉 برای استفاده:
1. نام فرد یا @username او را تایپ کنید
2. بنویسید: دلقک <نام>
3. منتظر پاسخ هوشمندانه ربات باشید!`

	case "locks":
		// نمایش وضعیت قفل‌ها و امکان تغییر
		clownEnabled, _ := r.storage.IsClownEnabled(chatID)
		linkEnabled, _ := r.storage.IsFeatureEnabled(chatID, "link")
		badwordEnabled, _ := r.storage.IsFeatureEnabled(chatID, "badword")

		clownIcon := "❌"
		if clownEnabled {
			clownIcon = "✅"
		}
		linkIcon := "❌"
		if linkEnabled {
			linkIcon = "✅"
		}

		locksKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🤡 دلقک "+clownIcon, "toggle_clown"),
				tgbotapi.NewInlineKeyboardButtonData("🔗 لینک "+linkIcon, "toggle_link"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🚫 فحش "+func() string {
					if badwordEnabled {
						return "✅"
					} else {
						return "❌"
					}
				}(), "toggle_badword"),
			),
		)
		msg.Text = "🔒 تنظیمات قفل‌ها:\n\nبا دکمه‌های زیر می‌توانید قفل‌ها را فعال/غیرفعال کنید."
		msg.ReplyMarkup = locksKeyboard

	case "toggle_clown":
		// تغییر وضعیت دلقک
		enabled, err := r.storage.IsClownEnabled(chatID)
		if err != nil {
			msg.Text = "❌ خطا در بررسی وضعیت دلقک"
			break
		}
		if err := r.storage.SetClownEnabled(chatID, !enabled); err != nil {
			msg.Text = "❌ خطا در تغییر وضعیت دلقک"
			break
		}

		// ساخت کیبورد بروز‌شده (هر دو قفل)
		clownEnabled, _ := r.storage.IsClownEnabled(chatID)
		linkEnabled, _ := r.storage.IsFeatureEnabled(chatID, "link")
		clownIcon := "❌"
		if clownEnabled {
			clownIcon = "✅"
		}
		linkIcon := "❌"
		if linkEnabled {
			linkIcon = "✅"
		}
		locksKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🤡 دلقک "+clownIcon, "toggle_clown"),
				tgbotapi.NewInlineKeyboardButtonData("🔗 لینک "+linkIcon, "toggle_link"),
			),
		)
		msg.Text = "وضعیت قفل‌ها بروز شد."
		msg.ReplyMarkup = locksKeyboard

	case "toggle_link":
		// تغییر وضعیت لینک
		enabled, err := r.storage.IsFeatureEnabled(chatID, "link")
		if err != nil {
			msg.Text = "❌ خطا در بررسی وضعیت لینک"
			break
		}
		if err := r.storage.SetFeatureEnabled(chatID, "link", !enabled); err != nil {
			msg.Text = "❌ خطا در تغییر وضعیت لینک"
			break
		}
		// ساخت کیبورد بروز‌شده (هر دو قفل)
		clownEnabled, _ := r.storage.IsClownEnabled(chatID)
		linkEnabled, _ := r.storage.IsFeatureEnabled(chatID, "link")
		clownIcon := "❌"
		if clownEnabled {
			clownIcon = "✅"
		}
		linkIcon := "❌"
		if linkEnabled {
			linkIcon = "✅"
		}
		locksKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🤡 دلقک "+clownIcon, "toggle_clown"),
				tgbotapi.NewInlineKeyboardButtonData("🔗 لینک "+linkIcon, "toggle_link"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🚫 فحش "+func() string {
					if enabled {
						return "✅"
					} else {
						return "❌"
					}
				}(), "toggle_badword"),
			),
		)
		msg.Text = "وضعیت قفل‌ها بروز شد."
		msg.ReplyMarkup = locksKeyboard

	case "toggle_badword":
		// تغییر وضعیت فحش
		enabled, err := r.storage.IsFeatureEnabled(chatID, "badword")
		if err != nil {
			msg.Text = "❌ خطا در بررسی وضعیت فحش"
			break
		}
		if err := r.storage.SetFeatureEnabled(chatID, "badword", !enabled); err != nil {
			msg.Text = "❌ خطا در تغییر وضعیت فحش"
			break
		}
		// ساخت کیبورد بروز‌شده
		clownEnabled, _ := r.storage.IsClownEnabled(chatID)
		linkEnabled, _ := r.storage.IsFeatureEnabled(chatID, "link")
		badwordEnabled, _ := r.storage.IsFeatureEnabled(chatID, "badword")
		clownIcon := "❌"
		if clownEnabled {
			clownIcon = "✅"
		}
		linkIcon := "❌"
		if linkEnabled {
			linkIcon = "✅"
		}
		badwordIcon := "❌"
		if badwordEnabled {
			badwordIcon = "✅"
		}
		locksKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🤡 دلقک "+clownIcon, "toggle_clown"),
				tgbotapi.NewInlineKeyboardButtonData("🔗 لینک "+linkIcon, "toggle_link"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🚫 فحش "+badwordIcon, "toggle_badword"),
			),
		)
		msg.Text = "وضعیت قفل‌ها بروز شد."
		msg.ReplyMarkup = locksKeyboard

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
