package dispatcher

import (
	"awesomeProject/src/config"
	"awesomeProject/src/dao"
	"awesomeProject/src/executor"
	"awesomeProject/src/handler"
	log "awesomeProject/src/logs"
	"awesomeProject/src/message"
	"awesomeProject/src/model"
	"awesomeProject/src/module"
	"awesomeProject/src/rabbitmq"
	"awesomeProject/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buptmiao/parallel"
	"github.com/go-redis/redis/v8"
	"github.com/ivpusic/grpool"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"time"
)

// 采集器多个对应多组mq
type MqProvider struct {
	mq       *rabbitmq.RabbitMQ     // mq
	exchange rabbitmq.QueueExchange // 交换机
	receiver rabbitmq.Receiver      // 接收器，用来转发，这里全部使用dispatcher
}

func NewProvider(exc rabbitmq.QueueExchange, rec rabbitmq.Receiver) *MqProvider {
	return &MqProvider{
		exchange: exc,
		receiver: rec,
	}
}

// 路由转发器
type MqDispatcher struct {
	mqProviders []*MqProvider
	pool        *grpool.Pool
}

var lock sync.Once
var instance *MqDispatcher

func GetInstance() *MqDispatcher {
	lock.Do(func() {
		instance = &MqDispatcher{}
	})
	return instance
}

// 停服
func (mqDispatcher *MqDispatcher) Stop() {
	log.Logger.Info("开始停服")
	// 释放 goroutine 池
	mqDispatcher.pool.Release()
	for _, provider := range mqDispatcher.mqProviders {
		if provider.mq != nil {
			log.Logger.WithField("queue", provider.exchange.QuName).Info("停止监听器")
			provider.mq.Close()
		} else {
			log.Logger.Error("没有可用监听器")
		}
	}
}

// 全局注册监听
func (mqDispatcher *MqDispatcher) RegisterMqProviders() {
	gConfig, err := config.GetConfigInstance().GlobalConfig()
	if err != nil {
		log.Logger.WithError(err).Error("load config error")
		return
	}
	mqConfig := gConfig.RabbitMQ
	mqExConfig := gConfig.RabbitMQ.Exchanges
	for k, v := range mqExConfig {
		ex := rabbitmq.QueueExchange{
			QuName: v.QueueName,
			RtKey:  v.RoutingKey,
			ExName: v.ExchangeName,
			ExType: v.ExchangeType,
		}
		provider := NewProvider(ex, mqDispatcher)
		mqDispatcher.RegisterMqProvider(provider, mqConfig)
		log.Logger.WithField("queue", k).Info("Load queue ready.")
	}
}

// 注册mq
func (mqDispatcher *MqDispatcher) RegisterMqProvider(provider *MqProvider, config config.RabbitMQConfig) {
	provider.mq = rabbitmq.New(&provider.exchange, &config)
	provider.mq.RegisterReceiver(provider.receiver)
	mqDispatcher.mqProviders = append(mqDispatcher.mqProviders, provider)
}

// 全部初始化
func (mqDispatcher *MqDispatcher) Start() {
	mqDispatcher.pool = grpool.NewPool(50, 33)
	for _, mqProvider := range mqDispatcher.mqProviders {
		mqProvider.mq.Start()
	}
}

// 消费监听
func (mqDispatcher *MqDispatcher) Consumer(body []byte) error {
	if body == nil || len(body) < 10 {
		log.Logger.Error("now body in msg")
		return errors.New("data body error")
	}
	data := make(map[string]interface{})
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Logger.WithError(err).Error("unmarshal body error!")
		return err
	}
	// 先用redis去重，key 为sp_items_加消息标识组成，内容key为消息md5
	appCode := int(data["app_code"].(float64))
	key := fmt.Sprintf("sp_items_%d", appCode)
	digest := utils.Md5(body)
	_, err = module.GetCacheProviderInstance().HGet(key, digest)
	if err == redis.Nil {
		log.Logger.WithFields(logrus.Fields{"key": key, "digest": digest}).Info("key未在Redis中，开始处理")
		err := module.GetCacheProviderInstance().HSet(key, digest, "1")
		if err != nil {
			return err
		}
		// 这里开始处理分发
		mqDispatcher.pool.JobQueue <- func() {
			mqDispatcher.Dispatch(data, body)
		}
	} else if err != nil {
		log.Logger.WithError(err).Error("获取Redis失败")
		return err
	} else {
		log.Logger.WithField("key", key).Info("已经存在，跳过此条数据->", digest)
		return nil
	}
	return nil
}

