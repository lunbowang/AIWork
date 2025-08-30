package visual

import (
	"ai/pkg/langchain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// 提取具体的数值 -》绘图

type BarTool struct {
	c        chains.Chain
	savePath string
}

func NewBarTool(llm llms.Model, savePath string) *BarTool {
	return &BarTool{
		c:        chains.NewLLMChain(llm, prompts.NewPromptTemplate(_defaultVBarPrompts, []string{"input"})),
		savePath: savePath,
	}
}

func (b *BarTool) Namespace() string {
	//TODO implement me
	panic("implement me")
}

func (b *BarTool) Destinations() string {
	//TODO implement me
	panic("implement me")
}

func (b *BarTool) Call(ctx context.Context, inputs map[string]string, options ...chains.ChainCallOption) (map[string]any, error) {
	// 获取数据的数值
	out, err := chains.Predict(ctx, b.c, map[string]any{
		langchain.Input: inputs[langchain.Input],
	})
	if err != nil {
		return nil, err
	}

	// 解析数据
	data, err := b.parseData(out)
	if err != nil {
		return nil, err
	}

	// 绘图
	path, err := b.plot(inputs, data)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		ImgUrl: path,
	}, nil
}

func (b *BarTool) plot(inputs map[string]string, data map[string]float64) (string, error) {
	// 创建一个新的图
	p := plot.New()

	// 设置标题和标签
	p.Title.Text = inputs[Title]
	p.X.Label.Text = inputs[X]
	p.Y.Label.Text = inputs[Y]

	var (
		values []float64
		labels []string
	)
	for k, v := range data {
		values = append(values, v)
		labels = append(labels, k)
	}
	// 创建柱状图
	bars, err := plotter.NewBarChart(plotter.Values(values), 10)
	if err != nil {
		return "", err
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(0)
	// 将柱状图添加到图表中
	p.Add(bars)
	// 设置X轴标签
	p.NominalX(labels...)

	fmt.Println("inputs : ", inputs, " values : ", values, " labels : ", labels)

	// 保存
	path := fmt.Sprintf("%s%v.png", b.savePath, time.Now().Unix())
	if err := p.Save(5*vg.Inch, 5*vg.Inch, path); err != nil {
		return "", err
	}

	return path, err

}

// {"语文": 900, ",数学": 700, "物理": 800}
func (b *BarTool) parseData(text string) (map[string]float64, error) {
	var parsed map[string]any

	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		// Remove the ```json that should be at the start of the text, and the ```
		// that should be at the end of the text.
		withoutJSONStart := strings.Split(text, "```json")
		if !(len(withoutJSONStart) > 1) {
			return nil, errors.New("no ```json at start of output")
		}

		withoutJSONEnd := strings.Split(withoutJSONStart[1], "```")
		if len(withoutJSONEnd) < 1 {
			return nil, errors.New("no ``` at end of output")
		}

		jsonString := withoutJSONEnd[0]

		err := json.Unmarshal([]byte(jsonString), &parsed)
		if err != nil {
			return nil, err
		}
	}
	res := make(map[string]float64)

	for k, v := range parsed {
		switch v.(type) {
		case float64:
			res[k] = v.(float64)
		case float32:
			res[k] = float64(v.(float32))
		default:
			// ...
		}
	}

	return res, nil
}
