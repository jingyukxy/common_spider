package module

import (
	config2 "awesomeProject/src/config"
	"awesomeProject/src/db"
	log "awesomeProject/src/logs"
	"sync"
)

var dbProviderInstance *DbProviderModule
var dbProviderOnceLock sync.Once

func GetDbProviderInstance() *DbProviderModule {
	dbProviderOnceLock.Do(func() {
		dbProviderInstance = &DbProviderModule{}
	})
	return dbProviderInstance
}

type DbProviderModule struct {
	poolConnection *db.PooledConnection
}

func (dbProviderModule *DbProviderModule) Destroy() {
	log.Logger.Info("destroy db provider")
	dbProviderModule.poolConnection.Close()
}

func (dbProviderModule *DbProviderModule) Init() error {
	config, err := config2.GetConfigInstance().GlobalConfig()
	if err != nil {
		log.Logger.WithError(err).Error("get config error!")
		return err
	}
	conn, err := db.InitDbConnection(&config.DbConfig)
	if err != nil {
		log.Logger.WithError(err).Error("Initial db failed!")
		return err
	}
	dbProviderModule.poolConnection = conn
	return err
}
