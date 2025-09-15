# 📚 مستندات API - Covo Bot

## 🎯 معرفی

این مستندات شامل جزئیات کامل API، ساختار کد، و نحوه استفاده از تمام قابلیت‌های Covo Bot است.

---

## 🏗️ ساختار کلی

### 📁 **معماری پروژه**

```
Covo_Chat/
├── main.go                 # نقطه ورود اصلی
├── config/                 # مدیریت تنظیمات
├── commands/               # دستورات بات
├── storage/                # لایه ذخیره‌سازی
├── ai/                     # ادغام هوش مصنوعی
├── limiter/                # محدودیت درخواست
├── scheduler/              # زمان‌بندی
└── jsonfile/               # فایل‌های داده
```

---

## 🔧 API Reference

### 🤖 **دستورات هوش مصنوعی**

#### `/covo <سوال>`
**توضیحات:** پرسش و پاسخ هوشمند با استفاده از DeepSeek AI

**پارامترها:**
- `سوال` (string, اجباری): سوال کاربر

**مثال:**
```go
// در commands/covo.go
func (r *CovoCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
    userID := update.Message.From.ID
    chatID := update.Message.Chat.ID
    
    // بررسی محدودیت درخواست
    if allowed, message := r.rateLimiter.CheckRateLimit(userID); !allowed {
        return tgbotapi.NewMessage(chatID, message)
    }
    
    // استخراج سوال از دستور
    text := update.Message.Text
    question := strings.TrimSpace(strings.TrimPrefix(text, "/covo"))
    
    // دریافت پاسخ از هوش مصنوعی
    response, err := r.aiClient.AskQuestion(question)
    // ...
}
```

**پاسخ:**
```
🤖 هوش مصنوعی کوو

[پاسخ هوش مصنوعی]
```

#### `/cj <موضوع>`
**توضیحات:** تولید جوک بر اساس موضوع

**پارامترها:**
- `موضوع` (string, اجباری): موضوع جوک

**مثال:**
```go
// در commands/rtj.go
func (r *CovoJokeCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
    text := update.Message.Text
    topic := strings.TrimSpace(strings.TrimPrefix(text, "/covoJoke"))
    
    // ساخت درخواست برای هوش مصنوعی
    prompt := fmt.Sprintf("هی، یک جوک خنده‌دار و مناسب خانواده درباره '%s' تولید کن و ارسال کن.", topic)
    
    // استفاده از AskQuestion برای ارسال درخواست
    joke, err := r.aiClient.AskQuestion(prompt)
    // ...
}
```

#### `/music`
**توضیحات:** پیشنهاد موسیقی بر اساس سلیقه

**نحوه استفاده:**
1. `/music` را ارسال کنید
2. روی پیام ریپلای کنید
3. سلیقه موسیقی خود را بنویسید

---

### 🎮 **بازی‌ها و سرگرمی**

#### **دلقک**
**توضیحات:** توهین هوشمند به اعضای گروه

**نحوه استفاده:**
```
دلقک علی
دلقک @username
```

**کد:**
```go
// در commands/clown.go
func (r *ClownCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
    // بررسی فعال بودن قابلیت دلقک در این گروه
    enabled, err := r.storage.IsClownEnabled(chatID)
    if !enabled {
        return tgbotapi.NewMessage(chatID, "🔒 قابلیت دلقک در این گروه غیرفعال است")
    }
    
    // استخراج نام مخاطب از دستور
    text := update.Message.Text
    cleaned := strings.TrimSpace(strings.TrimPrefix(text, "دلقک"))
    targetName := cleaned
    
    // انتخاب تصادفی یک فحش از لیست
    insult := r.randomInsult()
    formattedResponse := fmt.Sprintf("🤡 *دلقک به %s:*\n\n%s", targetName, insult)
    // ...
}
```

#### **کراش**
**توضیحات:** اعلام خودکار جفت‌های تصادفی

**دستورات:**
- `/crushon` - فعال‌سازی
- `/crushoff` - غیرفعال‌سازی
- `/کراشوضعیت` - نمایش وضعیت

**کد:**
```go
// در commands/crush.go
func (r *CrushCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
    text := update.Message.Text
    
    // بررسی دستور فعال‌سازی
    if text == "/crushon" {
        if err := r.storage.SetCrushEnabled(chatID, true); err != nil {
            return tgbotapi.NewMessage(chatID, "❌ خطا در فعال‌سازی قابلیت کراش")
        }
        // پیام تایید
    }
    // ...
}
```

#### **فال حافظ**
**توضیحات:** دریافت فال با تفسیر کامل

**دستور:**
```
/فال
فال
```

