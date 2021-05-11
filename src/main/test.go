package main

import (
	"awesomeProject/src/cache"
	"awesomeProject/src/config"
	"awesomeProject/src/db"
	log "awesomeProject/src/logs"
	"awesomeProject/src/rabbitmq"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type DbModel interface {
	GetTableName() string
}

type Basic struct {
	Id          int32     `DbType:"INT" ColumnName:"BASIC_ID" PK:"LastId"`
	Title       string    `DbType:"VARCHAR" ColumnName:"BASIC_TITLE"`
	Description string    `DbType:"TEXT" ColumnName:"BASIC_DESCRIPTION"`
	Thumbnails  string    `DbType:"VARCHAR" ColumnName:"BASIC_THUMBNAILS"`
	Hit         int64     `DbType:"BIGINT" ColumnName:"BASIC_HIT"`
	Sort        int32     `DbType:"INT" ColumnName:"BASIC_SORT"`
	AddTime     time.Time `DbType:"DATETIME" ColumnName:"BASIC_DATETIME"`
	UpdateTime  time.Time `DbType:"DATETIME" ColumnName:"BASIC_UPDATETIME"`
	PeopleId    int32     `DbType:"INT" ColumnName:"BASIC_PEOPLEID"`
	CategoryId  int32     `DbType:"INT" ColumnName:"BASIC_CATEGORYID"`
	AppId       int32     `DbType:"INT" ColumnName:"BASIC_APPID"`
	ModelId     int32     `DbType:"INT" ColumnName:"BASIC_MODELID"`
	Comment     int32     `DbType:"INT" ColumnName:"BASIC_COMMENT"`
	Collect     int32     `DbType:"INT" ColumnName:"BASIC_COLLECT"`
	Share       int32     `DbType:"INT" ColumnName:"BASIC_SHARE"`
	BasicType   string    `DbType:"VARCHAR" ColumnName:"BASIC_TYPE"`
}

func (basic *Basic) GetTableName() string {
	return "basic"
}

var globalConfig *config.Config

func SInit() {
	var etcFile = flag.String("c", "", "etc config file")
	flag.Parse()
	if *etcFile == "" {
		logrus.Fatal("etc file should not be empty!")
	}
	config, err := config.GetConfig(*etcFile)
	if err != nil {
		logrus.WithError(err).Fatal("load config err!")
	}
	globalConfig = config
	log.InitLogger(&config.Log)
	log.Logger.Info("load config =>", config.DbConfig.Database)
	_, err = db.InitDbConnection(&config.DbConfig)

	if err != nil {
		log.Logger.WithError(err).Error("init db error ")
	}
	log.Logger.Info("load db success")
	//cache.InitRedis(config.Redis)
	cache.NewRedisCache(config.Redis)
	log.Logger.Info("load redis success")
}

type Pc struct {
	Name string `name:"Name",id:name`
	Id   int32  `name:"Id",id:id`
	FC   *Fc
}

type Fc interface {
	GetName() string
}

func (pc *Pc) GetName(fc Fc) string {
	return fc.GetName()
}

type Ads struct {
	Id      int16  `DbType:"SMALLINT" ColumnName:"ads_id" PK:"LastId"`
	Name    string `DbType:"VARCHAR" ColumnName:"ads_name"`
	Content string `DbType:"TEXT" ColumnName:"ads_content"`
}

func testSql() chan int {
	var ch = make(chan int, 1)
	go func() {
		defer func() {
			ch <- 1
		}()
		var sqlSession *db.DefaultSqlSession
		var basic Basic
		//_, err := db.GetDbConnection()
		//if err != nil {
		//	log.Logger.WithError(err).Error("get db connection error!")
		//	return
		//}
		sqlSession = db.NewSqlSession()
		// 查多记录
		rows, err := sqlSession.SelectAll("SELECT * FROM basic limit 3", basic)
		if err != nil {
			log.Logger.WithError(err).Error("select all error!")
			return
		}
		log.Logger.Info("rows num ", len(rows))
		//for i := range rows {
		//	b := rows[i].(Basic)
		//	log.Logger.Info(b)
		//}
		// 单记录
		bc, err := sqlSession.SelectOne("select * from basic where BASIC_ID = ?", basic, 1)
		if err != nil {
			log.Logger.WithError(err).Error("select one error!!")
		}
		log.Logger.Info("get row=>", bc.(Basic).Id)
		//pq.Array()
		//log.Logger.Info(rows)
		// 更新
		affRow, err := sqlSession.Update("update basic set BASIC_PEOPLEID = ? where BASIC_ID=?", 300, 3)
		if err != nil {
			log.Logger.WithError(err).Error("update data error!!")
		}
		log.Logger.Info("update basic affectedRows ", affRow)
		// 添加
		affRow, retValue, err := sqlSession.Insert("insert into basic (BASIC_TITLE, BASIC_DESCRIPTION,BASIC_HIT)values(?,?,?)",
			basic, "Hell", "dfdsfdsfd", 10)
		log.Logger.Info("insert basic affectedRows , basic id ", affRow, retValue.(Basic).Id)
		// 切库
		sqlSession.SetDbFlag("vod")
		var ads Ads
		allAds, err := sqlSession.SelectAll("select * from ff_ads", ads)
		if err != nil {
			log.Logger.WithError(err).Error("from ads data error!!")
			return

		}
		log.Logger.Info("get ads all len ", len(allAds))
		//for i := range allAds {
		//	log.Logger.Info(allAds[i])
		//}
		// 切库
		sqlSession.SetDbFlag("default")
		// 重新查询
		affRow, retValue, err = sqlSession.Insert("insert into basic (BASIC_TITLE, BASIC_DESCRIPTION,BASIC_HIT)values(?,?,?)",
			basic, "Hell", "dfdsfdsfd", 10)
		if err != nil {
			log.Logger.WithError(err).Error("failed insert")
			return
		}
		log.Logger.Info("insert basic affectedRows , basic id ", affRow, retValue.(Basic).Id)
	}()
	return ch
}

func testRedis() {
	_, err := cache.Exec("set", "aft", "100", "EX", 1)
	if err != nil {
		log.Logger.Fatal(err)
	}
}
func testGoRedis() chan int {
	var rdChan = make(chan int)
	go func() {
		//err := cache.NewRedisCache()
		//if err != nil {
		//	log.Logger.WithError(err).Error("set tff error!")
		//}
		//tff, err := cache.Get("tff")
		//if err != nil {
		//	log.Logger.WithError(err).Error("get tff error!")
		//}
		//log.Logger.WithField("tff", tff).Info("get tff")
		////select {
		////case x := <-rdChan:
		////	log.Logger.WithField("x", x).Info("got that")
		////}
		//rdChan <- 100
		//close(rdChan)
	}()
	return rdChan
}

type TestPro struct {
	msgContent string
}

func (t *TestPro) MsgContent() string {
	return t.msgContent
}

func (t *TestPro) Consumer(dataByte []byte) error {
	//msg := string(dataByte)
	m := make(map[string]interface{})
	err := json.Unmarshal(dataByte, &m)
	if err != nil {
		log.Logger.WithError(err).Info("parse json error!")
	}
	log.Logger.Info(m)
	return nil
}

func testRb() {
	msgBody := "hello..."
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	q, err := ch.QueueDeclare(
		"testqueue", //Queue name
		true,        //durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	err = ch.Publish(
		"",     //exchange
		q.Name, //routing key(queue name)
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //Msg set as persistent
			ContentType:  "text/plain",
			Body:         []byte(msgBody),
		})
	if err != nil {
		panic(err)
	}
}

