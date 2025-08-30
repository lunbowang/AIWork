package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatType int

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)

type Chatlog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	ConversationId string   `bson:"conversationId"`
	SendId         string   `bson:"sendId"`
	RecvId         string   `bson:"recvId"`
	ChatType       ChatType `bson:"chatType"`
	MsgContent     string   `bson:"msgContent"`
	SendTime       int64    `bson:"sendTime"`

	UpdateAt int64 `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt int64 `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