**کد:**
```go
// در commands/hafez.go
func (r *HafezCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
    // دریافت فال
    text, err := r.getHafezFal()
    if err != nil {
        return tgbotapi.NewMessage(chatID, "❌ متأسفانه در دریافت فال خطایی رخ داد")
    }
    
    // اضافه کردن دکمه فال جدید به پیام
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🎲 فال جدید", "new_hafez"),
        ),
    )
    // ...
}
```

#### **چلنج روزانه**
**توضیحات:** بازی حدس ضرب‌المثل با ایموجی

**کد:**
```go
// در commands/daily_challenge.go
func (d *DailyChallengeCommand) PostDailyChallenge(groupID int64) {
    emojis, proverb, ok := getRandomZarb()
    if !ok {
        return
    }
    
    text := fmt.Sprintf("🧩 چلنج روزانه\n\n%s\n\nاز روی ایموجی ضرب‌المثل را حدس بزنید و روی همین پیام ریپلای کنید.\nاولین پاسخ صحیح لقب «باهوش‌ترین فرد گروه» را می‌گیرد!", emojis)
    msg := tgbotapi.NewMessage(groupID, text)
    // ...
}
```

#### **جرات یا حقیقت +18**
**توضیحات:** بازی تعاملی با دکمه‌های اینلاین

**کد:**
```go
// در commands/truthdare.go
func (r *TruthDareCommand) HandleStartWithoutSlash(update tgbotapi.Update) tgbotapi.MessageConfig {
    // initialize game room
    g := &tdGame{
        chatID:           chatID,
        starterUserID:    userID,
        isOpen:           true,
        participants:     []int64{},
        participantNames: map[int64]string{},
        currentIndex:     0,
        activeUserID:     0,
    }
    r.games[chatID] = g
    
    // announce room with join/close buttons
    text := "🎮 بازی جرات یا سوال +۱۸ شروع شد!\nاگر می‌خوای شرکت کنی، روی دکمه زیر بزن."
    joinBtn := tgbotapi.NewInlineKeyboardButtonData("➕ جوین شو", fmt.Sprintf("td_join:%d", chatID))
    closeBtn := tgbotapi.NewInlineKeyboardButtonData("🔒 بستن بازی (فقط ادمین)", fmt.Sprintf("td_close:%d", chatID))
    // ...
}
```

---

### 👥 **مدیریت گروه**

#### **حذف پیام**
**توضیحات:** حذف تکی یا دسته‌ای پیام‌ها

**دستورات:**
```
حذف 10          # حذف 10 پیام قبلی
حذف             # حذف پیام ریپلای شده
```

**کد:**
```go
// در commands/moderation.go
func (m *ModerationCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
    text := strings.TrimSpace(update.Message.Text)
    fields := strings.Fields(text)
    var count int
    if len(fields) > 1 {
        if n, err := strconv.Atoi(fields[1]); err == nil && n > 0 {
            if n > 300 {
                n = 300
            }
            count = n
        }
    }
    
    if count > 0 {
        // Bulk delete previous N messages
        m.bulkDeletePrev(chatID, update.Message.MessageID, count)
        return tgbotapi.MessageConfig{}
    }
    // ...
}
```

#### **سکوت کاربر**
**توضیحات:** سکوت موقت یا نامحدود

**دستورات:**
```
سکوت 2          # سکوت 2 ساعته
سکوت            # سکوت نامحدود
آزاد            # خارج کردن از سکوت
```

**کد:**
```go
// در commands/moderation.go
func (m *ModerationCommand) HandleMute(update tgbotapi.Update) tgbotapi.MessageConfig {
    // Parse optional hours
    var until int64 = 0
    fields := strings.Fields(strings.TrimSpace(update.Message.Text))
    if len(fields) > 1 {
        if hours, err := strconv.Atoi(fields[1]); err == nil && hours > 0 {
            until = time.Now().Add(time.Duration(hours) * time.Hour).Unix()
        }
    }
    
    restrictCfg := tgbotapi.RestrictChatMemberConfig{
        ChatMemberConfig: tgbotapi.ChatMemberConfig{
            ChatID: chatID,
            UserID: targetUserID,
        },
        Permissions: &tgbotapi.ChatPermissions{
            CanSendMessages:       false,
            CanSendMediaMessages:  false,
            // ... سایر مجوزها
        },
        UntilDate: until,
    }
    // ...
}
```

#### **بن کاربر**
**توضیحات:** اخراج دائمی اعضا

**دستور:**
```
بن              # بن کردن کاربر ریپلای شده
```

#### **تگ همه**
**توضیحات:** تگ کردن تمام اعضای گروه

**دستور:**
```
تگ              # تگ کردن تمام اعضا (روی پیام ریپلای)
```

