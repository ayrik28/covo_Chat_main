package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"redhat-bot/ai"
	"redhat-bot/commands"
	"redhat-bot/config"
	"redhat-bot/limiter"

	// "redhat-bot/scheduler"
	"redhat-bot/storage"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

type CovoBot struct {
	bot               *tgbotapi.BotAPI
	storage           *storage.MySQLStorage
	rateLimiter       *limiter.RateLimiter
	aiClient          *ai.DeepSeekClient
	covoCommand       *commands.CovoCommand
	covoJokeCommand   *commands.CovoJokeCommand
	musicCommand      *commands.MusicCommand
	crsCommand        *commands.CrsCommand
	clownCommand      *commands.ClownCommand
	crushCommand      *commands.CrushCommand
	gapCommand        *commands.GapCommand
	hafezCommand      *commands.HafezCommand
	adminCommand      *commands.AdminCommand
	moderationCommand *commands.ModerationCommand
	truthDareCommand  *commands.TruthDareCommand
	tagCommand        *commands.TagCommand
	dailyChallenge    *commands.DailyChallengeCommand
	// summaryScheduler *scheduler.DailySummaryScheduler
	cron *cron.Cron
}

// bad words cache
var loadBadWordsOnce sync.Once
var cachedBadWords []string
var cachedFinglishWords []string

type badWordsFile struct {
	FarsiWords    []string `json:"farsiWords"`
	FinglishWords []string `json:"finglishWords"`
}

func loadBadWords() {
	loadBadWordsOnce.Do(func() {
		filePath := filepath.Join("jsonfile", "badwords.json")
		f, err := os.Open(filePath)
		if err != nil {
			log.Printf("cannot open badwords file: %v", err)
			return
		}
		defer f.Close()

		var bw badWordsFile
		if err := json.NewDecoder(f).Decode(&bw); err != nil {
			log.Printf("cannot decode badwords file: %v", err)
			return
		}
		cachedBadWords = bw.FarsiWords
		cachedFinglishWords = bw.FinglishWords
	})
}

func containsBadWord(text string) bool {
	if text == "" {
		return false
	}
	loadBadWords()
	t := strings.ToLower(text)
	for _, w := range cachedBadWords {
		if w == "" {
			continue
		}
		if strings.Contains(t, strings.ToLower(w)) {
			return true
		}
	}
	for _, w := range cachedFinglishWords {
		if w == "" {
			continue
		}
		if strings.Contains(t, strings.ToLower(w)) {
			return true
		}
	}
	return false
}

func NewCovoBot() (*CovoBot, error) {
	// بارگذاری تنظیمات
	config.LoadConfig()

	// راه‌اندازی بات
	bot, err := tgbotapi.NewBotAPI(config.AppConfig.TelegramToken)
	if err != nil {
		return nil, err
	}

	// راه‌اندازی اتصال به دیتابیس
	storage, err := storage.NewMySQLStorage(
		config.AppConfig.MySQLHost,
		config.AppConfig.MySQLPort,
		config.AppConfig.MySQLUser,
		config.AppConfig.MySQLPassword,
		config.AppConfig.MySQLDatabase,
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing MySQL storage: %v", err)
	}

	// راه‌اندازی اجزا
	rateLimiter := limiter.NewRateLimiter(storage)
	aiClient := ai.NewDeepSeekClient()

	// راه‌اندازی دستورات
	covoCommand := commands.NewCovoCommand(aiClient, rateLimiter, bot)
	covoJokeCommand := commands.NewCovoJokeCommand(aiClient, rateLimiter, bot)
	musicCommand := commands.NewMusicCommand(aiClient, rateLimiter, bot)
	crsCommand := commands.NewCrsCommand(rateLimiter)
	clownCommand := commands.NewClownCommand(storage, rateLimiter, bot)
	crushCommand := commands.NewCrushCommand(storage, bot)
	hafezCommand := commands.NewHafezCommand(bot)
	adminCommand := commands.NewAdminCommand(bot, storage)
	gapCommand := commands.NewGapCommand(bot, storage, hafezCommand)
	moderationCommand := commands.NewModerationCommand(bot)
	truthDareCommand := commands.NewTruthDareCommand(bot, adminCommand)
	tagCommand := commands.NewTagCommand(bot, storage)

	// راه‌اندازی زمان‌بند
	// summaryScheduler := scheduler.NewDailySummaryScheduler(bot, storage, aiClient)

	// راه‌اندازی کران با تایم‌زون تهران
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		log.Printf("cannot load Asia/Tehran timezone, falling back to local: %v", err)
		loc = time.Local
	}
	cronJob := cron.New(cron.WithLocation(loc))

	return &CovoBot{
		bot:               bot,
		storage:           storage,
		rateLimiter:       rateLimiter,
		aiClient:          aiClient,
		covoCommand:       covoCommand,
		covoJokeCommand:   covoJokeCommand,
		musicCommand:      musicCommand,
		crsCommand:        crsCommand,
		clownCommand:      clownCommand,
		crushCommand:      crushCommand,
		gapCommand:        gapCommand,
		hafezCommand:      hafezCommand,
		adminCommand:      adminCommand,
		moderationCommand: moderationCommand,
		truthDareCommand:  truthDareCommand,
		tagCommand:        tagCommand,
		dailyChallenge:    commands.NewDailyChallengeCommand(storage, bot),
		// summaryScheduler: summaryScheduler,
		cron: cronJob,
	}, nil
}

