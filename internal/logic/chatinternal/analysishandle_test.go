/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package chatinternal

import (
	"ai/internal/domain"
	"ai/pkg/langchain"
	"ai/pkg/langchain/memoryx"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/schema"
)

func TestNewAnalysisHandle_DB(t *testing.T) {
	t.Log(svcTest.Config.MysqlDns)
	memory := memoryx.NewMemoryx(func() schema.Memory {
		m := memoryx.NewSummaryBuffer(svcTest.LLMs, 50, memoryx.WithCallback(svcTest.Callbacks))

		m.InputKey = langchain.Input
		return m
	})
	a := NewAnalysisHandle(svcTest, memory)
	ctx := context.Background()

	res, err := chains.Call(ctx, a.Chains(), map[string]any{
		"input": "分析财务预算中待发生、进行中、已完成的各项金额",
	})

	t.Log(res, " ", err)
}

func TestNewAnalysisHandle_File(t *testing.T) {
	t.Log(svcTest.Config.MysqlDns)

	memory := memoryx.NewMemoryx(func() schema.Memory {
		m := memoryx.NewSummaryBuffer(svcTest.LLMs, 50, memoryx.WithCallback(svcTest.Callbacks))

		m.InputKey = langchain.Input
		return m
	})

	ctx := context.WithValue(context.Background(), langchain.ChatId, "999")
	testUploadFile(ctx, memory, t)

	// 输出
	t.Log(memory.LoadMemoryVariables(ctx, map[string]any{}))

	a := NewAnalysisHandle(svcTest, memory)
	res, err := chains.Call(ctx, a.Chains(), map[string]any{
		"input": "请根据上传的财务预算报表分析财务预算中待发生、进行中、已完成的各项金额",
	})

	t.Log(res, " ", err)
}

func testUploadFile(ctx context.Context, memory schema.Memory, t *testing.T) {
	b, err := json.Marshal([]*domain.ChatFile{
		{
			Path: svcTest.Config.Upload.SavePath + "1.xlsx",
			Name: "财务预算报表",
			Time: time.Now(),
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	err = memory.SaveContext(ctx, map[string]any{
		langchain.Input: string(b),
	}, map[string]any{
		langchain.Input: "upload file",
	})

	if err != nil {
		t.Fatal(err)
	}
}
