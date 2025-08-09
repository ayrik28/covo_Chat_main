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
	if err := db.AutoMigrate(&UserUsage{}, &GroupMessage{}, &GroupMember{}, &FeatureSetting{}); err != nil {
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