**کد:**
```go
// در commands/tag.go
func (t *TagCommand) HandleTagAllOnReply(update tgbotapi.Update) tgbotapi.MessageConfig {
    // Load members
    members, err := t.storage.GetGroupMembers(chatID)
    if err != nil || len(members) == 0 {
        return tgbotapi.NewMessage(chatID, "❌ لیست اعضای گروه پیدا نشد")
    }
    
    // Build and send in chunks
    const chunkSize = 20
    var batch []string
    for _, m := range members {
        // Build HTML mention using tg://user?id
        escaped := html.EscapeString(displayName)
        mention := fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", m.UserID, escaped)
        batch = append(batch, mention)
        if len(batch) >= chunkSize {
            flush()
        }
    }
    // ...
}
```

---

### 📊 **آمار و گزارش**

#### **دسترسی از پنل:**
- **پنل** - نمایش منوی اصلی
- **آمار پیام** - آمار 24 ساعته
- **کاربران برتر** - 10 کاربر فعال
- **آمار من** - آمار شخصی

**کد:**
```go
// در commands/gap.go
case "show_stats":
    // چک فعال بودن قابلیت
    enabled, _ := r.storage.IsFeatureEnabled(chatID, "stats")
    if !enabled {
        msg.Text = "ℹ️ آمار پیام‌ها غیر فعال است. ابتدا آن را فعال کنید."
        break
    }
    // دریافت ۱۰ کاربر برتر
    top, err := r.storage.GetTopActiveUsersLast24h(chatID, 10)
    if err != nil {
        msg.Text = "❌ خطا در دریافت آمار"
        break
    }
    // ...
```

---

### 🔒 **امنیت و قفل‌ها**

#### **قفل‌های موجود:**
- **دلقک** - فعال/غیرفعال کردن توهین
- **لینک** - حذف پیام‌های حاوی لینک
- **فحش** - حذف پیام‌های نامناسب

**کد:**
```go
// در main.go
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
```

---

## 🗄️ پایگاه داده

### 📊 **جداول اصلی**

#### **user_usage**
```sql
CREATE TABLE user_usage (
    user_id BIGINT PRIMARY KEY,
    requests_today INT DEFAULT 0,
    last_reset TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_request TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### **group_messages**
```sql
CREATE TABLE group_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    username VARCHAR(255),
    message TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_group_timestamp (group_id, timestamp)
);
```

#### **group_members**
```sql
CREATE TABLE group_members (
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    name VARCHAR(255),
    PRIMARY KEY (group_id, user_id),
    UNIQUE INDEX idx_group_user (group_id, user_id)
);
```

#### **feature_settings**
```sql
CREATE TABLE feature_settings (
    group_id BIGINT NOT NULL,
    feature_name VARCHAR(255) NOT NULL,
    enabled BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (group_id, feature_name)
);
```

#### **daily_challenges**
```sql
CREATE TABLE daily_challenges (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_id BIGINT NOT NULL,
    message_id INT NOT NULL,
    proverb TEXT,
    emojis TEXT,
    answered BOOLEAN DEFAULT FALSE,
    winner_id BIGINT,
    winner_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_group (group_id),
    INDEX idx_created_at (created_at)
);
```

#### **bot_channels**
```sql
CREATE TABLE bot_channels (
    id INT AUTO_INCREMENT PRIMARY KEY,
    chat_id BIGINT UNIQUE NOT NULL,
    title VARCHAR(255),
    username VARCHAR(255),
    is_admin BOOLEAN DEFAULT FALSE,
    member_count INT DEFAULT 0,
    date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_check TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_date_added (date_added)
);
```

#### **required_channels**
```sql
CREATE TABLE required_channels (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_id BIGINT NOT NULL,
    title VARCHAR(255),
    link TEXT,
    channel_username VARCHAR(255),
    channel_id BIGINT,
    chat_id BIGINT,
    bot_joined BOOLEAN DEFAULT FALSE,
    member_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_group (group_id),
    INDEX idx_channel_username (channel_username),
    INDEX idx_channel_id (channel_id)
);
```

#### **user_onboarding**
```sql
CREATE TABLE user_onboarding (
    user_id BIGINT PRIMARY KEY,
    promo_sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_sent_at (sent_at)
);
```

---

## 🔧 توسعه

### 🛠️ **نحوه افزودن دستور جدید**

#### 1️⃣ **ایجاد فایل دستور**
```go
// commands/new_command.go
package commands

