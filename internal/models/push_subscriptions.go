package models

import "time"

type PushSubscription struct {
	ID        int64  `gorm:"primaryKey"`
	UserID    int64  `gorm:"index"`
	Endpoint  string `gorm:"uniqueIndex"`
	P256dh    string
	Auth      string
	CreatedAt time.Time
}

func (*PushSubscription) TableName() string {
	return "push_subscriptions"
}
