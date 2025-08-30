package ebook

import (
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/schema"
)

type Options func(opt *Option)

const DefaultMaxIterations = 10

type Option struct {
	CallbacksHandler callbacks.Handler
	memory           schema.Memory
	MaxIterations    int
}

func newOption(opts ...Options) *Option {
	o := Option{
		MaxIterations: DefaultMaxIterations,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &o
}

func WithMemory(memory schema.Memory) Options {
	return func(options *Option) {
		options.memory = memory
	}
}

func WithCallbacksHandler(callbacks callbacks.Handler) Options {
	return func(options *Option) {
		options.CallbacksHandler = callbacks
	}
}

func WithMaxIterations(maxIterations int) Options {
	return func(opt *Option) {
		opt.MaxIterations = maxIterations
	}
}
