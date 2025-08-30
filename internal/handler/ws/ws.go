package ws

import (
	"ai/internal/domain"
	"ai/internal/logic"
	"ai/internal/model"
	"ai/internal/svc"
	"ai/token"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"gitee.com/dn-jinmin/tlog"

	"github.com/gorilla/websocket"
)

type Ws struct {
	websocket.Upgrader //WebSocket升级器，用于将HTTP连接升级为WebSocket连接
	sync.RWMutex       // 读写锁，用于安全地操作共享的连接映射

	svc         *svc.ServiceContext // 服务上下文，包含配置等共享资源
	tokenparser *token.Parse        // 令牌解析器，用于验证用户身份
	chat        logic.Chat          // 聊天业务逻辑处理组件

	uidToConn map[string]*websocket.Conn // 用户ID到WebSocket连接的映射
	ConnToUid map[*websocket.Conn]string // WebSocket连接到用户ID的映射
}

// NewWs 创建一个新的WebSocket服务实例
func NewWs(svc *svc.ServiceContext) *Ws {
	// 初始化日志配置
	tlog.Init(
		tlog.WithLoggerWriter(tlog.NewLoggerWriter()),
		tlog.WithLabel(svc.Config.Tlog.Label),
		tlog.WithMode(svc.Config.Tlog.Mode),
	)

	return &Ws{
		// 初始化WebSocket升级器，允许所有来源的连接
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			}, // 允许跨域请求
		},
		svc:         svc,
		chat:        logic.NewChat(svc),                         // 初始化聊天业务逻辑
		tokenparser: token.NewTokenParse(svc.Config.Jwt.Secret), // 初始化令牌解析器

		uidToConn: make(map[string]*websocket.Conn), // 初始化用户ID到连接的映射
		ConnToUid: make(map[*websocket.Conn]string), // 初始化连接到用户ID的映射
	}
}

// Run 启动WebSocket服务
func (s *Ws) Run() {
	// 注册WebSocket处理函数，路径为"/ws"
	http.HandleFunc("/ws", s.ServerWs)
	// 打印启动信息并开始监听指定地址
	fmt.Println("启动websocket服务", s.svc.Config.Ws.Addr)
	http.ListenAndServe(s.svc.Config.Ws.Addr, nil)
}

// ServerWs 处理WebSocket连接请求
func (s *Ws) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			tlog.ErrorCtx(r.Context(), "serverWs", e)
		}
	}()

	// 对连接进行鉴权，获取用户ID和令牌
	uid, token, err := s.auth(r)
	if err != nil {
		tlog.ErrorfCtx(r.Context(), "serverWs", "auth fail %v", err.Error())
		return
	}

	// 请求升级，设置WebSocket响应头，包含协议信息
	respHeader := http.Header{
		"sec-websocket-protocol": []string{token},
	}

	// 将HTTP连接升级为WebSocket连接
	c, err := s.Upgrade(w, r, respHeader)
	if err != nil {
		tlog.ErrorfCtx(r.Context(), "serverWs", "Upgrade fail %v", err.Error())
		return
	}

	// 记录新建立的连接
	s.addConn(c, uid)

	// 启动goroutine处理该连接的消息
	go s.handlerConn(c, uid, token)
}

// handlerConn 处理单个WebSocket连接的消息循环
func (s *Ws) handlerConn(conn *websocket.Conn, uid, tok string) {
	for {
		// 读取客户端发送的消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			tlog.Errorf("serverWs", "conn.ReadMessage fail %v, uid %v", err.Error(), uid)
			s.closeConn(conn)
			return
		}

		// 创建包含用户信息的上下文
		ctx := s.context(uid, tok)

		// 解析消息为Message结构体
		var req domain.Message
		if err := json.Unmarshal(msg, &req); err != nil {
			tlog.ErrorfCtx(ctx, "handlerConn", "json.Unmarshal fail %v", err.Error())
			return
		}
		req.SendId = uid

		// 根据消息类型分发处理
		switch model.ChatType(req.ChatType) {
		case model.SingleChatType:
			err = s.privateChat(ctx, conn, &req)
		case model.GroupChatType:
			err = s.groupChat(ctx, conn, &req)
		}

		if err != nil {
			tlog.ErrorfCtx(ctx, "handlerConn", "message handle fail %v, msg %v", err.Error(), req)
			return
		}
	}
}

// context 创建包含用户身份信息和日志追踪的上下文
func (s *Ws) context(uid, tok string) context.Context {
	ctx := context.WithValue(context.Background(), token.Identify, uid)
	ctx = context.WithValue(ctx, token.Authorization, tok)

	return tlog.TraceStart(ctx)
}

// addConn 将新连接添加到映射中，线程安全
func (s *Ws) addConn(conn *websocket.Conn, uid string) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	if c := s.uidToConn[uid]; c != nil {
		c.Close()
	}

	s.uidToConn[uid] = conn
	s.ConnToUid[conn] = uid
}

// closeConn 关闭连接并从映射中移除，线程安全
func (s *Ws) closeConn(conn *websocket.Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.ConnToUid[conn]
	if uid == "" {
		return
	}

	fmt.Printf("关闭 %s 连接\n", uid)

	delete(s.ConnToUid, conn)
	delete(s.uidToConn, uid)

	conn.Close()
}

// send 向指定连接发送消息
func (s *Ws) send(ctx context.Context, conn *websocket.Conn, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		tlog.ErrorCtx(ctx, "conn.send", err.Error())
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, b)
}

// sendByUids 向指定用户列表发送消息，支持广播（uids为空时）
func (s *Ws) sendByUids(ctx context.Context, msg interface{}, uids ...string) error {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	if len(uids) == 0 {
		for i, _ := range s.uidToConn {
			if err := s.send(ctx, s.uidToConn[i], msg); err != nil {
				tlog.ErrorCtx(ctx, "sendByUids.all.send", err.Error())
				return err
			}
		}
	}

	for _, uid := range uids {
		c, ok := s.uidToConn[uid]
		if !ok {
			continue
		}
		if err := s.send(ctx, c, msg); err != nil {
			tlog.ErrorfCtx(ctx, "sendByUids.one.send err %v, uid %v", err.Error(), uid)
			return err
		}
	}
	return nil
}

// auth 验证WebSocket连接的身份
func (s *Ws) auth(r *http.Request) (uid string, tokenStr string, err error) {
	tok := r.Header.Get("sec-websocket-protocol")
	if tok == "" {
		return "", "", errors.New("没有登入，不存在访问权限")
	}

	claims, tokenStr, err := s.tokenparser.ParseToken(tok)
	if err != nil {
		return "", "", err
	}

	return claims[token.Identify].(string), tokenStr, nil
}
