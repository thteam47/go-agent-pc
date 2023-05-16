package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Survey struct {
	Id           *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	DomainId     string              `bson:"domain_id,omitempty" json:"domain_id,omitempty"`
	SurveyId     string              `json:"survey_id,omitempty" bson:"survey_id,omitempty"`
	Name         string              `json:"name,omitempty" bson:"name,omitempty"`
	Status       string              `json:"status,omitempty" bson:"status,omitempty"`
	CreatedTime  int32               `json:"created_time" bson:"created_time,omitempty"`
	UpdatedTime  int32               `json:"updated_time" bson:"updated_time,omitempty"`
	Type         string              `bson:"type,omitempty" json:"type,omitempty"`
	UserIdCreate string              `bson:"user_id_create,omitempty" json:"user_id_create,omitempty"`
	UserIdVerify string              `bson:"user_id_verify,omitempty" json:"user_id_verify,omitempty"`
	UserIdJoin   []string            `bson:"user_id_join,omitempty" json:"user_id_join,omitempty"`
	Questions    []Question          `bson:"questions,omitempty" json:"questions,omitempty"`
}

type Question struct {
	CategoryId string   `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Message    string   `json:"message,omitempty" bson:"message,omitempty"`
	Position   int32    `json:"position,omitempty" bson:"position,omitempty"`
	Type       string   `bson:"type,omitempty" json:"type,omitempty"`
	Answers    []string `bson:"answers,omitempty" json:"answers,omitempty"`
}
