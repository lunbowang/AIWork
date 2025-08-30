package toolx

import (
	"ai/internal/domain"
	"ai/internal/svc"
	"ai/pkg/curl"
	"ai/pkg/langchain/outputparserx"
	"ai/token"
	"context"
	"encoding/json"

	"github.com/tmc/langchaingo/callbacks"
)

type TodoAdd struct {
	svc          *svc.ServiceContext
	callback     callbacks.Handler
	outputparser outputparserx.Structured
}

func NewTodoAdd(svc *svc.ServiceContext) *TodoAdd {
	return &TodoAdd{
		svc:      svc,
		callback: svc.Callbacks,
		outputparser: outputparserx.NewStructured([]outputparserx.ResponseSchema{
			{
				Name:        "title",
				Description: "todo title",
			}, {
				Name:        "deadlineAt",
				Description: "calculate the final deadline based on the time information entered by the user and combined with today's time. a Unix time",
				Type:        "int64",
			}, {
				Name:        "desc",
				Description: "todo description",
			}, {
				Name:        "executeIds",
				Description: "list of participating users in the backlog. the data type is a set of string ids. none is empty",
				Type:        "[]string",
			},
		}),
	}
}

func (t *TodoAdd) Name() string {
	return "todo_add"
}

func (t *TodoAdd) Description() string {
	template := `
	a todo add interface.
	use when you need to create a todo.
	keep Chinese output.
` + t.outputparser.GetFormatInstructions()

	return template
}

func (t *TodoAdd) Call(ctx context.Context, input string) (string, error) {
	if t.callback != nil {
		t.callback.HandleText(ctx, "todo add start : "+input)
	}

	data, err := t.outputparser.Parse(input)
	if err != nil {
		return "", err
	}

	res, err := curl.PostRequest(token.GetTokenStr(ctx), t.svc.Config.Host+"/v1/todo", data)

	var idResp domain.IdRespInfo
	if err := json.Unmarshal(res, &idResp); err != nil {
		return "", err
	}

	return Success + "\ncreated todo id : " + idResp.Data.Id, nil
}