func (r *CovoBot) Start() error {
	log.Printf("🤖 بات کوو در حال راه‌اندازی است...")
	log.Printf("👤 نام کاربری بات: @%s", r.bot.Self.UserName)

	// تنظیم کار کران برای خلاصه‌های روزانه (ساعت ۹ صبح هر روز)
	_, err := r.cron.AddFunc("0 9 * * *", func() {
		// r.summaryScheduler.RunDailySummary()
	})
	if err != nil {
		return err
	}

	// کران چلنج روزانه ساعت ۱۰ به وقت ایران (با کرانی که روی Asia/Tehran تنظیم شده)
	if _, err := r.cron.AddFunc("0 10 * * *", func() {
		r.dailyChallenge.RunDailyForEnabledGroups()
	}); err != nil {
		return err
	}

	r.cron.Start()
	log.Println("⏰ زمان‌بندها راه‌اندازی شد (خلاصه ۹:۰۰، چلنج ~۱۰:۳۰ تهران)")

	// راه‌اندازی کراش scheduler
	r.crushCommand.StartCrushScheduler()
	log.Println("💘 کراش scheduler راه‌اندازی شد (هر 10 ساعت)")

	// تنظیم کانال به‌روزرسانی
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := r.bot.GetUpdatesChan(updateConfig)

	// پردازش به‌روزرسانی‌ها
	for update := range updates {
		go r.handleUpdate(update)
	}

	return nil
}

