goctl-gin model mongo --type user --dir ./internal/model/
goctl-gin model mongo --type department --dir ./internal/model/
goctl-gin model mongo --type departmentUser --dir ./internal/model/
goctl-gin model mongo --type todo --dir ./internal/model/
goctl-gin model mongo --type Approval  --dir ./internal/model/
goctl-gin model mongo --type UserTodo   --dir ./internal/model/
goctl-gin model mongo --type usertodomodel  --dir ./internal/model/
goctl-gin model mongo --type chatlog  --dir ./internal/model/

goctl-gin api go -api ./doc/api.api -dir ./  go-zero