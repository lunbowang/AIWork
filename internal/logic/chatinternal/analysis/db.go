package analysis

import (
	"ai/internal/domain"
	"ai/internal/svc"
	"ai/pkg/langchain"
	"ai/pkg/langchain/visual"
	"ai/pkg/xerr"
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/outputparser"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools/sqldatabase"
	"github.com/tmc/langchaingo/tools/sqldatabase/mysql"
)

type DB struct {
	svc            *svc.ServiceContext
	c              chains.Chain
	outputParser   outputparser.Structured
	analysisChains chains.Chain
	visual         chains.Chain
	callbacks      callbacks.Handler
}

func NewDB(svc *svc.ServiceContext) *DB {
	return &DB{
		svc:          svc,
		outputParser: _outputParser,
		c: chains.NewLLMChain(svc.LLMs, prompts.PromptTemplate{
			Template:       _defaultDBQuery,
			InputVariables: []string{langchain.Input},
			TemplateFormat: prompts.TemplateFormatGoTemplate,
			PartialVariables: map[string]any{
				"formatting": _outputParser.GetFormatInstructions(),
			},
		}),
		visual: visual.NewVisualChains(svc.LLMs, svc.Config.Upload.SavePath),
	}
}

func (d *DB) Name() string {
	return "db"
}

func (d *DB) Description() string {
	return `
	使用条件：在没有指定文件进行数据分析的时候
	使用限制: 只支持[-财务预算:budget]`
}

func (d *DB) Chains() chains.Chain {
	return chains.NewTransform(d.transform, nil, nil)
}

func (d *DB) transform(ctx context.Context, inputs map[string]any, options ...chains.ChainCallOption) (map[string]any, error) {
	// ？
	fmt.Println("analysis --- db --- start")

	if err := d.initDB(); err != nil {
		return nil, err
	}

	// 选择表 和 提示词 ： 优化
	out, err := chains.Predict(ctx, d.c, inputs, options...)
	if err != nil {
		return nil, err
	}
	fmt.Printf("analysis db chains.Predict out %v\n", out)

	query, err := d.extractQuery(ctx, out)
	if err != nil {
		return nil, err
	}
	fmt.Printf("analysis db d.extractQuery query %v\n", query)

	// 根据SQLDatabaseChains分析用户的问题
	res, err := d.analysisChains.Call(ctx, query, options...)
	if err != nil {
		return nil, err
	}
	fmt.Printf("analysis d.analysisChains.Call res %v\n", res)

	// 进行数据的可视化
	v, err := chains.Call(ctx, d.visual, map[string]any{
		langchain.Input: res["result"],
	})

	var resp domain.ChatResp
	if err != nil {
		fmt.Println("visual : ", err)
		resp = domain.ChatResp{
			ChatType: 0,
			Data:     res["result"].(string),
		}
	} else {
		resp = domain.ChatResp{
			ChatType: domain.ImgAndText,
			Data: domain.ImgAndTextResp{
				Text: res["result"].(string),
				Url:  d.svc.Config.Upload.Host + v[visual.ImgUrl].(string),
			},
		}
	}
	body, err := json.Marshal(&resp)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		langchain.OutPut: string(body),
	}, nil
}

// 解析查询的条件
func (d *DB) extractQuery(ctx context.Context, out string) (map[string]any, error) {
	res, err := d.outputParser.Parse(out)
	if err != nil {
		return nil, err
	}

	// callback

	data := res.(map[string]string)
	db, ok := data[_destinations]
	if !ok {
		return nil, xerr.WithMessage(err, "")
	}

	return map[string]any{
		"query":              data[langchain.Input],
		"table_names_to_use": []string{db},
	}, err
}

func (d *DB) initDB() (err error) {
	if d.analysisChains != nil {
		return nil
	}

	fmt.Println("analysis ---initDB")

	engine, err := mysql.NewMySQL(d.svc.Config.MysqlDns)
	if err != nil {
		return err
	}

	db, err := sqldatabase.NewSQLDatabase(engine, map[string]struct{}{
		"admin": {},
	})
	if err != nil {
		return err
	}

	d.analysisChains = chains.NewSQLDatabaseChain(d.svc.LLMs, 100, db)

	return nil
}