func (r *CovoBot) handleUpdate(update tgbotapi.Update) {
	// Handle my_chat_member updates (bot status in chats/channels changes)
	if update.MyChatMember != nil {
		r.handleMyChatMember(update)
		return
	}
	// Handle callback queries from inline keyboard
	if update.CallbackQuery != nil {
		var callback tgbotapi.CallbackConfig

		// بررسی نوع callback
		// گیت عضویت برای تمام کال‌بک‌ها به جز موارد ادمین و بررسی عضویت
		if !(strings.HasPrefix(update.CallbackQuery.Data, "admin_") || strings.HasPrefix(update.CallbackQuery.Data, "admin_check_join")) {
			if ok, prompt := r.checkRequiredMembershipAndPromptUser(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.ID); !ok {
				if prompt.ChatID != 0 {
					_, _ = r.bot.Send(prompt)
				}
				// ارسال ack کوتاه
				callback = tgbotapi.NewCallback(update.CallbackQuery.ID, "برای استفاده، ابتدا عضو کانال‌ها شوید")
				if _, err := r.bot.Request(callback); err != nil {
					log.Printf("Error handling callback: %v", err)
				}
				return
			}
		}
		switch {
		case strings.HasPrefix(update.CallbackQuery.Data, "admin_check_join"):
			// دکمه «بررسی عضویت» از پیام عضویت اجباری
			// تلاش مجدد برای بررسی
			dummy := tgbotapi.NewCallback(update.CallbackQuery.ID, "در حال بررسی...")
			if _, err := r.bot.Request(dummy); err != nil {
				log.Printf("callback ack error: %v", err)
			}
			// در PM ممکن است کاربر بخواهد مستقیم شروع کند
			if ok, prompt := r.checkRequiredMembershipAndPromptUser(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.ID); ok {
				// تایید عضویت
				notice := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "✅ عضویت شما تایید شد. حالا می‌توانید از دستورات استفاده کنید.")
				_, _ = r.bot.Send(notice)
				return
			} else {
				if prompt.ChatID != 0 {
					_, _ = r.bot.Send(prompt)
				}
				return
			}
		case strings.HasPrefix(update.CallbackQuery.Data, "admin_"):
			callback = r.adminCommand.HandleCallback(update)
		case strings.HasPrefix(update.CallbackQuery.Data, "td_"):
			// گیت عضویت برای کلیک‌های بازی
			if ok, prompt := r.checkRequiredMembershipAndPromptUser(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.ID); !ok {
				if prompt.ChatID != 0 {
					_, _ = r.bot.Send(prompt)
				}
				// ارسال ack کوتاه
				callback = tgbotapi.NewCallback(update.CallbackQuery.ID, "برای استفاده، ابتدا عضو کانال‌ها شوید")
				break
			}
			callback = r.truthDareCommand.HandleCallback(update)
		default:
			callback = r.gapCommand.HandleCallback(update)
		}

		if _, err := r.bot.Request(callback); err != nil {
			log.Printf("Error handling callback: %v", err)
		}
		return
	}

	if update.Message == nil {
		return
	}

	message := update.Message
	text := message.Text

	// بررسی عضو شدن بات در گروه جدید
	if message.NewChatMembers != nil {
		for _, user := range message.NewChatMembers {
			if user.UserName == r.bot.Self.UserName {
				// بات به گروه جدید اضافه شده
				if err := r.storage.SetFeatureEnabled(message.Chat.ID, "crush", false); err != nil {
					log.Printf("Error initializing crush feature: %v", err)
				}
				if err := r.storage.SetFeatureEnabled(message.Chat.ID, "clown", false); err != nil {
					log.Printf("Error initializing clown feature: %v", err)
				}
				if err := r.storage.SetFeatureEnabled(message.Chat.ID, "hafez", false); err != nil {
					log.Printf("Error initializing hafez feature: %v", err)
				}
				if err := r.storage.SetFeatureEnabled(message.Chat.ID, "badword", false); err != nil {
					log.Printf("Error initializing badword feature: %v", err)
				}
				welcomeMsg := tgbotapi.NewMessage(message.Chat.ID, `🤖 *سلام! من بات covo هستم!*

من دستیار هوشمند شما با قابلیت‌های جالب هستم:

💡 *دستورات:*
• /covo <سوال> - هر سوالی دارید بپرسید!
• /cj <موضوع> - جوک خنده‌دار درباره هر موضوعی تولید کن
• /music - پیشنهاد موسیقی بر اساس سلیقه شما
• دلقک <نام> - توهین به شخص مورد نظر
• /crushon - فعال‌سازی قابلیت کراش
• /فال - دریافت فال حافظ
• /crs - بررسی وضعیت بات
• پنل - نمایش دستورات مخصوص گروه
• /covog - نمایش راهنما
• /help - نمایش راهنما

🎯 *ویژگی‌ها:*
• درخواست‌های نامحدود
• خلاصه‌های هوشمند روزانه گروه‌ها
• تولید جوک بر اساس موضوع
• پیشنهاد موسیقی هوشمند
• قابلیت دلقک برای توهین هوشمند
• قابلیت کراش خودکار هر 15 ساعت

بیایید شروع کنیم! با /covo <سوال شما> چیزی از من بپرسید 🚀`)
				welcomeMsg.ParseMode = tgbotapi.ModeMarkdown
				_, err := r.bot.Send(welcomeMsg)
				if err != nil {
					log.Printf("خطا در ارسال پیام خوش‌آمدگویی: %v", err)
				}
				log.Printf("بات به گروه جدید اضافه شد: %s (ID: %d)", message.Chat.Title, message.Chat.ID)
				return
			}
		}
	}

	// بررسی خروج بات از گروه
	if message.LeftChatMember != nil && message.LeftChatMember.UserName == r.bot.Self.UserName {
		log.Printf("بات از گروه خارج شد: %s (ID: %d)", message.Chat.Title, message.Chat.ID)
		return
	}

	// پردازش پیام‌های گروه (برای ثبت)
	if message.Chat.Type == "group" || message.Chat.Type == "supergroup" {
		// افزودن گروه به زمان‌بند اگر قبلاً اضافه نشده
		// r.summaryScheduler.AddGroup(message.Chat.ID, message.Chat.Title)

		// ثبت پیام برای خلاصه روزانه
		username := message.From.UserName
		if username == "" {
			username = message.From.FirstName
		}
		if err := r.storage.AddGroupMessage(message.Chat.ID, message.From.ID, username, text); err != nil {
			log.Printf("Error adding group message: %v", err)
		}

		// اضافه کردن کاربر به لیست اعضای گروه (برای قابلیت کراش)
		userName := message.From.FirstName
		if message.From.UserName != "" {
			userName = "@" + message.From.UserName
		}
		if err := r.storage.AddGroupMember(message.Chat.ID, message.From.ID, userName); err != nil {
			log.Printf("Error adding group member: %v", err)
		}
		// اگر قفل لینک فعال است، پیام‌های حاوی لینک حذف شوند
		if enabled, err := r.storage.IsFeatureEnabled(message.Chat.ID, "link"); err == nil && enabled {
			if containsLink(text) {
				_, _ = r.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: message.Chat.ID, MessageID: message.MessageID})
				return
			}
		}

		// اگر قفل فحش فعال است، پیام‌های حاوی کلمات بد حذف شوند
		if enabled, err := r.storage.IsFeatureEnabled(message.Chat.ID, "badword"); err == nil && enabled {
			if containsBadWord(text) {
				_, _ = r.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: message.Chat.ID, MessageID: message.MessageID})
				return
			}
		}
	}

	// اگر پیام خصوصی از ادمین و در حالت افزودن لینک بود، قبل از هرچیز آن را هندل کن
	if message.Chat.Type == "private" && r.adminCommand.IsAdmin(message.From.ID) && r.adminCommand.HasPendingAdd(message.From.ID) {
		resp := r.adminCommand.HandlePrivateTextInput(update)
		if resp.ChatID != 0 {
			_, _ = r.bot.Send(resp)
		}
		return
	}

	// اگر پاسخ به چلنج روزانه است، اول رسیدگی شود
	if resp := r.dailyChallenge.HandleAnswer(update); resp.ChatID != 0 {
		if _, err := r.bot.Send(resp); err != nil {
			log.Printf("خطا در اعلام برنده چلنج: %v", err)
		}
		return
	}

	// پردازش دستورات
	if !strings.HasPrefix(text, "/") {
		// اگر یکی از تریگرهای اکشن بود، ابتدا گیت عضویت را بررسی کن
		trimmed := strings.TrimSpace(text)
		if trimmed == "پنل" || trimmed == "بازی" || trimmed == "توقف بازی" || trimmed == "کراش" || trimmed == "فال" || trimmed == "تگ" || strings.HasPrefix(trimmed, "دلقک") || strings.HasPrefix(trimmed, "سکوت") || trimmed == "ازاد" || strings.HasPrefix(trimmed, "حذف") {
			if ok, prompt := r.checkRequiredMembershipAndPromptUser(message.Chat.ID, message.From.ID); !ok {
				if prompt.ChatID != 0 {
					_, _ = r.bot.Send(prompt)
				}
				return
			}
		}
		// «بازی» بدون اسلش -> شروع روم (فقط ادمین)
		if strings.TrimSpace(text) == "بازی" {
			response := r.truthDareCommand.HandleStartWithoutSlash(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام بازی: %v", err)
				}
			}
			return
		}

		// «توقف بازی» بدون اسلش -> توقف کامل (فقط ادمین)
		if strings.TrimSpace(text) == "توقف بازی" {
			response := r.truthDareCommand.HandleStopWithoutSlash(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام توقف بازی: %v", err)
				}
			}
			return
		}
		// «کراش» بدون اسلش -> نمایش وضعیت
		if strings.TrimSpace(text) == "کراش" {
			status := r.crushCommand.BuildStatusMessage(message.Chat.ID)
			if status.ChatID != 0 {
				_, err := r.bot.Send(status)
				if err != nil {
					log.Printf("خطا در ارسال وضعیت کراش: %v", err)
				}
			}
			return
		}
		// «فال» بدون اسلش (در صورت فعال بودن قابلیت)
		if strings.TrimSpace(text) == "فال" {
			if enabled, err := r.storage.IsFeatureEnabled(message.Chat.ID, "hafez"); err == nil && enabled {
				response := r.hafezCommand.Handle(update)
				if response.ChatID != 0 {
					_, err := r.bot.Send(response)
					if err != nil {
						log.Printf("خطا در ارسال پیام فال: %v", err)
					}
				}
			} else {
				notice := tgbotapi.NewMessage(message.Chat.ID, "❌ قابلیت فال در این گروه غیرفعال است")
				_, _ = r.bot.Send(notice)
			}
			return
		}
		// پشتیبانی از «بن» روی ریپلای بدون اسلش
		if strings.TrimSpace(text) == "بن" {
			response := r.moderationCommand.HandleBanOnReply(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام: %v", err)
				}
			}
			return
		}

		// پشتیبانی از «سکوت [n]» روی ریپلای بدون اسلش (n = ساعت)
		if strings.HasPrefix(strings.TrimSpace(text), "سکوت") {
			response := r.moderationCommand.HandleMute(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام: %v", err)
				}
			}
			return
		}

		// پشتیبانی از «آزاد» روی ریپلای بدون اسلش
		if strings.TrimSpace(text) == "ازاد" {
			response := r.moderationCommand.HandleUnmute(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام: %v", err)
				}
			}
			return
		}

		// پشتیبانی از «حذف [n]» بدون اسلش
		if strings.HasPrefix(text, "حذف") {
			response := r.moderationCommand.Handle(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام: %v", err)
				}
			}
			return
		}

		// پشتیبانی از «پنل» بدون اسلش
		if strings.TrimSpace(text) == "پنل" {
			// گیت عضویت اجباری پیش از نمایش پنل
			if ok, prompt := r.checkRequiredMembershipAndPromptUser(message.Chat.ID, message.From.ID); !ok {
				if prompt.ChatID != 0 {
					_, _ = r.bot.Send(prompt)
				}
				return
			}
			response := r.gapCommand.Handle(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام: %v", err)
				}
			}
			return
		}

		// «تگ» بدون اسلش روی ریپلای -> تگ همه اعضا (فقط ادمین)
		if strings.TrimSpace(text) == "تگ" {
			response := r.tagCommand.HandleTagAllOnReply(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام: %v", err)
				}
			}
			return
		}

		// پشتیبانی از «دلقک <نام>» بدون اسلش
		if strings.HasPrefix(strings.TrimSpace(text), "دلقک") {
			response := r.clownCommand.Handle(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام: %v", err)
				}
			}
			return
		}

		// پشتیبانی از «دلقک <نام>» بدون اسلش
		if strings.HasPrefix(strings.TrimSpace(text), "دلقک") {
			response := r.clownCommand.Handle(update)
			if response.ChatID != 0 {
				_, err := r.bot.Send(response)
				if err != nil {
					log.Printf("خطا در ارسال پیام: %v", err)
				}
			}
			return
		}

		// بررسی ریپلای به دستور موسیقی
		if message.ReplyToMessage != nil && message.ReplyToMessage.Text != "" {
			replyText := message.ReplyToMessage.Text
			if strings.Contains(replyText, "پیشنهاد موسیقی") || strings.Contains(replyText, "چه نوع آهنگی") {
				response := r.musicCommand.Handle(update)
				if response.ChatID != 0 {
					_, err := r.bot.Send(response)
					if err != nil {
						log.Printf("خطا در ارسال پاسخ موسیقی: %v", err)
					}
				}
				return
			}
		}
		return
	}

	var response tgbotapi.MessageConfig

	// برای همه دستورات اسلش‌دار گیت عضویت
	if strings.HasPrefix(text, "/") {
		if ok, prompt := r.checkRequiredMembershipAndPromptUser(message.Chat.ID, message.From.ID); !ok {
			if prompt.ChatID != 0 {
				_, _ = r.bot.Send(prompt)
			}
			return
		}
	}

	switch {
	case strings.HasPrefix(text, "/covo"):
		if ok, prompt := r.checkRequiredMembershipAndPromptUser(message.Chat.ID, message.From.ID); !ok {
			if prompt.ChatID != 0 {
				_, _ = r.bot.Send(prompt)
			}
			return
		}
		response = r.covoCommand.Handle(update)
	case strings.HasPrefix(text, "/cj"):
		response = r.covoJokeCommand.Handle(update)
	case strings.HasPrefix(text, "/music"):
		response = r.musicCommand.Handle(update)
	case strings.HasPrefix(text, "/crs"):
		response = r.crsCommand.Handle(update)
	case strings.HasPrefix(text, "/clown"):
		response = r.clownCommand.Handle(update)
	case strings.HasPrefix(text, "/crushon"), strings.HasPrefix(text, "/crushoff"), strings.HasPrefix(text, "/کراشوضعیت"):
		response = r.crushCommand.Handle(update)
	case strings.HasPrefix(text, "/start"):
		response = r.handleStartCommand(update)
	case strings.HasPrefix(text, "/covog"):
		response = r.handleStartCommand(update)
	case strings.HasPrefix(text, "/help"):
		response = r.handleHelpCommand(update)
		// پشتیبانی از /gap حذف شد؛ از «پنل» بدون اسلش استفاده کنید
	case strings.HasPrefix(text, "/فال"):
		if ok, prompt := r.checkRequiredMembershipAndPromptUser(message.Chat.ID, message.From.ID); !ok {
			if prompt.ChatID != 0 {
				_, _ = r.bot.Send(prompt)
			}
			return
		}
		response = r.hafezCommand.Handle(update)
	case strings.HasPrefix(text, "/admin"):
		response = r.adminCommand.Handle(update)
	case strings.HasPrefix(text, "/showusers"):
		response = r.adminCommand.HandleShowUsers(update)
	case strings.HasPrefix(text, "/showgroups"):
		response = r.adminCommand.HandleShowGroups(update)
	case strings.HasPrefix(text, "/del"):
		response = r.moderationCommand.Handle(update)
	default:
		return // نادیده گرفتن دستورات ناشناخته
	}

	// ارسال پاسخ (اگر پیام خالی نباشد)
	if response.ChatID != 0 {
		_, err := r.bot.Send(response)
		if err != nil {
			log.Printf("خطا در ارسال پیام: %v", err)
		}
	}
}

