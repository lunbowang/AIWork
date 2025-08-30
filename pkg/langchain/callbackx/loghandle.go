package callbackx

import (
	"context"
	"encoding/json"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"

	"gitee.com/dn-jinmin/tlog"
)

type LogHandle struct {
	tlog.Logger
}

func NewLogHandler(logger tlog.Logger) *LogHandle {
	return &LogHandle{Logger: logger}
}

func (l *LogHandle) HandleText(ctx context.Context, text string) {
	l.InfoCtx(ctx, "text", text)
}

func (l *LogHandle) HandleLLMStart(ctx context.Context, prompts []string) {
	l.InfoCtx(ctx, "llm_start", prompts)
}

func (l *LogHandle) HandleLLMGenerateContentStart(ctx context.Context, ms []llms.MessageContent) {
	l.InfoCtx(ctx, "llm_generate_content_start", l.mustJsonMarshal(ms))
}

func (l *LogHandle) HandleLLMGenerateContentEnd(ctx context.Context, res *llms.ContentResponse) {
	l.InfoCtx(ctx, "llm_generate_content_end", l.mustJsonMarshal(res))
}

func (l *LogHandle) HandleLLMError(ctx context.Context, err error) {
	l.ErrorCtx(ctx, "llm_error", err.Error())
}

func (l *LogHandle) HandleChainStart(ctx context.Context, inputs map[string]any) {
	l.InfoCtx(ctx, "chain_start", l.mustJsonMarshal(inputs))
}

func (l *LogHandle) HandleChainEnd(ctx context.Context, outputs map[string]any) {
	l.InfoCtx(ctx, "chain_end", l.mustJsonMarshal(outputs))
}

func (l *LogHandle) HandleChainError(ctx context.Context, err error) {
	l.ErrorCtx(ctx, "chain_error", err.Error())
}

func (l *LogHandle) HandleToolStart(ctx context.Context, input string) {
	l.InfoCtx(ctx, "tool_start", input)
}

func (l *LogHandle) HandleToolEnd(ctx context.Context, output string) {
	l.InfoCtx(ctx, "tool_end", output)
}

func (l *LogHandle) HandleToolError(ctx context.Context, err error) {
	l.ErrorCtx(ctx, "tool_error", err.Error())
}

func (l *LogHandle) HandleAgentAction(ctx context.Context, action schema.AgentAction) {
	l.InfoCtx(ctx, "agent_action", l.mustJsonMarshal(action))
}

func (l *LogHandle) HandleAgentFinish(ctx context.Context, finish schema.AgentFinish) {
	l.InfoCtx(ctx, "agent_finish", l.mustJsonMarshal(finish))
}

func (l *LogHandle) HandleRetrieverStart(ctx context.Context, query string) {
	l.InfoCtx(ctx, "retriever_start", query)
}

func (l *LogHandle) HandleRetrieverEnd(ctx context.Context, query string, documents []schema.Document) {
	l.InfofCtx(ctx, "retriever_end", "query %s, documents %s", query, l.mustJsonMarshal(documents))
}

func (l *LogHandle) HandleStreamingFunc(ctx context.Context, chunk []byte) {
	//l.InfoCtx(ctx, "streaming_func", string(chunk))
}

func (l *LogHandle) mustJsonMarshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
