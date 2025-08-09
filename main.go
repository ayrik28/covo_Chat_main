package main

import (
	"fmt"
	"log"
	"redhat-bot/ai"
	"redhat-bot/commands"
	"redhat-bot/config"
	"redhat-bot/limiter"

	// "redhat-bot/scheduler"
	"redhat-bot/storage"
	"strings"

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
	// summaryScheduler *scheduler.DailySummaryScheduler
	cron *cron.Cron
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
	clownCommand := commands.NewClownCommand(aiClient, rateLimiter, bot)
	crushCommand := commands.NewCrushCommand(storage, bot)
	hafezCommand := commands.NewHafezCommand(bot)
	adminCommand := commands.NewAdminCommand(bot, storage)
	gapCommand := commands.NewGapCommand(bot, storage, hafezCommand)
	moderationCommand := commands.NewModerationCommand(bot)

	// راه‌اندازی زمان‌بند
	// summaryScheduler := scheduler.NewDailySummaryScheduler(bot, storage, aiClient)

	// راه‌اندازی کران
	cronJob := cron.New()

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

	r.cron.Start()
	log.Println("⏰ زمان‌بند خلاصه روزانه راه‌اندازی شد (ساعت ۹ صبح هر روز)")

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
	// Handle callback queries from inline keyboard
	if update.CallbackQuery != nil {
		var callback tgbotapi.CallbackConfig

		// بررسی نوع callback
		switch {
		case strings.HasPrefix(update.CallbackQuery.Data, "admin_"):
			callback = r.adminCommand.HandleCallback(update)
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
				welcomeMsg := tgbotapi.NewMessage(message.Chat.ID, `🤖 *سلام! من بات covo هستم!*

من دستیار هوشمند شما با قابلیت‌های جالب هستم:

💡 *دستورات:*
• /covo <سوال> - هر سوالی دارید بپرسید!
• /cj <موضوع> - جوک خنده‌دار درباره هر موضوعی تولید کن
• /music - پیشنهاد موسیقی بر اساس سلیقه شما
• /clown <نام> - توهین هوشمند به شخص مورد نظر
• /crushon - فعال‌سازی قابلیت کراش
• /فال - دریافت فال حافظ
• /crs - بررسی وضعیت بات
• /gap - نمایش دستورات مخصوص گروه
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
	}

	// پردازش دستورات
	if !strings.HasPrefix(text, "/") {
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

	switch {
	case strings.HasPrefix(text, "/covo"):
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
	case strings.HasPrefix(text, "/gap"):
		response = r.gapCommand.Handle(update)
	case strings.HasPrefix(text, "/فال"):
		response = r.hafezCommand.Handle(update)
	case strings.HasPrefix(text, "/admin"):
		response = r.adminCommand.Handle(update)
	case strings.HasPrefix(text, "/showusers"):
		response = r.adminCommand.HandleShowUsers(update)
	case strings.HasPrefix(text, "/showgroups"):
		response = r.adminCommand.HandleShowGroups(update)
	case strings.HasPrefix(text, "/del"):
		response = r.moderationCommand.HandleDelete(update)
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

func (r *CovoBot) handleStartCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	chatType := update.Message.Chat.Type
	userID := update.Message.From.ID

	var response string

	// تشخیص نوع چت و ارسال پیام مناسب
	if chatType == "private" {
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
• /clown <نام> - توهین هوشمند به شخص مورد نظر
• /crushon - فعال‌سازی قابلیت کراش
• /فال - دریافت فال حافظ
• /crs - بررسی وضعیت بات
• /gap - نمایش دستورات مخصوص گروه
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
• /clown <نام> - توهین هوشمند به شخص مورد نظر
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
• دستور /clown را با نام شخص تایپ کنید
• مثال: /clown علی یا /clown @username
• بات با هوش مصنوعی به آن شخص توهین می‌کند

💘 *قابلیت کراش:*
• با /crushon قابلیت را فعال کنید
• هر 10 ساعت یک جفت کراش جدید اعلام می‌شود

• با /کراشوضعیت وضعیت را بررسی کنید

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
