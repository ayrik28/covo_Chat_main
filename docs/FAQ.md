# ❓ سوالات متداول - Covo Bot

## 📋 فهرست مطالب

- [🤖 سوالات عمومی](#-سوالات-عمومی)
- [⚙️ نصب و راه‌اندازی](#️-نصب-و-راه‌اندازی)
- [🔧 پیکربندی](#-پیکربندی)
- [🐛 مشکلات رایج](#-مشکلات-رایج)
- [💡 نکات و ترفندها](#-نکات-و-ترفندها)
- [🔒 امنیت](#-امنیت)
- [📊 عملکرد](#-عملکرد)
- [🤝 مشارکت](#-مشارکت)

---

## 🤖 سوالات عمومی

### ❓ **Covo Bot چیست؟**

Covo Bot یک بات تلگرام پیشرفته و چندمنظوره است که با زبان Go نوشته شده و قابلیت‌های متنوعی از جمله:

- **هوش مصنوعی** - پرسش و پاسخ هوشمند با DeepSeek
- **بازی‌های تعاملی** - جرات یا حقیقت، چلنج روزانه، فال حافظ
- **مدیریت گروه** - حذف پیام، سکوت، بن، تگ
- **آمارگیری** - آمار پیام‌ها، کاربران فعال
- **امنیت** - فیلتر کلمات، محدودیت درخواست

### ❓ **آیا استفاده از بات رایگان است؟**

بله، استفاده از بات کاملاً رایگان است. فقط نیاز به:
- سرور برای اجرای بات
- توکن تلگرام (رایگان)
- کلید API DeepSeek (رایگان)

### ❓ **بات از چه زبان برنامه‌نویسی استفاده می‌کند؟**

بات با **Go 1.24.5** نوشته شده است که یک زبان برنامه‌نویسی مدرن، سریع و قابل اعتماد است.

### ❓ **آیا می‌توانم کد را تغییر دهم؟**

بله، این پروژه تحت مجوز MIT منتشر شده است و شما می‌توانید:
- کد را تغییر دهید
- از آن استفاده تجاری کنید
- آن را توزیع کنید
- نسخه‌های اصلاح شده ایجاد کنید

---

## ⚙️ نصب و راه‌اندازی

### ❓ **حداقل سیستم مورد نیاز چیست؟**

**سیستم عامل:**
- Linux (Ubuntu 20.04+, CentOS 7+, Debian 10+)
- macOS (10.15+)
- Windows (10+)

**سخت‌افزار:**
- RAM: حداقل 512MB، توصیه 1GB+
- CPU: 1 هسته، توصیه 2 هسته+
- فضای دیسک: 100MB

**نرم‌افزار:**
- Go 1.24.5+
- MySQL 5.7+
- Git

### ❓ **چگونه بات را نصب کنم؟**

```bash
# 1. کلون کردن پروژه
git clone https://github.com/your-username/covo-bot.git
cd covo-bot

# 2. نصب وابستگی‌ها
go mod tidy

# 3. تنظیم پایگاه داده
mysql -u root -p
CREATE DATABASE myappdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'covouser'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON myappdb.* TO 'covouser'@'localhost';
FLUSH PRIVILEGES;

# 4. تنظیم متغیرهای محیطی
cp .env.example .env
# فایل .env را ویرایش کنید

# 5. اجرای بات
go run main.go
```

### ❓ **چگونه توکن تلگرام دریافت کنم؟**

1. به [@BotFather](https://t.me/botfather) در تلگرام پیام دهید
2. دستور `/newbot` را ارسال کنید
3. نام بات را وارد کنید
4. نام کاربری بات را وارد کنید
5. توکن را کپی کنید

### ❓ **چگونه کلید API DeepSeek دریافت کنم؟**

1. به [OpenRouter](https://openrouter.ai/) بروید
2. حساب کاربری ایجاد کنید
3. به بخش API Keys بروید
4. کلید جدید ایجاد کنید

---

## 🔧 پیکربندی

### ❓ **چگونه تنظیمات بات را تغییر دهم؟**

تنظیمات در فایل `.env` قرار دارد:

```env
# تنظیمات تلگرام
TELEGRAM_TOKEN=your_bot_token

# تنظیمات هوش مصنوعی
DEEPSEEK_TOKEN=your_api_key

# تنظیمات پایگاه داده
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=covouser
MYSQL_PASSWORD=your_password
MYSQL_DATABASE=myappdb

# تنظیمات محدودیت
MAX_REQUESTS_PER_DAY=1000
COOLDOWN_SECONDS=5
```

### ❓ **چگونه محدودیت درخواست‌ها را تغییر دهم؟**

در فایل `.env`:
```env
MAX_REQUESTS_PER_DAY=500    # حداکثر درخواست روزانه
COOLDOWN_SECONDS=10         # فاصله زمانی بین درخواست‌ها
```

### ❓ **چگونه قابلیت‌ها را فعال/غیرفعال کنم؟**

از دستور **پنل** در گروه استفاده کنید:
1. `پنل` را تایپ کنید
2. روی **قابلیت‌ها** کلیک کنید
3. قابلیت مورد نظر را فعال/غیرفعال کنید

---

## 🐛 مشکلات رایج

### ❓ **خطای "Error connecting to MySQL" دریافت می‌کنم**

**علت:** MySQL در حال اجرا نیست یا تنظیمات اشتباه است.

**راه حل:**
```bash
# بررسی وضعیت MySQL
sudo systemctl status mysql

# شروع MySQL
sudo systemctl start mysql

# بررسی تنظیمات
mysql -u covouser -p myappdb
```

### ❓ **بات پیام‌ها را دریافت نمی‌کند**

**علت:** توکن تلگرام اشتباه است یا بات در گروه‌ها دسترسی ندارد.

**راه حل:**
```bash
# بررسی توکن
curl "https://api.telegram.org/bot<YOUR_TOKEN>/getMe"

# بررسی دسترسی‌ها
curl "https://api.telegram.org/bot<YOUR_TOKEN>/getChatMember?chat_id=<CHAT_ID>&user_id=<BOT_USER_ID>"
```

### ❓ **خطای "401 Unauthorized" از DeepSeek API**

**علت:** کلید API اشتباه است یا منقضی شده.

**راه حل:**
```bash
# بررسی کلید API
curl -H "Authorization: Bearer <YOUR_API_KEY>" \
     -H "Content-Type: application/json" \
     https://openrouter.ai/api/v1/models
```

### ❓ **بات کند پاسخ می‌دهد**

**علت:** محدودیت سرور یا پایگاه داده.

**راه حل:**
```bash
# بررسی استفاده از منابع
top -p $(pgrep covo-bot)

# بهینه‌سازی پایگاه داده
mysql -u root -p
OPTIMIZE TABLE group_messages;
OPTIMIZE TABLE user_usage;
```

### ❓ **پیام‌های گروه ثبت نمی‌شوند**

**علت:** بات در گروه دسترسی ندارد یا قابلیت آمار غیرفعال است.

**راه حل:**
1. مطمئن شوید بات در گروه عضو است
2. از دستور **پنل** → **قابلیت‌ها** → **آمار پیام** را فعال کنید
3. بات را ادمین کنید

---

## 💡 نکات و ترفندها

### ❓ **چگونه بات را ادمین کنم؟**

1. به گروه بروید
2. روی نام بات کلیک کنید
3. **Admin** را انتخاب کنید
4. مجوزهای لازم را اعطا کنید

### ❓ **چگونه لاگ‌ها را ببینم؟**

```bash
# اجرا با لاگ‌های تفصیلی
go run main.go 2>&1 | tee bot.log

# مانیتورینگ real-time
tail -f bot.log

# جستجو در لاگ‌ها
grep "ERROR" bot.log
```

### ❓ **چگونه بات را در پس‌زمینه اجرا کنم؟**

```bash
# با nohup
nohup go run main.go > bot.log 2>&1 &

# با screen
screen -S covo-bot
go run main.go
# Ctrl+A, D برای خروج

# با systemd (Linux)
sudo systemctl start covo-bot
```

### ❓ **چگونه از بات بکاپ بگیرم؟**

```bash
# بکاپ از پایگاه داده
mysqldump -u covouser -p myappdb > backup_$(date +%Y%m%d_%H%M%S).sql

# بکاپ از کد
tar -czf covo-bot-backup-$(date +%Y%m%d_%H%M%S).tar.gz .

# بکاپ از تنظیمات
cp .env .env.backup
```

---

## 🔒 امنیت

### ❓ **چگونه امنیت بات را افزایش دهم؟**

1. **رمزهای قوی استفاده کنید:**
```env
MYSQL_PASSWORD=YourVeryStrongPassword123!
```

2. **فایل .env را محافظت کنید:**
```bash
chmod 600 .env
```

3. **فایروال تنظیم کنید:**
```bash
# فقط پورت‌های لازم را باز کنید
sudo ufw allow 3306  # MySQL
sudo ufw allow 22    # SSH
```

4. **به‌روزرسانی‌های امنیتی:**
```bash
# به‌روزرسانی سیستم
sudo apt update && sudo apt upgrade

# به‌روزرسانی وابستگی‌ها
go get -u ./...
```

### ❓ **چگونه دسترسی ادمین را محدود کنم؟**

در فایل `commands/admin.go`:
```go
var adminUsers = map[int64]string{
    7853092812: "مهشید",
    990475046:  "هانتر",
    // فقط ID های مجاز را اضافه کنید
}
```

### ❓ **چگونه کلمات نامناسب را فیلتر کنم؟**

فایل `jsonfile/badwords.json` را ویرایش کنید:
```json
{
    "farsiWords": [
        "کلمه1", "کلمه2"
    ],
    "finglishWords": [
        "word1", "word2"
    ]
}
```

---

## 📊 عملکرد

### ❓ **چگونه عملکرد بات را بهبود دهم؟**

1. **تنظیمات Go:**
```bash
export GOMAXPROCS=4
export GOGC=100
export GOMEMLIMIT=2GiB
```

2. **بهینه‌سازی پایگاه داده:**
```sql
-- ایجاد ایندکس‌های اضافی
CREATE INDEX idx_group_messages_timestamp ON group_messages(timestamp);
CREATE INDEX idx_user_usage_last_request ON user_usage(last_request);
```

3. **تنظیمات MySQL:**
```ini
# /etc/mysql/mysql.conf.d/mysqld.cnf
[mysqld]
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
max_connections = 200
```

### ❓ **چگونه مانیتورینگ کنم؟**

```bash
# مانیتورینگ سیستم
htop

# مانیتورینگ پایگاه داده
mysql -u root -p
SHOW PROCESSLIST;
SHOW STATUS;

# مانیتورینگ لاگ‌ها
tail -f bot.log | grep -E "(ERROR|WARN)"
```

### ❓ **چگونه بات را مقیاس‌پذیر کنم؟**

1. **استفاده از Load Balancer**
2. **تقسیم پایگاه داده (Sharding)**
3. **استفاده از Redis برای Cache**
4. **استفاده از Docker Swarm یا Kubernetes**

---

## 🤝 مشارکت

### ❓ **چگونه در پروژه مشارکت کنم؟**

1. **Fork کنید**
2. **Branch جدید ایجاد کنید:**
```bash
git checkout -b feature/amazing-feature
```
3. **تغییرات را commit کنید:**
```bash
git commit -m 'Add amazing feature'
```
4. **Push کنید:**
```bash
git push origin feature/amazing-feature
```
5. **Pull Request ایجاد کنید**

### ❓ **چگونه باگ گزارش دهم؟**

1. به [GitHub Issues](https://github.com/your-username/covo-bot/issues) بروید
2. Issue جدید ایجاد کنید
3. توضیح کامل مشکل را بنویسید
4. مراحل تکرار باگ را ذکر کنید
5. لاگ‌های مربوطه را ضمیمه کنید

### ❓ **چگونه درخواست قابلیت جدید دهم؟**

1. به [GitHub Issues](https://github.com/your-username/covo-bot/issues) بروید
2. Issue جدید ایجاد کنید
3. برچسب "enhancement" اضافه کنید
4. توضیح کامل قابلیت را بنویسید
5. مزایای آن را ذکر کنید

---

## 📞 پشتیبانی

### ❓ **کجا می‌توانم کمک بگیرم؟**

- **GitHub Issues:** [ایجاد Issue](https://github.com/your-username/covo-bot/issues)
- **Email:** support@covo-bot.com
- **Telegram:** [@CovoBotSupport](https://t.me/CovoBotSupport)
- **Discord:** [سرور Discord](https://discord.gg/covo-bot)

### ❓ **چگونه مستندات را به‌روزرسانی کنم؟**

1. فایل‌های مستندات را در پوشه `docs/` ویرایش کنید
2. تغییرات را commit کنید
3. Pull Request ایجاد کنید

---

## 🔄 به‌روزرسانی

### ❓ **چگونه بات را به‌روزرسانی کنم؟**

```bash
# دریافت آخرین تغییرات
git pull origin main

# نصب وابستگی‌های جدید
go mod tidy

# کامپایل مجدد
go build -ldflags="-s -w" -o covo-bot main.go

# راه‌اندازی مجدد
sudo systemctl restart covo-bot
```

### ❓ **چگونه از به‌روزرسانی بکاپ بگیرم؟**

```bash
# بکاپ از پایگاه داده
mysqldump -u covouser -p myappdb > backup_before_update.sql

# بکاپ از کد
git tag v1.0.0
git push origin v1.0.0
```

---

<div align="center">

**❓ سوال دیگری دارید؟**

[🔝 بازگشت به بالا](#-سوالات-متداول---covo-bot)

</div>
