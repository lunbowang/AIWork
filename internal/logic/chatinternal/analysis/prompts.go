package analysis

import (
	"ai/pkg/langchain"

	"github.com/tmc/langchaingo/outputparser"
)

const (
	_destinations = "destinations"

	_defaultDBQuery = `
	# 你是一个助手，你将理解用户输入的信息，根据问题回答

	## 工作
	1. 你需要先优化用户输入的提示词，让它表达更准确
	2. 分析用户的问题是选择那个数据库，并输出数据库名
	3. 用户输入的提示词中如果涉及金额，在提示词后加上；"最终结果除100"

	## 数据库
	- budget: 用于分析财务预算相关问题的,如果有金额相关则提示词后应追加

	## 输出
	{{.formatting}}

	## 输入
	{{.input}}
`
	FilePath = "filePath"
	FileName = "fileName"
	FileId   = "fileId"
	System   = "system"

	_defaultFileQuery = `
	## 历史对话
	{{.history}}

	## 角色
	你是一个文件分析助理，你根据用户的输入信息和历史对话中的文件信息

	## 工作
	1. 你需要理解用户输入的信息并查阅历史对话
	2. 你需要选择一个与问题相关的文件并返回
	3. 你需要描述对该问题处理的时候系统的角色和工作内容

	## 输出
	{{.formatting}}

	## 输入
	{{.input}}
`
)

var (
	_outputParser = outputparser.NewStructured([]outputparser.ResponseSchema{
		{
			Name:        _destinations,
			Description: "name of the prompt to use or \"DEFAULT\"",
		},
		{
			Name:        langchain.Input,
			Description: "对用户输入的原始提示词进行优化",
		},
	})

	_defaultFileOutputParser = outputparser.NewStructured([]outputparser.ResponseSchema{
		{
			Name:        FilePath,
			Description: "与问题符合的一个文件地址",
		}, {
			Name:        FileName,
			Description: "选择的文件命名",
		}, {
			Name:        FileId,
			Description: `选择的文件id; 如果不存在fileId就 输出 "1"`,
		}, {
			Name:        System,
			Description: "系统的角色和工作的描述",
		},
	})
)
