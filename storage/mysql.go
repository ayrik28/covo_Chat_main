package storage

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Models
type UserUsage struct {
	UserID        int64 `gorm:"primaryKey"`
	RequestsToday int
	LastReset     time.Time
	LastRequest   time.Time
}

type GroupMessage struct {
	ID        uint  `gorm:"primaryKey"`
	GroupID   int64 `gorm:"index:idx_group_timestamp"`
	UserID    int64
	Username  string
	Message   string
	Timestamp time.Time `gorm:"index:idx_group_timestamp"`
}

type GroupMember struct {
	GroupID int64 `gorm:"primaryKey;uniqueIndex:idx_group_user"`
	UserID  int64 `gorm:"primaryKey;uniqueIndex:idx_group_user"`
	Name    string
}

type FeatureSetting struct {
	GroupID     int64  `gorm:"primaryKey"`
	FeatureName string `gorm:"primaryKey"`
	Enabled     bool
}

// DailyChallenge نگهداری وضعیت چالش روزانه در هر گروه
type DailyChallenge struct {
	ID         uint      `gorm:"primaryKey"`
	GroupID    int64     `gorm:"index"`
	MessageID  int       `gorm:"index"`
	Proverb    string    `gorm:"type:text"`
	Emojis     string    `gorm:"type:text"`
	Answered   bool      `gorm:"default:false"`
	WinnerID   int64     `gorm:"index"`
	WinnerName string    `gorm:"type:varchar(255)"`
	CreatedAt  time.Time `gorm:"index"`
	UpdatedAt  time.Time `gorm:"index"`
}

// UserOnboarding برای ارسال یک‌باره پیام عضویت در اولین استارت
type UserOnboarding struct {
	UserID    int64     `gorm:"primaryKey"`
	PromoSent bool      `gorm:"default:false"`
	SentAt    time.Time `gorm:"index"`
}

// BotChannel نگهداری اطلاعات کانال‌هایی که ربات در آن‌ها حضور دارد
type BotChannel struct {
	ID          uint  `gorm:"primaryKey"`
	ChatID      int64 `gorm:"uniqueIndex;column:chat_id"`
	Title       string
	Username    string    `gorm:"index"`
	IsAdmin     bool      `gorm:"default:false;column:is_admin"`
	MemberCount int       `gorm:"default:0"`
	DateAdded   time.Time `gorm:"index;column:date_added"`
	LastCheck   time.Time `gorm:"index;column:last_check"`
}

// RequiredChannel نگهداری لینک/کانال‌های الزام عضویت (جهانی یا بر اساس گروه)
type RequiredChannel struct {
	ID              uint  `gorm:"primaryKey"`
	GroupID         int64 `gorm:"index"` // 0: سراسری
	Title           string
	Link            string
	ChannelUsername string `gorm:"index"`                // مثال: mychannel بدون @
	ChannelID       int64  `gorm:"index"`                // اگر شناسه عددی موجود است
	ChatID          int64  `gorm:"index;column:chat_id"` // chat_id عددی کانال/گروه
	BotJoined       bool   `gorm:"default:false"`
	MemberCount     int    `gorm:"default:0"`
	CreatedAt       time.Time
	LastChecked     time.Time
}

type MySQLStorage struct {
	db *gorm.DB
}

func NewMySQLStorage(host, port, user, password, dbname string) (*MySQLStorage, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to MySQL: %v", err)
	}

	// پاکسازی داده‌های تکراری قبل از اضافه کردن ایندکس یونیک
	if err := cleanupDuplicateMembers(db); err != nil {
		return nil, fmt.Errorf("error cleaning up duplicate members: %v", err)
	}

	// Auto Migrate the schemas
	if err := db.AutoMigrate(&UserUsage{}, &GroupMessage{}, &GroupMember{}, &FeatureSetting{}, &RequiredChannel{}, &UserOnboarding{}, &BotChannel{}, &DailyChallenge{}); err != nil {
		return nil, fmt.Errorf("error migrating database: %v", err)
	}

	return &MySQLStorage{db: db}, nil
}

