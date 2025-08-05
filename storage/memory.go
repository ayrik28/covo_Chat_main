package storage

// import (
// 	"sync"
// 	"time"
// )

// type UserUsage struct {
// 	RequestsToday int
// 	LastReset     time.Time
// 	LastRequest   time.Time
// }

// type GroupMessage struct {
// 	UserID    int64
// 	Username  string
// 	Message   string
// 	Timestamp time.Time
// }

// type GroupMember struct {
// 	UserID int64
// 	Name   string
// }

// type MemoryStorage struct {
// 	mu            sync.RWMutex
// 	userUsage     map[int64]*UserUsage
// 	groupMessages map[int64][]GroupMessage
// 	clownEnabled  map[int64]bool
// 	crushEnabled  map[int64]bool
// 	groupMembers  map[int64][]GroupMember
// }

// func NewMemoryStorage() *MemoryStorage {
// 	return &MemoryStorage{
// 		userUsage:     make(map[int64]*UserUsage),
// 		groupMessages: make(map[int64][]GroupMessage),
// 		clownEnabled:  make(map[int64]bool),
// 		crushEnabled:  make(map[int64]bool),
// 		groupMembers:  make(map[int64][]GroupMember),
// 	}
// }

// func (m *MemoryStorage) GetUserUsage(userID int64) *UserUsage {
// 	m.mu.RLock()
// 	defer m.mu.RUnlock()

// 	if usage, exists := m.userUsage[userID]; exists {
// 		// بررسی نیاز به بازنشانی شمارنده روزانه
// 		if time.Since(usage.LastReset) >= 24*time.Hour {
// 			usage.RequestsToday = 0
// 			usage.LastReset = time.Now()
// 		}
// 		return usage
// 	}

// 	// ایجاد استفاده جدید کاربر
// 	usage := &UserUsage{
// 		RequestsToday: 0,
// 		LastReset:     time.Now(),
// 		LastRequest:   time.Time{},
// 	}
// 	m.userUsage[userID] = usage
// 	return usage
// }

// func (m *MemoryStorage) IncrementUserUsage(userID int64) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	usage := m.GetUserUsage(userID)
// 	usage.RequestsToday++
// 	usage.LastRequest = time.Now()
// }

// func (m *MemoryStorage) AddGroupMessage(groupID int64, userID int64, username, message string) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	// پاک کردن پیام‌های قدیمی (بیش از ۲۴ ساعت)
// 	m.cleanOldMessages(groupID)

// 	groupMsg := GroupMessage{
// 		UserID:    userID,
// 		Username:  username,
// 		Message:   message,
// 		Timestamp: time.Now(),
// 	}

// 	m.groupMessages[groupID] = append(m.groupMessages[groupID], groupMsg)
// }

// func (m *MemoryStorage) GetGroupMessages(groupID int64) []GroupMessage {
// 	m.mu.RLock()
// 	defer m.mu.RUnlock()

// 	m.cleanOldMessages(groupID)
// 	return m.groupMessages[groupID]
// }

// func (m *MemoryStorage) cleanOldMessages(groupID int64) {
// 	messages := m.groupMessages[groupID]
// 	var validMessages []GroupMessage

// 	cutoff := time.Now().Add(-24 * time.Hour)

// 	for _, msg := range messages {
// 		if msg.Timestamp.After(cutoff) {
// 			validMessages = append(validMessages, msg)
// 		}
// 	}

// 	m.groupMessages[groupID] = validMessages
// }

// func (m *MemoryStorage) ClearGroupMessages(groupID int64) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	delete(m.groupMessages, groupID)
// }

// // IsClownEnabled checks if the clown feature is enabled for a chat
// func (m *MemoryStorage) IsClownEnabled(chatID int64) bool {
// 	m.mu.RLock()
// 	defer m.mu.RUnlock()
// 	return m.clownEnabled[chatID]
// }

// // SetClownEnabled sets the clown feature status for a chat
// func (m *MemoryStorage) SetClownEnabled(chatID int64, enabled bool) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	m.clownEnabled[chatID] = enabled
// }

// // IsCrushEnabled checks if the crush feature is enabled for a chat
// func (m *MemoryStorage) IsCrushEnabled(chatID int64) bool {
// 	m.mu.RLock()
// 	defer m.mu.RUnlock()
// 	return m.crushEnabled[chatID]
// }

// // SetCrushEnabled sets the crush feature status for a chat
// func (m *MemoryStorage) SetCrushEnabled(chatID int64, enabled bool) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	m.crushEnabled[chatID] = enabled
// }

// // GetCrushEnabledGroups returns all groups where crush feature is enabled
// func (m *MemoryStorage) GetCrushEnabledGroups() []int64 {
// 	m.mu.RLock()
// 	defer m.mu.RUnlock()

// 	var enabledGroups []int64
// 	for groupID, enabled := range m.crushEnabled {
// 		if enabled {
// 			enabledGroups = append(enabledGroups, groupID)
// 		}
// 	}
// 	return enabledGroups
// }

// // AddGroupMember adds a member to a group
// func (m *MemoryStorage) AddGroupMember(groupID int64, userID int64, name string) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	// بررسی اینکه آیا کاربر قبلاً اضافه شده
// 	for _, member := range m.groupMembers[groupID] {
// 		if member.UserID == userID {
// 			return // کاربر قبلاً وجود دارد
// 		}
// 	}

// 	member := GroupMember{
// 		UserID: userID,
// 		Name:   name,
// 	}
// 	m.groupMembers[groupID] = append(m.groupMembers[groupID], member)
// }

// // GetGroupMembers returns all members of a group
// func (m *MemoryStorage) GetGroupMembers(groupID int64) []GroupMember {
// 	m.mu.RLock()
// 	defer m.mu.RUnlock()
// 	return m.groupMembers[groupID]
// }
