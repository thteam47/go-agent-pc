package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ResultCard struct {
	UUID           string    `gorm:"primaryKey"`
	DomainId       string    `json:"domain_id,omitempty" bson:"domain_id,omitempty"`
	UserId         string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	PositionUser   int32     `json:"position_user,omitempty" bson:"position_user,omitempty"`
	CategoryId     string    `json:"category_id,omitempty" bson:"category_id,omitempty"`
	PositionItem   int32     `json:"position_item,omitempty" bson:"position_item,omitempty"`
	Option         string    `json:"option,omitempty" bson:"option,omitempty"`
	PositionOption int32     `json:"position_option,omitempty" bson:"position_option,omitempty"`
	SurveyId       string    `json:"survey_id,omitempty" bson:"survey_id,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (resultCard *ResultCard) BeforeCreate(tx *gorm.DB) {
	resultCard.UUID = uuid.New().String()
	return
}
