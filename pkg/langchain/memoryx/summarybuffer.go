package memoryx

import (
	"context"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
)

// outParser 定义输出解析函数类型，用于对AI输出内容进行自定义处理
type outParser func(ctx context.Context, input string) string

// SummaryBuffer 对话摘要缓存器，结合普通对话缓存与增量摘要功能
type SummaryBuffer struct {
	outParser
	*memory.ConversationBuffer // 正常聊天的消息
	chain                      chains.Chain

	MaxTokenLimit int
	callback      callbacks.Handler

	buffer llms.ChatMessage //总结后的消息内容
}

func NewSummaryBuffer(llms llms.Model, maxTokenLimit int, opts ...Option) *SummaryBuffer {
	opt := newOption(opts...)

	return &SummaryBuffer{
		ConversationBuffer: memory.NewConversationBuffer(),
		chain:              chains.NewLLMChain(llms, createSummaryPrompt(), chains.WithCallback(opt.callback)),
		callback:           opt.callback,
		MaxTokenLimit:      maxTokenLimit,
		buffer:             nil,
		outParser:          opt.outParser,
	}
}

// GetMemoryKey getter for memory key.
func (s *SummaryBuffer) GetMemoryKey(ctx context.Context) string {
	return s.ConversationBuffer.GetMemoryKey(ctx)
}

// MemoryVariables Input keys this memory class will load dynamically.
func (s *SummaryBuffer) MemoryVariables(ctx context.Context) []string {
	return s.ConversationBuffer.MemoryVariables(ctx)
}

// LoadMemoryVariables 加载内存变量，返回当前完整的对话上下文（摘要+未总结消息）
func (s *SummaryBuffer) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	var (
		res []llms.ChatMessage // 存储完整对话上下文（摘要+普通消息）
		err error
	)

	// 若存在已生成的摘要，先将摘要（系统消息）加入上下文
	if s.buffer != nil {
		res = append(res, s.buffer)
	}

	// 读取当前未总结的普通对话消息（用户消息+AI消息）
	message, err := s.ChatHistory.Messages(ctx)
	if err != nil {
		return nil, err
	}

	// 将普通对话消息追加到上下文后
	res = append(res, message...)

	// 将对话消息列表转换为字符串（按"人类前缀+消息内容"、"AI前缀+消息内容"格式拼接）
	// s.HumanPrefix：人类消息前缀（如"Human:"），s.AIPrefix：AI消息前缀（如"Ai:"）
	bufferString, err := llms.GetBufferString(res, s.HumanPrefix, s.AIPrefix)

	// 返回完整的上下文字符串，键为内存键
	return map[string]any{
		s.MemoryKey: bufferString,
	}, nil
}

// SaveContext 保存对话上下文，核心逻辑：
// 1. 先存储最新的用户输入和AI输出；
// 2. 检查当前对话令牌数是否超过上限，超过则生成摘要并清空普通消息。
func (s *SummaryBuffer) SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error {
	// 1. 提取并存储用户最新输入消息
	// 从输入数据中获取用户输入（s.InputKey为用户输入的键名，如"input"）
	userInputValue, err := memory.GetInputValue(inputs, s.InputKey)
	if err != nil {
		return err
	}
	// 将用户输入添加到对话历史（以人类消息形式存储）
	err = s.ChatHistory.AddUserMessage(ctx, userInputValue)
	if err != nil {
		return err
	}

	// 2. 提取、处理并存储AI最新输出消息
	// 从输出数据中获取AI输出（s.OutputKey为AI输出的键名，如"output"）
	aiOutPutValue, err := memory.GetInputValue(outputs, s.OutputKey)
	if err != nil {
		return err
	}

	// 若配置了输出解析器，先对AI输出进行自定义处理（如格式清洗、内容过滤）
	if s.outParser != nil {
		aiOutPutValue = s.outParser(ctx, aiOutPutValue)
	}

	err = s.ChatHistory.AddAIMessage(ctx, aiOutPutValue)
	if err != nil {
		return err
	}

	// 3. 读取当前所有普通对话消息，计算令牌数
	// 获取当前未总结的普通消息列表
	messages, err := s.ChatHistory.Messages(ctx)
	if err != nil {
		return err
	}
	// 将消息列表转换为字符串，用于计算令牌数
	bufferString, err := llms.GetBufferString(messages, s.ConversationBuffer.HumanPrefix,
		s.ConversationBuffer.AIPrefix)
	if err != nil {
		return err
	}

	// 4. 检查令牌数是否超过上限，未超过则直接返回（无需生成摘要）
	// llms.CountTokens：计算字符串对应的令牌数（第一个参数为模型名，空字符串表示使用默认计算逻辑）
	if llms.CountTokens("", bufferString) <= s.MaxTokenLimit {
		//没有超过上限
		return nil
	}

	// 5. 令牌数超过上限，生成新的对话摘要
	// 若已有历史摘要，将其内容作为"已有总结"传入；若无则为空字符串
	var newLines string
	if s.buffer != nil {
		newLines = s.buffer.GetContent()
	}

	// 调用摘要LLM链，生成合并后的新摘要
	// 输入参数："summary"为当前普通消息的字符串，"new_lines"为历史摘要（如有）
	newSummary, err := chains.Predict(ctx, s.chain, map[string]any{
		"summary":   bufferString,
		"new_lines": newLines,
	})

	if err != nil {
		return err
	}

	// 6. 保存新摘要并清空普通对话消息
	// 将新摘要以系统消息形式存储（SystemChatMessage），便于后续加载时区分
	s.buffer = &llms.SystemChatMessage{
		Content: newSummary,
	}

	// 清空普通对话消息（已总结的内容无需再保留原始消息）
	return s.ChatHistory.SetMessages(ctx, nil)
}

// Clear memory contents.
func (s *SummaryBuffer) Clear(ctx context.Context) error {
	s.buffer = nil
	return s.ConversationBuffer.Clear(ctx)
}
