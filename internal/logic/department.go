package logic

import (
	"ai/internal/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"ai/internal/domain"
	"ai/internal/svc"

	"github.com/jinzhu/copier"
)

type Department interface {
	Soa(ctx context.Context) (resp *domain.DepartmentSoaResp, err error)
	Info(ctx context.Context, req *domain.IdPathReq) (resp *domain.Department, err error)
	Create(ctx context.Context, req *domain.Department) (err error)
	Edit(ctx context.Context, req *domain.Department) (err error)
	Delete(ctx context.Context, req *domain.IdPathReq) (err error)
	SetDepUsers(ctx context.Context, req *domain.SetDepUser) (err error)
	DepUserInfo(ctx context.Context, req *domain.IdPathReq) (resp *domain.Department, err error)
}

type department struct {
	svcCtx *svc.ServiceContext
}

func NewDepartment(svcCtx *svc.ServiceContext) Department {
	return &department{
		svcCtx: svcCtx,
	}
}

// Soa 获取部门的树形结构数据（SOA服务接口风格）
func (l *department) Soa(ctx context.Context) (resp *domain.DepartmentSoaResp, err error) {
	deps, err := l.svcCtx.DepartmentModel.All(ctx)
	if err != nil {
		return nil, err
	}

	groupDep := make(map[string][]*domain.Department, len(deps))
	rootDep := make([]*domain.Department, 0)

	for i := range deps {
		if len(deps[i].ParentPath) == 0 {
			rootDep = append(rootDep, deps[i].ToDepartment())
			continue
		}
		groupDep[deps[i].ParentPath] = append(groupDep[deps[i].ParentPath], deps[i].ToDepartment())
	}

	l.buildTree(rootDep, groupDep)

	return &domain.DepartmentSoaResp{
		Child: rootDep,
	}, nil
}

func (l *department) buildTree(rootDep []*domain.Department, groupDep map[string][]*domain.Department) {
	for i := range rootDep {
		path := model.DepartmentParentPath(rootDep[i].ParentPath, rootDep[i].Id)

		data, ok := groupDep[path]
		if !ok || len(data) == 0 {
			continue
		}

		l.buildTree(data, groupDep)

		rootDep[i].Child = data
	}
}

// Info 获取部门信息
func (l *department) Info(ctx context.Context, req *domain.IdPathReq) (resp *domain.Department, err error) {
	dep, err := l.svcCtx.DepartmentModel.FindOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserModel.FindOne(ctx, dep.LeaderId)
	if err != nil {
		return nil, err
	}

	var res domain.Department
	copier.Copy(&res, dep)

	res.Leader = user.Name

	return &res, nil
}

// Create 创建部门
func (l *department) Create(ctx context.Context, req *domain.Department) (err error) {
	dep, err := l.svcCtx.DepartmentModel.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, model.ErrDepNotFound) {
		return err
	}
	if dep != nil {
		return errors.New("已存在该部门")
	}

	var parentPath string
	if len(req.ParentId) > 0 {
		pdep, err := l.svcCtx.DepartmentModel.FindOne(ctx, req.ParentId)
		if err != nil && !errors.Is(err, model.ErrDepNotFound) {
			return err
		}
		parentPath = model.DepartmentParentPath(pdep.ParentPath, req.ParentId)
	}

	depId := primitive.NewObjectID()

	err = l.svcCtx.DepartmentModel.Insert(ctx, &model.Department{
		ID:         depId,
		Name:       req.Name,
		ParentId:   req.ParentId,
		ParentPath: parentPath,
		Level:      req.Level,
		LeaderId:   req.LeaderId,
		CreateAt:   time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	// 将部门主管也添加到部门中
	return l.svcCtx.DepartmentUserModel.Insert(ctx, &model.DepartmentUser{
		DepId:  depId.Hex(),
		UserId: req.LeaderId,
	})
}

// Edit 修改部门
func (l *department) Edit(ctx context.Context, req *domain.Department) (err error) {
	dep, err := l.svcCtx.DepartmentModel.FindOne(ctx, req.Id)
	if err != nil {
		return err
	}

	dep2, err := l.svcCtx.DepartmentModel.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, model.ErrDepNotFound) {
		return err
	}
	if dep2 != nil && dep2.ID.Hex() != dep.ID.Hex() {
		return errors.New("已存在该部门")
	}

	return l.svcCtx.DepartmentModel.Update(ctx, &model.Department{
		ID:       dep.ID,
		Name:     req.Name,
		ParentId: req.ParentId,
		Level:    req.Level,
		LeaderId: req.LeaderId,
	})
}

// Delete 删除部门
func (l *department) Delete(ctx context.Context, req *domain.IdPathReq) (err error) {
	dep, err := l.svcCtx.DepartmentModel.FindOne(ctx, req.Id)
	if err != nil {
		if errors.Is(err, model.ErrDepNotFound) {
			return nil
		}
		return err
	}

	depUser, err := l.svcCtx.DepartmentUserModel.List(ctx, &domain.DepartmentListReq{DepId: req.Id})
	if err != nil {
		return err
	}

	if len(depUser) == 0 {
		return l.svcCtx.DepartmentModel.Delete(ctx, req.Id)
	}

	if len(depUser) > 1 || depUser[0].UserId != dep.LeaderId {
		return errors.New("该部门下还存在用户，不能删除该部门")
	}

	return l.svcCtx.DepartmentModel.Delete(ctx, req.Id)
}

// SetDepUsers 设置部门成员
func (l *department) SetDepUsers(ctx context.Context, req *domain.SetDepUser) (err error) {
	_, err = l.svcCtx.DepartmentModel.FindOne(ctx, req.DepId)
	if err != nil {
		return err
	}

	err = l.svcCtx.DepartmentUserModel.DeleteByDepId(ctx, req.DepId)
	if err != nil {
		return err
	}

	return l.svcCtx.DepartmentUserModel.Inserts(ctx, req.DepId, req.UserIds)
}

// DepUserInfo 获取部门成员信息
func (l *department) DepUserInfo(ctx context.Context, req *domain.IdPathReq) (resp *domain.Department, err error) {
	dep, err := l.svcCtx.DepartmentModel.FindOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if len(dep.ParentPath) == 0 {
		return dep.ToDepartment(), err
	}
	parentIds := model.ParseParentPath(dep.ParentPath)
	pdeps, err := l.svcCtx.DepartmentModel.ListToMap(ctx, &domain.DepartmentListReq{
		DepIds: parentIds,
	})

	if err != nil {
		return nil, err
	}

	var root *domain.Department
	var node *domain.Department
	for _, id := range parentIds {
		if _, ok := pdeps[id]; !ok {
			continue
		}

		if root == nil {
			root = pdeps[id].ToDepartment()
			node = root
			continue
		}
		tmp := pdeps[id].ToDepartment()
		node.Child = append(node.Child, tmp)
		node = tmp
	}
	if node != nil {
		node.Child = append(node.Child, dep.ToDepartment())
	}

	return root, nil
}
