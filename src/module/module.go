package module

import (
	log "awesomeProject/src/logs"
	"sync"
)

var managerInstance *DefaultModuleManager
var managerLock sync.Once

func GetManagerInstance() *DefaultModuleManager {
	managerLock.Do(func() {
		managerInstance = &DefaultModuleManager{}
	})
	return managerInstance
}

// 模块
type IModule interface {
	Init() error
	Destroy()
}

// 模块管理器
type IModuleManager interface {
	Init()
	RegisterModule(module IModule)
	Destroy()
}

type DefaultModuleManager struct {
	modules []IModule
}

func (moduleManager *DefaultModuleManager) RegisterModule(module IModule) {
	moduleManager.modules = append(moduleManager.modules, module)
}

func (moduleManager *DefaultModuleManager) Destroy() {
	for _, module := range moduleManager.modules {
		module.Destroy()
	}
}

func (moduleManager *DefaultModuleManager) Init() {
	// 先添加进来
	moduleManager.RegisterModules()
	for _, module := range moduleManager.modules {
		if err := module.Init(); err != nil {
			log.Logger.WithError(err).Fatal("module start error!")
		}
	}
}

func (moduleManager *DefaultModuleManager) RegisterModules() {
	// 先添加db
	moduleManager.RegisterModule(GetDbProviderInstance())
	// db provider
	moduleManager.RegisterModule(GetDbProcessorInstance())
	// cache provider
	moduleManager.RegisterModule(GetCacheProviderInstance())
}
