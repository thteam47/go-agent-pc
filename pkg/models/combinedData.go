package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type CombinedData struct {
	UUID           string            `gorm:"primaryKey"`
	TenantId1      string            `bson:"tenant_id_1,omitempty" json:"tenant_id_1,omitempty"`
	TenantId2      string            `json:"tenant_id_2,omitempty" bson:"tenant_id_2,omitempty"`
	Meta           map[string]string `bson:"meta,omitempty" json:"meta,omitempty"`
	Status         bool              `json:"status,omitempty" bson:"status,omitempty"`
	CreatedTime    int32             `json:"created_time" bson:"created_time,omitempty"`
	UpdatedTime    int32             `json:"updated_time" bson:"updated_time,omitempty"`
	CombinedDataId string            `bson:"combined_data_id,omitempty" json:"combined_data_id,omitempty"`
	NumberItem1    int32             `bson:"number_item_1,omitempty" json:"number_item_1,omitempty"`
	NumberItem2    int32             `bson:"number_item_2,omitempty" json:"number_item_2,omitempty"`
	STwoPart       int32             `bson:"s_two_part,omitempty" json:"s_two_part,omitempty"`
	SkTwoPart      int32             `bson:"sk_two_part,omitempty" json:"sk_two_part,omitempty"`
	NkTwoPart      int32             `bson:"nk_two_part,omitempty" json:"nk_two_part,omitempty"`
	SOnePart1      int32             `bson:"s_one_part_1,omitempty" json:"s_one_part_1,omitempty"`
	SOnePart2      int32             `bson:"s_one_part_2,omitempty" json:"s_one_part_2,omitempty"`
	NkOnePart1     int32             `bson:"nk_one_part_1,omitempty" json:"nk_one_part_1,omitempty"`
	NkOnePart2     int32             `bson:"nk_one_part_2,omitempty" json:"nk_one_part_2,omitempty"`
	NumberUser     int32             `bson:"number_user,omitempty" json:"number_user,omitempty"`
}

func (keyInfo *CombinedData) BeforeCreate(tx *gorm.DB) {
	keyInfo.UUID = uuid.New().String()
	return
}
