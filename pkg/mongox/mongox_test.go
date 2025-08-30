package mongox

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"sort"
	"testing"
	"time"
)

func TestUpdate_one(t *testing.T) {
	db, _ := MongodbDatabase(&MongodbConfig{
		Hosts:    []string{"127.0.0.1"},
		Port:     27017,
		Database: "aiworkc",
	})

	col := db.Collection("chat_log")
	type ChatLog struct {
		ConversationId string   `bson:"conversationId"`
		SendId         string   `bson:"sendId"`
		RecvId         string   `bson:"recvId"`
		ChatType       ChatType `bson:"chatType"`
		MsgContent     string   `bson:"msgContent"`
		SendTime       int64    `bson:"sendTime"`
		UpdateAt       int64    `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
		CreateAt       int64    `bson:"createAt,omitempty" json:"createAt,omitempty"`
	}
	ts := time.Now().Unix()
	uid1 := "66af51b514fd35e240c9ab92"
	uid2 := "66af525714fd35e240c9ab93"
	conversationId := GenerateUniqueID(uid1, uid2)
	data := []ChatLog{
		{
			ConversationId: conversationId,
			SendId:         uid1,
			RecvId:         uid2,
			ChatType:       SingleChatType,
			MsgContent:     "最近的项目进展如何？有没有遇到什么问题？",
			SendTime:       ts,
			CreateAt:       ts,
		}, {
			ConversationId: conversationId,
			SendId:         uid2,
			RecvId:         uid1,
			ChatType:       SingleChatType,
			MsgContent:     "项目进展还算顺利，但在集成测试阶段，发现了几个接口不兼容的问题。",
			SendTime:       ts + 20,
			CreateAt:       ts + 20,
		}, {
			ConversationId: conversationId,
			SendId:         uid1,
			RecvId:         uid2,
			ChatType:       SingleChatType,
			MsgContent:     "具体是哪些接口？我们需不需要调整一下技术方案？",
			SendTime:       ts + 100,
			CreateAt:       ts + 100,
		}, {
			ConversationId: conversationId,
			SendId:         uid2,
			RecvId:         uid1,
			ChatType:       SingleChatType,
			MsgContent:     "主要是前端与后端的数据格式不一致，暂时可以通过一些临时解决方案绕过，但我建议尽快进行一些调整，以免影响后续的开发。",
			SendTime:       ts + 110,
			CreateAt:       ts + 110,
		}, {
			ConversationId: conversationId,
			SendId:         uid1,
			RecvId:         uid2,
			ChatType:       SingleChatType,
			MsgContent:     "好的。这部分你预计需要多长时间来解决？",
			SendTime:       ts + 130,
			CreateAt:       ts + 130,
		}, {
			ConversationId: conversationId,
			SendId:         uid2,
			RecvId:         uid1,
			ChatType:       SingleChatType,
			MsgContent:     "如果能够尽快确认接口规范，预计一周内能完全解决。",
			SendTime:       ts + 170,
			CreateAt:       ts + 170,
		}, {
			ConversationId: conversationId,
			SendId:         uid1,
			RecvId:         uid2,
			ChatType:       SingleChatType,
			MsgContent:     "一周的时间可以接受，但我希望你能在下周三之前给我一个明确的解决方案和进度更新。",
			SendTime:       ts + 200,
			CreateAt:       ts + 200,
		}, {
			ConversationId: conversationId,
			SendId:         uid2,
			RecvId:         uid1,
			ChatType:       SingleChatType,
			MsgContent:     "明白了，我会尽快确认接口规范，并在周三之前给您反馈。",
			SendTime:       ts + 220,
			CreateAt:       ts + 220,
		}, {
			ConversationId: conversationId,
			SendId:         uid1,
			RecvId:         uid2,
			ChatType:       SingleChatType,
			MsgContent:     "好的。我会安排一个小组会议，让大家落实接口规范的问题。这样可以做到快速反馈。",
			SendTime:       ts + 290,
			CreateAt:       ts + 290,
		}, {
			ConversationId: conversationId,
			SendId:         uid2,
			RecvId:         uid1,
			ChatType:       SingleChatType,
			MsgContent:     "谢谢！另外，我想询问一下关于下周的发布计划，是否已经定下来？",
			SendTime:       ts + 330,
			CreateAt:       ts + 330,
		}, {
			ConversationId: conversationId,
			SendId:         uid1,
			RecvId:         uid2,
			ChatType:       SingleChatType,
			MsgContent:     "是的，发布计划已经初步确定在下周五。我们需要加快测试进度，确保没有重大bug。",
			SendTime:       ts + 350,
			CreateAt:       ts + 350,
		}, {
			ConversationId: conversationId,
			SendId:         uid2,
			RecvId:         uid1,
			ChatType:       SingleChatType,
			MsgContent:     "明白了，我这边会督促测试团队，提高效率。",
			SendTime:       ts + 410,
			CreateAt:       ts + 410,
		}, {
			ConversationId: conversationId,
			SendId:         uid1,
			RecvId:         uid2,
			ChatType:       SingleChatType,
			MsgContent:     "很好，还有其他需要讨论的事项吗？",
			SendTime:       ts + 450,
			CreateAt:       ts + 450,
		}, {
			ConversationId: conversationId,
			SendId:         uid2,
			RecvId:         uid1,
			ChatType:       SingleChatType,
			MsgContent:     "目前没有了。如果有新的问题，我会及时反馈。",
			SendTime:       ts + 470,
			CreateAt:       ts + 470,
		}, {
			ConversationId: conversationId,
			SendId:         uid1,
			RecvId:         uid2,
			ChatType:       SingleChatType,
			MsgContent:     "好的，那就先这样，辛苦了。",
			SendTime:       ts + 500,
			CreateAt:       ts + 500,
		},
	}

	for _, datum := range data {
		t.Log(col.InsertOne(context.Background(), datum))
	}
}

func TestUpdate_group(t *testing.T) {
	db, _ := MongodbDatabase(&MongodbConfig{
		Hosts:    []string{"127.0.0.1"},
		Port:     27017,
		Database: "aiworkc",
	})

	col := db.Collection("chat_log")
	type ChatLog struct {
		ConversationId string   `bson:"conversationId"`
		SendId         string   `bson:"sendId"`
		RecvId         string   `bson:"recvId"`
		ChatType       ChatType `bson:"chatType"`
		MsgContent     string   `bson:"msgContent"`
		SendTime       int64    `bson:"sendTime"`
		UpdateAt       int64    `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
		CreateAt       int64    `bson:"createAt,omitempty" json:"createAt,omitempty"`
	}
	ts := time.Now().Unix() + 10000
	// root 山药
	shanyao := "66af51b514fd35e240c9ab92"
	// 小闵
	xiaomin := "66c349356c772061c5daf208"
	// 木兮
	muxi := "66c349666c772061c5daf20a"
	xiaoxiao := "66c3493e6c772061c5daf209"
	cara := "66c349776c772061c5daf20b"

	cid := "PU0zXtG2ePJzJpkfbE5gm/"
	data := []ChatLog{
		{
			ConversationId: cid,
			SendId:         shanyao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "我们今天主要讨论一下项目的进展情况，以及接下来要解决的几个问题。小闵，你先说说后端的进展吧。",
			SendTime:       ts,
			CreateAt:       ts,
		}, {
			ConversationId: cid,
			SendId:         xiaomin,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "好的。后端这边的功能开发大致完成了，现在正在进行一些接口的调试和优化。预计这周三能完成所有接口的开发和测试。",
			SendTime:       ts + 20,
			CreateAt:       ts + 20,
		}, {
			ConversationId: cid,
			SendId:         shanyao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "很好。但你们的接口文档是否更新了？上次会议上提到的文档问题解决了吗？",
			SendTime:       ts + 100,
			CreateAt:       ts + 100,
		}, {
			ConversationId: cid,
			SendId:         xiaomin,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "我已经更新了，昨天发给了大家。请大家尽快查看，特别是小小，你负责测试，能否在周四之前给出反馈？",
			SendTime:       ts + 110,
			CreateAt:       ts + 110,
		}, {
			ConversationId: cid,
			SendId:         xiaoxiao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "没问题，我会尽快查阅。之前我在测试中也发现了一些小问题，这两天会尽量整理出来。",
			SendTime:       ts + 130,
			CreateAt:       ts + 130,
		}, {
			ConversationId: cid,
			SendId:         shanyao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "小小，关于你之前提到的那几个问题，能否在明天下午之前把整理好的列表发给大家？",
			SendTime:       ts + 170,
			CreateAt:       ts + 170,
		}, {
			ConversationId: cid,
			SendId:         xiaoxiao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "可以的，我会争取在明天下午三点之前发给你们。",
			SendTime:       ts + 200,
			CreateAt:       ts + 200,
		}, {
			ConversationId: cid,
			SendId:         shanyao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "好的，那木兮，前端这边有什么进展？",
			SendTime:       ts + 220,
			CreateAt:       ts + 220,
		}, {
			ConversationId: cid,
			SendId:         muxi,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "前端的界面设计已经完成，现在在进行组件的实现。原定周五完成，但由于有几处细节需要改动，可能会延迟到下周一。",
			SendTime:       ts + 290,
			CreateAt:       ts + 290,
		}, {
			ConversationId: cid,
			SendId:         shanyao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "下周一必须完成，如果时间紧，是否需要增加人手协助？",
			SendTime:       ts + 330,
			CreateAt:       ts + 330,
		}, {
			ConversationId: cid,
			SendId:         muxi,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "如果能调配到一名同事帮忙的话，进度会快很多。",
			SendTime:       ts + 350,
			CreateAt:       ts + 350,
		}, {
			ConversationId: cid,
			SendId:         shanyao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "那我联系一下其他部门，看能否借调一位同事。还有你提到的请假问题，木兮，你大概什么时候请假？",
			SendTime:       ts + 410,
			CreateAt:       ts + 410,
		}, {
			ConversationId: cid,
			SendId:         muxi,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "我打算下周四请一天假，去参加一个前端开发的研讨会。会议的内容会对我们的项目有所帮助。",
			SendTime:       ts + 450,
			CreateAt:       ts + 450,
		}, {
			ConversationId: cid,
			SendId:         shanyao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "好的，尽量确保请假前的工作完成，不要影响项目进度。cara，你这边能否在周四之前把相关接口梳理清楚，方便木兮对接？",
			SendTime:       ts + 470,
			CreateAt:       ts + 470,
		}, {
			ConversationId: cid,
			SendId:         cara,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "没问题，明天一早我会把接口的使用说明整理好，确保木兮能顺利接入。",
			SendTime:       ts + 500,
			CreateAt:       ts + 500,
		}, {
			ConversationId: cid,
			SendId:         shanyao,
			RecvId:         cid,
			ChatType:       GroupChatType,
			MsgContent:     "很好。今天的讨论到此为止，周四我们再开一个会，检查进展情况。大家每天有问题随时沟通，确保按时完成各自的工作。",
			SendTime:       ts + 560,
			CreateAt:       ts + 560,
		},
	}

	for _, datum := range data {
		col.InsertOne(context.Background(), datum)
	}
}

type ChatType int

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)

func GenerateUniqueID(id1, id2 string) string {
	// 将两个 ID 放入切片中
	ids := []string{id1, id2}

	// 对 IDs 切片进行排序
	sort.Strings(ids)

	// 将排序后的 ID 组合起来
	combined := ids[0] + ids[1]

	// 创建 SHA-256 哈希对象
	hasher := sha256.New()

	// 写入合并后的字符串
	hasher.Write([]byte(combined))

	// 计算哈希值
	hash := hasher.Sum(nil)

	// 返回哈希值的十六进制字符串表示
	return base64.RawStdEncoding.EncodeToString(hash)[:22] // 可以选择更短的长度
}
