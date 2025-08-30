package model

import (
	"ai/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Name     string `bson:"name"`
	Password string `bson:"password"`
	Status   int    `bson:"status"`
	IsSystem bool   `bson:"isSystem"`

	UpdateAt int64 `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt int64 `bson:"createAt,omitempty" json:"createAt,omitempty"`
}

func (u *User) ToDomainUser() *domain.User {
	return &domain.User{
		Id:     u.ID.Hex(),
		Name:   u.Name,
		Status: u.Status,
	}
}
