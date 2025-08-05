package commands

import (
	"fmt"
	"log"
	"math/rand"
	"redhat-bot/storage"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CrushCommand struct {
	storage *storage.MySQLStorage
	bot     *tgbotapi.BotAPI
}

func NewCrushCommand(storage *storage.MySQLStorage, bot *tgbotapi.BotAPI) *CrushCommand {
	return &CrushCommand{
		storage: storage,
		bot:     bot,
	}
}

func (r *CrushCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// بررسی دستور فعال‌سازی
	if text == "/کراشفعال" {
		if err := r.storage.SetCrushEnabled(chatID, true); err != nil {
			log.Printf("Error enabling crush: %v", err)
			return tgbotapi.NewMessage(chatID, "❌ خطا در فعال‌سازی قابلیت کراش")
		}
		msg := tgbotapi.NewMessage(chatID, "💘 قابلیت کراش فعال شد! ✅\n\nهر 15 ساعت یک جفت کراش جدید اعلام می‌شود! 😂")
		return msg
	}

	// بررسی دستور غیرفعال‌سازی
	if text == "/کراشغیرفعال" {
		if err := r.storage.SetCrushEnabled(chatID, false); err != nil {
			log.Printf("Error disabling crush: %v", err)
			return tgbotapi.NewMessage(chatID, "❌ خطا در غیرفعال‌سازی قابلیت کراش")
		}
		msg := tgbotapi.NewMessage(chatID, "💘 قابلیت کراش غیرفعال شد! ❌\n\nدیگر اعلام کراش نخواهد شد.")
		return msg
	}

	// بررسی دستور وضعیت
	if text == "/کراشوضعیت" {
		isEnabled, err := r.storage.IsCrushEnabled(chatID)
		if err != nil {
			log.Printf("Error checking crush status: %v", err)
			return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی وضعیت کراش")
		}

		var status string
		if isEnabled {
			status = "فعال ✅"
		} else {
			status = "غیرفعال ❌"
		}

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("💘 *وضعیت قابلیت کراش:*\n\nوضعیت فعلی: %s\n\nدستورات:\n`/کراشفعال` - فعال‌سازی\n`/کراشغیرفعال` - غیرفعال‌سازی\n`/کراشدستی` - اعلام کراش دستی", status))
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	}

	// بررسی دستور کراش دستی
	if text == "/کراشدستی" {
		go r.announceRandomCrush(chatID)
		msg := tgbotapi.NewMessage(chatID, "💘 در حال انتخاب جفت کراش... 😂")
		return msg
	}

	// اگر دستور نامعتبر بود
	msg := tgbotapi.NewMessage(chatID, "💘 *دستورات کراش:*\n\n`/کراشفعال` - فعال‌سازی قابلیت\n`/کراشغیرفعال` - غیرفعال‌سازی\n`/کراشوضعیت` - نمایش وضعیت\n`/کراشدستی` - اعلام کراش دستی")
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

// تابع اعلام کراش تصادفی
func (r *CrushCommand) announceRandomCrush(chatID int64) {
	// دریافت لیست کاربران گروه
	users, err := r.storage.GetGroupMembers(chatID)
	if err != nil {
		log.Printf("Error getting group members: %v", err)
		msg := tgbotapi.NewMessage(chatID, "❌ خطا در دریافت لیست اعضای گروه")
		r.bot.Send(msg)
		return
	}

	if len(users) < 2 {
		msg := tgbotapi.NewMessage(chatID, "💘 تعداد اعضای گروه برای اعلام کراش کافی نیست! 😅")
		r.bot.Send(msg)
		return
	}

	// انتخاب دو کاربر تصادفی
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(users))
	j := rand.Intn(len(users))
	for i == j {
		j = rand.Intn(len(users))
	}

	user1 := users[i]
	user2 := users[j]

	// ساخت پیام کراش
	crushMessages := []string{
		fmt.Sprintf("💘 امروز %s کراش %s هست! 😂", user1.Name, user2.Name),
		fmt.Sprintf("💕 %s عاشق %s شده! 🥰", user1.Name, user2.Name),
		fmt.Sprintf("🔥 %s و %s عاشق همدیگه شدن! 💕", user1.Name, user2.Name),
		fmt.Sprintf("💘 %s دلش برای %s می‌تپه! 💓", user1.Name, user2.Name),
		fmt.Sprintf("🥰 %s عاشقانه به %s نگاه می‌کنه! 😍", user1.Name, user2.Name),
	}

	selectedMessage := crushMessages[rand.Intn(len(crushMessages))]

	// ارسال پیام
	msg := tgbotapi.NewMessage(chatID, selectedMessage)
	_, err = r.bot.Send(msg)
	if err != nil {
		log.Printf("خطا در ارسال پیام کراش: %v", err)
	}
}

// تابع شروع کرون جاب برای اعلام خودکار کراش
func (r *CrushCommand) StartCrushScheduler() {
	go func() {
		for {
			time.Sleep(15 * time.Hour) // هر 15 ساعت

			// دریافت تمام گروه‌هایی که قابلیت کراش فعال دارند
			enabledGroups, err := r.storage.GetCrushEnabledGroups()
			if err != nil {
				log.Printf("Error getting crush enabled groups: %v", err)
				continue
			}

			for _, groupID := range enabledGroups {
				go r.announceRandomCrush(groupID)
				time.Sleep(5 * time.Second) // فاصله بین اعلام‌ها
			}
		}
	}()
}
