package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TruthDareCommand پیاده‌سازی بازی «جرات یا سوال +۱۸»
type TruthDareCommand struct {
	bot        *tgbotapi.BotAPI
	admin      *AdminCommand
	mu         sync.Mutex
	games      map[int64]*tdGame // per chat
	dareItems  []string
	truthItems []string
	loadOnce   sync.Once
}

type tdGame struct {
	chatID           int64
	starterUserID    int64
	isOpen           bool
	participants     []int64
	participantNames map[int64]string
	currentIndex     int
	activeUserID     int64
}

func NewTruthDareCommand(bot *tgbotapi.BotAPI, admin *AdminCommand) *TruthDareCommand {
	return &TruthDareCommand{
		bot:   bot,
		admin: admin,
		games: make(map[int64]*tdGame),
	}
}

// HandleStartWithoutSlash شروع بازی با متن «بازی» توسط ادمین
func (r *TruthDareCommand) HandleStartWithoutSlash(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	if !r.admin.IsAdmin(userID) {
		msg := tgbotapi.NewMessage(chatID, "❌ فقط ادمین می‌تواند بازی را شروع کند")
		return msg
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.games[chatID]; exists {
		msg := tgbotapi.NewMessage(chatID, "ℹ️ یک بازی در حال حاضر فعال است. برای بستن ثبت‌نام از دکمه استفاده کنید یا «توقف بازی» بزنید.")
		return msg
	}

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
	text := "🎮 بازی جرات یا سوال +۱۸ شروع شد!\nاگر می‌خوای شرکت کنی، روی دکمه زیر بزن.\n\nپس از پایان ثبت‌نام، ادمین می‌تونه بازی رو ببنده و شروع کنه."
	joinBtn := tgbotapi.NewInlineKeyboardButtonData("➕ جوین شو", fmt.Sprintf("td_join:%d", chatID))
	closeBtn := tgbotapi.NewInlineKeyboardButtonData("🔒 بستن بازی (فقط ادمین)", fmt.Sprintf("td_close:%d", chatID))
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(joinBtn),
		tgbotapi.NewInlineKeyboardRow(closeBtn),
	)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = kb
	return msg
}

// HandleStopWithoutSlash توقف بازی با متن «توقف بازی» توسط ادمین
func (r *TruthDareCommand) HandleStopWithoutSlash(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	if !r.admin.IsAdmin(userID) {
		return tgbotapi.NewMessage(chatID, "❌ فقط ادمین می‌تواند بازی را متوقف کند")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.games[chatID]; !ok {
		return tgbotapi.NewMessage(chatID, "ℹ️ در حال حاضر بازی فعالی وجود ندارد.")
	}
	delete(r.games, chatID)
	return tgbotapi.NewMessage(chatID, "🛑 بازی متوقف شد و اتاق بسته شد.")
}

// HandleCallback پردازش کال‌بک‌های اینلاین
func (r *TruthDareCommand) HandleCallback(update tgbotapi.Update) tgbotapi.CallbackConfig {
	cq := update.CallbackQuery
	data := cq.Data

	// join
	if strings.HasPrefix(data, "td_join:") {
		return r.handleJoin(update)
	}
	// close registration
	if strings.HasPrefix(data, "td_close:") {
		return r.handleClose(update)
	}
	// pick dare/truth
	if strings.HasPrefix(data, "td_pick:") {
		return r.handlePick(update)
	}
	// done answering -> next
	if strings.HasPrefix(data, "td_done:") {
		return r.handleDone(update)
	}

	// default ack
	return tgbotapi.NewCallback(cq.ID, "")
}

func (r *TruthDareCommand) handleJoin(update tgbotapi.Update) tgbotapi.CallbackConfig {
	cq := update.CallbackQuery
	chatID := cq.Message.Chat.ID
	userID := cq.From.ID
	username := cq.From.FirstName
	if cq.From.UserName != "" {
		username = "@" + cq.From.UserName
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.games[chatID]
	if !ok || !g.isOpen {
		return tgbotapi.NewCallback(cq.ID, "ثبت‌نام بسته است")
	}
	// prevent duplicates
	for _, uid := range g.participants {
		if uid == userID {
			return tgbotapi.NewCallback(cq.ID, "قبلاً جوین شدی!")
		}
	}
	g.participants = append(g.participants, userID)
	g.participantNames[userID] = username

	// update message with participant list
	names := make([]string, 0, len(g.participants))
	for _, uid := range g.participants {
		names = append(names, g.participantNames[uid])
	}
	newText := fmt.Sprintf("🎮 بازی جرات یا سوال +۱۸\nشرکت‌کنندگان (%d):\n%s\n\nبرای پیوستن دکمه را بزنید. ادمین می‌تواند ثبت‌نام را ببندد.", len(names), strings.Join(names, "\n"))
	joinBtn := tgbotapi.NewInlineKeyboardButtonData("➕ جوین شو", fmt.Sprintf("td_join:%d", chatID))
	closeBtn := tgbotapi.NewInlineKeyboardButtonData("🔒 بستن بازی (فقط ادمین)", fmt.Sprintf("td_close:%d", chatID))
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(joinBtn),
		tgbotapi.NewInlineKeyboardRow(closeBtn),
	)
	edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, cq.Message.MessageID, newText, kb)
	if _, err := r.bot.Request(edit); err != nil {
		log.Printf("td: failed to edit join message: %v", err)
	}

	return tgbotapi.NewCallback(cq.ID, "به بازی اضافه شدی ✅")
}

