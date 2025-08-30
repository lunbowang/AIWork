package ebook

import (
	"ai/pkg/langchain"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

type Tool interface {
	Call(ctx context.Context, input, filepath string) (string, error)
}

type Ebook struct {
	plan          chains.Chain
	maxPlanRetry  int
	savePath      string
	subject       Tool
	content       Tool
	MaxIterations int
}

func NewEbook(llm llms.Model, savePath string, opts ...Options) *Ebook {
	opt := newOption(opts...)

	return &Ebook{
		plan:          chains.NewLLMChain(llm, CreateEbookWritePlan()),
		savePath:      savePath,
		subject:       NewSubjectTool(llm),
		content:       NewContentTool(llm),
		MaxIterations: opt.MaxIterations,
		maxPlanRetry:  1,
	}
}

func (e *Ebook) Call(ctx context.Context, inputs map[string]any, options ...chains.ChainCallOption) (map[string]any, error) {
	filepath := fmt.Sprintf("%s/_file_%d.txt", e.savePath, time.Now().Unix())

	// 生成主题、章节列表、摘要
	subjectStr, err := e.subject.Call(ctx, inputs[langchain.Input].(string), filepath)
	if err != nil {
		return nil, err
	}

	// 记录
	inputs[langchain.Input] = subjectStr

	step := Step{
		Subject: subjectStr,
		Record:  "",
	}
	// 生成内容
	var finish map[string]any
	for i := 0; i < e.MaxIterations; i++ {
		// gen content
		step, finish, err = e.doIteration(ctx, inputs, step, filepath)
		if err != nil || finish != nil {
			// 结束
			return finish, err
		}
	}

	return map[string]any{
		langchain.OutPut: filepath,
	}, nil
}

// 执行计划生成内容
func (e *Ebook) doIteration(ctx context.Context, inputs map[string]any, step Step, filepath string) (Step,
	map[string]any, error) {

	// 确定计划
	action, final, err := e.doPlan(ctx, inputs, step, filepath)
	if err != nil || final != nil {
		return step, final, err
	}

	// 内容生成
	fmt.Println("start ----- 执行计划 --------- : ", action.Input)

	// 执行
	_, err = e.content.Call(ctx, action.Input+"\n## 文章信息"+step.Subject, filepath)
	if err != nil {
		return step, nil, err
	}

	fmt.Println("end ----- 执行计划 --------- : ", action.Input)
	step.Record += "\n" + action.Input + " : 撰写完成"

	return step, nil, nil
}

// 判断下一步的计划 或者 完成的验证
func (e *Ebook) doPlan(ctx context.Context, inputs map[string]any, step Step, filepath string) (*Action, map[string]any, error) {
	// 设置撰写记录
	inputs["history"] = step.Record

	do := func() (*Action, map[string]any, error) {
		output, err := chains.Predict(ctx, e.plan, inputs)
		if err != nil {
			return nil, nil, err
		}

		// 验证，是否完成
		if strings.Contains(output, _finalAnswerAction) {
			// 完成了
			return nil, map[string]any{
				langchain.Input: filepath,
			}, err
		}

		// 解析
		r := regexp.MustCompile(`计划：\s*(.+)\s*`)
		matches := r.FindStringSubmatch(output)
		if len(matches) == 0 {
			return nil, nil, fmt.Errorf("不存在信息 ：%s", output)
		}

		return &Action{
			Input: strings.TrimSpace(matches[1]),
		}, nil, err
	}

	var rerr error
	for i := 0; i < e.maxPlanRetry; i++ {
		res, final, err := do()
		if err == nil {
			return res, final, nil
		}
		rerr = err
	}

	return nil, nil, rerr
}

func (e *Ebook) GetMemory() schema.Memory {
	return memory.NewSimple()
}

func (e *Ebook) GetInputKeys() []string {
	return []string{}
}

func (e *Ebook) GetOutputKeys() []string {
	return []string{}
}

func SaveContent(content, filepath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// 如果文件不存在，先创建文件
		return os.WriteFile(filepath, []byte(content), 0644)
	}

	// 如果文件已存在，在文件内容尾部增加内容
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
