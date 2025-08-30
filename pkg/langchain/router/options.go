package router

import (
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
)

// Options 结构体定义了执行器的配置选项集合
type Options struct {
	prompt       prompts.PromptTemplate // 用于生成提示词的模板
	memory       schema.Memory          // 用于存储对话历史或状态的内存组件
	callback     callbacks.Handler      // 用于处理执行过程中回调事件的处理器
	emptyHandler Handler                // 当没有匹配的处理器时使用的默认处理器
}

// Option 函数类型定义了配置选项的设置方式
// 采用函数式选项模式，允许灵活地配置Options结构体
type Option func(options *Options)

// executorDefaultOptions 生成执行器的默认配置选项
func executorDefaultOptions(handler []Handler) Options {
	return Options{
		prompt: createPrompt(handler),
		memory: memory.NewSimple(),
	}
}

// WithMemory 用是一个选项配置函数，于自定义内存组件
func WithMemory(m schema.Memory) Option {
	return func(options *Options) {
		options.memory = m
	}
}

// WithEmptyHandler 是一个选项配置函数，用于设置空处理器
func WithEmptyHandler(emptyHandler Handler) Option {
	return func(options *Options) {
		options.emptyHandler = emptyHandler
	}
}

// Withcallback 是一个选项配置函数，用于设置回调处理器
func Withcallback(callback callbacks.Handler) Option {
	return func(options *Options) {
		options.callback = callback
	}
}
