package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	JTI string `gorm:"uniqueIndex;not null"`

	ExpiresAt time.Time

	Revoked bool `gorm:"default:false"`

	CreatedAt time.Time
}