// containsLink بررسی وجود لینک در متن پیام
func containsLink(text string) bool {
	t := strings.ToLower(text)
	if strings.Contains(t, "http://") || strings.Contains(t, "https://") {
		return true
	}
	if strings.Contains(t, "t.me/") || strings.Contains(t, "telegram.me/") {
		return true
	}
	// تشخیص ساده دامنه‌ها مانند example.com
	if strings.Contains(t, ".com") || strings.Contains(t, ".ir") || strings.Contains(t, ".org") || strings.Contains(t, ".net") {
		return true
	}
	return false
}

// handleMyChatMember ثبت/آپدیت اطلاعات چت/کانالی که وضعیت ربات در آن تغییر کرده است
func (r *CovoBot) handleMyChatMember(update tgbotapi.Update) {
	m := update.MyChatMember
	if m == nil {
		return
	}
	chat := m.Chat
	status := strings.ToLower(m.NewChatMember.Status)
	isAdmin := status == "administrator" || status == "creator"
	title := chat.Title
	username := chat.UserName

	// تلاش برای ذخیره/آپدیت رکورد کانال/گروه
	if err := r.storage.UpsertBotChannel(chat.ID, title, username, isAdmin, 0); err != nil {
		log.Printf("UpsertBotChannel error: %v", err)
	}
}

