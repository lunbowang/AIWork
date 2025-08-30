package memoryx

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"

	"github.com/tmc/langchaingo/chains"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/memory"
)

type Summary struct {
	*memory.ConversationBuffer
	callback callbacks.Handler

	chain chains.Chain
}

func NewSummary(llm llms.Model, opts ...Option) *Summary {
	opt := newOption(opts...)
	return &Summary{
		callback:           opt.callback,
		ConversationBuffer: memory.NewConversationBuffer(),
		chain:              chains.NewLLMChain(llm, createSummaryPrompt(), chains.WithCallback(opt.callback)),
	}
}

// GetMemoryKey getter for memory key.
func (s *Summary) GetMemoryKey(ctx context.Context) string {
	return s.ConversationBuffer.GetMemoryKey(ctx)
}

// MemoryVariables Input keys this memory class will load dynamically.
func (s *Summary) MemoryVariables(ctx context.Context) []string {
	return s.ConversationBuffer.MemoryVariables(ctx)
}

// LoadMemoryVariables Return key-value pairs given the text input to the chain.
// If None, return all memories
func (s *Summary) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	return s.ConversationBuffer.LoadMemoryVariables(ctx, inputs)
}

// SaveContext Save the context of this model run to memory.
func (s *Summary) SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error {
	// 加载已有的对话摘要（从内存中获取历史总结信息）
	// message 是一个键值对映射，包含内存中存储的对话摘要
	message, err := s.LoadMemoryVariables(ctx, inputs)
	if err != nil {
		return err
	}
	// 从内存数据中提取当前的对话摘要（s.MemoryKey 是存储摘要的键名）
	summary := message[s.MemoryKey]

	// 从输入数据中提取用户的最新提问内容
	// s.InputKey 是输入数据中用户提问的键名（如"input"）
	userInputValue, err := memory.GetInputValue(inputs, s.InputKey)
	if err != nil {
		return err
	}

	// 从输出数据中提取AI的最新回答内容
	// s.OutputKey 是输出数据中AI回答的键名（如"output"）
	userOutPutValue, err := memory.GetInputValue(outputs, s.OutputKey)
	if err != nil {
		return err
	}

	// 格式化新的对话内容（用户提问 + AI回答）
	// 采用固定格式便于后续LLM理解和总结
	newLines := fmt.Sprintf("Homan:%s \nAi: %s", userInputValue, userOutPutValue)

	// 调用LLM链生成更新后的对话摘要
	// 将现有摘要和新对话内容传入，让LLM生成合并后的新摘要
	newSummary, err := chains.Predict(ctx, s.chain, map[string]any{
		"summary":   summary,
		"new_lines": newLines,
	})

	if err != nil {
		return err
	}

	// 将更新后的对话摘要保存到聊天历史中
	// 以系统消息的形式存储（llms.SystemChatMessage），便于后续加载
	return s.ChatHistory.SetMessages(ctx, []llms.ChatMessage{
		llms.SystemChatMessage{Content: newSummary},
	})
}

// Clear memory contents.
func (s *Summary) Clear(ctx context.Context) error {
	return s.ConversationBuffer.Clear(ctx)
}