// پاکسازی داده‌های تکراری از جدول group_members
func cleanupDuplicateMembers(db *gorm.DB) error {
	// ایجاد جدول موقت برای نگهداری آخرین رکورد هر ترکیب group_id و user_id
	if err := db.Exec(`
		CREATE TEMPORARY TABLE temp_members AS
		SELECT group_id, user_id, MAX(name) as name
		FROM group_members
		GROUP BY group_id, user_id
	`).Error; err != nil {
		return fmt.Errorf("error creating temporary table: %v", err)
	}

	// پاک کردن تمام داده‌های جدول اصلی
	if err := db.Exec("DELETE FROM group_members").Error; err != nil {
		return fmt.Errorf("error deleting from group_members: %v", err)
	}

	// انتقال داده‌های پاکسازی شده به جدول اصلی
	if err := db.Exec(`
		INSERT INTO group_members (group_id, user_id, name)
		SELECT group_id, user_id, name FROM temp_members
	`).Error; err != nil {
		return fmt.Errorf("error inserting cleaned data: %v", err)
	}

	// حذف جدول موقت
	if err := db.Exec("DROP TEMPORARY TABLE IF EXISTS temp_members").Error; err != nil {
		return fmt.Errorf("error dropping temporary table: %v", err)
	}

	// چک کردن وجود ایندکس
	var hasIndex bool
	err := db.Raw(`SELECT 1 FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = 'group_members' 
		AND index_name = 'idx_group_user' LIMIT 1`).Scan(&hasIndex).Error
	if err != nil {
		return fmt.Errorf("error checking index existence: %v", err)
	}

	// اگر ایندکس وجود نداشت، اضافه کن
	if !hasIndex {
		if err := db.Exec("ALTER TABLE group_members ADD UNIQUE INDEX idx_group_user (group_id, user_id)").Error; err != nil {
			return fmt.Errorf("error adding unique index: %v", err)
		}
	}

	return nil
}

// User Usage Methods
func (m *MySQLStorage) GetUserUsage(userID int64) (*UserUsage, error) {
	var usage UserUsage
	result := m.db.FirstOrCreate(&usage, UserUsage{UserID: userID})
	if result.Error != nil {
		return nil, result.Error
	}

	// Check if reset needed
	if time.Since(usage.LastReset) >= 24*time.Hour {
		usage.RequestsToday = 0
		usage.LastReset = time.Now()
		if err := m.db.Save(&usage).Error; err != nil {
			return nil, err
		}
	}

	return &usage, nil
}

func (m *MySQLStorage) IncrementUserUsage(userID int64) error {
	usage, err := m.GetUserUsage(userID)
	if err != nil {
		return err
	}

	usage.RequestsToday++
	usage.LastRequest = time.Now()
	return m.db.Save(usage).Error
}

// Group Messages Methods
func (m *MySQLStorage) AddGroupMessage(groupID int64, userID int64, username, message string) error {
	// Clean old messages first
	if err := m.cleanOldMessages(groupID); err != nil {
		log.Printf("Error cleaning old messages: %v", err)
	}

	msg := GroupMessage{
		GroupID:   groupID,
		UserID:    userID,
		Username:  username,
		Message:   message,
		Timestamp: time.Now(),
	}

	return m.db.Create(&msg).Error
}

func (m *MySQLStorage) GetGroupMessages(groupID int64) ([]GroupMessage, error) {
	if err := m.cleanOldMessages(groupID); err != nil {
		log.Printf("Error cleaning old messages: %v", err)
	}

	var messages []GroupMessage
	err := m.db.Where("group_id = ?", groupID).
		Order("timestamp desc").
		Find(&messages).Error

	return messages, err
}

func (m *MySQLStorage) cleanOldMessages(groupID int64) error {
	cutoff := time.Now().Add(-24 * time.Hour)
	return m.db.Where("group_id = ? AND timestamp < ?", groupID, cutoff).
		Delete(&GroupMessage{}).Error
}

func (m *MySQLStorage) ClearGroupMessages(groupID int64) error {
	return m.db.Where("group_id = ?", groupID).Delete(&GroupMessage{}).Error
}

// Stats and Analytics (24h)
// UserMessageCount holds aggregated count of messages for a user in last 24 hours
type UserMessageCount struct {
	UserID   int64
	Username string
	Count    int64
}

// GetUserMessageCountLast24h returns the number of messages a specific user has sent
// in the specified group during the last 24 hours.
func (m *MySQLStorage) GetUserMessageCountLast24h(groupID int64, userID int64) (int64, error) {
	cutoff := time.Now().Add(-24 * time.Hour)
	var count int64
	err := m.db.Model(&GroupMessage{}).
		Where("group_id = ? AND user_id = ? AND timestamp >= ?", groupID, userID, cutoff).
		Count(&count).Error
	return count, err
}

// GetTopActiveUsersLast24h returns top N users with most messages in the last 24 hours for a group
func (m *MySQLStorage) GetTopActiveUsersLast24h(groupID int64, limit int) ([]UserMessageCount, error) {
	cutoff := time.Now().Add(-24 * time.Hour)
	var results []UserMessageCount
	// Use MAX(username) to pick a representative username for the user
	err := m.db.Table("group_messages").
		Select("user_id as user_id, MAX(username) as username, COUNT(*) as count").
		Where("group_id = ? AND timestamp >= ?", groupID, cutoff).
		Group("user_id").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}

