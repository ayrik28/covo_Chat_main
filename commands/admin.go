package commands

import (
	"fmt"
	"log"

	"redhat-bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminCommand struct {
	bot     *tgbotapi.BotAPI
	storage *storage.MySQLStorage
}

// لیست ادمین‌های مجاز
var adminUsers = map[int64]string{
	7853092812: "مهشید",
	990475046:  "هانتر",
}

func NewAdminCommand(bot *tgbotapi.BotAPI, storage *storage.MySQLStorage) *AdminCommand {
	return &AdminCommand{
		bot:     bot,
		storage: storage,
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
	case 7853092812: // مهشید
		return fmt.Sprintf(`🌟 *سلام %s عزیز!* 🌟


💝 مرسی که یا ایده هات منو به اینجا رسوندی و این همه قابلیت جالب اضافه کردی!

🛠️ *دستورات ادمین:*
• /showusers - نمایش لیست تمام کاربران
• /showgroups - نمایش لیست تمام گروه‌ها
• /admin - بازگشت به منوی ادمین

✨ از اینکه منو ساختی ممنونم! 💖`, name)

	case 990475046: // هانتر
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

	switch update.CallbackQuery.Data {
	case "admin_showusers":
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

	case "admin_showgroups":
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
