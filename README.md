# 🤖 Covo Bot - Telegram AI Assistant
این پروژه برای من فقط حکم سرگرمی داشته ولی بخاطر بزرگتر شدن اسکیلش گفتم بد نیست اوپن سورس باشه ❤️

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.24.5-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Telegram](https://img.shields.io/badge/Telegram-Bot-blue.svg)
![AI](https://img.shields.io/badge/AI-DeepSeek-purple.svg)
![Database](https://img.shields.io/badge/Database-MySQL-orange.svg)

**یک بات تلگرام پیشرفته و چندمنظوره با قابلیت‌های هوش مصنوعی، بازی‌های تعاملی و مدیریت گروه**

[🚀 Quick Start](#-quick-start) • [📚 Documentation](#-documentation) • [⚙️ Configuration](#️-configuration) • [🎮 Features](#-features)

</div>

---

## 📋 فهرست مطالب

- [🎯 معرفی](#-معرفی)
- [✨ ویژگی‌ها](#-ویژگی‌ها)
- [🏗️ معماری](#️-معماری)
- [🚀 راه‌اندازی سریع](#-راه‌اندازی-سریع)
- [⚙️ تنظیمات](#️-تنظیمات)
- [📚 مستندات کامل](#-مستندات-کامل)
- [🔧 توسعه](#-توسعه)
- [📊 آمار پروژه](#-آمار-پروژه)
- [🤝 مشارکت](#-مشارکت)
- [📄 مجوز](#-مجوز)

---

## 🎯 معرفی

**Covo Bot** یک بات تلگرام پیشرفته و چندمنظوره است که با زبان Go نوشته شده و قابلیت‌های متنوعی از جمله هوش مصنوعی، بازی‌های تعاملی، مدیریت گروه و آمارگیری را ارائه می‌دهد.

### 🌟 چرا Covo Bot؟

- **🤖 هوش مصنوعی پیشرفته** - استفاده از DeepSeek برای پاسخ‌های هوشمند
- **🎮 بازی‌های تعاملی** - جرات یا حقیقت، چلنج روزانه، فال حافظ
- **📊 آمارگیری کامل** - آمار پیام‌ها، کاربران فعال، گزارش‌های روزانه
- **🔒 امنیت بالا** - سیستم محدودیت، فیلتر کلمات، مدیریت دسترسی
- **⚡ عملکرد عالی** - Go + MySQL + Cron Jobs
- **🎨 رابط کاربری زیبا** - کیبوردهای اینلاین و پیام‌های فرمت‌شده

---

## ✨ ویژگی‌ها

### 🤖 **دستورات هوش مصنوعی**
- **`/covo <سوال>`** - پرسش و پاسخ هوشمند با DeepSeek
- **`/cj <موضوع>`** - تولید جوک بر اساس موضوع
- **`/music`** - پیشنهاد موسیقی (با ریپلای)

### 🎮 **بازی‌ها و سرگرمی**
- **دلقک** - توهین هوشمند به اعضای گروه
- **کراش** - اعلام خودکار جفت‌های تصادفی هر 10 ساعت
- **فال حافظ** - دریافت فال با تفسیر کامل
- **چلنج روزانه** - بازی حدس ضرب‌المثل با ایموجی
- **جرات یا حقیقت +18** - بازی تعاملی با دکمه‌های اینلاین

### 👥 **مدیریت گروه**
- **حذف پیام** - حذف تکی یا دسته‌ای پیام‌ها
- **سکوت کاربر** - سکوت موقت یا نامحدود
- **بن کاربر** - اخراج دائمی اعضا
- **تگ همه** - تگ کردن تمام اعضای گروه

### 📊 **آمار و گزارش**
- **آمار پیام‌ها** - نمایش آمار 24 ساعته
- **کاربران فعال** - لیست کاربران پرکار
- **آمار شخصی** - آمار فردی هر کاربر
- **خلاصه روزانه** - تحلیل هوشمند پیام‌های گروه

### 🔒 **امنیت و قفل‌ها**
- **قفل لینک** - حذف خودکار پیام‌های حاوی لینک
- **قفل فحش** - حذف پیام‌های حاوی کلمات نامناسب
- **عضویت اجباری** - اجبار عضویت در کانال‌های مشخص
- **محدودیت درخواست** - 1000 درخواست در روز + 5 ثانیه فاصله

---

## 🏗️ معماری

### 📁 **ساختار پروژه**

```
Covo_Chat/
├── 📁 ai/                    # ادغام هوش مصنوعی
│   └── openai.go            # کلاینت DeepSeek
├── 📁 commands/              # دستورات بات
│   ├── admin.go             # مدیریت ادمین
│   ├── clown.go             # قابلیت دلقک
│   ├── covo.go              # دستور اصلی AI
│   ├── crush.go             # قابلیت کراش
│   ├── daily_challenge.go   # چلنج روزانه
│   ├── gap.go               # پنل مدیریت
│   ├── hafez.go             # فال حافظ
│   ├── moderation.go        # مدیریت گروه
│   ├── music.go             # پیشنهاد موسیقی
│   ├── redhat.go            # دستور اصلی (legacy)
│   ├── rrs.go               # وضعیت بات
│   ├── rtj.go               # تولید جوک
│   ├── tag.go               # تگ کردن
│   └── truthdare.go         # بازی جرات یا حقیقت
├── 📁 config/                # تنظیمات
│   └── env.go               # مدیریت متغیرهای محیطی
├── 📁 jsonfile/              # فایل‌های داده
│   ├── badwords.json        # کلمات نامناسب
│   ├── clown.json           # متن‌های دلقک
│   ├── dare.json            # چالش‌های جرات
│   ├── fal.json             # فال‌های حافظ
│   ├── truth+18.json        # سوالات +18
│   └── zarb.json            # ضرب‌المثل‌ها
├── 📁 limiter/               # محدودیت درخواست
│   └── rate_limiter.go      # سیستم Rate Limiting
├── 📁 scheduler/             # زمان‌بندی
│   └── daily_summary.go     # خلاصه روزانه
├── 📁 storage/               # ذخیره‌سازی
│   ├── mysql.go             # پایگاه داده MySQL
│   └── memory.go            # ذخیره‌سازی حافظه (غیرفعال)
├── 📄 main.go                # نقطه ورود اصلی
├── 📄 go.mod                 # وابستگی‌های Go
├── 📄 go.sum                 # چک‌سام وابستگی‌ها
└── 📄 README.md              # این فایل
```

### 🗄️ **پایگاه داده**

#### **جداول اصلی:**

| جدول | توضیحات | کلید اصلی |
|------|---------|-----------|
| `user_usage` | آمار استفاده کاربران | `user_id` |
| `group_messages` | پیام‌های گروه‌ها (24 ساعت) | `id` |
| `group_members` | اعضای گروه‌ها | `group_id, user_id` |
| `feature_settings` | تنظیمات قابلیت‌ها | `group_id, feature_name` |
| `daily_challenges` | چالش‌های روزانه | `id` |
| `bot_channels` | کانال‌های ربات | `chat_id` |
| `required_channels` | کانال‌های الزامی | `id` |
| `user_onboarding` | پیگیری عضویت | `user_id` |

---

## 🚀 راه‌اندازی سریع

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
```bash
# ایجاد فایل .env
cp .env.example .env
```

```env
# فایل .env
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

## ⚙️ تنظیمات

### 🔧 **متغیرهای محیطی**

| متغیر | پیش‌فرض | توضیحات |
|-------|---------|---------|
| `TELEGRAM_TOKEN` | - | توکن بات تلگرام (اجباری) |
| `DEEPSEEK_TOKEN` | - | کلید API DeepSeek (اجباری) |
| `MYSQL_HOST` | `localhost` | آدرس سرور MySQL |
| `MYSQL_PORT` | `3306` | پورت MySQL |
| `MYSQL_USER` | `covouser` | نام کاربری MySQL |
| `MYSQL_PASSWORD` | - | رمز عبور MySQL (اجباری) |
| `MYSQL_DATABASE` | `myappdb` | نام پایگاه داده |
| `MAX_REQUESTS_PER_DAY` | `1000` | حداکثر درخواست روزانه |
| `COOLDOWN_SECONDS` | `5` | فاصله زمانی بین درخواست‌ها |

### ⏰ **تنظیمات زمان‌بندی**

```go
// در main.go - تنظیم زمان‌بندی Cron Jobs
_, err := r.cron.AddFunc("0 9 * * *", func() {
    // خلاصه روزانه ساعت 9 صبح
})

_, err := r.cron.AddFunc("0 10 * * *", func() {
    // چلنج روزانه ساعت 10 صبح
})
```

### 🎛️ **تنظیمات قابلیت‌ها**

هر گروه می‌تواند قابلیت‌های زیر را فعال/غیرفعال کند:

- **کراش** - اعلام خودکار جفت‌های تصادفی
- **فال** - قابلیت فال حافظ
- **آمار** - آمار پیام‌ها
- **چلنج روزانه** - بازی ضرب‌المثل
- **دلقک** - قابلیت توهین
- **قفل لینک** - حذف پیام‌های حاوی لینک
- **قفل فحش** - حذف پیام‌های نامناسب

---

## 📚 مستندات کامل

### 🤖 **دستورات هوش مصنوعی**

#### `/covo <سوال>`
پرسش و پاسخ هوشمند با استفاده از DeepSeek AI

**مثال:**
```
/covo بهترین زبان برنامه‌نویسی چیست؟
```

**پاسخ:**
```
🤖 هوش مصنوعی کوو

بهترین زبان برنامه‌نویسی بستگی به نیاز و هدف شما دارد:

• Python: برای هوش مصنوعی و تحلیل داده
• JavaScript: برای توسعه وب
• Go: برای سیستم‌های backend
• Rust: برای برنامه‌نویسی سیستم
• Java: برای برنامه‌نویسی enterprise
```

#### `/cj <موضوع>`
تولید جوک بر اساس موضوع

**مثال:**
```
/cj برنامه‌نویسی
```

**پاسخ:**
```
😄 تولیدکننده جوک کوو

موضوع: برنامه‌نویسی

چرا برنامه‌نویس‌ها همیشه پنجره‌ها را باز می‌گذارند؟
چون می‌خواهند ببینند چه اتفاقی در console می‌افتد! 🖥️
```

#### `/music`
پیشنهاد موسیقی بر اساس سلیقه

**نحوه استفاده:**
1. `/music` را ارسال کنید
2. روی پیام ریپلای کنید
3. سلیقه موسیقی خود را بنویسید

**مثال:**
```
ریپلای: موسیقی غمگین و آرام
```

### 🎮 **بازی‌ها و سرگرمی**

#### **دلقک**
توهین هوشمند به اعضای گروه

**نحوه استفاده:**
```
دلقک علی
دلقک @username
```

#### **کراش**
اعلام خودکار جفت‌های تصادفی

**دستورات:**
- `/crushon` - فعال‌سازی
- `/crushoff` - غیرفعال‌سازی
- `/کراشوضعیت` - نمایش وضعیت

#### **فال حافظ**
دریافت فال با تفسیر کامل

**دستور:**
```
/فال
فال
```

#### **چلنج روزانه**
بازی حدس ضرب‌المثل با ایموجی

**نحوه بازی:**
1. هر روز ساعت 10 صبح ایموجی ارسال می‌شود
2. روی پیام ریپلای کنید
3. ضرب‌المثل را حدس بزنید
4. اولین پاسخ صحیح برنده می‌شود

#### **جرات یا حقیقت +18**
بازی تعاملی با دکمه‌های اینلاین

**نحوه بازی:**
1. ادمین: `بازی` را تایپ کند
2. اعضا روی دکمه "جوین شو" کلیک کنند
3. ادمین: "بستن بازی" را کلیک کند
4. هر نفر نوبت خود را انتخاب کند
5. سوال/چالش را پاسخ دهد

### 👥 **مدیریت گروه**

#### **حذف پیام**
```
حذف 10          # حذف 10 پیام قبلی
حذف             # حذف پیام ریپلای شده
```

#### **سکوت کاربر**
```
سکوت 2          # سکوت 2 ساعته
سکوت            # سکوت نامحدود
آزاد            # خارج کردن از سکوت
```

#### **بن کاربر**
```
بن              # بن کردن کاربر ریپلای شده
```

#### **تگ همه**
```
تگ              # تگ کردن تمام اعضا (روی پیام ریپلای)
```

### 📊 **آمار و گزارش**

#### **دسترسی از پنل:**
- **پنل** - نمایش منوی اصلی
- **آمار پیام** - آمار 24 ساعته
- **کاربران برتر** - 10 کاربر فعال
- **آمار من** - آمار شخصی

### 🔒 **امنیت و قفل‌ها**

#### **قفل‌های موجود:**
- **دلقک** - فعال/غیرفعال کردن توهین
- **لینک** - حذف پیام‌های حاوی لینک
- **فحش** - حذف پیام‌های نامناسب

#### **عضویت اجباری:**
- ادمین‌ها می‌توانند کانال‌های الزامی تعریف کنند
- کاربران باید عضو کانال‌ها باشند تا از بات استفاده کنند

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

## 📊 آمار پروژه

### 📈 **آمار کلی**
- **زبان برنامه‌نویسی:** Go 1.24.5
- **تعداد فایل‌های Go:** 15+
- **تعداد دستورات:** 20+
- **تعداد جداول DB:** 8
- **تعداد فایل‌های JSON:** 6
- **خطوط کد:** ~3000+
- **پوشش تست:** در حال توسعه

### 🏗️ **معماری**
- **Backend:** Go + GORM
- **Database:** MySQL 5.7+
- **AI:** DeepSeek via OpenRouter
- **Scheduling:** Cron Jobs
- **API:** Telegram Bot API

### 📦 **وابستگی‌های اصلی**
```go
require (
    github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
    github.com/joho/godotenv v1.5.1
    github.com/robfig/cron/v3 v3.0.1
    gorm.io/driver/mysql v1.6.0
    gorm.io/gorm v1.30.1
)
```

---

## 🤝 مشارکت

### 🔄 **نحوه مشارکت**

1. **Fork** کنید
2. **Branch** جدید ایجاد کنید (`git checkout -b feature/amazing-feature`)
3. **Commit** کنید (`git commit -m 'Add amazing feature'`)
4. **Push** کنید (`git push origin feature/amazing-feature`)
5. **Pull Request** ایجاد کنید

### 📋 **راهنمای مشارکت**

- کد را تمیز و قابل خواندن بنویسید
- کامنت‌های مناسب اضافه کنید
- تست‌های لازم را بنویسید
- مستندات را به‌روزرسانی کنید
- از نام‌گذاری مناسب استفاده کنید

### 🐛 **گزارش باگ**

برای گزارش باگ:
1. Issue جدید ایجاد کنید
2. توضیح کامل مشکل را بنویسید
3. مراحل تکرار باگ را ذکر کنید
4. لاگ‌های مربوطه را ضمیمه کنید

---

## 📄 مجوز

این پروژه تحت مجوز **MIT** منتشر شده است.

```
MIT License

Copyright (c) 2024 Covo Bot

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 🙏 تشکر

- **Go Team** برای زبان برنامه‌نویسی عالی
- **Telegram** برای API قدرتمند
- **DeepSeek** برای هوش مصنوعی پیشرفته
- **GORM** برای ORM عالی
- **OpenRouter** برای دسترسی آسان به AI

---

<div align="center">

**⭐ اگر این پروژه برایتان مفید بود، ستاره بدهید!**

[🔝 بازگشت به بالا](#-covo-bot---telegram-ai-assistant)

</div>
