# 🚀 راهنمای نصب و پیکربندی - Covo Bot

## 📋 فهرست مطالب

- [🎯 پیش‌نیازها](#-پیش‌نیازها)
- [⚙️ نصب و راه‌اندازی](#️-نصب-و-راه‌اندازی)
- [🔧 پیکربندی](#-پیکربندی)
- [🗄️ تنظیم پایگاه داده](#️-تنظیم-پایگاه-داده)
- [🤖 ایجاد بات تلگرام](#-ایجاد-بات-تلگرام)
- [🧠 تنظیم هوش مصنوعی](#-تنظیم-هوش-مصنوعی)
- [🚀 اجرای بات](#-اجرای-بات)
- [🔍 عیب‌یابی](#-عیب‌یابی)
- [📊 مانیتورینگ](#-مانیتورینگ)

---

## 🎯 پیش‌نیازها

### 💻 **سیستم عامل**
- **Linux** (Ubuntu 20.04+, CentOS 7+, Debian 10+)
- **macOS** (10.15+)
- **Windows** (10+)

### 🔧 **نرم‌افزارهای مورد نیاز**

#### **Go**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# CentOS/RHEL
sudo yum install golang

# macOS
brew install go

# Windows
# دانلود از https://golang.org/dl/
```

**بررسی نصب:**
```bash
go version
# باید Go 1.24.5 یا بالاتر نمایش دهد
```

#### **MySQL**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server

# CentOS/RHEL
sudo yum install mysql-server

# macOS
brew install mysql

# Windows
# دانلود از https://dev.mysql.com/downloads/mysql/
```

**بررسی نصب:**
```bash
mysql --version
# باید MySQL 5.7 یا بالاتر نمایش دهد
```

#### **Git**
```bash
# Ubuntu/Debian
sudo apt install git

# CentOS/RHEL
sudo yum install git

# macOS
brew install git

# Windows
# دانلود از https://git-scm.com/download/win
```

---

## ⚙️ نصب و راه‌اندازی

### 1️⃣ **کلون کردن پروژه**

```bash
# کلون کردن از GitHub
git clone https://github.com/your-username/covo-bot.git
cd covo-bot

# یا اگر از ZIP استفاده می‌کنید
wget https://github.com/your-username/covo-bot/archive/main.zip
unzip main.zip
cd covo-bot-main
```

### 2️⃣ **نصب وابستگی‌ها**

```bash
# نصب وابستگی‌های Go
go mod tidy

# بررسی وابستگی‌ها
go mod verify
```

### 3️⃣ **ساخت پروژه**

```bash
# کامپایل پروژه
go build -o covo-bot main.go

# یا برای production
go build -ldflags="-s -w" -o covo-bot main.go
```

---

## 🔧 پیکربندی

### 📄 **ایجاد فایل تنظیمات**

```bash
# ایجاد فایل .env
cp .env.example .env
```

### 🔐 **تنظیم متغیرهای محیطی**

```env
# فایل .env
# ===========================================
# تنظیمات تلگرام
# ===========================================
TELEGRAM_TOKEN=your_telegram_bot_token_here

# ===========================================
# تنظیمات هوش مصنوعی
# ===========================================
DEEPSEEK_TOKEN=your_deepseek_api_key_here

# ===========================================
# تنظیمات پایگاه داده
# ===========================================
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=covouser
MYSQL_PASSWORD=your_secure_password_here
MYSQL_DATABASE=myappdb

# ===========================================
# تنظیمات محدودیت
# ===========================================
MAX_REQUESTS_PER_DAY=1000
COOLDOWN_SECONDS=5
```

### 🔒 **امنیت فایل .env**

```bash
# تنظیم مجوزهای امنیتی
chmod 600 .env

# اطمینان از عدم commit شدن
echo ".env" >> .gitignore
```

---

## 🗄️ تنظیم پایگاه داده

### 1️⃣ **راه‌اندازی MySQL**

```bash
# شروع سرویس MySQL
sudo systemctl start mysql
sudo systemctl enable mysql

# یا در macOS
brew services start mysql
```

### 2️⃣ **ایجاد پایگاه داده**

```sql
-- اتصال به MySQL
mysql -u root -p

-- ایجاد پایگاه داده
CREATE DATABASE myappdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- ایجاد کاربر
CREATE USER 'covouser'@'localhost' IDENTIFIED BY 'your_secure_password_here';

-- اعطای مجوزها
GRANT ALL PRIVILEGES ON myappdb.* TO 'covouser'@'localhost';
FLUSH PRIVILEGES;

-- خروج
EXIT;
```

### 3️⃣ **تست اتصال**

```bash
# تست اتصال با کاربر جدید
mysql -u covouser -p myappdb

# در صورت موفقیت، خروج
EXIT;
```

### 4️⃣ **مایگریشن خودکار**

```bash
# اجرای بات برای ایجاد جداول
./covo-bot

# یا در حالت توسعه
go run main.go
```

**جداول ایجاد شده:**
- `user_usage` - آمار استفاده کاربران
- `group_messages` - پیام‌های گروه‌ها
- `group_members` - اعضای گروه‌ها
- `feature_settings` - تنظیمات قابلیت‌ها
- `daily_challenges` - چالش‌های روزانه
- `bot_channels` - کانال‌های ربات
- `required_channels` - کانال‌های الزامی
- `user_onboarding` - پیگیری عضویت

---

## 🤖 ایجاد بات تلگرام

### 1️⃣ **ایجاد بات جدید**

1. به [@BotFather](https://t.me/botfather) در تلگرام پیام دهید
2. دستور `/newbot` را ارسال کنید
3. نام بات را وارد کنید (مثال: `Covo Bot`)
4. نام کاربری بات را وارد کنید (مثال: `covo_ai_bot`)
5. توکن بات را کپی کنید

### 2️⃣ **تنظیم دستورات بات**

```bash
# ارسال دستورات به BotFather
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setMyCommands" \
-H "Content-Type: application/json" \
-d '{
  "commands": [
    {"command": "start", "description": "شروع کار با بات"},
    {"command": "help", "description": "راهنمای استفاده"},
    {"command": "covo", "description": "پرسش و پاسخ هوشمند"},
    {"command": "cj", "description": "تولید جوک"},
    {"command": "music", "description": "پیشنهاد موسیقی"},
    {"command": "crs", "description": "وضعیت بات"}
  ]
}'
```

### 3️⃣ **تنظیم توکن در .env**

```env
TELEGRAM_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
```

---

## 🧠 تنظیم هوش مصنوعی

### 1️⃣ **ایجاد حساب کاربری OpenRouter**

1. به [OpenRouter](https://openrouter.ai/) بروید
2. حساب کاربری ایجاد کنید
3. به بخش API Keys بروید
4. کلید API جدید ایجاد کنید

### 2️⃣ **تنظیم کلید در .env**

```env
DEEPSEEK_TOKEN=sk-or-v1-your-api-key-here
```

### 3️⃣ **تست اتصال**

```bash
# اجرای تست
go run main.go

# در لاگ‌ها باید پیام موفقیت ببینید
# "🤖 بات کوو در حال راه‌اندازی است..."
```

---

## 🚀 اجرای بات

### 🖥️ **حالت توسعه**

```bash
# اجرای مستقیم
go run main.go

# یا با لاگ‌های تفصیلی
go run main.go 2>&1 | tee bot.log
```

### 🏭 **حالت Production**

#### **روش 1: اجرای مستقیم**

```bash
# کامپایل
go build -ldflags="-s -w" -o covo-bot main.go

# اجرا
./covo-bot
```

#### **روش 2: Systemd Service (Linux)**

```bash
# ایجاد فایل سرویس
sudo nano /etc/systemd/system/covo-bot.service
```

```ini
[Unit]
Description=Covo Bot
After=network.target mysql.service

[Service]
Type=simple
User=covouser
WorkingDirectory=/path/to/covo-bot
ExecStart=/path/to/covo-bot/covo-bot
Restart=always
RestartSec=5
Environment=GOMAXPROCS=2

[Install]
WantedBy=multi-user.target
```

```bash
# فعال‌سازی سرویس
sudo systemctl daemon-reload
sudo systemctl enable covo-bot
sudo systemctl start covo-bot

# بررسی وضعیت
sudo systemctl status covo-bot
```

#### **روش 3: Docker**

```dockerfile
# Dockerfile
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o covo-bot main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/covo-bot .
COPY --from=builder /app/jsonfile ./jsonfile
CMD ["./covo-bot"]
```

```bash
# ساخت و اجرای Docker
docker build -t covo-bot .
docker run -d --name covo-bot --env-file .env covo-bot
```

---

## 🔍 عیب‌یابی

### 🐛 **مشکلات رایج**

#### **خطای اتصال به پایگاه داده**
```
Error connecting to MySQL: dial tcp [::1]:3306: connect: connection refused
```

**راه حل:**
```bash
# بررسی وضعیت MySQL
sudo systemctl status mysql

# شروع MySQL
sudo systemctl start mysql

# بررسی تنظیمات
mysql -u covouser -p -e "SELECT 1"
```

#### **خطای Telegram API**
```
Error sending message: 400 Bad Request
```

**راه حل:**
```bash
# بررسی توکن
curl "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getMe"

# بررسی دسترسی‌ها
curl "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getChatMember?chat_id=<CHAT_ID>&user_id=<BOT_USER_ID>"
```

#### **خطای DeepSeek API**
```
Error making request: 401 Unauthorized
```

**راه حل:**
```bash
# بررسی کلید API
curl -H "Authorization: Bearer <YOUR_API_KEY>" \
     -H "Content-Type: application/json" \
     https://openrouter.ai/api/v1/models
```

### 📝 **لاگ‌ها**

#### **فعال‌سازی لاگ‌های تفصیلی**

```go
// در main.go
log.SetFlags(log.LstdFlags | log.Lshortfile)
log.SetOutput(os.Stdout)
```

#### **مانیتورینگ لاگ‌ها**

```bash
# مانیتورینگ real-time
tail -f bot.log

# جستجو در لاگ‌ها
grep "ERROR" bot.log
grep "WARN" bot.log
```

---

## 📊 مانیتورینگ

### 📈 **آمار عملکرد**

#### **مانیتورینگ پایگاه داده**

```sql
-- آمار استفاده کاربران
SELECT user_id, requests_today, last_request 
FROM user_usage 
ORDER BY requests_today DESC 
LIMIT 10;

-- آمار پیام‌های گروه‌ها
SELECT group_id, COUNT(*) as message_count 
FROM group_messages 
WHERE timestamp > DATE_SUB(NOW(), INTERVAL 24 HOUR)
GROUP BY group_id 
ORDER BY message_count DESC;

-- آمار قابلیت‌های فعال
SELECT feature_name, COUNT(*) as enabled_count 
FROM feature_settings 
WHERE enabled = true 
GROUP BY feature_name;
```

#### **مانیتورینگ سیستم**

```bash
# بررسی استفاده از CPU و RAM
top -p $(pgrep covo-bot)

# بررسی استفاده از دیسک
df -h

# بررسی اتصالات شبکه
netstat -tulpn | grep covo-bot
```

### 🔔 **هشدارها**

#### **تنظیم هشدار برای خطاها**

```bash
# اسکریپت مانیتورینگ
#!/bin/bash
# monitor.sh

LOG_FILE="bot.log"
ERROR_COUNT=$(grep -c "ERROR" $LOG_FILE)

if [ $ERROR_COUNT -gt 10 ]; then
    echo "High error count detected: $ERROR_COUNT"
    # ارسال هشدار (ایمیل، تلگرام، etc.)
fi
```

#### **Cron Job برای مانیتورینگ**

```bash
# اضافه کردن به crontab
crontab -e

# هر 5 دقیقه چک کردن
*/5 * * * * /path/to/monitor.sh
```

---

## 🔧 بهینه‌سازی

### ⚡ **بهینه‌سازی عملکرد**

#### **تنظیمات Go**

```bash
# تنظیم متغیرهای محیطی
export GOMAXPROCS=4
export GOGC=100
export GOMEMLIMIT=2GiB
```

#### **بهینه‌سازی پایگاه داده**

```sql
-- ایجاد ایندکس‌های اضافی
CREATE INDEX idx_group_messages_timestamp ON group_messages(timestamp);
CREATE INDEX idx_user_usage_last_request ON user_usage(last_request);

-- بهینه‌سازی جداول
OPTIMIZE TABLE group_messages;
OPTIMIZE TABLE user_usage;
```

#### **تنظیمات MySQL**

```ini
# /etc/mysql/mysql.conf.d/mysqld.cnf
[mysqld]
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
max_connections = 200
query_cache_size = 64M
query_cache_type = 1
```

---

## 🔄 به‌روزرسانی

### 📦 **به‌روزرسانی کد**

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

### 🗄️ **به‌روزرسانی پایگاه داده**

```bash
# بکاپ از پایگاه داده
mysqldump -u covouser -p myappdb > backup_$(date +%Y%m%d_%H%M%S).sql

# اجرای مایگریشن‌های جدید
./covo-bot
```

---

## 🆘 پشتیبانی

### 📞 **راه‌های ارتباطی**

- **GitHub Issues:** [ایجاد Issue](https://github.com/your-username/covo-bot/issues)
- **Email:** support@covo-bot.com
- **Telegram:** [@CovoBotSupport](https://t.me/CovoBotSupport)

### 📚 **منابع مفید**

- [مستندات کامل API](API_DOCUMENTATION.md)
- [راهنمای توسعه](DEVELOPMENT_GUIDE.md)
- [FAQ](FAQ.md)

---

<div align="center">

**🚀 موفق باشید!**

[🔝 بازگشت به بالا](#-راهنمای-نصب-و-پیکربندی---covo-bot)

</div>
