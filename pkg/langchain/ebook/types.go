package ebook

type (
	Action struct {
		Tool  string
		Input string
	}

	Step struct {
		Subject string // 主题、章节列表、摘要
		Record  string // 记录文件的撰写进度
	}
)
