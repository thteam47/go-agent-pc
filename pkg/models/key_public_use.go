package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type KeyPublicUse struct {
	Id             *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	DomainId       string              `bson:"domain_id,omitempty" json:"domain_id,omitempty"`
	CreatedTime    int32               `json:"created_time" bson:"created_time,omitempty"`
	UpdatedTime    int32               `json:"updated_time" bson:"updated_time,omitempty"`
	KeyPublicUseId string              `bson:"key_public_use_id,omitempty" json:"key_public_use_id,omitempty"`
	CategoryId     string              `bson:"category_id,omitempty" json:"category_id,omitempty"`
	PositionItem   int32               `bson:"position_item,omitempty" json:"position_item,omitempty"`
	Part           int32               `bson:"part,omitempty" json:"part,omitempty"`
	KeyPublic      string              `json:"key_public,omitempty" bson:"key_public,omitempty"`
}
