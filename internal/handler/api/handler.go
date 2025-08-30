package api

import (
	"ai/internal/handler"
	"ai/internal/middleware"
	"ai/internal/svc"
	"ai/pkg/httpx"
	"time"

	"gitee.com/dn-jinmin/tlog"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

// Handler 定义API处理器接口，所有具体的API处理器都需要实现InitRegister方法
// 用于统一注册路由的方式
type Handler interface {
	// InitRegister 初始化并注册路由到Gin引擎
	InitRegister(*gin.Engine)
}

// handle 实现了API服务的处理器结构体
// 包含Gin引擎实例和服务监听地址
type handle struct {
	srv  *gin.Engine // Gin框架的引擎实例，用于路由管理和HTTP处理
	addr string      // 服务监听的地址，格式为"IP:端口"
}

// NewHandle 创建一个新的API处理器实例
// 参数：svc 是服务上下文，包含配置、数据库连接等资源
// 返回值：新创建的handle实例
func NewHandle(svc *svc.ServiceContext) *handle {
	h := &handle{
		srv:  gin.Default(),
		addr: "0.0.0.0:8080",
	}

	// 如果配置中指定了服务地址，则使用配置中的地址
	if len(svc.Config.Addr) > 0 {
		h.addr = svc.Config.Addr
	}

	// 日志初始化
	tlog.Init(
		tlog.WithLoggerWriter(tlog.NewLoggerWriter()),
		tlog.WithLabel(svc.Config.Tlog.Label),
		tlog.WithMode(svc.Config.Tlog.Mode),
	)

	h.srv.Use(middleware.NewLog().Handler)
	h.srv.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:63342",
			"http://127.0.0.1:8888",
			"http://localhost:8888",
			"http://127.0.0.1:8088", // 反代静态页时常用
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 设置HTTP工具包的错误处理器为自定义的错误处理函数
	httpx.SetErrorHandler(handler.ErrorHandler)

	// 初始化所有API处理器（具体实现可能在initHandler函数中）
	handlers := initHandler(svc)

	// 遍历所有处理器，调用它们的InitRegister方法注册路由
	for _, handler := range handlers {
		handler.InitRegister(h.srv)
	}

	return h
}

// Run 启动HTTP服务，开始监听并处理请求
// 返回值：启动服务过程中可能出现的错误
func (h *handle) Run() error {
	// 调用Gin引擎的Run方法，在指定地址启动服务
	return h.srv.Run(h.addr)
}
