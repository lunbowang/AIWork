package model

import (
	"ai/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1 我提交、2 我审核
type ApprovalOptionType int

const (
	ApprovalSubmit ApprovalOptionType = iota + 1
	ApprovalAudit
)

// ApprovalType 审批类型
// 1.通用，2.请假，3.补卡，4.外出，5.报销，6.付款，7.采购，8.收款
type ApprovalType int

const (
	UniversalApproval     ApprovalType = 1  // 通用
	LeaveApproval         ApprovalType = 2  // 请假
	MakeCardApproval      ApprovalType = 3  // 补卡
	GoOutApproval         ApprovalType = 4  // 外出
	ReimburseApproval     ApprovalType = 5  // 报销
	PaymentApproval       ApprovalType = 6  // 付款
	BuyerApproval         ApprovalType = 7  // 采购
	ProceedsApproval      ApprovalType = 8  // 收款
	PositiveApproval      ApprovalType = 9  // 转正
	DimissionApproval     ApprovalType = 10 // 离职
	OvertimeApproval      ApprovalType = 11 // 加班
	BuyerContractApproval ApprovalType = 12 // 合同
)

func (t ApprovalType) ToString() string {
	switch t {
	case LeaveApproval:
		return "请假审批"
	case MakeCardApproval:
		return "补卡审批"
	case GoOutApproval:
		return "外出审批"
	case PaymentApproval:
		return "付款审批"
	case ProceedsApproval:
		return "收款审批"
	case BuyerApproval:
		return "采购审批"
	case ReimburseApproval:
		return "报销审批"
	case PositiveApproval:
		return "转正审批"
	case OvertimeApproval:
		return "加班审批"
	case BuyerContractApproval:
		return "采购合同审批"
	case UniversalApproval:
		return "通用审批"
	default:
		return ""
	}
}

// 0. 没有开始 ，1. 进行中 2. 完成-通过 ，3. 撤销， 4. 拒绝
// ApprovalStatus 审批状态
type ApprovalStatus int

const (
	Notstarted ApprovalStatus = iota //未开始
	Processed                        //处理中
	Pass                             //通过
	Refuse                           //拒绝
	Cancel                           //撤销
	AutoPass                         //自动通过
)

// LeaveType 请假类型
// 0.事假, 1.调休, 2.病假, 3.年假, 4.产假, 5.陪产假, 6.婚假, 7.丧假, 8.哺乳假
type LeaveType int

const (
	Matter        LeaveType = iota + 1 //事假
	Rest                               //调休
	Fall                               //病假
	Annual                             //年假
	Maternity                          //产假
	Paternity                          //陪产假
	Marriage                           //婚假
	Funeral                            //丧假
	Breastfeeding                      //哺乳假
)

func (t LeaveType) ToString() string {
	switch t {
	case Matter:
		return "事假"
	case Rest:
		return "调休"
	case Fall:
		return "病假"
	case Annual:
		return "年假"
	case Maternity:
		return "产假"
	case Paternity:
		return "陪产假"
	case Marriage:
		return "婚假"
	case Funeral:
		return "丧假"
	case Breastfeeding:
		return "哺乳假"
	}
	return ""
}

// WorkCheckType 打卡类型
// 1. 上班卡, 2. 下班卡
type WorkCheckType int

const (
	OnWorkCheck  WorkCheckType = 1 // 上班
	OffWorkCheck WorkCheckType = 2 // 下班
)

// 1. 小时， 2. 天，3. 半天，4. 上半天， 5. 下半天
type TimeFormatType int

const (
	HourTimeFormatType TimeFormatType = iota + 1
	DayTimeFormatType
)

type (
	Approval struct {
		ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		// TODO: Fill your own fields
		UserId   string         `bson:"userId,omitempty" json:"userId,omitempty"`
		No       string         `bson:"no,omitempty" json:"no,omitempty"`
		Type     ApprovalType   `bson:"type,omitempty" json:"type,omitempty"`
		Status   ApprovalStatus `bson:"status,omitempty" json:"status,omitempty"`
		Title    string         `bson:"title,omitempty" json:"title,omitempty"`
		Abstract string         `bson:"abstract,omitempty" json:"abstract,omitempty"`
		Reason   string         `bson:"reason,omitempty" json:"reason,omitempty"`

		ApprovalId    string      `bson:"approvalId,omitempty"`
		ApprovalIdx   int         `bson:"approvalIdx,omitempty"`
		Approvers     []*Approver `base:"approvers,omitempty"`
		CopyPersons   []*Approver `base:"copyPersons,omitempty"`
		Participation []string    `bson:"participation,omitempty"`

		FinishAt    int64 `bson:"finishAt,omitempty" json:"finishAt,omitempty"`
		FinishDay   int64 `bson:"finishDay,omitempty" json:"finishDay,omitempty"`
		FinishMonth int64 `bson:"finishMonth,omitempty" json:"finishMonth,omitempty"`
		FinishYeas  int64 `bson:"finishYeas,omitempty" json:"finishYeas,omitempty"`

		MakeCard *MakeCard `bson:"makeCard,omitempty" json:"makeCard,omitempty"`
		Leave    *Leave    `bson:"leave,omitempty" json:"leave,omitempty"`
		GoOut    *GoOut    `bson:"goOut,omitempty" json:"goOut,omitempty"`

		UpdateAt int64 `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
		CreateAt int64 `bson:"createAt,omitempty" json:"createAt,omitempty"`
	}
	Approver struct {
		UserId   string         `bson:"userId,omitempty"`
		UserName string         `bson:"userName,omitempty"`
		Status   ApprovalStatus `bson:"status,omitempty"`
		Reason   string         `bson:"reason,omitempty"`
	}

	// MakeCard 补卡
	MakeCard struct {
		Date      int64         `bson:"date,omitempty"`          //补卡时间
		Reason    string        `bson:"reason,omitempty"`        //补卡理由
		Day       int64         `bson:"day,omitempty"`           //补卡日期(20221011)
		CheckType WorkCheckType `bson:"workCheckType,omitempty"` //补卡类型
	}

	// Leave 请假
	Leave struct {
		Type      LeaveType      `bson:"type,omitempty"`      //请假类型
		StartTime int64          `bson:"startTime,omitempty"` //开始时间
		EndTime   int64          `bson:"endTime,omitempty"`   //结束时间
		Reason    string         `bson:"reason,omitempty"`    //请假原由
		TimeType  TimeFormatType `bson:"timeType,omitempty"`  //请假类型  1=小时 2=天
	}

	// GoOut 外出
	GoOut struct {
		StartTime int64  `bson:"startTime,omitempty"` //开始时间
		EndTime   int64  `bson:"endTime,omitempty"`   //结束时间
		Reason    string `bson:"reason,omitempty"`    //请假原由
	}

	// ..
)

func (m *Approval) ToDomainApprovalInfo() *domain.ApprovalInfoResp {
	res := &domain.ApprovalInfoResp{
		Id:          m.ID.Hex(),
		No:          m.No,
		Type:        int(m.Type),
		Status:      int(m.Status),
		Title:       m.Title,
		Abstract:    m.Abstract,
		Reason:      m.Reason,
		Approvers:   nil,
		FinishAt:    m.FinishAt,
		FinishDay:   m.FinishDay,
		FinishMonth: m.FinishMonth,
		FinishYeas:  m.FinishYeas,
		UpdateAt:    m.UpdateAt,
		CreateAt:    m.CreateAt,
	}

	switch ApprovalType(res.Type) {
	case LeaveApproval:
		res.Leave = &domain.Leave{
			Type:      int(m.Leave.Type),
			StartTime: m.Leave.StartTime,
			EndTime:   m.Leave.EndTime,
			Reason:    m.Leave.Reason,
			TimeType:  int(m.Leave.TimeType),
		}
	case MakeCardApproval:
		res.MakeCard = &domain.MakeCard{
			Date:      m.MakeCard.Date,
			Reason:    m.MakeCard.Reason,
			Day:       m.MakeCard.Day,
			CheckType: int(m.MakeCard.CheckType),
		}
	case GoOutApproval:
		res.GoOut = &domain.GoOut{
			StartTime: m.GoOut.StartTime,
			EndTime:   m.GoOut.EndTime,
			Reason:    m.GoOut.Reason,
		}
	}

	return res
}

func (m *Approval) ToDomainApprovalList() *domain.ApprovalList {
	return &domain.ApprovalList{
		Id:       m.ApprovalId,
		Type:     int(m.Type),
		Status:   int(m.Status),
		Title:    m.Title,
		Abstract: m.Abstract,
	}
}
