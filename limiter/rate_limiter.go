package limiter

import (
	"redhat-bot/storage"
	"time"
)

type RateLimiter struct {
	storage *storage.MySQLStorage
}

func NewRateLimiter(storage *storage.MySQLStorage) *RateLimiter {
	return &RateLimiter{
		storage: storage,
	}
}

func (r *RateLimiter) CheckRateLimit(userID int64) (bool, string) {
	usage, err := r.storage.GetUserUsage(userID)
	if err != nil {
		// در صورت خطا، اجازه درخواست بده
		return true, ""
	}

	// بررسی محدودیت درخواست‌ها
	if usage.RequestsToday >= 1000 {
		return false, "⚠️ شما به محدودیت درخواست روزانه رسیده‌اید. لطفاً فردا دوباره تلاش کنید."
	}

	// بررسی فاصله زمانی بین درخواست‌ها
	if !usage.LastRequest.IsZero() && time.Since(usage.LastRequest) < 5*time.Second {
		return false, "⚠️ لطفاً بین درخواست‌ها کمی صبر کنید."
	}

	return true, ""
}

func (r *RateLimiter) IncrementUsage(userID int64) {
	if err := r.storage.IncrementUserUsage(userID); err != nil {
		// خطا را لاگ کن اما اجازه ادامه بده
		// log.Printf("Error incrementing usage: %v", err)
	}
}

func (r *RateLimiter) GetRemainingRequests(userID int64) (int, time.Time) {
	usage, err := r.storage.GetUserUsage(userID)
	if err != nil {
		return 999, time.Now().Add(24 * time.Hour)
	}

	remaining := 1000 - usage.RequestsToday
	if remaining < 0 {
		remaining = 0
	}

	return remaining, usage.LastReset.Add(24 * time.Hour)
}
