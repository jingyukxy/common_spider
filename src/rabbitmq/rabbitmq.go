package rabbitmq

import (
	"awesomeProject/src/config"
	log "awesomeProject/src/logs"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

var mqConn *amqp.Connection
var mqChan *amqp.Channel

// 定义生产者
type Producer interface {
	MsgContent() string
}

// 定义接收者
type Receiver interface {
	Consumer([]byte) error
}

type RabbitMQ struct {
	rabbitConfig *config.RabbitMQConfig
	connection   *amqp.Connection
	channel      *amqp.Channel
	queueName    string // 队列名
	routingKey   string // key
	exchangeName string // exchange name
	exchangeType string // exchange type
	producerList []Producer
	receiverList []Receiver
	mu           sync.RWMutex
	quitChan     chan struct{}
}

type QueueExchange struct {
	QuName string
	RtKey  string
	ExName string
	ExType string
}

// 连接
func (r *RabbitMQ) mqConnect() (err error) {
	rabbitUrl := r.rabbitConfig.RabbitMQURL
	mqConn, err = amqp.DialConfig(rabbitUrl, amqp.Config{
		Vhost:      r.rabbitConfig.VHost,
		ChannelMax: r.rabbitConfig.ChannelMax,
		Locale:     r.rabbitConfig.Locale,
		Heartbeat:  10 * time.Second,
	})
	r.connection = mqConn
	if err != nil {
		log.Logger.WithError(err).Error("open mq connect error")
		return
	}
	mqChan, err = mqConn.Channel()
	if err != nil {
		log.Logger.WithError(err).Error("open mq channel error!")
		return
	}
	r.channel = mqChan
	log.Logger.WithFields(logrus.Fields{
		"url":   r.rabbitConfig.RabbitMQURL,
		"vhost": r.rabbitConfig.VHost,
	}).Info("连接Rabbit启动服务成功")
	return nil
}

func (r *RabbitMQ) Close() {
	r.mqClose()
}

func (r *RabbitMQ) mqClose() {
	close(r.quitChan)
	err := r.channel.Close()
	if err != nil {
		log.Logger.WithError(err).Error("close mq channel failed")
	}
	err = r.connection.Close()
	if err != nil {
		log.Logger.WithError(err).Error("close mq connection failed!")
	}
}

//创建新的操作对象
func New(q *QueueExchange, config *config.RabbitMQConfig) *RabbitMQ {
	// 需要定义成属性
	return &RabbitMQ{
		rabbitConfig: config,
		queueName:    q.QuName,
		routingKey:   q.RtKey,
		exchangeName: q.ExName,
		exchangeType: q.ExType,
		quitChan:     make(chan struct{}, 1),
	}
}

// 启动RabbitMq客户端，并初始化
func (r *RabbitMQ) Start() {
	// 开启监听发送
	for _, producer := range r.producerList {
		go r.listenProducer(producer)
	}
	// 开启接收
	for _, receiver := range r.receiverList {
		go r.listenReceiver(receiver)
	}
}

// 注册发送指定队列指定路由的生产者
func (r *RabbitMQ) RegisterProducer(producer Producer) {
	r.producerList = append(r.producerList, producer)
}

func (r *RabbitMQ) IsClosed() bool {
	return r.connection.IsClosed()
}

// 发送任务
func (r *RabbitMQ) listenProducer(producer Producer) {
	// 验证是否正常，不正常则重连
	if r.channel == nil || r.connection.IsClosed() {
		err := r.mqConnect()
		if err != nil {
			log.Logger.WithError(err).Error("MQ init connection error!")
			return
		}
	}
	// 用于检测队列是否存在，已经存在则不需要重复声明
	_, err := r.channel.QueueDeclarePassive(r.queueName, true, false, false, false, nil)
	if err != nil {
		r.channel, err = r.connection.Channel()
		// 队列不存在声明队列
		// name :队列名称; durable:是否持久化,队列存盘，true服务重启后信息不会丢失,影响性能,autoDelete:是否自动删除;noWait:是否是非阻塞
		// true为是，不等待RMQ返回信息；args:参数，传nil即可,exclusive:是否设置排他性
		_, err = r.channel.QueueDeclare(r.queueName, true, false, false, false, nil)
		log.Logger.Info("declare query ", r.queueName)
		if err != nil {
			log.Logger.WithError(err).Error("MQ Declare Queue Failed")
			return
		}
	}
	// 绑定队列
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true, nil)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"queueName":    r.queueName,
			"routingKey":   r.routingKey,
			"exchangeName": r.exchangeName,
		}).WithError(err).Error("MQ binding Queue failed!")
		return
	}

	//检查exchange 是否已存在，如存在则不需重新声明
	err = r.channel.ExchangeDeclarePassive(r.exchangeName, r.exchangeType, true, false, false, false, nil)
	if err != nil {
		// 上面channel可能会被断开
		r.channel, err = r.connection.Channel()
		if err != nil {
			log.Logger.WithError(err).Error("channel error")
			return
		}
		log.Logger.Info("exchange declare")
		//不存在则注册exchange
		// name:exchange name, kind:exchange type;durable:持久化;noWait:是否阻塞;autoDelete:是否自动删除;internal是否为内部
		err = r.channel.ExchangeDeclare(r.exchangeName, r.exchangeType, true, false, false, false, nil)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"exchangeType": r.exchangeType,
				"exchangeName": r.exchangeName,
			}).WithError(err).Error("MQ Register Queue failed!")
			return
		}
	}
	// 发送任务消息
	err = r.channel.Publish(r.exchangeName, r.routingKey, false, false, amqp.Publishing{
		ContentType:     "text/plain",
		ContentEncoding: "UTF-8",
		Body:            []byte(producer.MsgContent()),
	})
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"body": producer.MsgContent(),
		}).WithError(err).Error("MQ Send Task failed!")
		return
	}
	log.Logger.Info("publish success.....")
}

