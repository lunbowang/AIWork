package domain

type Message struct {
	ConversationId string `json:"conversationId"`

	RecvId string `json:"recvId"`
	SendId string `json:"sendId"`

	ChatType    int    `json:"chatType"`
	Content     string `json:"content"`
	ContentType int    `json:"contentType"`
}
