package scheduler

// import (
// 	"fmt"
// 	"log"
// 	"strings"
// 	"time"
// 	"redhat-bot/ai"
// 	"redhat-bot/storage"
// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// type DailySummaryScheduler struct {
// 	bot       *tgbotapi.BotAPI
// 	storage   *storage.MemoryStorage
// 	aiClient  *ai.DeepSeekClient
// 	groups    map[int64]string // groupID -> group name
// }

// func NewDailySummaryScheduler(bot *tgbotapi.BotAPI, storage *storage.MemoryStorage, aiClient *ai.DeepSeekClient) *DailySummaryScheduler {
// 	return &DailySummaryScheduler{
// 		bot:      bot,
// 		storage:  storage,
// 		aiClient: aiClient,
// 		groups:   make(map[int64]string),
// 	}
// }

// func (d *DailySummaryScheduler) AddGroup(groupID int64, groupName string) {
// 	d.groups[groupID] = groupName
// }

// func (d *DailySummaryScheduler) RunDailySummary() {
// 	log.Println("🔄 Starting daily group analysis...")

// 	for groupID, groupName := range d.groups {
// 		d.analyzeAndPostSummary(groupID, groupName)
// 		// Small delay between groups to avoid rate limiting
// 		time.Sleep(2 * time.Second)
// 	}

// 	log.Println("✅ Daily group analysis completed!")
// }

// func (d *DailySummaryScheduler) analyzeAndPostSummary(groupID int64, groupName string) {
// 	// Get messages from the last 24 hours
// 	messages := d.storage.GetGroupMessages(groupID)

// 	if len(messages) == 0 {
// 		log.Printf("No messages found for group %s (%d)", groupName, groupID)
// 		return
// 	}

// 	// Prepare messages for AI analysis
// 	var messageTexts []string
// 	for _, msg := range messages {
// 		// Skip bot commands and very short messages
// 		if len(msg.Message) < 3 || strings.HasPrefix(msg.Message, "/") {
// 			continue
// 		}

// 		formattedMsg := fmt.Sprintf("%s: %s", msg.Username, msg.Message)
// 		messageTexts = append(messageTexts, formattedMsg)
// 	}

// 	if len(messageTexts) == 0 {
// 		log.Printf("No valid messages found for group %s (%d)", groupName, groupID)
// 		return
// 	}

// 	// Generate AI summary
// 	summary, err := d.aiClient.SummarizeGroupChat(messageTexts)
// 	if err != nil {
// 		log.Printf("Error generating summary for group %s (%d): %v", groupName, groupID, err)
// 		return
// 	}

// 	// Format and post summary
// 	formattedSummary := fmt.Sprintf("🧠 *Daily Group Summary*\n\n"+
// 		"📅 *Date:* %s\n"+
// 		"👥 *Group:* %s\n"+
// 		"💬 *Messages Analyzed:* %d\n\n"+
// 		"%s\n\n"+
// 		"Powered by Redhat Bot 🤖",
// 		time.Now().Format("January 2, 2006"),
// 		groupName,
// 		len(messageTexts),
// 		summary)

// 	msg := tgbotapi.NewMessage(groupID, formattedSummary)
// 	msg.ParseMode = tgbotapi.ModeMarkdown

// 	_, err = d.bot.Send(msg)
// 	if err != nil {
// 		log.Printf("Error posting summary to group %s (%d): %v", groupName, groupID, err)
// 		return
// 	}

// 	log.Printf("✅ Posted daily summary to group %s (%d)", groupName, groupID)

// 	// Clear old messages after posting summary
// 	d.storage.ClearGroupMessages(groupID)
// }
