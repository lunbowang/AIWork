package memoryx

import "github.com/tmc/langchaingo/callbacks"

// Options 定义了组件的配置选项结构体
// 目前主要用于配置回调处理器，可根据需求扩展其他配置项
type Options struct {
	outParser
	callback callbacks.Handler // 回调处理器，用于处理组件生命周期中的各类事件
}

// Option 函数类型定义了配置选项的设置方式
// 采用函数式选项模式，允许灵活地配置Options结构体
type Option func(options *Options)

func newOption(opts ...Option) *Options {
	opt := &Options{
		callback: nil,
	}

	// 应用所有传入的配置选项，覆盖默认值
	for _, o := range opts {
		o(opt)
	}
	return opt
}

// WithCallback 是一个配置选项函数，用于设置回调处理器
func WithCallback(handler callbacks.Handler) Option {
	return func(options *Options) {
		options.callback = handler
	}
}

func WithOutParser(outParser outParser) Option {
	return func(options *Options) {
		options.outParser = outParser
	}
}
