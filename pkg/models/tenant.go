package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tenant struct {
	Id         *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	TenantId   string              `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	Name       string              `json:"name,omitempty" bson:"name,omitempty"`
	Domain     string              `json:"domain,omitempty" bson:"domain,omitempty"`
	CustomerId string              `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	Meta       map[string]string   `bson:"meta,omitempty" json:"meta,omitempty"`
	CreateTime int32               `json:"create_time" bson:"create_time,omitempty"`
	UpdateTime int32               `json:"update_time" bson:"update_time,omitempty"`
}
