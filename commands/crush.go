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
	if text == "/crushon" {
		if err := r.storage.SetCrushEnabled(chatID, true); err != nil {
			log.Printf("Error enabling crush: %v", err)
			return tgbotapi.NewMessage(chatID, "❌ خطا در فعال‌سازی قابلیت کراش")
		}
		msg := tgbotapi.NewMessage(chatID, "💘 *قابلیت کراش با موفقیت فعال شد!* ✅\n\n🔥 از این لحظه هر 15 ساعت یک بار، دو نفر از اعضای گروه به صورت تصادفی به عنوان کراش انتخاب می‌شوند!\n\n👀 منتظر اعلام اولین جفت کراش باشید...")
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	}

	// بررسی دستور غیرفعال‌سازی
	if text == "/crushoff" {
		if err := r.storage.SetCrushEnabled(chatID, false); err != nil {
			log.Printf("Error disabling crush: %v", err)
			return tgbotapi.NewMessage(chatID, "❌ خطا در غیرفعال‌سازی قابلیت کراش")
		}
		msg := tgbotapi.NewMessage(chatID, "💔 *قابلیت کراش غیرفعال شد!* ❌\n\n🚫 دیگر اعلام خودکار کراش در این گروه انجام نخواهد شد.\n\n✅ برای فعال‌سازی مجدد از دستور `/crushon` استفاده کنید.")
		msg.ParseMode = tgbotapi.ModeMarkdown
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

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("💘 *وضعیت قابلیت کراش:*\n\nوضعیت فعلی: %s\n\nدستورات:\n`/crushon` - فعال‌سازی\n`/crushoff` - غیرفعال‌سازی", status))
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	}

	// دستور کراش دستی حذف شد

	// اگر دستور نامعتبر بود
	msg := tgbotapi.NewMessage(chatID, "💘 *دستورات کراش:*\n\n`/crushon` - فعال‌سازی قابلیت\n`/crushoff` - غیرفعال‌سازی\n`/کراشوضعیت` - نمایش وضعیت")
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

// تابع اعلام کراش تصادفی
func (r *CrushCommand) announceRandomCrush(chatID int64) {
	// دریافت لیست کاربران گروه مستقیماً از دیتابیس
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

	// انتخاب دو کاربر تصادفی با استفاده از rand.Shuffle برای انتخاب تصادفی بهتر
	rand.Seed(time.Now().UnixNano())

	// برای اطمینان از انتخاب تصادفی واقعی، لیست را شافل می‌کنیم
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})

	// انتخاب دو کاربر اول از لیست شافل شده
	user1 := users[0]
	user2 := users[1]

	// ساخت پیام کراش
	crushMessages := []string{
		fmt.Sprintf("💘 *کراش امروز:* %s روی %s کراش زده! 😍", user1.Name, user2.Name),
		fmt.Sprintf("💕 *خبر داغ:* %s عاشق %s شده! 🥰", user1.Name, user2.Name),
		fmt.Sprintf("🔥 *اعلام رسمی:* %s و %s عاشق همدیگه شدن! 💕", user1.Name, user2.Name),
		fmt.Sprintf("💘 *قلب‌ها به تپش افتادند:* %s دلش برای %s می‌تپه! 💓", user1.Name, user2.Name),
		fmt.Sprintf("🥰 *نگاه‌های عاشقانه:* %s مخفیانه به %s نگاه می‌کنه! 😍", user1.Name, user2.Name),
		fmt.Sprintf("💐 *عشق در هوا موج می‌زند:* %s می‌خواهد به %s نزدیک شود! 💝", user1.Name, user2.Name),
		fmt.Sprintf("💌 *پیام عاشقانه:* %s برای %s نامه‌ای پر از عشق نوشته! 📝", user1.Name, user2.Name),
		fmt.Sprintf("🌹 *گل سرخ عشق:* %s برای %s گل فرستاده! 🌷", user1.Name, user2.Name),
		fmt.Sprintf("🍫 *شیرینی عشق:* %s برای %s شکلات خریده! 🍬", user1.Name, user2.Name),
		fmt.Sprintf("🎵 *آهنگ عاشقانه:* %s برای %s آهنگ می‌خواند! 🎤", user1.Name, user2.Name),
	}

	selectedMessage := crushMessages[rand.Intn(len(crushMessages))]

	// ارسال پیام
	msg := tgbotapi.NewMessage(chatID, selectedMessage)
	msg.ParseMode = tgbotapi.ModeMarkdown
	_, err = r.bot.Send(msg)
	if err != nil {
		log.Printf("خطا در ارسال پیام کراش: %v", err)
	}
}

// تابع شروع کرون جاب برای اعلام خودکار کراش
func (r *CrushCommand) StartCrushScheduler() {
	go func() {
		for {
			time.Sleep(15 * time.Hour) // هر 15 ساعت یکبار

			// دریافت تمام گروه‌هایی که قابلیت کراش فعال دارند
			enabledGroups, err := r.storage.GetCrushEnabledGroups()
			if err != nil {
				log.Printf("Error getting crush enabled groups: %v", err)
				continue
			}

			log.Printf("Sending crush announcements to %d enabled groups", len(enabledGroups))

			for _, groupID := range enabledGroups {
				go r.announceRandomCrush(groupID)
				time.Sleep(5 * time.Minute) // فاصله کوتاه بین اعلام‌ها
			}
		}
	}()

	log.Println("💘 Crush scheduler started - announcing crushes every 15 hours")
}
