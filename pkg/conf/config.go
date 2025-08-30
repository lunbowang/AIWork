/**
 * @Author: dn-jinmin
 * @File:  gogen
 * @Version: 1.0.0
 * @Date: 2024/4/16
 * @Description:
 */

package conf

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type Loadhandler func(string, any) error

var (
	loaders = map[string]Loadhandler{
		".yaml": LoadFromYamlBytes,
	}
)

func MustLoad(file string, v any) {
	Load(file, v)
}

// Load 从指定文件加载配置到v中，支持.json、.yaml和.yml格式
// 参数：file为配置文件路径，v为接收配置的结构体指针
// 返回值：加载过程中可能出现的错误
func Load(file string, v any) error {
	// 根据扩展名获取对应的加载器
	loader, ok := loaders[strings.ToLower(path.Ext(file))]
	if !ok {
		return fmt.Errorf("unrecognized file type: %s", file)
	}

	return loader(file, v)
}

// LoadFromYamlBytes 从YAML文件加载配置到v中
// 参数：file为YAML文件路径，v为接收配置的结构体指针
// 返回值：加载过程中可能出现的错误
func LoadFromYamlBytes(file string, v any) error {
	viper.SetConfigType("yaml")
	// 将文件路径中的反斜杠替换为斜杠，处理跨平台路径问题
	file = strings.Replace(file, "\\", "/", -1)
	// 提取配置文件所在的目录，并添加到viper的配置路径中
	viper.AddConfigPath(file[:strings.LastIndex(file, "/")+1])
	// 提取配置文件的文件名（不含路径和扩展名）
	filename := file[strings.LastIndex(file, "/")+1 : strings.LastIndex(file, ".")]
	// 设置viper要读取的配置文件名
	viper.SetConfigName(filename)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	// 将配置文件内容解析到目标结构体v中
	if err := viper.Unmarshal(v); err != nil {
		return err
	}

	return nil
}
