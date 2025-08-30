package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotUser         = errors.New("查询不到该用户")
	ErrDepNotFound     = errors.New("不存在该部门")
	ErrNotFound        = mongo.ErrNoDocuments
	ErrInvalidObjectId = errors.New("invalid objectId")
)