func (r *TruthDareCommand) handleClose(update tgbotapi.Update) tgbotapi.CallbackConfig {
	cq := update.CallbackQuery
	chatID := cq.Message.Chat.ID
	userID := cq.From.ID

	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.games[chatID]
	if !ok {
		return tgbotapi.NewCallback(cq.ID, "بازی‌ای پیدا نشد")
	}
	if userID != g.starterUserID {
		return tgbotapi.NewCallback(cq.ID, "فقط شروع‌کننده می‌تواند ببندد")
	}
	if !g.isOpen {
		return tgbotapi.NewCallback(cq.ID, "ثبت‌نام قبلاً بسته شده")
	}
	if len(g.participants) == 0 {
		return tgbotapi.NewCallback(cq.ID, "کسی جوین نشده")
	}
	g.isOpen = false
	g.currentIndex = 0
	g.activeUserID = g.participants[0]

	// announce start
	names := make([]string, 0, len(g.participants))
	for _, uid := range g.participants {
		names = append(names, g.participantNames[uid])
	}
	startText := fmt.Sprintf("🚀 بازی شروع شد!\nنوبت‌ها به ترتیب شرکت‌کنندگان است.\n\nترتیب: %s", strings.Join(names, "، "))
	msg := tgbotapi.NewMessage(chatID, startText)
	if _, err := r.bot.Send(msg); err != nil {
		log.Printf("td: failed to announce start: %v", err)
	}

	// prompt first player to choose
	r.promptPickLocked(g)
	return tgbotapi.NewCallback(cq.ID, "ثبت‌نام بسته شد و بازی شروع شد")
}

func (r *TruthDareCommand) promptPickLocked(g *tdGame) {
	// assumes r.mu locked
	r.ensureLoaded()
	currentID := g.activeUserID
	name := g.participantNames[currentID]
	text := fmt.Sprintf("نوبت %s هست. انتخاب کن:\n👉 جرات یا سوال +۱۸؟", name)

	dareBtn := tgbotapi.NewInlineKeyboardButtonData("🔥 جرات", fmt.Sprintf("td_pick:dare:%d", currentID))
	truthBtn := tgbotapi.NewInlineKeyboardButtonData("🫣 سوال +۱۸", fmt.Sprintf("td_pick:truth:%d", currentID))
	kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(dareBtn, truthBtn))
	msg := tgbotapi.NewMessage(g.chatID, text)
	msg.ReplyMarkup = kb
	if _, err := r.bot.Send(msg); err != nil {
		log.Printf("td: failed to prompt pick: %v", err)
	}
}