func testRabbitMQ() {
	msg := fmt.Sprintf("Test Job Task")
	t := &TestPro{
		msg,
	}
	queueExchange := rabbitmq.QueueExchange{
		"items_qqMovie",
		"qq_movie_item",
		"amq.direct",
		"direct",
	}
	//pmq := rabbitmq.New(&queueExchange, &globalConfig.RabbitMQ)
	//pmq.RegisterProducer(t)
	cmq := rabbitmq.New(&queueExchange, &globalConfig.RabbitMQ)
	cmq.RegisterReceiver(t)
	//pmq.Start()
	time.Sleep(1 * time.Second)
	cmq.Start()
}

type fake interface {
	Get() int
}

var _ fake = (*fakeImpl)(nil)

//func
type fk func(p string)

func (f fk) Get() int {
	f("ddd")
	return 0
}

type fakeImpl struct {
	fk
}

func tt() {
	p := fakeImpl {
		fk: func(d string) {
			fmt.Println(d)
		},
	}
	p.Get()
}

func testGoRoutine() {
	var chs []chan int
	for i := 0; i < 100; i++ {
		rdc := testSql()
		chs = append(chs, rdc)
	}
	for i := range chs {
		x := <-chs[i]
		log.Logger.WithField("XX", x).Info("Get Value from cha")
		close(chs[i])
		//chs[i] <- 2
	}
}