// 通过代码获取采集信息,先从redis中获取，找不到再去数据库中找，找到后存入redis，失效时间一小时
func (mqDispatcher *MqDispatcher) getSpiderInfo(code int) (*model.SpiderInfo, error) {
	var spiderInfo = model.SpiderInfo{}
	key := fmt.Sprintf("spider_info_code_%d", code)
	result, err := module.GetCacheProviderInstance().Get(key)
	// redis没结果
	if err == redis.Nil {
		// 不在redis中 去数据库中查
		spiderDao := dao.NewSpiderDao()
		spiderInfo, err = spiderDao.GetSpiderInfoByCode(code)
		// 数据库中也没有找到，直接返回
		if err != nil {
			log.Logger.WithError(err).Error("get spider info error!")
			return nil, err
		}
		// 查找后把结果保存在redis中
		data, err := json.Marshal(spiderInfo)
		if err != nil {
			log.Logger.WithField("spiderInfo", spiderInfo).WithError(err).Error("marshal spider info json error!!")
		} else {
			err = module.GetCacheProviderInstance().Set(key, data, 60*time.Minute)
			// 存储redis失败
			if err != nil {
				log.Logger.WithError(err).Error("save redis cache error!")
			}
		}
		return &spiderInfo, nil
	} else if err != nil {
		log.Logger.WithError(err).Error("from redis error!")
		return nil, err
	} else {
		// 有结果
		buff := []byte(result)
		err = json.Unmarshal(buff, &spiderInfo)
		if err != nil {
			log.Logger.WithError(err).Error("unmarshal json from redis error!")
			return nil, err
		} else {
			return &spiderInfo, nil
		}
	}
}

// 合并数据
func (mqDispatcher *MqDispatcher) CombineMsg(transMsg message.IMessage, downMsg message.IMessage, info *model.SpiderInfo) message.IMessage {
	// 如里不需要下载也不需要翻译，返回空
	if info.NeedDownLoad == 0 && info.NeedTranslate == 0 {
		return nil
	}
	// 不需要翻译，则返回下载内容
	if info.NeedTranslate == 0 {
		return downMsg
	}
	// 只需要翻译则返回 翻译内容
	if info.NeedDownLoad == 0 {
		return transMsg
	}
	// 如同时操作则需要通过反射把下载后的结果合并到翻译中，创建新的消息返回
	vType := reflect.TypeOf(transMsg).Elem()
	newValue := reflect.New(vType).Elem()
	oldValue := reflect.ValueOf(transMsg).Elem()
	downValue := reflect.ValueOf(downMsg).Elem()
	// 创建新的message后进行数据拷贝
	reflect.Copy(newValue, oldValue)
	fieldName := downMsg.GetDownField()
	// 开始合并赋值
	downResult := downValue.FieldByName(fieldName).String()
	newValue.FieldByName(fieldName).Set(reflect.ValueOf(downResult))
	msg := newValue.Interface().(message.IMessage)
	return msg
}

// 核心转发器
func (mqDispatcher *MqDispatcher) Dispatch(data map[string]interface{}, body []byte) {
	log.Logger.WithField("app_code", data["app_code"]).Info("dispatch body")
	appCode := int(data["app_code"].(float64))
	processor := handler.GetAssistantInstance().GetProcessor(appCode)
	if processor == nil {
		log.Logger.Info("没有找到相应的处理器")
		return
	}
	// 先转码
	err := processor.Message.DoDecoder(body)
	// 转码出错
	if err != nil {
		log.Logger.WithError(err).WithFields(logrus.Fields{
			"app": data,
		}).Info("message decoder error!")
		return
	}
	// 从mysql中取出spiderDao,然后执行handle方法
	// 先执行filter中的方法，然后再进行下一步操作
	spiderInfo, err := mqDispatcher.getSpiderInfo(processor.Message.GetMsgCode())
	// 获取失败则直接返回
	if err != nil {
		log.Logger.WithError(err).Error("get spider info error!")
		return
	}
	if spiderInfo == nil {
		log.Logger.Info("Can not find the spider info in persistent")
		return
	}
	// 先看是不是最新的数据
	isNew, err := processor.Handler.IsNewMessage(processor.Message)
	if err != nil {
		log.Logger.WithError(err).Error("查询数据失败，队列处理返回")
		return
	}
	if isNew {
		// 先同步处理下载和翻译，然后再处理打标签
		downloadExecutor := executor.DownloadExecutor{}
		transExecutor := executor.TranslateExecutor{}
		var downMsg message.IMessage
		var tranMsg message.IMessage
		p := parallel.NewParallel()
		p.Register(downloadExecutor.Execute, spiderInfo, processor.Message).SetReceivers(&downMsg)
		p.Register(transExecutor.Execute, spiderInfo, processor.Message).SetReceivers(&tranMsg)
		p.Run()
		// 合并数据
		newMsg := mqDispatcher.CombineMsg(tranMsg, downMsg, spiderInfo)
		if newMsg == nil {
			newMsg = processor.Message
		}
		// 接下来处理切词
		if spiderInfo.NeedTag == 1 {
			segExecutor := executor.SegmentExecutor{}
			newMsg = segExecutor.Execute(spiderInfo, newMsg)
		}
		// 最后异步处理存储
		err = processor.Handler.Handle(spiderInfo, newMsg)
		if err != nil {
			log.Logger.WithError(err).Error("处理数据失败")
		}
	} else {
		// 不是最新的，直接处理 不下载，不处理
		err = processor.Handler.Handle(spiderInfo, processor.Message)
		if err != nil {
			log.Logger.WithError(err).Error("处理数据失败")
		}
	}
}
