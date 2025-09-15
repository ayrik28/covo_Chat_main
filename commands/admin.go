package commands

import (
	"fmt"
	"log"
	"strings"

	"redhat-bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminCommand struct {
	bot     *tgbotapi.BotAPI
	storage *storage.MySQLStorage
	// وضعیت موقت برای دریافت ورودی لینک جدید از ادمین‌ها (کلاینت خصوصی)
	pendingAdd map[int64]bool // key: admin user id
}

// لیست ادمین‌های مجاز
var adminUsers = map[int64]string{
	1234567890: "x",
	2345678901: "y",
}

func NewAdminCommand(bot *tgbotapi.BotAPI, storage *storage.MySQLStorage) *AdminCommand {
	return &AdminCommand{
		bot:        bot,
		storage:    storage,
		pendingAdd: make(map[int64]bool),
	}
}

// بررسی اینکه آیا کاربر ادمین است یا نه
func (r *AdminCommand) IsAdmin(userID int64) bool {
	_, exists := adminUsers[userID]
	return exists
}

// نمایش پیام خوش‌آمدگویی برای ادمین‌ها
func (r *AdminCommand) GetAdminWelcome(userID int64) string {
	name, exists := adminUsers[userID]
	if !exists {
		return ""
	}

	switch userID {
	case 1234567890:
		return fmt.Sprintf(`🌟 *سلام %s عزیز!* 🌟


💝 مرسی که یا ایده هات منو به اینجا رسوندی و این همه قابلیت جالب اضافه کردی!

🛠️ *دستورات ادمین:*
• /showusers - نمایش لیست تمام کاربران
• /showgroups - نمایش لیست تمام گروه‌ها
• /admin - بازگشت به منوی ادمین

✨ از اینکه منو ساختی ممنونم! 💖`, name)

	case 2345678901:
		return fmt.Sprintf(`🌟 *سلام %s عزیز!* 🌟

🎯 خوش اومدی به پنل ادمین!

🛠️ *دستورات ادمین:*
• /showusers - نمایش لیست تمام کاربران
• /showgroups - نمایش لیست تمام گروه‌ها
• /admin - بازگشت به منوی ادمین

✨ آماده خدمت‌رسانی هستم! 💪`, name)

	default:
		return ""
	}
}

func (r *AdminCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	// فقط در چت خصوصی کار می‌کند
	if update.Message.Chat.Type != "private" {
		return tgbotapi.NewMessage(chatID, "❌ این دستور فقط در چت خصوصی با بات قابل استفاده است.")
	}

	// بررسی دسترسی ادمین
	if !r.IsAdmin(userID) {
		return tgbotapi.NewMessage(chatID, "❌ شما دسترسی ادمین ندارید.")
	}

	// نمایش منوی ادمین
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 نمایش کاربران", "admin_showusers"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏢 نمایش گروه‌ها", "admin_showgroups"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📣 تبلیغات / عضویت اجباری", "admin_ads"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, r.GetAdminWelcome(userID))
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard
	return msg
}