// checkRequiredMembershipAndPrompt بررسی می‌کند کاربر عضو همه کانال‌های لازم است یا خیر
// اگر عضو نبود، پیام راهنما با دکمه‌های جوین و دکمه «بررسی عضویت» ارسال می‌کند
func (r *CovoBot) checkRequiredMembershipAndPromptUser(chatID int64, userID int64) (bool, tgbotapi.MessageConfig) {

	// لینک‌ها را از scope سراسری می‌خوانیم (GroupID=0)
	channels, err := r.storage.ListRequiredChannels(0)
	if err != nil || len(channels) == 0 {
		// چیزی برای بررسی نیست
		return true, tgbotapi.MessageConfig{}
	}

	// بررسی عضویت برای هر کانال
	notJoined := 0
	for _, ch := range channels {
		// اولویت با ChannelID (گروه/کانال عمومی)؛ اگر نبود، سعی با Username
		var targetChatID int64
		if ch.ChannelID != 0 {
			targetChatID = ch.ChannelID
		} else if ch.ChannelUsername != "" {
			// در GetChatMember باید chatID کانال عددی یا @username باشد؛ کتابخانه فقط int64 می‌گیرد
			// پس این مورد را نمی‌توان مستقیم چک کرد؛ از این‌رو در چنین حالتی صرفاً نمایش لینک می‌دهیم
			continue
		} else {
			continue
		}

		cfg := tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{ChatID: targetChatID, UserID: userID}}
		member, err := r.bot.GetChatMember(cfg)
		if err != nil {
			notJoined++
			continue
		}
		// Joined if member/admin/creator/restricted
		status := strings.ToLower(member.Status)
		if !(status == "member" || status == "administrator" || status == "creator" || status == "restricted") {
			notJoined++
		}
	}

	if notJoined == 0 {
		return true, tgbotapi.MessageConfig{}
	}

	// ساخت پیام و کیبورد جوین
	text := "برای استفاده از ربات، لطفاً در کانال‌های زیر عضو شوید:"
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, ch := range channels {
		link := ch.Link
		if link == "" && ch.ChannelUsername != "" {
			link = "https://t.me/" + ch.ChannelUsername
		}
		if link == "" {
			continue
		}
		title := ch.Title
		if title == "" {
			title = "Join"
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(title, link),
		))
	}
	// دکمه بررسی عضویت
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ عضو شدم، بررسی کن", "admin_check_join"),
	))
	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = kb
	return false, msg
}

