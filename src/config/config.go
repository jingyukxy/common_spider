package config

import (
	//"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

func GetConfig(file string) (*Config, error) {
	config := new(Config)
	viperLoader := viper.New()
	viperLoader.SetConfigFile(file)
	viperLoader.SetConfigType("toml")
	err := viperLoader.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viperLoader.Unmarshal(config)
	//_, err := toml.DecodeFile(file, &config)
	if err != nil {
		return nil, err
	}
	return config, err
}

type DownloadConfig struct {
	DownloadPath   string
	DownloadSuffix string
	PrefixUrl      string
}

type AppConfig struct {
	AppName        string
	InternalServer string
}

type Config struct {
	AppConfig AppConfig
	Redis     RedisConfig
	DbConfig  DataSourceConfig
	RabbitMQ  RabbitMQConfig
	Log       LogConfig
	Download  DownloadConfig
}

type DataSourceConfig struct {
	Database map[string]DataConfig
}

type LogConfig struct {
	MaxCount int
	AppId    int
	FileName string
	ApName   string
}

type RedisConfig struct {
	Database           int
	MaxIdle            int
	MaxActive          int
	IdleTimeout        int
	Network            string
	Address            string
	DialReadTimeout    int
	DialWriteTimeout   int
	DialConnectTimeout int
}

type DataConfig struct {
	Driver             string
	ConnectionString   string
	MaxOpenConnections uint8
	MaxIdleConnections uint8
	ConnId             uint8
	ConnName           string
	ConnMaxLifetime    int
}

type RabbitMQConfig struct {
	RabbitMQURL string
	VHost       string
	Locale      string
	ChannelMax  int
	Exchanges   map[string]RabbitMQExchange
}

type RabbitMQExchange struct {
	QueueName    string
	RoutingKey   string
	ExchangeName string
	ExchangeType string
}
