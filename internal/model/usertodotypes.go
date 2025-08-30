package model

import (
	"ai/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTodo struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	UserId     string     `bson:"userId,omitempty"`
	TodoId     string     `bson:"todoId,omitempty"`
	TodoStatus TodoStatus `bson:"todoStatus,omitempty"`
	// TODO: Fill your own fields
	UpdateAt int64 `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt int64 `bson:"createAt,omitempty" json:"createAt,omitempty"`
}

func (m *UserTodo) ToDomain(username string) *domain.UserTodo {
	return &domain.UserTodo{
		ID:         m.ID.Hex(),
		UserId:     m.UserId,
		UserName:   username,
		TodoId:     m.TodoId,
		TodoStatus: int(m.TodoStatus),
	}
}
