package commands

import (
	"fmt"
	"html"
	"strings"

	"redhat-bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TagCommand struct {
	bot     *tgbotapi.BotAPI
	storage *storage.MySQLStorage
}

func NewTagCommand(bot *tgbotapi.BotAPI, storage *storage.MySQLStorage) *TagCommand {
	return &TagCommand{bot: bot, storage: storage}
}

// HandleTagAllOnReply tags all known group members when user replies a message and sends "تگ" (without slash)
// Sends messages in chunks to avoid Telegram limits. Only admins can use it.
func (t *TagCommand) HandleTagAllOnReply(update tgbotapi.Update) tgbotapi.MessageConfig {
	chat := update.Message.Chat
	chatID := chat.ID

	if chat.Type != "group" && chat.Type != "supergroup" {
		return tgbotapi.NewMessage(chatID, "❌ این قابلیت فقط در گروه‌ها قابل استفاده است")
	}

	if update.Message.ReplyToMessage == nil {
		return tgbotapi.NewMessage(chatID, "لطفاً روی یک پیام ریپلای کنید و بنویسید: تگ")
	}

	// Only group admins can use
	isAdmin, err := t.isUserAdmin(chatID, update.Message.From.ID)
	if err != nil {
		return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی دسترسی ادمین")
	}
	if !isAdmin {
		return tgbotapi.NewMessage(chatID, "❌ فقط ادمین‌های گروه می‌توانند همه را تگ کنند")
	}

	// Load members
	members, err := t.storage.GetGroupMembers(chatID)
	if err != nil || len(members) == 0 {
		return tgbotapi.NewMessage(chatID, "❌ لیست اعضای گروه پیدا نشد")
	}

	// Build and send in chunks
	const chunkSize = 20
	replyTo := update.Message.ReplyToMessage.MessageID

	var batch []string
	flush := func() {
		if len(batch) == 0 {
			return
		}
		text := strings.Join(batch, " \u2063") // invisible separator to avoid formatting merges
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyToMessageID = replyTo
		_, _ = t.bot.Send(msg)
		batch = batch[:0]
	}

	for _, m := range members {
		// Prefer saved name; sanitize
		displayName := strings.TrimSpace(m.Name)
		if displayName == "" {
			displayName = fmt.Sprintf("User %d", m.UserID)
		}
		// Build HTML mention using tg://user?id
		escaped := html.EscapeString(displayName)
		mention := fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", m.UserID, escaped)
		batch = append(batch, mention)
		if len(batch) >= chunkSize {
			flush()
		}
	}
	flush()

	// Try to delete the command message for cleanliness (ignore error)
	_, _ = t.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: update.Message.MessageID})

	// Return empty config; we already sent messages
	return tgbotapi.MessageConfig{}
}

func (t *TagCommand) isUserAdmin(chatID int64, userID int64) (bool, error) {
	cfg := tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{ChatID: chatID, UserID: userID}}
	member, err := t.bot.GetChatMember(cfg)
	if err != nil {
		return false, err
	}
	return member.IsAdministrator() || member.IsCreator(), nil
}
