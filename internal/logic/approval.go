package logic

import (
	"ai/internal/model"
	"ai/pkg/timex"
	"ai/token"
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"ai/internal/domain"
	"ai/internal/svc"
)

type Approval interface {
	Info(ctx context.Context, req *domain.IdPathReq) (resp *domain.ApprovalInfoResp, err error)
	Create(ctx context.Context, req *domain.Approval) (resp *domain.IdResp, err error)
	Dispose(ctx context.Context, req *domain.DisposeReq) (err error)
	List(ctx context.Context, req *domain.ApprovalListReq) (resp *domain.ApprovalListResp, err error)
}

type approval struct {
	svcCtx *svc.ServiceContext
}

func NewApproval(svcCtx *svc.ServiceContext) Approval {
	return &approval{
		svcCtx: svcCtx,
	}
}

func (l *approval) Info(ctx context.Context, req *domain.IdPathReq) (resp *domain.ApprovalInfoResp, err error) {
	approval, err := l.svcCtx.ApprovalModel.FindOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	resp = approval.ToDomainApprovalInfo()
	users, err := l.svcCtx.UserModel.ListToMaps(ctx, &domain.UserListReq{
		Ids: approval.Participation,
	})
	if err != nil || len(users) == 0 {
		return resp, err
	}
	resp.User = &domain.Approver{
		UserId:   users[approval.UserId].ID.Hex(),
		UserName: users[approval.UserId].Name,
	}
	resp.Approver = &domain.Approver{
		UserId:   users[approval.ApprovalId].ID.Hex(),
		UserName: users[approval.ApprovalId].Name,
	}
	for _, approver := range approval.Approvers {
		resp.Approvers = append(resp.Approvers, &domain.Approver{
			UserId:   users[approver.UserId].ID.Hex(),
			UserName: users[approver.UserId].Name,
			Status:   int(approval.Status),
			Reason:   approver.Reason,
		})
	}

	return
}

func (l *approval) Create(ctx context.Context, req *domain.Approval) (resp *domain.IdResp, err error) {
	uid := token.GetUId(ctx)
	req.UserId = uid
	approval := l.newApproval(req)

	var abstract string
	switch model.ApprovalType(req.Type) {
	case model.LeaveApproval:
		approval.Leave = &model.Leave{
			Type:      model.LeaveType(req.Leave.Type),
			StartTime: req.Leave.StartTime,
			EndTime:   req.Leave.EndTime,
			Reason:    req.Leave.Reason,
			TimeType:  model.TimeFormatType(req.Leave.TimeType),
		}
		abstract = fmt.Sprintf("【%s】: 【%s】-【%s】", model.LeaveType(req.Leave.Type).ToString(),
			timex.Format(req.Leave.StartTime), timex.Format(req.Leave.EndTime))
		approval.Reason = req.Leave.Reason
	case model.GoOutApproval:
		approval.GoOut = &model.GoOut{
			StartTime: req.GoOut.StartTime,
			EndTime:   req.GoOut.EndTime,
			Reason:    req.GoOut.Reason,
		}
		abstract = fmt.Sprintf("【%s】-【%s】", timex.Format(req.GoOut.StartTime), timex.Format(req.GoOut.EndTime))
		approval.Reason = req.GoOut.Reason
	case model.MakeCardApproval:
		approval.MakeCard = &model.MakeCard{
			Date:      req.MakeCard.Date,
			Reason:    req.MakeCard.Reason,
			Day:       req.MakeCard.Day,
			CheckType: model.WorkCheckType(req.MakeCard.CheckType),
		}
		abstract = fmt.Sprintf("【%s】【%s】", timex.Format(req.MakeCard.Date), req.MakeCard.Reason)
		approval.Reason = req.MakeCard.Reason
	default:
		// ...
	}

	user, err := l.svcCtx.UserModel.FindOne(ctx, uid)
	if err != nil {
		return
	}
	approval.Title = fmt.Sprintf("%s 提交的 %s", user.Name, model.ApprovalType(req.Type).ToString())
	approval.Abstract = abstract

	// 审批人
	depUser, err := l.svcCtx.DepartmentUserModel.FindByUserId(ctx, user.ID.Hex())
	if err != nil {
		return
	}
	dep, err := l.svcCtx.DepartmentModel.FindOne(ctx, depUser.DepId)
	if err != nil {
		return
	}

	parentIds := model.ParseParentPath(dep.ParentPath)
	pdeps, err := l.svcCtx.DepartmentModel.ListToMap(ctx, &domain.DepartmentListReq{
		DepId:  "",
		DepIds: parentIds,
	})
	var (
		approvals      []*model.Approver
		participations []string
	)
	approvals = append(approvals, &model.Approver{
		UserId: dep.LeaderId,
		Status: model.Processed,
	})
	participations = append(participations, dep.LeaderId, uid)

	for i := len(parentIds) - 1; i > 0; i-- {
		if _, ok := pdeps[parentIds[i]]; !ok {
			continue
		}
		approvals = append(approvals, &model.Approver{
			UserId: pdeps[parentIds[i]].LeaderId,
		})
		participations = append(participations, pdeps[parentIds[i]].LeaderId)
	}

	approval.Approvers = approvals
	approval.Participation = participations
	approval.ApprovalId = dep.LeaderId
	approval.UserId = uid

	if err = l.svcCtx.ApprovalModel.Insert(ctx, approval); err != nil {
		return
	}

	return &domain.IdResp{
		Id: approval.ID.Hex(),
	}, nil
}

