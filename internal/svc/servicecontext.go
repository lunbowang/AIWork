package svc

import (
	"ai/internal/config"
	"ai/internal/middleware"
	"ai/internal/model"
	"ai/pkg/langchain/callbackx"
	"ai/pkg/mongox"
	"ai/token"
	"context"
	"errors"

	"gitee.com/dn-jinmin/tlog"

	openaiSdk "github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms/openai"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrAuth = errors.New("不具有权限")

type ServiceContext struct {
	*middleware.Jwt

	Config config.Config

	Mongo *mongo.Database

	model.UserModel
	model.DepartmentModel
	model.DepartmentUserModel
	model.UserTodoModel
	model.TodoModel
	model.ApprovalModel
	model.ChatlogModel

	LLMs           *openai.LLM
	AliProxyOpenai *openaiSdk.Client
	OpenaiClient   *openaiSdk.Client
	Callbacks      callbacks.Handler

	Auth func(ctx context.Context) error
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {

	// 初始化MongoDB数据库连接
	mongoDb, err := mongox.MongodbDatabase(&mongox.MongodbConfig{
		User:     c.Mongo.User,
		Password: c.Mongo.Password,
		Hosts:    c.Mongo.Hosts,
		Port:     c.Mongo.Port,
		Database: c.Mongo.Database,
	})
	if err != nil {
		return nil, err
	}

	logger := tlog.NewLogger()
	callbacks := callbacks.CombiningHandler{
		Callbacks: []callbacks.Handler{
			callbackx.NewLogHandler(logger),
			callbackx.NewTitTokenHandle(logger),
		},
	}

	opts := []openai.Option{
		openai.WithBaseURL(c.Langchain.Url),
		openai.WithToken(c.Langchain.ApiKey),
		openai.WithCallback(callbacks),
		openai.WithEmbeddingModel("text-embedding-ada-002"),
		openai.WithModel("gpt-3.5-turbo"),
	}

	llm, err := openai.New(opts...)
	if err != nil {
		return nil, err
	}

	aliGPTCfg := openaiSdk.DefaultConfig(c.AliGPT.ApiKey)
	aliGPTCfg.BaseURL = c.AliGPT.Url
	aliProxyOpenAi := openaiSdk.NewClientWithConfig(aliGPTCfg)
	openaiGPTCfg := openaiSdk.DefaultConfig(c.Langchain.ApiKey)
	openaiGPTCfg.BaseURL = c.Langchain.Url
	openaiGPT := openaiSdk.NewClientWithConfig(aliGPTCfg)

	// 创建用户数据模型实例，传入MongoDB连接
	userModel := model.NewUserModel(mongoDb)
	svc := &ServiceContext{
		Config:              c,
		Jwt:                 middleware.NewJwt(c.Jwt.Secret),
		Mongo:               mongoDb,
		UserModel:           userModel,
		DepartmentUserModel: model.NewDepartmentUserModel(mongoDb),
		DepartmentModel:     model.NewDepartmentModel(mongoDb),
		UserTodoModel:       model.NewUserTodoModel(mongoDb),
		TodoModel:           model.NewTodoModel(mongoDb),
		ApprovalModel:       model.NewApprovalModel(mongoDb),
		ChatlogModel:        model.NewChatlogModel(mongoDb),

		LLMs:           llm,
		Callbacks:      callbacks,
		AliProxyOpenai: aliProxyOpenAi,
		OpenaiClient:   openaiGPT,

		Auth: func(ctx context.Context) error {
			uid := token.GetUId(ctx)
			if uid == "" {
				return ErrAuth
			}

			user, err := userModel.FindOne(ctx, uid)
			if err != nil {
				return err
			}

			if !user.IsSystem {
				return ErrAuth
			}

			return nil
		},
	}

	return svc, initUser(svc)
}

// initUser 初始化用户数据，确保系统存在默认的管理员用户
// 参数：svc 是服务上下文实例
// 返回值：初始化过程中可能出现的错误
func initUser(svc *ServiceContext) error {
	ctx := context.Background()

	// 尝试查询系统用户
	systemUser, err := svc.UserModel.FindSysStemUser(ctx)
	if err != nil && err != model.ErrNotUser {
		return err
	}
	if systemUser != nil {
		return nil
	}

	// 如果系统用户不存在，则创建默认的root用户
	// 密码是经过bcrypt加密的"000000"（示例）
	return svc.UserModel.Insert(ctx, &model.User{
		Name:     "root",
		Password: "$2a$10$ddIvqt7U6zNA9poys.FNCuEZTJY6V.axWy4P7A44TuT9KBegGZlD6",
		Status:   0,
		IsSystem: true,
	})
}