// HandleCallback handles admin callback queries
func (r *AdminCommand) HandleCallback(update tgbotapi.Update) tgbotapi.CallbackConfig {
	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Chat.ID

	// بررسی دسترسی ادمین
	if !r.IsAdmin(userID) {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "❌ دسترسی ندارید")
		return callback
	}

	data := update.CallbackQuery.Data
	switch {
	case data == "admin_ads":
		// نمایش لیست کانال‌های اجباری و دکمه‌های مدیریت
		channels, err := r.storage.ListRequiredChannels(0)
		if err != nil {
			log.Printf("Error getting required channels: %v", err)
			return tgbotapi.NewCallback(update.CallbackQuery.ID, "❌ خطا در دریافت لینک‌ها")
		}

		text := "📣 لینک‌های عضویت اجباری:\n\n"
		if len(channels) == 0 {
			text += "فعلاً لینکی ثبت نشده است."
		} else {
			for _, ch := range channels {
				link := ch.Link
				if link == "" && ch.ChannelUsername != "" {
					link = "https://t.me/" + ch.ChannelUsername
				}
				if ch.Title == "" {
					text += fmt.Sprintf("#%d — %s\n", ch.ID, link)
				} else {
					text += fmt.Sprintf("#%d — %s (%s)\n", ch.ID, ch.Title, link)
				}
			}
		}

		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ اضافه لینک جدید", "admin_ads_add"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🗑️ پاک کردن لینک", "admin_ads_del_menu"),
			),
		)
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = kb
		r.bot.Send(msg)
		return tgbotapi.NewCallback(update.CallbackQuery.ID, "")

	case data == "admin_ads_add":
		r.pendingAdd[userID] = true
		prompt := "لطفاً لینک کانال را ارسال کنید.\n\nفرمت‌های قابل قبول:\n• لینک عمومی: https://t.me/<username> | عنوان دلخواه\n• لینک خصوصی: https://t.me/+joincode | عنوان دلخواه\n(می‌توانید عنوان را ننویسید)"
		r.bot.Send(tgbotapi.NewMessage(chatID, prompt))
		return tgbotapi.NewCallback(update.CallbackQuery.ID, "")

	case data == "admin_ads_del_menu":
		channels, err := r.storage.ListRequiredChannels(0)
		if err != nil {
			log.Printf("Error getting required channels: %v", err)
			return tgbotapi.NewCallback(update.CallbackQuery.ID, "❌ خطا در دریافت لینک‌ها")
		}
		if len(channels) == 0 {
			r.bot.Send(tgbotapi.NewMessage(chatID, "📭 لینکی برای حذف وجود ندارد."))
			return tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		}
		// ساخت دکمه‌های حذف به‌صورت چند ردیفه
		var rows [][]tgbotapi.InlineKeyboardButton
		for _, ch := range channels {
			label := fmt.Sprintf("حذف #%d", ch.ID)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, fmt.Sprintf("admin_ads_del:%d", ch.ID)),
			))
		}
		kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg := tgbotapi.NewMessage(chatID, "یک لینک را برای حذف انتخاب کنید:")
		msg.ReplyMarkup = kb
		r.bot.Send(msg)
		return tgbotapi.NewCallback(update.CallbackQuery.ID, "")

	case strings.HasPrefix(data, "admin_ads_del:"):
		var id uint
		if _, err := fmt.Sscanf(data, "admin_ads_del:%d", &id); err != nil {
			return tgbotapi.NewCallback(update.CallbackQuery.ID, "شناسه نامعتبر")
		}
		if err := r.storage.RemoveRequiredChannel(id); err != nil {
			log.Printf("Error removing required channel: %v", err)
			return tgbotapi.NewCallback(update.CallbackQuery.ID, "❌ خطا در حذف")
		}
		r.bot.Send(tgbotapi.NewMessage(chatID, "✅ لینک حذف شد"))
		return tgbotapi.NewCallback(update.CallbackQuery.ID, "")

	case data == "admin_showusers":
		users, err := r.storage.GetAllUsers()
		if err != nil {
			log.Printf("Error getting users: %v", err)
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "❌ خطا در دریافت کاربران")
			return callback
		}

		// ارسال لیست کاربران
		userList := "👥 *لیست تمام کاربران:*\n\n"
		for i, user := range users {
			userList += fmt.Sprintf("%d. ID: %d | نام: %s\n", i+1, user.UserID, user.Name)
		}

		msg := tgbotapi.NewMessage(chatID, userList)
		msg.ParseMode = tgbotapi.ModeMarkdown
		r.bot.Send(msg)

	case data == "admin_showgroups":
		groups, err := r.storage.GetAllGroups()
		if err != nil {
			log.Printf("Error getting groups: %v", err)
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "❌ خطا در دریافت گروه‌ها")
			return callback
		}

		// ارسال لیست گروه‌ها
		groupList := "🏢 *لیست تمام گروه‌ها:*\n\n"
		for i, group := range groups {
			groupList += fmt.Sprintf("%d. ID: %d | نام: %s\n", i+1, group.GroupID, group.GroupName)
		}

		msg := tgbotapi.NewMessage(chatID, groupList)
		msg.ParseMode = tgbotapi.ModeMarkdown
		r.bot.Send(msg)
	}

	return tgbotapi.NewCallback(update.CallbackQuery.ID, "✅")
}

