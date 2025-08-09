package commands

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ModerationCommand struct {
	bot *tgbotapi.BotAPI
}

func NewModerationCommand(bot *tgbotapi.BotAPI) *ModerationCommand {
	return &ModerationCommand{bot: bot}
}

// HandleDelete deletes a replied message if the requester is a group admin
func (m *ModerationCommand) HandleDelete(update tgbotapi.Update) tgbotapi.MessageConfig {
	chat := update.Message.Chat
	chatID := chat.ID

	// Only in groups
	if chat.Type != "group" && chat.Type != "supergroup" {
		return tgbotapi.NewMessage(chatID, "❌ این دستور فقط در گروه‌ها قابل استفاده است")
	}

	// Require reply
	if update.Message.ReplyToMessage == nil {
		return tgbotapi.NewMessage(chatID, "لطفاً روی یک پیام ریپلای کنید تا حذف شود")
	}

	// Only group admins can use
	isAdmin, err := m.isUserAdmin(chatID, update.Message.From.ID)
	if err != nil {
		log.Printf("getChatMember error: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی دسترسی ادمین")
	}
	if !isAdmin {
		return tgbotapi.NewMessage(chatID, "❌ فقط ادمین‌های گروه می‌توانند پیام حذف کنند")
	}

	// Delete the replied message
	targetMsgID := update.Message.ReplyToMessage.MessageID
	if _, err := m.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: targetMsgID}); err != nil {
		log.Printf("deleteMessage target error: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ حذف پیام انجام نشد. مطمئن شوید ربات دسترسی حذف دارد و پیام خیلی قدیمی نیست")
	}

	// Try to delete the command message for cleanliness (ignore error)
	_, _ = m.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: update.Message.MessageID})

	// No further message to send
	return tgbotapi.MessageConfig{}
}

func (m *ModerationCommand) isUserAdmin(chatID int64, userID int64) (bool, error) {
	cfg := tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{ChatID: chatID, UserID: userID}}
	member, err := m.bot.GetChatMember(cfg)
	if err != nil {
		return false, err
	}
	return member.IsAdministrator() || member.IsCreator(), nil
}
