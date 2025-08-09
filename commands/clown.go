package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"redhat-bot/limiter"
	"redhat-bot/storage"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ClownCommand struct {
	rateLimiter *limiter.RateLimiter
	bot         *tgbotapi.BotAPI
	storage     *storage.MySQLStorage
	insults     []string
}

func NewClownCommand(storage *storage.MySQLStorage, rateLimiter *limiter.RateLimiter, bot *tgbotapi.BotAPI) *ClownCommand {
	c := &ClownCommand{
		rateLimiter: rateLimiter,
		bot:         bot,
		storage:     storage,
		insults:     nil,
	}

	// تلاش برای بارگذاری لیست فحش‌ها از فایل
	c.loadInsultsFromFile("jsonfile/clown.json")
	return c
}

func (r *ClownCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	// بررسی فعال بودن قابلیت دلقک در این گروه
	enabled, err := r.storage.IsClownEnabled(chatID)
	if err != nil {
		return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی وضعیت دلقک")
	}
	if !enabled {
		return tgbotapi.NewMessage(chatID, "🔒 قابلیت دلقک در این گروه غیرفعال است. از منوی پنل → قفل‌ها آن را فعال کنید.")
	}

	// بررسی محدودیت درخواست
	if allowed, message := r.rateLimiter.CheckRateLimit(userID); !allowed {
		return tgbotapi.NewMessage(chatID, message)
	}

	// استخراج نام مخاطب از دستور (پشتیبانی از «دلقک» و «/clown»)
	text := update.Message.Text
	cleaned := strings.TrimSpace(strings.TrimPrefix(text, "/clown"))
	cleaned = strings.TrimSpace(strings.TrimPrefix(cleaned, "دلقک"))
	targetName := cleaned

	if targetName == "" {
		msg := tgbotapi.NewMessage(chatID, "🤡 *دستور دلقک*\n\nنحوه استفاده: دلقک <نام مخاطب>\n\nمثال: دلقک علی یا دلقک @username")
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	}

	// انتخاب تصادفی یک فحش از لیست
	insult := r.randomInsult()
	if insult == "" {
		insult = "فعلاً چیزی برای گفتن ندارم!"
	}

	formattedResponse := fmt.Sprintf("🤡 *دلقک به %s:*\n\n%s", targetName, insult)
	msg := tgbotapi.NewMessage(chatID, formattedResponse)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

// loadInsultsFromFile بارگذاری فهرست فحش‌ها از فایل JSON
func (r *ClownCommand) loadInsultsFromFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("خطا در خواندن فایل فحش‌ها: %v", err)
		return
	}
	var items []struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(data, &items); err != nil {
		log.Printf("خطا در parse فایل فحش‌ها: %v", err)
		return
	}
	insults := make([]string, 0, len(items))
	for _, it := range items {
		text := strings.TrimSpace(it.Description)
		if text != "" {
			insults = append(insults, text)
		}
	}
	r.insults = insults
}

// randomInsult انتخاب تصادفی از فهرست
func (r *ClownCommand) randomInsult() string {
	if len(r.insults) == 0 {
		// تلاش دوباره برای بارگذاری در صورت خالی بودن
		r.loadInsultsFromFile("jsonfile/clown.json")
	}
	if len(r.insults) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return r.insults[rand.Intn(len(r.insults))]
}
