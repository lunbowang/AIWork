package api

import (
	"ai/internal/domain"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"

	"ai/internal/logic"
	"ai/internal/svc"
	"ai/pkg/httpx"
)

type Upload struct {
	svcCtx *svc.ServiceContext
	chat   logic.Chat
}

func NewUpload(svcCtx *svc.ServiceContext, chat logic.Chat) *Upload {
	return &Upload{
		svcCtx: svcCtx,
		chat:   chat,
	}
}

func (h *Upload) InitRegister(engine *gin.Engine) {
	g := engine.Group("v1/upload", h.svcCtx.Jwt.Handler)
	g.POST("/file", h.File)
	g.POST("/multiplefiles", h.Multiplefiles)
}

// File 处理文件上传的方法
// 接收HTTP请求中的文件，保存到本地，并返回文件信息
func (h *Upload) File(ctx *gin.Context) {
	// 从请求中获取名为"file"的上传文件
	// file: 文件内容的读取流
	// header: 文件的元信息（如文件名、大小等）
	// err: 可能出现的错误（如没有上传文件、文件过大等）
	file, header, err := ctx.Request.FormFile("file")
	var (
		filename string                 // 用于存储生成的唯一文件名
		buf      = bytes.NewBuffer(nil) // 缓冲区，用于临时存储文件内容
	)
	// 延迟关闭文件流，确保资源释放
	defer file.Close()

	if _, err := io.Copy(buf, file); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	// 生成唯一文件名：使用ksuid生成唯一ID + 保留原文件的扩展名
	// ksuid是一种分布式唯一ID生成算法，确保文件名不重复
	filename = ksuid.New().String() + filepath.Ext(header.Filename)

	// 创建新文件：保存路径由配置文件指定，文件名使用上面生成的唯一名称
	newFile, err := os.Create(h.svcCtx.Config.Upload.SavePath + filename)
	if err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}
	// 延迟关闭新创建的文件，确保内容写入完成
	defer newFile.Close()

	// 将缓冲区中的文件内容写入到新创建的文件中
	if _, err := newFile.Write(buf.Bytes()); err != nil {
		httpx.FailWithErr(ctx, err)
		return
	}

	// 构建文件响应信息
	resp := domain.FileResp{
		Host:     h.svcCtx.Config.Host,
		File:     fmt.Sprintf("%s%s", h.svcCtx.Config.Upload.SavePath, filename),
		Filename: filename,
	}

	// 从请求表单中获取"chat"字段的值
	// 如果该字段存在且不为空，说明需要将上传的文件关联到聊天功能
	chat := ctx.Request.FormValue("chat")
	if len(chat) > 0 {
		// 调用聊天服务的File方法，处理文件与聊天的关联（如解析文件内容到知识库等）
		h.chat.File(ctx.Request.Context(), []*domain.FileResp{
			&resp,
		})
	}

	// 返回最终处理结果
	if err != nil {
		httpx.FailWithErr(ctx, err)
	} else {
		httpx.OkWithData(ctx, resp)
	}
}

func (h *Upload) Multiplefiles(ctx *gin.Context) {
}