// 注册指定队列指定路由的消费者
func (r *RabbitMQ) RegisterReceiver(receiver Receiver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.receiverList = append(r.receiverList, receiver)
}

// 监听接收者接收任务
func (r *RabbitMQ) listenReceiver(receiver Receiver) {
	log.Logger.Info("receiver start")
	// 连接不存在需要重连
	if r.channel == nil || r.connection.IsClosed() {
		// 需要连接
		err := r.mqConnect()
		if err != nil {
			log.Logger.WithError(err).Error("MQ init connect error")
			return
		}
	}
	// 检查队列是否存在
	_, err := r.channel.QueueDeclarePassive(r.queueName, true, false, false, false, nil)
	if err != nil {
		r.channel, err = r.connection.Channel()
		if err != nil {
			log.Logger.WithError(err).Error("channel error")
			return
		}
		log.Logger.Info("declareQueue......")
		//队列不存在则声明
		_, err = r.channel.QueueDeclare(r.queueName, true, false, false, false, nil)
		if err != nil {
			//声明失败
			log.Logger.WithError(err).Error("MQ Declare queue failed")
		}
	}

	//绑定任务
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true, nil)
	if err != nil {
		log.Logger.WithError(err).Error("MQ Binding Queue failed!")
		return
	}
	// 获取消费通道，确保MQ一个一个发送消息
	err = r.channel.Qos(1, 0, true)
	if err != nil {
		log.Logger.WithError(err).Error("Qos error")
		return
	}
	msgList, err := r.channel.Consume(r.queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Logger.WithError(err).Error("MQ Fetch Consumer Channel Failed!")
		return
	}
	// 注册被关闭的情况
	closeCh := make(chan bool, 1)
	go func(chan<- bool) {
		cc := make(chan *amqp.Error)
		e := <-r.channel.NotifyClose(cc)
		log.Logger.WithError(e).Error("RabbitMQ channel closed")
		closeCh <- true
	}(closeCh)

	go func(<-chan bool) {
		for {
			select {
			case <-r.quitChan:
				log.Logger.Info("RabbitMQ服务正常退出")
				return
			case msg := <-msgList:
				err := receiver.Consumer(msg.Body)
				if err != nil {
					err = msg.Ack(true)
					if err != nil {
						log.Logger.WithError(err).Error("MQ Ack Unfinished Msg Failed!")
					}
				} else {
					// 确认消息，必须为 false
					err = msg.Ack(false)
					if err != nil {
						log.Logger.WithError(err).Error("MQ Ack Finished Msg Failed!")
					}
				}
			case <-closeCh:
				// 断线重连
				if r.connection.IsClosed() {
					err = r.mqConnect()
					if err != nil {
						log.Logger.WithError(err).Error("Reconnect error!")
						return
					}
					log.Logger.Info("Reconnect with connection")
				} else {
					r.channel, err = r.connection.Channel()
					if err != nil {
						log.Logger.WithError(err).Error("Reopen channel error!")
						return
					}
					log.Logger.Info("Reconnect with channel")
				}
			}
		}
	}(closeCh)
	//for msg := range msgList {
	//	err := receiver.Consumer(msg.Body)
	//	if err != nil {
	//		err = msg.Ack(true)
	//		if err != nil {
	//			log.Logger.WithError(err).Error("MQ Ack Unfinished Msg Failed!")
	//			return
	//		}
	//	} else {
	//		// 确认消息，必须为 false
	//		err = msg.Ack(false)
	//		if err != nil {
	//			log.Logger.WithError(err).Error("MQ Ack Finished Msg Failed!")
	//			return
	//		}
	//	}
	//}
}
