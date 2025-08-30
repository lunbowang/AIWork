package memoryx

import (
	"ai/pkg/langchain"
	"context"
	"sync"

	"github.com/tmc/langchaingo/schema"
)

type Memoryx struct {
	sync.Mutex
	getMemory     func() schema.Memory
	memorys       map[string]schema.Memory
	defaultMemory schema.Memory
}

func NewMemoryx(handle func() schema.Memory) *Memoryx {
	return &Memoryx{
		getMemory:     handle,
		memorys:       make(map[string]schema.Memory),
		defaultMemory: handle(),
	}
}

// GetMemoryKey getter for memory key.
func (s *Memoryx) GetMemoryKey(ctx context.Context) string {
	return s.memory(ctx).GetMemoryKey(ctx)
}

// MemoryVariables Input keys this memory class will load dynamically.
func (s *Memoryx) MemoryVariables(ctx context.Context) []string {
	return s.memory(ctx).MemoryVariables(ctx)
}

// LoadMemoryVariables 加载内存变量，返回当前完整的对话上下文（摘要+未总结消息）
func (s *Memoryx) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	return s.memory(ctx).LoadMemoryVariables(ctx, inputs)
}

// SaveContext 保存对话上下文，核心逻辑：
// 1. 先存储最新的用户输入和AI输出；
// 2. 检查当前对话令牌数是否超过上限，超过则生成摘要并清空普通消息。
func (s *Memoryx) SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error {
	return s.memory(ctx).SaveContext(ctx, inputs, outputs)
}

// Clear memory contents.
func (s *Memoryx) Clear(ctx context.Context) error {
	return s.memory(ctx).Clear(ctx)
}

func (s *Memoryx) memory(ctx context.Context) schema.Memory {
	s.Lock()
	defer s.Unlock()

	var chatId string
	v := ctx.Value(langchain.ChatId)
	if v == nil {
		return s.defaultMemory
	}

	chatId = v.(string)
	memory, ok := s.memorys[chatId]
	if !ok {
		memory = s.getMemory()
		s.memorys[chatId] = memory
	}

	return memory
}
