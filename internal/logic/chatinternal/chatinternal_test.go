package chatinternal

import (
	"ai/internal/config"
	"ai/internal/svc"
	"ai/pkg/conf"
	"fmt"
	"path/filepath"
)

var svcTest *svc.ServiceContext

func init() {
	var c config.Config
	conf.MustLoad(filepath.Join("../../../etc/api.yaml"), &c)

	fmt.Println(c)

	svc, err := svc.NewServiceContext(c)
	if err != nil {
		panic(err)
	}

	svcTest = svc
}