func (r *TruthDareCommand) handlePick(update tgbotapi.Update) tgbotapi.CallbackConfig {
	cq := update.CallbackQuery
	chatID := cq.Message.Chat.ID
	fromID := cq.From.ID
	parts := strings.Split(cq.Data, ":") // td_pick:<kind>:<expectedUserID>
	if len(parts) < 3 {
		return tgbotapi.NewCallback(cq.ID, "درخواست نامعتبر")
	}
	kind := parts[1]
	// parse expected user id (ignore error, compare string)
	expected := parts[2]
	if fmt.Sprint(fromID) != expected {
		return tgbotapi.NewCallback(cq.ID, "این نوبت شما نیست")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.games[chatID]
	if !ok || g.activeUserID != fromID {
		return tgbotapi.NewCallback(cq.ID, "این نوبت شما نیست")
	}

	r.ensureLoaded()
	var q string
	switch kind {
	case "dare":
		if len(r.dareItems) == 0 {
			return tgbotapi.NewCallback(cq.ID, "بانک جرات خالی است")
		}
		q = r.dareItems[rand.Intn(len(r.dareItems))]
	case "truth":
		if len(r.truthItems) == 0 {
			return tgbotapi.NewCallback(cq.ID, "بانک سوال خالی است")
		}
		q = r.truthItems[rand.Intn(len(r.truthItems))]
	default:
		return tgbotapi.NewCallback(cq.ID, "انتخاب نامعتبر")
	}

	// send question with "done" button
	question := fmt.Sprintf("❓ %s\n\n%s لطفاً پاسخ را به همین پیام ریپلای کن و پس از پاسخ، دکمه زیر را بزن.", q, r.displayName(g, fromID))
	doneBtn := tgbotapi.NewInlineKeyboardButtonData("✅ جواب دادم", fmt.Sprintf("td_done:%d", fromID))
	kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(doneBtn))
	msg := tgbotapi.NewMessage(chatID, question)
	msg.ReplyMarkup = kb
	if _, err := r.bot.Send(msg); err != nil {
		log.Printf("td: failed to send question: %v", err)
	}
	return tgbotapi.NewCallback(cq.ID, "سوال ارسال شد")
}

func (r *TruthDareCommand) handleDone(update tgbotapi.Update) tgbotapi.CallbackConfig {
	cq := update.CallbackQuery
	chatID := cq.Message.Chat.ID
	fromID := cq.From.ID
	parts := strings.Split(cq.Data, ":") // td_done:<expectedUserID>
	if len(parts) < 2 {
		return tgbotapi.NewCallback(cq.ID, "درخواست نامعتبر")
	}
	expected := parts[1]
	if fmt.Sprint(fromID) != expected {
		return tgbotapi.NewCallback(cq.ID, "این دکمه فقط برای فردی است که نوبت اوست")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.games[chatID]
	if !ok || g.activeUserID != fromID {
		return tgbotapi.NewCallback(cq.ID, "نوبت معتبر نیست")
	}

	// advance to next participant (circular)
	if len(g.participants) == 0 {
		return tgbotapi.NewCallback(cq.ID, "شرکت‌کننده‌ای وجود ندارد")
	}
	g.currentIndex = (g.currentIndex + 1) % len(g.participants)
	g.activeUserID = g.participants[g.currentIndex]

	// prompt next
	r.promptPickLocked(g)
	return tgbotapi.NewCallback(cq.ID, "نوبت بعدی")
}

func (r *TruthDareCommand) ensureLoaded() {
	r.loadOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
		// load dare.json
		darePath := filepath.Join("jsonfile", "dare.json")
		if data, err := os.ReadFile(darePath); err == nil {
			var arr []struct {
				ID   int    `json:"id"`
				Dare string `json:"dare"`
			}
			if err := json.Unmarshal(data, &arr); err == nil {
				for _, it := range arr {
					s := strings.TrimSpace(it.Dare)
					if s != "" {
						r.dareItems = append(r.dareItems, s)
					}
				}
			} else {
				log.Printf("td: cannot parse dare.json: %v", err)
			}
		} else {
			log.Printf("td: cannot read dare.json: %v", err)
		}

		// load truth+18.json
		truthPath := filepath.Join("jsonfile", "truth+18.json")
		if data, err := os.ReadFile(truthPath); err == nil {
			var arr []struct {
				ID       int    `json:"id"`
				Question string `json:"question"`
			}
			if err := json.Unmarshal(data, &arr); err == nil {
				for _, it := range arr {
					s := strings.TrimSpace(it.Question)
					if s != "" {
						r.truthItems = append(r.truthItems, s)
					}
				}
			} else {
				log.Printf("td: cannot parse truth+18.json: %v", err)
			}
		} else {
			log.Printf("td: cannot read truth+18.json: %v", err)
		}
	})
}

func (r *TruthDareCommand) displayName(g *tdGame, userID int64) string {
	if name, ok := g.participantNames[userID]; ok && name != "" {
		return name
	}
	return fmt.Sprintf("کاربر %d", userID)
}
