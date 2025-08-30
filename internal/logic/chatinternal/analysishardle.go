package chatinternal

import (
	"ai/internal/logic/chatinternal/analysis"
	"ai/internal/svc"
	"ai/pkg/langchain/router"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/schema"
)

type AnalysisHandle struct {
	svc *svc.ServiceContext
	c   chains.Chain
}

func NewAnalysisHandle(svc *svc.ServiceContext, memory schema.Memory) *AnalysisHandle {
	return &AnalysisHandle{
		svc: svc,

		c: router.NewRouter(svc.LLMs, []router.Handler{
			analysis.NewDB(svc),
			analysis.NewFile(svc),
		}, router.WithMemory(memory)),
	}
}

func (a *AnalysisHandle) Name() string {
	return "analysis"
}

func (a *AnalysisHandle) Description() string {
	return "use when you need to analyze certain data"
}

func (a *AnalysisHandle) Chains() chains.Chain {
	return a.c
}