// GetAllActiveUsersLast24h returns all users with message counts in the last 24 hours for a group
func (m *MySQLStorage) GetAllActiveUsersLast24h(groupID int64) ([]UserMessageCount, error) {
	cutoff := time.Now().Add(-24 * time.Hour)
	var results []UserMessageCount
	err := m.db.Table("group_messages").
		Select("user_id as user_id, MAX(username) as username, COUNT(*) as count").
		Where("group_id = ? AND timestamp >= ?", groupID, cutoff).
		Group("user_id").
		Order("count DESC").
		Scan(&results).Error
	return results, err
}

// Feature Settings Methods
func (m *MySQLStorage) IsFeatureEnabled(chatID int64, feature string) (bool, error) {
	var setting FeatureSetting
	err := m.db.Where("group_id = ? AND feature_name = ?", chatID, feature).
		First(&setting).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return setting.Enabled, nil
}

func (m *MySQLStorage) SetFeatureEnabled(chatID int64, feature string, enabled bool) error {
	setting := FeatureSetting{
		GroupID:     chatID,
		FeatureName: feature,
		Enabled:     enabled,
	}

	return m.db.Save(&setting).Error
}

// GetEnabledGroupsForFeature returns all group IDs that have a specific feature enabled
func (m *MySQLStorage) GetEnabledGroupsForFeature(feature string) ([]int64, error) {
	var settings []FeatureSetting
	if err := m.db.Where("feature_name = ? AND enabled = ?", feature, true).Find(&settings).Error; err != nil {
		return nil, err
	}
	groups := make([]int64, 0, len(settings))
	for _, s := range settings {
		groups = append(groups, s.GroupID)
	}
	return groups, nil
}

