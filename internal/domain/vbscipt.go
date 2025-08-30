package domain

import "time"

type AiChatType int

const (
	DefaultHandler = iota
	TodoFind
	TodoAdd

	ApprovalFind

	ChatLog

	ImgAndText // 图片+文本
	File       //
)

type ChatFile struct {
	Path string
	Name string
	Time time.Time
}

type ImgAndTextResp struct {
	Text string `json:"text"`
	Url  string `json:"url"`
}
type ChatFileResp struct {
	Url string `json:"url"`
}
