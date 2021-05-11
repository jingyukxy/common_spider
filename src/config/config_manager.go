package config

import (
	"errors"
	"log"
	"sync"
)

var instance *Manager
var onceLock sync.Once

func GetConfigInstance() *Manager {
	onceLock.Do(func() {
		instance = &Manager{}
	})
	return instance
}

type Manager struct {
	*Config
}

func (configManager *Manager) Init(fileName string) error {
	gConfig, err := GetConfig(fileName)
	if err != nil {
		log.Printf("Initial config error! %v", err)
		panic(err)
	}
	configManager.Config = gConfig
	return err
}

func (configManager *Manager) GlobalConfig() (globalConfig *Config, err error) {
	globalConfig = configManager.Config
	if globalConfig == nil {
		return nil, errors.New("CM global config is not init")
	}
	return
}

func (configManager *Manager) GetRabbitMQConfig() *RabbitMQConfig {
	if configManager.Config != nil {
		return &configManager.RabbitMQ
	}
	return nil
}

func (configManager *Manager) GetRedisConfig() *RedisConfig {
	if configManager.Config != nil {
		return &configManager.Config.Redis
	}
	return nil
}

func (configManager *Manager) GetDownloadConfig() *DownloadConfig {
	if configManager.Config != nil {
		return &configManager.Download
	}
	return nil
}

func (configManager *Manager) GetLogConfig() *LogConfig {
	if configManager.Config != nil {
		return &configManager.Log
	}
	return nil
}

func (configManager *Manager) GetDataSourceConfig() *DataSourceConfig {
	if configManager.Config != nil {
		return &configManager.DbConfig
	}
	return nil
}
