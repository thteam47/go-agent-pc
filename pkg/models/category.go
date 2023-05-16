package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Id          *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	DomainId    string              `bson:"domain_id,omitempty" json:"domain_id,omitempty"`
	CategoryId  string              `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Name        string              `json:"name,omitempty" bson:"name,omitempty"`
	Position    int32               `json:"position,omitempty" bson:"position,omitempty"`
	CreatedTime int32               `json:"created_time" bson:"created_time,omitempty"`
	UpdatedTime int32               `json:"updated_time" bson:"updated_time,omitempty"`
	Type        string              `bson:"type,omitempty" json:"type,omitempty"`
}
