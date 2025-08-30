package logic

import (
	"ai/internal/model"
	"ai/token"
	"context"
	"errors"
	"time"

	"gitee.com/dn-jinmin/tlog"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"ai/internal/domain"
	"ai/internal/svc"
)

type Todo interface {
	Info(ctx context.Context, req *domain.IdPathReq) (resp *domain.TodoInfoResp, err error)
	Create(ctx context.Context, req *domain.Todo) (resp *domain.IdResp, err error)
	Edit(ctx context.Context, req *domain.Todo) (err error)
	Delete(ctx context.Context, req *domain.IdPathReq) (err error)
	Finish(ctx context.Context, req *domain.FinishedTodoReq) (err error)
	CreateRecord(ctx context.Context, req *domain.TodoRecord) (err error)
	List(ctx context.Context, req *domain.TodoListReq) (resp *domain.TodoListResp, err error)
}

type todo struct {
	svcCtx *svc.ServiceContext
}

func NewTodo(svcCtx *svc.ServiceContext) Todo {
	return &todo{
		svcCtx: svcCtx,
	}
}

func (l *todo) Info(ctx context.Context, req *domain.IdPathReq) (resp *domain.TodoInfoResp, err error) {
	t, err := l.svcCtx.TodoModel.FindOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 获取到关联的用户列表
	uids := make([]string, 0, len(t.Executes))
	for i, _ := range t.Executes {
		uids = append(uids, t.Executes[i].UserId)
	}
	// 关联用户的用户信息
	users, err := l.svcCtx.UserModel.ListToMaps(ctx, &domain.UserListReq{Ids: uids})
	if err != nil {
		return nil, err
	}
	userTodoDomains := make([]*domain.UserTodo, 0, len(t.Executes))
	for i, _ := range t.Executes {
		u, ok := users[t.Executes[i].UserId]
		if !ok {
			u = new(model.User)
		}
		userTodoDomains = append(userTodoDomains, t.Executes[i].ToDomain(u.Name))
	}

	if time.Now().Unix() > t.DeadlineAt {
		t.TodoStatus = model.TodoTimeout
	}

	// 创建人
	creator, ok := users[t.CreatorId]
	if !ok {
		return nil, errors.New("用户信息查询失败")
	}

	return &domain.TodoInfoResp{
		ID:          t.ID.Hex(),
		CreatorId:   t.CreatorId,
		CreatorName: creator.Name,
		Title:       t.Title,
		DeadlineAt:  t.DeadlineAt,
		Desc:        t.Desc,
		Records:     t.ToDomainTodoRecords(),
		Status:      int(t.TodoStatus),
		ExecuteIds:  userTodoDomains,
		TodoStatus:  int(t.TodoStatus),
	}, nil
}

func (l *todo) Create(ctx context.Context, req *domain.Todo) (resp *domain.IdResp, err error) {

	tlog.InfoCtx(ctx, "create todo ", req)
	uid := token.GetUId(ctx)
	executes := make([]*model.UserTodo, 0, len(req.ExecuteIds))
	for _, id := range req.ExecuteIds {
		executes = append(executes, &model.UserTodo{
			UserId:     id,
			TodoStatus: model.TodoInProgress,
		})
	}

	if len(executes) == 0 {
		executes = append(executes, &model.UserTodo{
			UserId:     uid,
			TodoStatus: model.TodoInProgress,
		})
	}

	var records []*model.TodoRecord
	copier.Copy(&records, req.Records)

	tlog.InfoCtx(ctx, "create todo insert", req)

	id := primitive.NewObjectID()
	err = l.svcCtx.TodoModel.Insert(ctx, &model.Todo{
		ID:         id,
		CreatorId:  uid,
		Title:      req.Title,
		DeadlineAt: req.DeadlineAt,
		Desc:       req.Desc,
		Records:    records,
		Executes:   executes,
		TodoStatus: model.TodoInProgress,
		CreateAt:   time.Now().Unix(),
		UpdateAt:   time.Now().Unix(),
	})
	if err != nil {
		return
	}

	return &domain.IdResp{
		Id: id.Hex(),
	}, nil
}

func (l *todo) Edit(ctx context.Context, req *domain.Todo) (err error) {
	// todo
	return
}

func (l *todo) Delete(ctx context.Context, req *domain.IdPathReq) (err error) {
	uid := token.GetUId(ctx)

	todo, err := l.svcCtx.TodoModel.FindOne(ctx, req.Id)
	if err != nil {
		return err
	}

	if uid != todo.CreatorId {
		return errors.New("你不能删除该待办事项")
	}

	return l.svcCtx.TodoModel.Delete(ctx, req.Id)
}

func (l *todo) Finish(ctx context.Context, req *domain.FinishedTodoReq) (err error) {
	todo, err := l.svcCtx.TodoModel.FindOne(ctx, req.TodoId)
	if err != nil {
		return err
	}

	for i, _ := range todo.Executes {
		if todo.Executes[i].UserId != req.UserId {
			continue
		}

		todo.Executes[i].TodoStatus = model.TodoFinish
	}

	isAllFinished := true
	for i, _ := range todo.Executes {
		if todo.Executes[i].TodoStatus != model.TodoFinish {
			isAllFinished = false
			break
		}
	}

	return l.svcCtx.TodoModel.UpdateFinished(ctx, todo, isAllFinished)
}

func (l *todo) CreateRecord(ctx context.Context, req *domain.TodoRecord) (err error) {

	req.UserId = token.GetUId(ctx)

	todo, err := l.svcCtx.TodoModel.FindOne(ctx, req.TodoId)
	if err != nil {
		return err
	}

	var record model.TodoRecord
	copier.Copy(&record, req)

	todo.Records = append(todo.Records, &record)

	return l.svcCtx.TodoModel.Update(ctx, todo)
}

func (l *todo) List(ctx context.Context, req *domain.TodoListReq) (resp *domain.TodoListResp, err error) {
	data, count, err := l.svcCtx.TodoModel.List(ctx, req)
	if err != nil {
		return nil, err
	}

	var todoDomains []*domain.Todo
	for i, _ := range data {
		if time.Now().Unix() > data[i].DeadlineAt {
			data[i].TodoStatus = model.TodoTimeout
		}
		todoDomains = append(todoDomains, data[i].ToDomainTodo())
	}

	return &domain.TodoListResp{
		Count: count,
		List:  todoDomains,
	}, nil
}
