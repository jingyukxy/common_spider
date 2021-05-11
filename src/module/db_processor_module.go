package module

import (
	log "awesomeProject/src/logs"
	"sync"
)

var dbProcessorInstance *DbProcessorModule
var dbProcessorOnceLock sync.Once

func GetDbProcessorInstance() *DbProcessorModule {
	dbProcessorOnceLock.Do(func() {
		dbProcessorInstance = &DbProcessorModule{}
		dbProcessorInstance.asyncDataCh = make(chan AsyncCall)
		dbProcessorInstance.cancel = make(chan int, 1)
	})
	return dbProcessorInstance
}

type AsyncCall func()

type DbProcessorModule struct {
	asyncDataCh chan AsyncCall
	cancel      chan int
}

func (dbProcessorModule *DbProcessorModule) Destroy() {
	dbProcessorModule.cancel <- 1
	close(dbProcessorModule.cancel)
	close(dbProcessorModule.asyncDataCh)
}

func (dbProcessorModule *DbProcessorModule) AddAsyncCall(asyncCall AsyncCall) {
	dbProcessorModule.asyncDataCh <- asyncCall
}

func (dbProcessorModule *DbProcessorModule) Init() (err error) {
	go func() {
		for {
			select {
			case asyncCall := <-dbProcessorModule.asyncDataCh:
				go asyncCall()
			case <-dbProcessorModule.cancel:
				log.Logger.Info("db processor normal exit")
				return
			}
		}
	}()
	return nil
}
