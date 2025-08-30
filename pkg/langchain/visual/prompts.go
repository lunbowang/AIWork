package visual

import (
	"ai/pkg/langchain"

	"github.com/tmc/langchaingo/outputparser"
)

// 获取分析的数据字段、数据值 选择的工具
const (
	ImgUrl = "imgUrl"
	X      = "x-title"
	Y      = "y-title"
	Title  = "title"

	_destinations = "destinations"
	_formatting   = "formatting"

	_defaultVisualParsePrompts = `
	# 角色
	你是一个数据分析信息提取助手，你负责分析用户输入的信息，然后根据要求提取相关的信息并格式化输出
	
	## 工作
	1. 你需要先理解用户输入的信息内容
	2. 你将确定用户的分析主要信息是什么
	3. 你要分析用户输入的信息适合那种输出情况并从 [bar] 选择一个
	4. 在用户提供的信息中确定数据分析坐标图中 x、y 的含义是什么
	
	## 禁止
	1. 你不需要对用户输入的信息做调整
	2. 你不需要对用户输入的信息进行数据分析
	
	## 类型
	- bar: 适合在x坐标轴上只有单一类型数据值的时候使用，主要绘制柱状图
	
	## 输出
	{{.formatting}}
	
	## 输入
	{{.input}}
`
	_defaultVBarPrompts = `
	# 角色
	你是一个数据分析的助手，你的工作是根据用户输入的信息提取出用户需要分析的信息并以要求的格式输出

	## 工作
	1. 理解用户输入的信息
	2. 根据数据坐标x轴标题和y轴标题提取信息

	## 禁止
	1. 你不需要对用户输入的信息做调整
	2. 你不需要对用户输入的信息进行数据分析

	## 输出
	1. 以标题为k，以数值为y
	2. 标题应简洁，多个标题不能存在相同的内容
	3. 请将标题转为英文
	4. 基于json的方式输出如：{"语文": 100, ",数学": 90, "物理": 90}
	
	## 示例
	用户输入：在这次考试里木兮的语文有80分、数学有90分、物理有90分
	最终结果：{"语文": 900, ",数学": 700, "物理": 800}

	## 输入
	{{.input}}
`
)

var (
	_outputParserVisual = outputparser.NewStructured([]outputparser.ResponseSchema{
		{
			Name:        _destinations,
			Description: "name of the prompt to use or \"DEFAULT\"",
		}, {
			Name:        X,
			Description: "数据坐标x轴的标题",
		}, {
			Name:        Y,
			Description: "数据坐标y轴的标题",
		}, {
			Name:        Title,
			Description: "数据分析的主标题",
		}, {
			Name:        langchain.Input,
			Description: "用户原始输入",
		},
	})
)
