package ws

import (
	"ai/internal/domain"
	"context"

	"github.com/gorilla/websocket"
)

func (s *Ws) privateChat(ctx context.Context, conn *websocket.Conn, req *domain.Message) error {
	if err := s.chat.PrivateChat(ctx, req); err != nil {
		return err
	}

	return s.sendByUids(ctx, req, req.RecvId)
}

func (s *Ws) groupChat(ctx context.Context, conn *websocket.Conn, req *domain.Message) error {
	if _, err := s.chat.GroupChat(ctx, req); err != nil {
		return err
	}

	return s.sendByUids(ctx, req)
}
