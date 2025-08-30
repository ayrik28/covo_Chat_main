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

	"redhat-bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DailyChallengeCommand struct {
	storage *storage.MySQLStorage
	bot     *tgbotapi.BotAPI
}

func NewDailyChallengeCommand(storage *storage.MySQLStorage, bot *tgbotapi.BotAPI) *DailyChallengeCommand {
	return &DailyChallengeCommand{storage: storage, bot: bot}
}

// ---------- Load proverbs (zarb.json) ----------

type emojifiedProverb struct {
	ID      int    `json:"id"`
	Proverb string `json:"proverb"`
	Emojis  string `json:"emojis"`
}

type zarbFile struct {
	Items []emojifiedProverb `json:"emojified_proverbs"`
}

var (
	loadZarbOnce    sync.Once
	cachedZarbItems []emojifiedProverb
)

func loadZarb() {
	loadZarbOnce.Do(func() {
		filePath := filepath.Join("jsonfile", "zarb.json")
		f, err := os.Open(filePath)
		if err != nil {
			log.Printf("cannot open zarb file: %v", err)
			return
		}
		defer f.Close()

		var z zarbFile
		if err := json.NewDecoder(f).Decode(&z); err != nil {
			log.Printf("cannot decode zarb file: %v", err)
			return
		}
		cachedZarbItems = z.Items
	})
}

func getRandomZarb() (emojis string, proverb string, ok bool) {
	loadZarb()
	if len(cachedZarbItems) == 0 {
		return "", "", false
	}
	rand.Seed(time.Now().UnixNano())
	it := cachedZarbItems[rand.Intn(len(cachedZarbItems))]
	return it.Emojis, it.Proverb, true
}

// ---------- Posting daily challenge ----------

func (d *DailyChallengeCommand) PostDailyChallenge(groupID int64) {
	emojis, proverb, ok := getRandomZarb()
	if !ok {
		log.Printf("daily challenge: zarb list is empty")
		return
	}

	text := fmt.Sprintf("🧩 چلنج روزانه\n\n%s\n\nاز روی ایموجی ضرب‌المثل را حدس بزنید و روی همین پیام ریپلای کنید.\nاولین پاسخ صحیح لقب «باهوش‌ترین فرد گروه» را می‌گیرد!", emojis)
	msg := tgbotapi.NewMessage(groupID, text)
	sent, err := d.bot.Send(msg)
	if err != nil {
		log.Printf("daily challenge: send error: %v", err)
		return
	}

	if err := d.storage.CreateDailyChallenge(groupID, sent.MessageID, proverb, emojis); err != nil {
		log.Printf("daily challenge: save state error: %v", err)
	}
}

// RunDailyForEnabledGroups posts the daily challenge to all enabled groups
func (d *DailyChallengeCommand) RunDailyForEnabledGroups() {
	groups, err := d.storage.GetEnabledGroupsForFeature("daily_challenge")
	if err != nil {
		log.Printf("daily challenge: cannot list enabled groups: %v", err)
		return
	}
	for _, gid := range groups {
		d.PostDailyChallenge(gid)
		time.Sleep(2 * time.Second)
	}
}

// ---------- Handling answers ----------

func normalizePersian(s string) string {
	if s == "" {
		return ""
	}
	s = strings.TrimSpace(s)
	// unify Arabic and Persian letters and spacing
	replacements := []struct{ old, new string }{
		{"ي", "ی"}, {"ك", "ک"}, {"ۀ", "ه"}, {"ة", "ه"},
		{"\u0640", ""},  // kashida
		{"\u200c", " "}, // ZWNJ -> space
	}
	for _, r := range replacements {
		s = strings.ReplaceAll(s, r.old, r.new)
	}
	// collapse spaces
	s = strings.Join(strings.Fields(s), " ")
	return s
}

// HandleAnswer checks if a message is a reply to the latest active challenge and, if correct, announces the winner
func (d *DailyChallengeCommand) HandleAnswer(update tgbotapi.Update) tgbotapi.MessageConfig {
	empty := tgbotapi.MessageConfig{}
	if update.Message == nil || update.Message.ReplyToMessage == nil {
		return empty
	}
	chatID := update.Message.Chat.ID
	replyToID := update.Message.ReplyToMessage.MessageID

	challenge, err := d.storage.GetActiveChallengeForGroup(chatID)
	if err != nil || challenge == nil {
		return empty
	}
	if challenge.MessageID != replyToID || challenge.Answered {
		return empty
	}

	answer := normalizePersian(update.Message.Text)
	correct := normalizePersian(challenge.Proverb)
	if answer == "" {
		return empty
	}

	// accept exact or containing match
	if !(answer == correct || strings.Contains(answer, correct)) {
		return empty
	}

	user := update.Message.From
	winnerName := user.FirstName
	if user.LastName != "" {
		winnerName = strings.TrimSpace(winnerName + " " + user.LastName)
	}

	// try to mark answered atomically
	ok, err := d.storage.TryMarkChallengeAnswered(challenge.ID, user.ID, winnerName)
	if err != nil || !ok {
		return empty
	}

	text := fmt.Sprintf("🎉 %s اولین نفر بود که ضرب‌المثل را درست حدس زد!\n\n✅ پاسخ صحیح: «%s»\n🧠 لقب امروز: باهوش‌ترین فرد گروه 👑", winnerName, challenge.Proverb)
	msg := tgbotapi.NewMessage(chatID, text)
	return msg
}
