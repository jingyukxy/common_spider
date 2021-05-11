package handler

import (
	"awesomeProject/src/dao"
	log "awesomeProject/src/logs"
	"awesomeProject/src/message"
	"awesomeProject/src/model"
	"awesomeProject/src/module"
	"awesomeProject/src/utils"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"reflect"
)

var _ Handler = (*QQMovieHandler)(nil)

type QQMovieHandler struct {
	movieDao *dao.MovieDao
}

func NewQQMovieHandler() *QQMovieHandler {
	return &QQMovieHandler{
		movieDao: dao.NewMovieDao(),
	}
}

// 判断是否为已经处理过的数据，只需要更新则不需要切词等操作
func (movie *QQMovieHandler) IsNewMessage(msg message.IMessage) (bool, error) {
	movieMsg := reflect.ValueOf(msg).Interface().(*message.MovieMessage)
	// 先从redis中取出来
	digest := utils.Md5WithString(movieMsg.Name)
	key := "movie_spider"
	_, err := module.GetCacheProviderInstance().RedisCache.HGet(key, digest)
	if err == redis.Nil {
		//redis中没到 去数据库中找
		rows, err := movie.movieDao.SelectMoviesByName(movieMsg.Name)
		if err != nil {
			return false, err
		}
		// 数据库中无此数据
		if len(rows) == 0 {
			return true, nil
		}
	} else if err != nil {
		return false, err
	}
	return false, nil
}

// 业务处理
func (movie *QQMovieHandler) Handle(spiderInfo *model.SpiderInfo, msg message.IMessage) (err error) {
	movieMsg := reflect.ValueOf(msg).Interface().(*message.MovieMessage)
	log.Logger.WithFields(logrus.Fields{"name": movieMsg.Name, "spider": spiderInfo.Name}).Info("start process")
	dbProcessor := module.GetDbProcessorInstance()
	var asyncSave module.AsyncCall = func() {
		err := movie.movieDao.SaveMovie(spiderInfo,*movieMsg)
		if err != nil {
			log.Logger.WithError(err).Error("async save movie error!!!!")
		} else {
			log.Logger.Info("async save movie finished")
			// 把数据存储在redis中
			err := module.GetCacheProviderInstance().HSet("movie_spider", utils.Md5WithString(movieMsg.Name), "1")
			if err != nil {
				log.Logger.WithError(err).Error("存储redis 去重数据失败")
			}
		}
	}
	dbProcessor.AddAsyncCall(asyncSave)
	return err
}
