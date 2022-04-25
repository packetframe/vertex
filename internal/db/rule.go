package db

import "time"

// Rule represents a filter rule
type Rule struct {
	ID string `gorm:"primaryKey,type:uuid;default:uuid_generate_v4()" json:"id"`

	Name      string `json:"name"`
	Filter    string `json:"filter"`
	ExpireStr string `json:"expire" gorm:"-"`

	Expire    time.Duration `json:"-"`
	CreatedAt time.Time     `json:"created"`
}
