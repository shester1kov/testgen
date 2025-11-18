package entity

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID     *uuid.UUID `json:"user_id,omitempty" gorm:"type:uuid;index"`
	Action     string     `json:"action" gorm:"type:varchar(255);not null"`
	EntityType string     `json:"entity_type,omitempty" gorm:"type:varchar(100)"`
	EntityID   *uuid.UUID `json:"entity_id,omitempty" gorm:"type:uuid"`
	IPAddress  string     `json:"ip_address,omitempty" gorm:"type:inet"`
	UserAgent  string     `json:"user_agent,omitempty" gorm:"type:text"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime;index"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (ActivityLog) TableName() string {
	return "activity_logs"
}
