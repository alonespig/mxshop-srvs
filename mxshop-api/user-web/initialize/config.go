package initialize

import (
	"fmt"
	"mxshop-api/user-web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)

	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()

	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

	zap.L().Info("配置文件加载成功", zap.Any("serverConfig", global.ServerConfig))
	zap.L().Info("配置文件加载成功", zap.String("name", global.ServerConfig.Name))

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.L().Info("配置文件修改了...", zap.String("e.Name", e.Name))
		v.ReadInConfig()
		v.Unmarshal(global.ServerConfig)
		zap.L().Info("配置文件重新加载...", zap.Any("serverConfig", global.ServerConfig))
	})
}
