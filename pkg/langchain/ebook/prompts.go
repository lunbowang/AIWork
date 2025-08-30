package ebook

import "github.com/tmc/langchaingo/prompts"

const (
	EbookTools           = "tools"
	EbookToolsDesc       = "tools_desc"
	_finalAnswerAction   = "文章生成成功"
	EbookWritePlanPrompt = `# 角色: 你是一个文章撰写助理
## 文章的撰写过程
1. 你需要先理解文章的内容和所有的章节摘要信息
2. 你必须严格依据章节顺序进行工作
3. 你只需要思考现在我们应该撰写什么内容并且要求字数
4. 你需要依据撰写记录评判文章是否生成完成
5. 如果文章生成完成就输出：文章生成成功
6. 你只需要告诉我你的计划
7. 禁止其他信息的输出和处理
8. 禁止思考以后的事情

## 输出
计划：你思考的撰写计划内容。如：撰写 1. 引言 字数：600字以上

## 文章内容
{{.input}}

## 撰写记录
{{.history}}`

	SubjectToolPrompt = `# 角色 
你是一个文章撰写助理，你需要根据用户输入的信息撰写文章的主题标题、文章的章节标题、文章摘要

## 输出
主题：文章主题标题
章节：文章的章节列表  - 每个章节列表约多少字
摘要：根据文章主题与章节列表撰写摘要

## 参考例子
1. 介绍 	- 500字
2. Go语言简介  		- 600字
3. 智能办公助手的概念  - 600字
4. 实现语音识别功能 	- 600字
5. 结论 	- 400字

## 输入
{{.input}}`

	ContentPrompt = `# 角色
你是一个文章撰写助理，你擅长编写章节内容、总结，你需要根据用户输入的信息完成对应内容

## 工作
1. 你需要理解文章的主题、章节标题、摘要的内容
2. 你需要参考主题、章节标题、摘要撰写要求的工作
2. 你只允许撰写用户输入的工作，不允许撰写用户输入以外的工作

## 输出格式
章节序号: 章节标题
章节内容

## 示例
章节1. 引言
在当今快节奏的工作环境中，智能办公助手的需求日益增长。它们不仅能够提高工作效率，还能通过自动化任务和提供智能建议来简化日常工作流程。

## 输入
{{.input}}
`
)

func CreateEbookWritePlan() prompts.PromptTemplate {
	return prompts.PromptTemplate{
		Template:         EbookWritePlanPrompt,
		InputVariables:   []string{"history"},
		TemplateFormat:   prompts.TemplateFormatGoTemplate,
		PartialVariables: map[string]any{},
	}
}