// CreateDailyChallenge inserts a new daily challenge row for a group
func (m *MySQLStorage) CreateDailyChallenge(groupID int64, messageID int, proverb string, emojis string) error {
	dc := DailyChallenge{
		GroupID:   groupID,
		MessageID: messageID,
		Proverb:   proverb,
		Emojis:    emojis,
		Answered:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return m.db.Create(&dc).Error
}

// GetActiveChallengeForGroup returns the latest not-expired challenge for a group (today)
func (m *MySQLStorage) GetActiveChallengeForGroup(groupID int64) (*DailyChallenge, error) {
	// limit to last 24 hours to ensure "روزانه" semantics
	cutoff := time.Now().Add(-24 * time.Hour)
	var dc DailyChallenge
	err := m.db.Where("group_id = ? AND created_at >= ?", groupID, cutoff).Order("id DESC").First(&dc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dc, nil
}

// TryMarkChallengeAnswered marks challenge answered if not already answered; returns true if succeeded
func (m *MySQLStorage) TryMarkChallengeAnswered(id uint, winnerID int64, winnerName string) (bool, error) {
	// optimistic update where answered=false
	res := m.db.Model(&DailyChallenge{}).
		Where("id = ? AND answered = ?", id, false).
		Updates(map[string]interface{}{
			"answered":    true,
			"winner_id":   winnerID,
			"winner_name": winnerName,
			"updated_at":  time.Now(),
		})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// Clown Feature Methods
func (m *MySQLStorage) IsClownEnabled(chatID int64) (bool, error) {
	return m.IsFeatureEnabled(chatID, "clown")
}

func (m *MySQLStorage) SetClownEnabled(chatID int64, enabled bool) error {
	return m.SetFeatureEnabled(chatID, "clown", enabled)
}

// Crush Feature Methods
func (m *MySQLStorage) IsCrushEnabled(chatID int64) (bool, error) {
	return m.IsFeatureEnabled(chatID, "crush")
}

func (m *MySQLStorage) SetCrushEnabled(chatID int64, enabled bool) error {
	return m.SetFeatureEnabled(chatID, "crush", enabled)
}

func (m *MySQLStorage) GetCrushEnabledGroups() ([]int64, error) {
	var settings []FeatureSetting
	err := m.db.Where("feature_name = ? AND enabled = ?", "crush", true).
		Find(&settings).Error
	if err != nil {
		return nil, err
	}

	groups := make([]int64, len(settings))
	for i, setting := range settings {
		groups[i] = setting.GroupID
	}
	return groups, nil
}

// Group Members Methods
func (m *MySQLStorage) AddGroupMember(groupID int64, userID int64, name string) error {
	// استفاده از Upsert برای اضافه کردن یا آپدیت کردن عضو گروه
	return m.db.Exec(`
		INSERT INTO group_members (group_id, user_id, name)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
		name = ?`,
		groupID, userID, name, name,
	).Error
}

func (m *MySQLStorage) GetGroupMembers(groupID int64) ([]GroupMember, error) {
	var members []GroupMember
	err := m.db.Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}

// Close database connection
// ساختار برای لیست کاربران
type UserInfo struct {
	UserID int64
	Name   string
}

// ساختار برای لیست گروه‌ها
type GroupInfo struct {
	GroupID   int64
	GroupName string
}

func (m *MySQLStorage) GetAllUsers() ([]UserInfo, error) {
	var users []UserInfo
	err := m.db.Table("group_members").
		Select("DISTINCT user_id, MAX(name) as name").
		Group("user_id").
		Find(&users).Error
	return users, err
}

func (m *MySQLStorage) GetAllGroups() ([]GroupInfo, error) {
	var groups []GroupInfo
	err := m.db.Table("group_members").
		Select("DISTINCT group_id, 'Group' as group_name").
		Group("group_id").
		Find(&groups).Error
	return groups, err
}

func (m *MySQLStorage) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// BotChannels methods

// UpsertBotChannel ثبت/به‌روزرسانی اطلاعات کانال ربات
func (m *MySQLStorage) UpsertBotChannel(chatID int64, title string, username string, isAdmin bool, memberCount int) error {
	var bc BotChannel
	if err := m.db.Where("chat_id = ?", chatID).First(&bc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			bc = BotChannel{
				ChatID:      chatID,
				Title:       title,
				Username:    username,
				IsAdmin:     isAdmin,
				MemberCount: memberCount,
				DateAdded:   time.Now(),
				LastCheck:   time.Now(),
			}
			return m.db.Create(&bc).Error
		}
		return err
	}
	// update
	bc.Title = title
	bc.Username = username
	bc.IsAdmin = isAdmin
	bc.MemberCount = memberCount
	bc.LastCheck = time.Now()
	return m.db.Save(&bc).Error
}

// ListBotChannels لیست تمام کانال‌هایی که ربات در آن‌ها حضور دارد
func (m *MySQLStorage) ListBotChannels() ([]BotChannel, error) {
	var list []BotChannel
	if err := m.db.Order("date_added ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// WasPromoSent اولین استارت را چک می‌کند (آیا پیام عضویت قبلاً برای کاربر ارسال شده؟)
func (m *MySQLStorage) WasPromoSent(userID int64) (bool, error) {
	var rec UserOnboarding
	if err := m.db.First(&rec, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return rec.PromoSent, nil
}

// MarkPromoSent علامت‌گذاری ارسال پیام عضویت برای کاربر
func (m *MySQLStorage) MarkPromoSent(userID int64) error {
	rec := UserOnboarding{UserID: userID, PromoSent: true, SentAt: time.Now()}
	return m.db.Save(&rec).Error
}

// Required membership methods

// AddRequiredChannel اضافه/آپدیت یک لینک الزام عضویت (بر اساس GroupID و ChannelUsername/Link)
func (m *MySQLStorage) AddRequiredChannel(groupID int64, title string, link string, channelUsername string, channelID int64) error {
	rc := RequiredChannel{
		GroupID:         groupID,
		Title:           title,
		Link:            link,
		ChannelUsername: channelUsername,
		ChannelID:       channelID,
		ChatID:          channelID,
	}
	return m.db.Create(&rc).Error
}

// RemoveRequiredChannel حذف بر اساس ID رکورد
func (m *MySQLStorage) RemoveRequiredChannel(id uint) error {
	return m.db.Delete(&RequiredChannel{}, id).Error
}

// ListRequiredChannels لیست لینک‌های الزام عضویت برای یک گروه (به‌همراه لینک‌های سراسری group_id=0)
func (m *MySQLStorage) ListRequiredChannels(groupID int64) ([]RequiredChannel, error) {
	var list []RequiredChannel
	if err := m.db.Where("group_id = ? OR group_id = 0", groupID).Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateRequiredChannelStatus به‌روزرسانی وضعیت عضویت ربات و تعداد اعضا برای یک کانال
func (m *MySQLStorage) UpdateRequiredChannelStatus(id uint, botJoined bool, memberCount int) error {
	return m.db.Model(&RequiredChannel{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"bot_joined":   botJoined,
			"member_count": memberCount,
			"last_checked": time.Now(),
		}).Error
}

// UpdateRequiredChannelResolved به‌روزرسانی متادیتای کانال (ChannelID/Username/Title)
func (m *MySQLStorage) UpdateRequiredChannelResolved(id uint, channelID int64, username string, title string) error {
	updates := map[string]interface{}{}
	if channelID != 0 {
		updates["channel_id"] = channelID
		updates["chat_id"] = channelID
	}
	if username != "" {
		updates["channel_username"] = username
	}
	if title != "" {
		updates["title"] = title
	}
	if len(updates) == 0 {
		return nil
	}
	return m.db.Model(&RequiredChannel{}).Where("id = ?", id).Updates(updates).Error
}