func (l *approval) Dispose(ctx context.Context, req *domain.DisposeReq) (err error) {
	approval, err := l.svcCtx.ApprovalModel.FindOne(ctx, req.ApprovalId)
	if err != nil {
		return err
	}
	uid := token.GetUId(ctx)
	// 撤销
	if model.ApprovalStatus(req.Status) == model.Cancel {
		if req.ApprovalId != approval.UserId {
			return errors.New("审核用户错误")
		}

		approval.Status = model.Cancel

		return l.svcCtx.ApprovalModel.Update(ctx, approval)
	}

	// 通过或拒绝

	if approval.ApprovalId != uid {
		return errors.New("审核用户错误")
	}
	switch approval.Status {
	case model.Cancel:
		return errors.New("该审核已撤销")
	case model.Pass:
		return errors.New("该审核已通过")
	case model.Refuse:
		return errors.New("该审核已拒绝")
	}

	// 当前用户审批
	approval.Approvers[approval.ApprovalIdx].Status = model.Pass
	approval.Approvers[approval.ApprovalIdx].Reason = req.Reason

	// 切换下个审批人

	if model.ApprovalStatus(req.Status) == model.Pass && len(approval.Approvers)-1 > approval.ApprovalIdx {
		approval.ApprovalIdx++
		approval.ApprovalId = approval.Approvers[approval.ApprovalIdx].UserId
	} else {
		// 通过
		isPass := true
		for _, approver := range approval.Approvers {
			if approver.Status != model.Pass {
				isPass = false
				break
			}
		}

		if isPass {
			approval.Status = model.Pass
		}
	}

	return l.svcCtx.ApprovalModel.Update(ctx, approval)
}

func (l *approval) List(ctx context.Context, req *domain.ApprovalListReq) (resp *domain.ApprovalListResp, err error) {

	data, count, err := l.svcCtx.ApprovalModel.List(ctx, req)
	if err != nil {
		return nil, err
	}

	var list []*domain.ApprovalList
	for i, _ := range data {
		list = append(list, data[i].ToDomainApprovalList())
	}

	return &domain.ApprovalListResp{
		List:  list,
		Count: count,
	}, nil
}

func (l *approval) newApproval(req *domain.Approval) *model.Approval {
	return &model.Approval{
		ID:     primitive.NewObjectID(),
		UserId: req.UserId,
		No:     GenRandomNo(11),
		Type:   model.ApprovalType(req.Type),
		Status: model.Processed,
		Reason: req.Reason,
	}
}

func GenRandomNo(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb bytes.Buffer
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
