package initialize

import (
	"encoding/json"
	"fmt"
	"mxshop-api/user-web/global"

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
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)

	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()

	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		zap.S().Fatalf("[config] 读取配置文件 %s 失败 %s", configFileName, err.Error())
	}
	zap.S().Infof("[config] 读取配置文件 %s 成功", configFileName)

	if err := v.Unmarshal(global.NacosConfig); err != nil {
		zap.S().Fatalf("[config] 解析配置文件 %s 失败 %s", configFileName, err.Error())
	}

	zap.S().Infof("[config] 解析配置文件成功")
	zap.S().Infof("[config] NacosConfig: %+v", *global.NacosConfig)

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
		zap.S().Fatalf("[nacos] 创建配置客户端失败 %s", err.Error())
	}
	zap.S().Infof("[nacos] 创建配置客户端成功")

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Fatalf("[nacos] 读取配置文件失败 %s", err.Error())
	}

	zap.S().Infof("[nacos] 读取配置文件成功 %s", content)

	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("[nacos] 反序列化配置文件失败 %s", err.Error())
	}

	zap.S().Infof("[nacos] 反序列化配置文件成功 %+v", global.ServerConfig)

	serverConfig := global.ServerConfig
	zap.S().Infof("ServerConfig:\nname:%s host:%s port:%d\n", serverConfig.Name, serverConfig.Host, serverConfig.Port)
	zap.S().Infof("UserSrvInfo:\n %+v\n", serverConfig.UserSrvInfo)
	zap.S().Infof("ConsulInfo:\n %+v\n", serverConfig.ConsulInfo)

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
		OnChange: func(namespace string, group string, dataId string, data string) {
			zap.S().Infof("[nacos] 配置文件发生变化")
			fmt.Println("gropu ", group, "dataId ", dataId, "data ", data)
			err = json.Unmarshal([]byte(data), global.ServerConfig)
			if err != nil {
				zap.S().Fatal("[nacos] 反序列化配置文件失败", zap.Error(err))
			}

			zap.S().Infof("[nacos] 配置文件发生变化", global.ServerConfig)
		},
	})
	if err != nil {
		zap.S().Fatal("[nacos] 监听配置文件失败", zap.Error(err))
		panic(err)
	}

}
