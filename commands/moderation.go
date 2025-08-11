package commands

import (
	"log"
	"strconv"
	"strings"
	"time"

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

// Handle processes both "/del [n]" and "حذف [n]". If n>0 deletes n previous messages; otherwise deletes replied message.
func (m *ModerationCommand) Handle(update tgbotapi.Update) tgbotapi.MessageConfig {
	chat := update.Message.Chat
	chatID := chat.ID

	if chat.Type != "group" && chat.Type != "supergroup" {
		return tgbotapi.NewMessage(chatID, "❌ این دستور فقط در گروه‌ها قابل استفاده است")
	}

	// Admin check
	isAdmin, err := m.isUserAdmin(chatID, update.Message.From.ID)
	if err != nil {
		log.Printf("getChatMember error: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی دسترسی ادمین")
	}
	if !isAdmin {
		return tgbotapi.NewMessage(chatID, "❌ فقط ادمین‌های گروه می‌توانند پیام حذف کنند")
	}

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
		// Try to delete command message too
		_, _ = m.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: update.Message.MessageID})
		return tgbotapi.MessageConfig{}
	}

	// No count provided: fallback to reply deletion
	if update.Message.ReplyToMessage != nil {
		return m.HandleDelete(update)
	}

	return tgbotapi.NewMessage(chatID, "برای حذف چند پیام بنویسید: حذف 10 (حداکثر 300)\nیا روی یک پیام ریپلای کنید و بنویسید: حذف")
}

func (m *ModerationCommand) bulkDeletePrev(chatID int64, fromMessageID int, count int) {
	// Delete up to count previous message IDs; ignore individual errors
	for i := 1; i <= count; i++ {
		target := fromMessageID - i
		if target <= 0 {
			break
		}
		_, _ = m.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: target})
	}
}

func (m *ModerationCommand) isUserAdmin(chatID int64, userID int64) (bool, error) {
	cfg := tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{ChatID: chatID, UserID: userID}}
	member, err := m.bot.GetChatMember(cfg)
	if err != nil {
		return false, err
	}
	return member.IsAdministrator() || member.IsCreator(), nil
}

// HandleBanOnReply bans the replied user permanently, only if requester is admin and target is not admin
func (m *ModerationCommand) HandleBanOnReply(update tgbotapi.Update) tgbotapi.MessageConfig {
	chat := update.Message.Chat
	chatID := chat.ID

	// Only in groups
	if chat.Type != "group" && chat.Type != "supergroup" {
		return tgbotapi.NewMessage(chatID, "❌ این دستور فقط در گروه‌ها قابل استفاده است")
	}

	// Require reply to a user's message
	if update.Message.ReplyToMessage == nil || update.Message.ReplyToMessage.From == nil {
		return tgbotapi.NewMessage(chatID, "لطفاً روی پیام کاربری که می‌خواهید بن شود ریپلای کنید و بنویسید: بن")
	}

	// Only group admins can use
	isAdmin, err := m.isUserAdmin(chatID, update.Message.From.ID)
	if err != nil {
		log.Printf("getChatMember error (requester): %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی دسترسی ادمین")
	}
	if !isAdmin {
		return tgbotapi.NewMessage(chatID, "❌ فقط ادمین‌های گروه می‌توانند کاربر را بن کنند")
	}

	// Target user
	targetUserID := update.Message.ReplyToMessage.From.ID

	// Prevent banning admins
	isTargetAdmin, err := m.isUserAdmin(chatID, targetUserID)
	if err != nil {
		log.Printf("getChatMember error (target): %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی نقش کاربر هدف")
	}
	if isTargetAdmin {
		return tgbotapi.NewMessage(chatID, "❌ امکان بن کردن ادمین یا صاحب گروه وجود ندارد")
	}

	// Perform ban (permanent)
	banCfg := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: targetUserID,
		},
		UntilDate: 0,
	}

	if _, err := m.bot.Request(banCfg); err != nil {
		log.Printf("banChatMember error: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ بن انجام نشد. مطمئن شوید ربات دسترسی بن دارد")
	}

	// Try to delete the command message for cleanliness (ignore error)
	_, _ = m.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: update.Message.MessageID})

	// Confirmation message
	return tgbotapi.NewMessage(chatID, "✅ کاربر موردنظر بن شد")
}

