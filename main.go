package main

import (
	"ai/internal/config"
	"ai/internal/handler/api"
	"ai/internal/handler/ws"
	"ai/internal/svc"
	"ai/pkg/conf"
	"flag"
	"sync"
)

type Serve interface {
	Run() error
}

const (
	Api = "api"

	// add other module
)

var (
	// 定义命令行参数：配置文件路径，默认值为./etc/api.yaml
	configFile = flag.String("f", "./etc/api.yaml", "the config file")
	// 定义命令行参数：服务运行模式，默认值为api
	modeType = flag.String("m", "api", "server run mod")

	// 等待组，用于等待所有服务goroutine执行完成
	sw sync.WaitGroup
)

func main() {
	// 解析命令行参数，将命令行输入的值绑定到上面定义的变量
	flag.Parse()

	// 声明配置结构体变量，用于存储加载的配置
	var cfg config.Config
	// 加载配置文件到cfg结构体中
	conf.MustLoad(*configFile, &cfg)

	// 向等待组添加一个计数，表示有一个goroutine需要等待
	sw.Add(1)
	// 启动一个goroutine来运行服务，避免阻塞主goroutine
	go func() {
		// 当goroutine退出时，减少等待组的计数
		defer sw.Done()
		// 创建服务上下文，初始化服务所需的各种资源（如数据库连接、缓存等）
		svc, err := svc.NewServiceContext(cfg)
		if err != nil {
			panic(err)
		}
		// 创建API处理器实例，传入服务上下文
		srv := api.NewHandle(svc)
		// 启动API服务
		srv.Run()
	}()

	sw.Add(1)
	go func() {
		defer sw.Done()
		svc, err := svc.NewServiceContext(cfg)
		if err != nil {
			panic(err)
		}
		srv := ws.NewWs(svc)
		srv.Run()
	}()

	sw.Wait()

	//var srv Serve
	//switch *modeType {
	//case Api:
	//	srv = api.NewHandle(svc)
	//// add other module case
	//default:
	//	panic("请指定正确的服务")
	//}
	//
	//srv.Run()
}
