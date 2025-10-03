package config

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ServerConfig struct {
	Name          string         `mapstructure:"name" json:"name"`
	Host          string         `mapstructure:"host" json:"host"`
	Port          int            `mapstructure:"port" json:"port"`
	Tags          []string       `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo  GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	UseropSrvInfo GoodsSrvConfig `mapstructure:"userop_srv" json:"userop_srv"`
	JwtInfo       JwtConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo    ConsulConfig   `mapstructure:"consul" json:"consul"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type CheckConfig struct {
	Interval                       string `mapstructure:"interval" json:"interval"`
	DeregisterCriticalServiceAfter string `mapstructure:"deregister_critical_service_after" json:"deregister_critical_service_after"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
}
