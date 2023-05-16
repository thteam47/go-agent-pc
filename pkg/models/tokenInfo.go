package models

type TokenInfo struct {
	AuthenticationDone bool                   `json:"authentication_done,omitempty" bson:"authentication_done,omitempty"`
	DomainId           string                 `json:"domain_id,omitempty" bson:"domain_id,omitempty"`
	Subject            string                 `bson:"subject,omitempty" json:"subject,omitempty"`
	Exp                int64                  `json:"exp,omitempty" bson:"exp,omitempty"`
	Roles              []string               `bson:"roles,omitempty" json:"roles,omitempty"`
	PermissionAll      bool                   `json:"permission_all,omitempty" bson:"permission_all,omitempty"`
	Permissions        []Permission           `json:"permissions,omitempty" bson:"permissions,omitempty"`
	MetaData           map[string]interface{} `bson:"meta_data,omitempty" json:"meta_data,omitempty"`
	initialized        bool
	subjectParsed      bool
	subjectID          string
	subjectType        string
}
