package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DepartmentUser struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	DepId  string `bson:"depId,omitempty"`
	UserId string `bson:"userId,omitempty"`

	// TODO: Fill your own fields
	UpdateAt int64 `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt int64 `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
