package router

import (
	"context"
	"errors"

	"github.com/tmc/langchaingo/llms"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/outputparser"
	"github.com/tmc/langchaingo/schema"
)

const Empty = "Default"

var ErrNotHandles = errors.New("不存在合适的handler")

// Router 路由处理器，用于根据输入选择合适的处理器(Handler)执行处理逻辑
type Router struct {
	handlers     map[string]Handler      // 存储处理器映射，键为处理器名称，值为处理器实例
	chain        chains.Chain            // LLM链，用于执行提示词生成和模型调用
	callbacks    callbacks.Handler       // 回调处理器，用于处理路由过程中的事件
	memory       schema.Memory           // 内存组件，用于存储对话状态或历史信息
	outputparser outputparser.Structured // 输出解析器，用于解析LLM的输出结果
	emptyHandle  Handler                 // 空处理器，当没有找到匹配的处理器时使用
}

// NewRouter 创建一个新的Router实例
func NewRouter(llm llms.Model, handler []Handler, opts ...Option) *Router {
	// 初始化默认配置
	opt := executorDefaultOptions(handler)
	// 应用自定义配置选项
	for _, o := range opts {
		o(&opt)
	}

	// 将处理器列表转换为映射，便于快速查找
	hs := make(map[string]Handler, len(handler))
	for _, h := range handler {
		hs[h.Name()] = h
	}

	// 创建并返回Router实例
	return &Router{
		handlers:     hs,
		chain:        chains.NewLLMChain(llm, opt.prompt), // 使用LLM和提示词模板创建链
		callbacks:    opt.callback,
		memory:       opt.memory,
		emptyHandle:  opt.emptyHandler,
		outputparser: _outputparser, // 使用全局输出解析器
	}
}

// Call 执行路由处理逻辑，根据输入选择合适的处理器并调用
func (r *Router) Call(ctx context.Context, inputs map[string]any, options ...chains.ChainCallOption) (map[string]any, error) {
	// 触发链开始事件的回调
	if r.callbacks != nil {
		r.callbacks.HandleChainStart(ctx, inputs)
	}

	// 当没有注册任何处理器时的处理逻辑
	if len(r.handlers) == 0 {
		if r.emptyHandle != nil {
			return chains.Call(ctx, r.emptyHandle.Chains(), inputs)
		} else {
			return nil, ErrNotHandles
		}
	}

	// 调用LLM链获取路由决策结果
	result, err := chains.Call(ctx, r.chain, inputs, options...)
	if err != nil {
		return nil, err
	}

	// 提取LLM输出的文本结果
	text, ok := result["text"]
	if !ok {
		return nil, chains.ErrNotFound // 文本结果不存在时返回错误
	}

	// 解析LLM输出为结构化数据
	out, err := r.outputparser.Parse(text.(string))
	if err != nil {
		return nil, err
	}

	// 触发链结束事件的回调
	if r.callbacks != nil {
		r.callbacks.HandleChainEnd(ctx, map[string]any{
			"out": out,
		})
	}

	// 将解析结果转换为字符串映射
	data := out.(map[string]string)
	// 获取目标处理器名称
	next, ok := data[_destinations]
	// 检查目标处理器是否有效，无效则使用空处理器或返回错误
	if !ok || next == Empty || r.handlers[next] == nil {
		if r.emptyHandle != nil {
			return chains.Call(ctx, r.emptyHandle.Chains(), inputs)
		} else {
			return nil, ErrNotHandles
		}
	}

	// 调用目标处理器处理输入
	return chains.Call(ctx, r.handlers[next].Chains(), inputs)
}

// GetMemory 返回路由使用的内存组件
func (r *Router) GetMemory() schema.Memory {
	//TODO implement me
	return r.memory
}

// GetInputKeys 返回输入键列表
func (r *Router) GetInputKeys() []string {
	return nil
}

// GetOutputKeys 返回输出键列表
func (r *Router) GetOutputKeys() []string {
	return nil
}
