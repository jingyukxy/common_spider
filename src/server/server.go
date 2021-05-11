package server

import (
	config2 "awesomeProject/src/config"
	"awesomeProject/src/dao"
	"awesomeProject/src/dispatcher"
	"awesomeProject/src/handler"
	log "awesomeProject/src/logs"
	"awesomeProject/src/module"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
)

func New() *Bootstrap {
	return &Bootstrap{}
}

// 启动结构体
type Bootstrap struct {
	status   bool // 状态
	ch       chan struct{}
	listener net.Listener
}

// 启动状态
func (bootstrap *Bootstrap) Status() bool {
	return bootstrap.status
}

// 启动后预先加载
func (bootstrap *Bootstrap) afterLoad() {
	// processor 绑定处理
	handler.PreLoad()
	// 绑定信息加载
	movieDao := dao.NewMovieDao()
	movieDao.PreloadCategory()
}

func (bootstrap *Bootstrap) Init(configFile string) {
	bootstrap.ch = make(chan struct{}, 1)
	configManager := config2.GetConfigInstance()
	err := configManager.Init(configFile)
	if err != nil {
		logrus.WithError(err).Fatal("加载日志文件失败，直接停服!")
	}
	log.InitLogger(configManager.GetLogConfig())
	log.Logger.Info("初始化配置文件日志完成")
	// 初始化模块
	module.GetManagerInstance().Init()
	// 预加载
	bootstrap.afterLoad()
	// 注册dispatcher
	dispatcher.GetInstance().RegisterMqProviders()
}

func (bootstrap *Bootstrap) startInternalServer() {
	appConfig := config2.GetConfigInstance().AppConfig
	listen, err := net.Listen("tcp", appConfig.InternalServer)
	if err != nil {
		log.Logger.WithError(err).Fatal("启动服务失败")
	}
	sf := fmt.Sprintf("Server:[%s]服务器启动成功,地址:[%s]", appConfig.AppName, appConfig.InternalServer)
	log.Logger.Info(sf)
	bootstrap.listener = listen
	//ch := make(chan struct{}, 1)
	for bootstrap.status {
		conn, err := listen.Accept()
		if err != nil {
			select {
			case <-bootstrap.ch:
				// 优雅退出
				bootstrap.Stop()
				log.Logger.Info("退出服务器")
				return
			}
			log.Logger.WithError(err).Error("接收客户端失败")
			continue
		}
		go bootstrap.process(conn)
	}
}

func (bootstrap *Bootstrap) process(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Logger.WithError(err).Error("获取信息失败!")
		}
		str := string(buf[:n])
		switch str {
		case "quit":
			bootstrap.status = false
			bootstrap.listener.Close()
			close(bootstrap.ch)
			conn.Write([]byte("Server process shutdown.."))
			return
		case "status":
			log.Logger.Info("status")
			return
		default:
			log.Logger.WithField("cmd", str).Info("no command!")
			return
		}
	}
}

// 启动方法
func (bootstrap *Bootstrap) Start(configFile string) {
	// 初始化
	bootstrap.Init(configFile)
	// 启动dispatcher
	dispatcher.GetInstance().Start()
	bootstrap.status = true
	bootstrap.startInternalServer()
}

func (bootstrap *Bootstrap) ShutDownServer() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		logrus.Info("服务器已经关闭,不能发送信息")
		return
	}
	q := "quit"
	_, err = conn.Write([]byte(q))
	if err != nil {
		logrus.Info("服务器写入处理失败")
	}
}

func (bootstrap *Bootstrap) Stop() {
	// 先把队列任务停了
	dispatcher.GetInstance().Stop()
	// 然后模块停止
	module.GetManagerInstance().Destroy()
	log.Logger.Info("停服完成...")
}