// HandleShowUsers command
func (r *AdminCommand) HandleShowUsers(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	// فقط در چت خصوصی کار می‌کند
	if update.Message.Chat.Type != "private" {
		return tgbotapi.NewMessage(chatID, "❌ این دستور فقط در چت خصوصی با بات قابل استفاده است.")
	}

	// بررسی دسترسی ادمین
	if !r.IsAdmin(userID) {
		return tgbotapi.NewMessage(chatID, "❌ شما دسترسی ادمین ندارید.")
	}

	users, err := r.storage.GetAllUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در دریافت لیست کاربران.")
	}

	if len(users) == 0 {
		return tgbotapi.NewMessage(chatID, "📭 هیچ کاربری یافت نشد.")
	}

	userList := "👥 *لیست تمام کاربران:*\n\n"
	for i, user := range users {
		userList += fmt.Sprintf("%d. ID: %d | نام: %s\n", i+1, user.UserID, user.Name)
	}

	msg := tgbotapi.NewMessage(chatID, userList)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

// HandleShowGroups command
func (r *AdminCommand) HandleShowGroups(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	// فقط در چت خصوصی کار می‌کند
	if update.Message.Chat.Type != "private" {
		return tgbotapi.NewMessage(chatID, "❌ این دستور فقط در چت خصوصی با بات قابل استفاده است.")
	}

	// بررسی دسترسی ادمین
	if !r.IsAdmin(userID) {
		return tgbotapi.NewMessage(chatID, "❌ شما دسترسی ادمین ندارید.")
	}

	groups, err := r.storage.GetAllGroups()
	if err != nil {
		log.Printf("Error getting groups: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در دریافت لیست گروه‌ها.")
	}

	if len(groups) == 0 {
		return tgbotapi.NewMessage(chatID, "📭 هیچ گروهی یافت نشد.")
	}

	groupList := "🏢 *لیست تمام گروه‌ها:*\n\n"
	for i, group := range groups {
		groupList += fmt.Sprintf("%d. ID: %d | نام: %s\n", i+1, group.GroupID, group.GroupName)
	}

	msg := tgbotapi.NewMessage(chatID, groupList)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

// HasPendingAdd آیا ادمین در حالت افزودن لینک است؟
func (r *AdminCommand) HasPendingAdd(userID int64) bool {
	return r.pendingAdd[userID]
}

// HandlePrivateTextInput پردازش متن ارسالی خصوصی وقتی حالت افزودن فعال است
func (r *AdminCommand) HandlePrivateTextInput(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	text := strings.TrimSpace(update.Message.Text)

	if !r.IsAdmin(userID) {
		return tgbotapi.NewMessage(chatID, "❌ شما دسترسی ادمین ندارید.")
	}
	if !r.HasPendingAdd(userID) {
		return tgbotapi.MessageConfig{}
	}

	// پارس ورودی: "<link> | <title?>"
	link := text
	title := ""
	if idx := strings.Index(text, "|"); idx >= 0 {
		link = strings.TrimSpace(text[:idx])
		title = strings.TrimSpace(text[idx+1:])
	}
	if link == "" {
		return tgbotapi.NewMessage(chatID, "❌ لینک نامعتبر است")
	}

	// استخراج یوزرنیم کانال از لینک t.me
	channelUsername := ""
	if strings.Contains(link, "t.me/") {
		// نمونه‌ها: https://t.me/username or https://t.me/+code
		parts := strings.Split(link, "/")
		if len(parts) > 0 {
			last := parts[len(parts)-1]
			if last != "" {
				if strings.HasPrefix(last, "+") {
					// لینک خصوصی؛ یوزرنیم ندارد
					channelUsername = ""
				} else {
					channelUsername = strings.TrimPrefix(last, "@")
				}
			}
		}
	}

	if err := r.storage.AddRequiredChannel(0, title, link, channelUsername, 0); err != nil {
		log.Printf("Error adding required channel: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در ذخیره لینک")
	}

	delete(r.pendingAdd, userID)
	return tgbotapi.NewMessage(chatID, "✅ لینک اضافه شد")
}
