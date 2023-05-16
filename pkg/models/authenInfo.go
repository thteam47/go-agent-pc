package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthenInfo struct {
	Id           *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	AuthenInfoId string              `json:"authen_info_id,omitempty" bson:"authen_info_id,omitempty"`
	UserId       string              `json:"user_id,omitempty" bson:"user_id,omitempty"`
	HashPassword string              `json:"hash_password,omitempty" bson:"hash_password,omitempty"`
	Mfas         []*Mfa              `json:"mfas,omitempty" bson:"mfas,omitempty"`
	CreateTime   int32               `json:"create_time" bson:"create_time,omitempty"`
	UpdateTime   int32               `json:"update_time" bson:"update_time,omitempty"`
}

type Mfa struct {
	Type       string `json:"type,omitempty" bson:"type,omitempty"`
	Enabled    bool   `json:"enabled,omitempty" bson:"enable,omitempty"`
	Secret     string `json:"secret,omitempty" bson:"secret,omitempty"`
	PublicData string `json:"public_data,omitempty" bson:"public_data,omitempty"`
	Configured bool   `json:"configured,omitempty" bson:"configured,omitempty"`
	Url        string `json:"url,omitempty" bson:"url,omitempty"`
}
