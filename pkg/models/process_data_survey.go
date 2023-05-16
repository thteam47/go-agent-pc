package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ProcessDataSurvey struct {
	UUID                  string    `gorm:"primaryKey"`
	DomainId              string    `json:"domain_id,omitempty" bson:"domain_id,omitempty"`
	UserId                string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	PositionUser          int32     `json:"position_user,omitempty" bson:"position_user,omitempty"`
	PositionItem          int32     `json:"position_item,omitempty" bson:"position_item,omitempty"`
	ProcessedData         int32     `json:"processed_data,omitempty" bson:"processed_data,omitempty"`
	PositionItemOriginal  int32     `json:"position_item_original,omitempty" bson:"position_item_original,omitempty"`
	PositionItemOriginal1 int32     `json:"position_item_original_1,omitempty" bson:"position_item_original_1,omitempty"`
	PositionItemOriginal2 int32     `json:"position_item_original_2,omitempty" bson:"position_item_original_2,omitempty"`
	CreatedAt             time.Time `gorm:"autoCreateTime"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime"`
}

func (resultCard *ProcessDataSurvey) BeforeCreate(tx *gorm.DB) {
	resultCard.UUID = uuid.New().String()
	return
}
