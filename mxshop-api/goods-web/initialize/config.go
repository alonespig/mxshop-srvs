package initialize

import (
	"encoding/json"
	"fmt"
	"mxshop-api/goods-web/global"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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
	configFileName := fmt.Sprintf("goods-web/%s-pro.yaml", configFilePrefix)

	if debug {
		configFileName = fmt.Sprintf("goods-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()

	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}

	zap.L().Info("配置文件加载成功", zap.Any("NacosConfig", global.NacosConfig))

	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(content)

	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		zap.L().Fatal("读取配置文件失败", zap.Error(err))
	}

	fmt.Println(global.ServerConfig)

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
		OnChange: func(namespace string, group string, dataId string, data string) {
			fmt.Println("配置文件发生变化")
			fmt.Println("gropu ", group, "dataId ", dataId, "data ", data)
			err = json.Unmarshal([]byte(data), global.ServerConfig)
			if err != nil {
				panic(err)
			}

			fmt.Println(global.ServerConfig)
		},
	})
	if err != nil {
		zap.L().Fatal("监听配置文件失败", zap.Error(err))
		panic(err)
	}

}