// buildJoinPromptWithoutCheck فقط بر اساس لینک‌های ثبت شده، پیام عضویت و دکمه‌ها را می‌سازد (بدون بررسی عضویت)
func (r *CovoBot) buildJoinPromptWithoutCheck(chatID int64) (bool, tgbotapi.MessageConfig) {
	channels, err := r.storage.ListRequiredChannels(0)
	if err != nil || len(channels) == 0 {
		return false, tgbotapi.MessageConfig{}
	}
	text := "برای استفاده از ربات، لطفاً در کانال‌های زیر عضو شوید:"
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, ch := range channels {
		link := ch.Link
		if link == "" && ch.ChannelUsername != "" {
			link = "https://t.me/" + ch.ChannelUsername
		}
		if link == "" {
			continue
		}
		title := ch.Title
		if title == "" {
			title = "Join"
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(title, link),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ عضو شدم، بررسی کن", "admin_check_join"),
	))
	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = kb
	return true, msg
}

func (r *CovoBot) handleStartCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	chatType := update.Message.Chat.Type
	userID := update.Message.From.ID

	var response string

	// تشخیص نوع چت و ارسال پیام مناسب
	if chatType == "private" {
		// فقط یک‌بار در اولین استارت پیام عضویت را بدون چک ارسال کن
		if sent, err := r.storage.WasPromoSent(userID); err == nil && !sent {
			if has, prompt := r.buildJoinPromptWithoutCheck(chatID); has {
				_ = r.storage.MarkPromoSent(userID)
				return prompt
			}
			_ = r.storage.MarkPromoSent(userID)
		}
		// بررسی اینکه آیا کاربر ادمین است
		if r.adminCommand.IsAdmin(userID) {
			// پیام مخصوص ادمین‌ها
			response = r.adminCommand.GetAdminWelcome(userID)
		} else {
			// پیام عادی برای کاربران
			response = `🤖 *به بات covo خوش آمدید!*

من دستیار هوشمند شما با قابلیت‌های جالب هستم:

💡 *دستورات:*
• /covo <سوال> - هر سوالی دارید بپرسید!
• /cj <موضوع> - جوک خنده‌دار درباره هر موضوعی تولید کن
• /music - پیشنهاد موسیقی بر اساس سلیقه شما
• /فال - دریافت فال حافظ
• /crs - بررسی وضعیت بات
• /help - نمایش این پیام راهنما

🎯 *ویژگی‌ها:*
• درخواست‌های نامحدود
• پشتیبانی از چت خصوصی و گروهی
• خلاصه‌های هوشمند روزانه گروه‌ها
• تولید جوک بر اساس موضوع
• پیشنهاد موسیقی هوشمند

بیایید شروع کنیم! با /covo <سوال شما> چیزی از من بپرسید 🚀`
		}
	} else {
		// پیام برای گروه‌ها
		response = `🤖 *سلام! من بات covo هستم!*

من دستیار هوشمند شما با قابلیت‌های جالب هستم:

💡 *دستورات:*
• /covo <سوال> - هر سوالی دارید بپرسید!
• /cj <موضوع> - جوک خنده‌دار درباره هر موضوعی تولید کن
• /music - پیشنهاد موسیقی بر اساس سلیقه شما
• دلقک <نام> - توهین به شخص مورد نظر
• /crushon - فعال‌سازی قابلیت کراش
• /فال - دریافت فال حافظ
• /crs - بررسی وضعیت بات
• پنل - نمایش دستورات مخصوص گروه
• /covog - نمایش راهنما
• /help - نمایش راهنما

🎯 *ویژگی‌ها:*
• درخواست‌های نامحدود
• خلاصه‌های هوشمند روزانه گروه‌ها
• تولید جوک بر اساس موضوع
• پیشنهاد موسیقی هوشمند
• قابلیت دلقک برای توهین هوشمند
• قابلیت کراش خودکار هر 15 ساعت

بیایید شروع کنیم! با /covo <سوال شما> چیزی از من بپرسید 🚀`
	}

	msg := tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

