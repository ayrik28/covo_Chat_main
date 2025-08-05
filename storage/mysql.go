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
	GroupID int64 `gorm:"primaryKey"`
	UserID  int64 `gorm:"primaryKey"`
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

	// Auto Migrate the schemas
	if err := db.AutoMigrate(&UserUsage{}, &GroupMessage{}, &GroupMember{}, &FeatureSetting{}); err != nil {
		return nil, fmt.Errorf("error migrating database: %v", err)
	}

	return &MySQLStorage{db: db}, nil
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
	member := GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Name:    name,
	}
	return m.db.Save(&member).Error
}

func (m *MySQLStorage) GetGroupMembers(groupID int64) ([]GroupMember, error) {
	var members []GroupMember
	err := m.db.Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}

// Close database connection
func (m *MySQLStorage) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
