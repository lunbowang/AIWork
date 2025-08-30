package config

import "gitee.com/dn-jinmin/tlog"

type Config struct {
	Name string
	Addr string
	Host string
	Jwt  struct {
		Secret string
		Expire int64
	}
	MysqlDns string
	Mongo    struct {
		User     string   //用户名
		Password string   //密码
		Hosts    []string //host地址
		Port     int      //端口
		Database string   //数据库名
		Params   string   //其他相关认证
		//MaxPoolSize uint64   //连接池最大数
	}

	Tlog struct {
		Mode  tlog.LogMod
		Label string
	}

	Ws struct {
		Addr string
	}

	Langchain struct {
		Url    string
		ApiKey string
	}
	AliGPT struct {
		Url    string
		ApiKey string
	}
	Upload struct {
		SavePath string
		Host     string
	}

	Redis struct {
		Addr string
	}
}