import (
    "redhat-bot/storage"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewCommand struct {
    bot     *tgbotapi.BotAPI
    storage *storage.MySQLStorage
}

func NewNewCommand(bot *tgbotapi.BotAPI, storage *storage.MySQLStorage) *NewCommand {
    return &NewCommand{
        bot:     bot,
        storage: storage,
    }
}

func (c *NewCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
    chatID := update.Message.Chat.ID
    // منطق دستور
    return tgbotapi.NewMessage(chatID, "پاسخ دستور")
}
```

#### 2️⃣ **اضافه کردن به main.go**
```go
// در NewCovoBot()
newCommand := commands.NewNewCommand(bot, storage)

// در handleUpdate()
case strings.HasPrefix(text, "/new"):
    response = r.newCommand.Handle(update)
```

### 🗄️ **نحوه افزودن جدول جدید**

#### 1️⃣ **تعریف مدل**
```go
// storage/mysql.go
type NewTable struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"type:varchar(255)"`
    CreatedAt time.Time
}
```

#### 2️⃣ **اضافه کردن به AutoMigrate**
```go
// در NewMySQLStorage()
if err := db.AutoMigrate(&NewTable{}); err != nil {
    return nil, fmt.Errorf("error migrating NewTable: %v", err)
}
```

### 🧪 **تست کردن**

```bash
# اجرای تست‌ها
go test ./...

# تست با coverage
go test -cover ./...

# تست یک پکیج خاص
go test ./commands
```

---

## 📦 وابستگی‌ها

### 🔧 **وابستگی‌های اصلی**

```go
require (
    github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
    github.com/joho/godotenv v1.5.1
    github.com/robfig/cron/v3 v3.0.1
    gorm.io/driver/mysql v1.6.0
    gorm.io/gorm v1.30.1
)
```

### 📋 **وابستگی‌های غیرمستقیم**

```go
require (
    filippo.io/edwards25519 v1.1.0
    github.com/go-sql-driver/mysql v1.9.3
    github.com/jinzhu/inflection v1.0.0
    github.com/jinzhu/now v1.1.5
    golang.org/x/text v0.20.0
)
```

---

## 🚀 راه‌اندازی

### 📋 **پیش‌نیازها**

- **Go** 1.24.5 یا بالاتر
- **MySQL** 5.7 یا بالاتر
- **Git** برای کلون کردن پروژه
- **Telegram Bot Token** از [@BotFather](https://t.me/botfather)
- **DeepSeek API Key** از [OpenRouter](https://openrouter.ai/)

### ⚡ **نصب و راه‌اندازی**

#### 1️⃣ **کلون کردن پروژه**
```bash
git clone https://github.com/your-username/covo-bot.git
cd covo-bot
```

#### 2️⃣ **نصب وابستگی‌ها**
```bash
go mod tidy
```

#### 3️⃣ **تنظیم پایگاه داده**
```sql
CREATE DATABASE myappdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'covouser'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON myappdb.* TO 'covouser'@'localhost';
FLUSH PRIVILEGES;
```

#### 4️⃣ **تنظیم متغیرهای محیطی**
```env
TELEGRAM_TOKEN=your_telegram_bot_token
DEEPSEEK_TOKEN=your_deepseek_api_key
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=covouser
MYSQL_PASSWORD=your_password
MYSQL_DATABASE=myappdb
MAX_REQUESTS_PER_DAY=1000
COOLDOWN_SECONDS=5
```

#### 5️⃣ **اجرای بات**
```bash
# حالت توسعه
go run main.go

# یا کامپایل و اجرا
go build -o covo-bot
./covo-bot
```

---

## 🔍 عیب‌یابی

### 🐛 **مشکلات رایج**

#### **خطای اتصال به پایگاه داده**
```
Error connecting to MySQL: dial tcp [::1]:3306: connect: connection refused
```
**راه حل:** مطمئن شوید MySQL در حال اجرا است و تنظیمات درست است.

#### **خطای Telegram API**
```
Error sending message: 400 Bad Request
```
**راه حل:** توکن بات را بررسی کنید و مطمئن شوید بات در گروه‌ها دسترسی دارد.

#### **خطای DeepSeek API**
```
Error making request: 401 Unauthorized
```
**راه حل:** کلید API DeepSeek را بررسی کنید.

### 📝 **لاگ‌ها**

```bash
# اجرا با لاگ‌های تفصیلی
go run main.go 2>&1 | tee bot.log

# یا در production
./covo-bot >> bot.log 2>&1 &
```

---

## 📚 منابع بیشتر

- [Telegram Bot API Documentation](https://core.telegram.org/bots/api)
- [Go Documentation](https://golang.org/doc/)
- [GORM Documentation](https://gorm.io/docs/)
- [DeepSeek API Documentation](https://platform.deepseek.com/api-docs/)

---

<div align="center">

**📚 این مستندات به‌روزرسانی می‌شود**

[🔝 بازگشت به بالا](#-مستندات-api---covo-bot)

</div>