func (r *CovoBot) handleHelpCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID

	response := `📚 *راهنمای باتcovo *

🤖 *دستورات دستیار هوشمند:*
• /covo <سوال> - هر سوالی دارید بپرسید! من پاسخ مفید می‌دهم
• /cj <موضوع> - جوک خنده‌دار و تمیز درباره هر موضوعی تولید کن
• /music - پیشنهاد موسیقی بر اساس سلیقه شما (با ریپلای)
• دلقک <نام> - توهین به شخص مورد نظر
• /crushon - فعال‌سازی قابلیت کراش
• /فال - دریافت فال حافظ با تفسیر
• /crs - بررسی وضعیت بات
• /gap - نمایش دستورات مخصوص گروه
• /covog - نمایش راهنما (در گروه‌ها)

📊 *استفاده:*
• درخواست‌های نامحدود
• بدون تأخیر بین درخواست‌ها
• بدون محدودیت روزانه

👥 *ویژگی‌های گروه:*
• در چت خصوصی و گروه‌ها کار می‌کند
• به طور خودکار پیام‌های گروه را ثبت می‌کند
• خلاصه‌های هوشمند روزانه ساعت ۹ صبح ارسال می‌کند

🤡 *قابلیت دلقک:*
• بنویسید: دلقک <نام>
• مثال: دلقک علی یا دلقک @username
• بات به‌صورت تصادفی یک پاسخ ارسال می‌کند

💘 *قابلیت کراش:*
• با /crushon قابلیت را فعال کنید
• هر 10 ساعت یک جفت کراش جدید اعلام می‌شود

• با /کراشوضعیت وضعیت را بررسی کنید

🎲 *بازی جرات یا حقیقت +۱۸:*
• «بازی» (بدون اسلش، فقط ادمین) — ایجاد روم و شروع ثبت‌نام با دکمه‌های اینلاین
• «توقف بازی» (بدون اسلش، فقط ادمین) — پایان بازی و بستن روم
• بعد از بستن ثبت‌نام: نوبت‌ها به‌ترتیب شرکت‌کنندگان است؛ هر نفر «جرات» یا «سوال +۱۸» را انتخاب می‌کند
• بات یک سؤال تصادفی می‌فرستد؛ کاربر باید ریپلای کند و سپس دکمه «✅ جواب دادم» را بزند تا نفر بعدی فعال شود

💡 *نکات:*
• برای پاسخ‌های بهتر، سوالات خود را دقیق مطرح کنید
• موضوعات مختلف را برای جوک امتحان کنید
• برای موسیقی، ابتدا /music بزنید، سپس ترجیحات خود را ریپلای کنید
• از /crs برای بررسی وضعیت بات استفاده کنید
• در گروه‌ها از /covog برای راهنما استفاده کنید

نیاز به کمک دارید؟ فقط بپرسید! 😊`

	msg := tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

func main() {
	bot, err := NewCovoBot()
	if err != nil {
		log.Fatal("خطا در ایجاد بات:", err)
	}

	log.Println("🚀 راه‌اندازی بات کوو...")

	if err := bot.Start(); err != nil {
		log.Fatal("خطا در راه‌اندازی بات:", err)
	}
}
