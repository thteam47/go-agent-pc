package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	FullName      string              `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Email         string              `json:"email,omitempty" bson:"email,omitempty"`
	Username      string              `json:"username,omitempty" bson:"username,omitempty"`
	PermissionAll bool                `json:"permission_all,omitempty" bson:"permission_all,omitempty"`
	Role          string              `json:"role,omitempty" bson:"role,omitempty"`
	Permissions   []*Permission       `json:"permissions,omitempty" bson:"permissions,omitempty"`
	UserId        string              `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreateTime    int32               `json:"create_time" bson:"create_time,omitempty"`
	UpdateTime    int32               `json:"update_time" bson:"update_time,omitempty"`
}

type Permission struct {
	Privilege string   `json:"privilege,omitempty" bson:"privilege,omitempty"`
	Actions   []string `json:"actions,omitempty" bson:"actions,omitempty"`
}

type UserContext struct {
	UserId string `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type Context struct {
	AccessToken string `json:"access_token,omitempty" bson:"access_token,omitempty"`
}