// HandleMute mutes a replied user. Supports optional hours: "سکوت [n]" where n is hours. Without n -> indefinite.
func (m *ModerationCommand) HandleMute(update tgbotapi.Update) tgbotapi.MessageConfig {
	chat := update.Message.Chat
	chatID := chat.ID

	if chat.Type != "group" && chat.Type != "supergroup" {
		return tgbotapi.NewMessage(chatID, "❌ این دستور فقط در گروه‌ها قابل استفاده است")
	}

	// Only group admins can use
	isAdmin, err := m.isUserAdmin(chatID, update.Message.From.ID)
	if err != nil {
		log.Printf("getChatMember error (requester): %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی دسترسی ادمین")
	}
	if !isAdmin {
		return tgbotapi.NewMessage(chatID, "❌ فقط ادمین‌های گروه می‌توانند کاربر را سکوت کنند")
	}

	if update.Message.ReplyToMessage == nil || update.Message.ReplyToMessage.From == nil {
		return tgbotapi.NewMessage(chatID, "لطفاً روی پیام کاربری که می‌خواهید سکوت شود ریپلای کنید و بنویسید: سکوت [ساعت]")
	}

	targetUserID := update.Message.ReplyToMessage.From.ID

	// Prevent muting admins
	isTargetAdmin, err := m.isUserAdmin(chatID, targetUserID)
	if err != nil {
		log.Printf("getChatMember error (target): %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی نقش کاربر هدف")
	}
	if isTargetAdmin {
		return tgbotapi.NewMessage(chatID, "❌ امکان سکوت کردن ادمین یا صاحب گروه وجود ندارد")
	}

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
			CanSendPolls:          false,
			CanSendOtherMessages:  false,
			CanAddWebPagePreviews: false,
			CanChangeInfo:         false,
			CanInviteUsers:        false,
			CanPinMessages:        false,
		},
		UntilDate: until,
	}

	if _, err := m.bot.Request(restrictCfg); err != nil {
		log.Printf("restrictChatMember error: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ سکوت انجام نشد. مطمئن شوید ربات دسترسی مناسب دارد")
	}

	// Try to delete the command message for cleanliness (ignore error)
	_, _ = m.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: update.Message.MessageID})

	if until == 0 {
		return tgbotapi.NewMessage(chatID, "✅ کاربر موردنظر به‌صورت نامحدود سکوت شد")
	}
	return tgbotapi.NewMessage(chatID, "✅ کاربر موردنظر سکوت شد")
}

// HandleUnmute lifts mute restrictions from a replied user: "آزاد" on reply.
func (m *ModerationCommand) HandleUnmute(update tgbotapi.Update) tgbotapi.MessageConfig {
	chat := update.Message.Chat
	chatID := chat.ID

	if chat.Type != "group" && chat.Type != "supergroup" {
		return tgbotapi.NewMessage(chatID, "❌ این دستور فقط در گروه‌ها قابل استفاده است")
	}

	// Only group admins can use
	isAdmin, err := m.isUserAdmin(chatID, update.Message.From.ID)
	if err != nil {
		log.Printf("getChatMember error (requester): %v", err)
		return tgbotapi.NewMessage(chatID, "❌ خطا در بررسی دسترسی ادمین")
	}
	if !isAdmin {
		return tgbotapi.NewMessage(chatID, "❌ فقط ادمین‌های گروه می‌توانند کاربر را از سکوت خارج کنند")
	}

	if update.Message.ReplyToMessage == nil || update.Message.ReplyToMessage.From == nil {
		return tgbotapi.NewMessage(chatID, "لطفاً روی پیام کاربری که می‌خواهید از سکوت خارج شود ریپلای کنید و بنویسید: آزاد")
	}

	targetUserID := update.Message.ReplyToMessage.From.ID

	// Lift restrictions by allowing messaging-related permissions
	unrestrictCfg := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: targetUserID,
		},
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       true,
			CanSendMediaMessages:  true,
			CanSendPolls:          true,
			CanSendOtherMessages:  true,
			CanAddWebPagePreviews: true,
			CanChangeInfo:         false,
			CanInviteUsers:        false,
			CanPinMessages:        false,
		},
		UntilDate: 0,
	}

	if _, err := m.bot.Request(unrestrictCfg); err != nil {
		log.Printf("restrictChatMember (unmute) error: %v", err)
		return tgbotapi.NewMessage(chatID, "❌ آزاد کردن انجام نشد. مطمئن شوید ربات دسترسی مناسب دارد")
	}

	// Try to delete the command message (ignore error)
	_, _ = m.bot.Request(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: update.Message.MessageID})

	return tgbotapi.NewMessage(chatID, "✅ کاربر از سکوت خارج شد")
}
