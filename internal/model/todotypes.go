package model

import (
	"ai/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoStatus int

const (
	TodoInProgress TodoStatus = iota + 1
	TodoFinish
	TodoCancel
	TodoTimeout
)

type (
	Todo struct {
		ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

		CreatorId  string        `bson:"creatorId"`
		Title      string        `bson:"title"`
		DeadlineAt int64         `bson:"deadlineAt"`
		Desc       string        `bson:"desc"`
		Records    []*TodoRecord `bson:"records"`
		Executes   []*UserTodo   `bson:"executes"`
		TodoStatus `bson:"todo_status"`

		// TODO: Fill your own fields
		UpdateAt int64 `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
		CreateAt int64 `bson:"createAt,omitempty" json:"createAt,omitempty"`
	}
	TodoRecord struct {
		UserId   string `json:"userId,omitempty"`
		UserName string `json:"userName,omitempty"`
		Content  string `json:"content,omitempty"`
		Image    string `json:"image,omitempty"`
		CreateAt int64  `json:"createAt,omitempty"`
	}
)

func (m *Todo) ToDomainTodoRecords() []*domain.TodoRecord {
	res := make([]*domain.TodoRecord, 0, len(m.Records))
	for _, record := range m.Records {
		res = append(res, record.ToDomainTodoRecord())
	}
	return res
}

func (m *Todo) ToDomainTodo() *domain.Todo {
	return &domain.Todo{
		ID:         m.ID.Hex(),
		CreatorId:  m.CreatorId,
		Title:      m.Title,
		DeadlineAt: m.DeadlineAt,
		Desc:       m.Desc,
		ExecuteIds: nil,
		TodoStatus: int(m.TodoStatus),
	}
}

func (m *TodoRecord) ToDomainTodoRecord() *domain.TodoRecord {
	return &domain.TodoRecord{
		UserId:   m.UserId,
		UserName: m.UserName,
		Content:  m.Content,
		Image:    m.Image,
		CreateAt: m.CreateAt,
	}
}
