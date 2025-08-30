package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func entityList(ctx context.Context, col *mongo.Collection, query interface{}, v any, opt ...*options.FindOptions) error {
	cur, err := col.Find(ctx, query, opt...)
	if err != nil {
		return err
	}

	return cur.All(ctx, v)
}

func entityUpdateOrInsert(ctx context.Context, col *mongo.Collection, filter interface{}, update interface{}) error {
	_, err := col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

// Pagination 处理分页
func Pagination(page, count int) (setLimit, setSkip *int64) {
	setLimit = new(int64)
	setSkip = new(int64)

	if page <= 0 {
		page = 0
	}

	if count <= 0 {
		count = 1000
	}

	*setLimit = int64(count)
	*setSkip = int64(count * page)

	return setLimit, setSkip
}
