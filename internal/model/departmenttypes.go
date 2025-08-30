package model

import (
	"ai/internal/domain"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Department struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Name       string `bson:"name,omitempty"`
	ParentId   string `bson:"parentId,omitempty"`
	ParentPath string `bson:"parent_path"`
	Level      int    `bson:"level,omitempty"`
	LeaderId   string `bson:"leaderId,omitempty"`

	UpdateAt int64 `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt int64 `bson:"createAt,omitempty" json:"createAt,omitempty"`
}

// Id Level3
// parent1Id:parent2Id:...
func DepartmentParentPath(path string, id string) string {
	return path + ":" + id
}

func ParseParentPath(parentPath string) []string {
	res := strings.Split(parentPath, ":")
	return res[1:]
}

func (d *Department) ToDepartment() *domain.Department {
	return &domain.Department{
		Id:         d.ID.Hex(),
		Name:       d.Name,
		ParentId:   d.ParentId,
		Level:      d.Level,
		LeaderId:   d.LeaderId,
		ParentPath: d.ParentPath,
	}
}
