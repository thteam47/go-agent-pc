package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type KeyItem struct {
	UUID               string    `gorm:"primaryKey"`
	DomainId           string    `json:"domain_id,omitempty" bson:"domain_id,omitempty"`
	UserId             string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	PositionUserId     int32     `json:"position_user_id,omitempty" bson:"position_user_id,omitempty"`
	CategoryId         string    `json:"category_id,omitempty" bson:"category_id,omitempty"`
	PositionCategoryId int32     `json:"position_category_id,omitempty" bson:"position_category_id,omitempty"`
	KeyPublic          string    `json:"key_public,omitempty" bson:"key_public,omitempty"`
	KeyPrivate         string    `json:"key_private,omitempty" bson:"key_private,omitempty"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (keyInfo *KeyItem) BeforeCreate(tx *gorm.DB) {
	keyInfo.UUID = uuid.New().String()
	return
}
