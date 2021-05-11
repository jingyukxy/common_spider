package db

import (
	"awesomeProject/src/config"
	log "awesomeProject/src/logs"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// 当前数据连接池
var currentConnection *PooledConnection

// 全局锁
var globalMutex sync.Mutex

// 初始化数据库连接 需要数据库配置，从配置文件中获取,此处有些侵入性，可以改为传入Configuration类,或全局声明Configuration
// 全局初始化一次即可，不需要多次初始化
func InitDbConnection(sourceConfig *config.DataSourceConfig) (*PooledConnection, error) {
	// 如多开，需要加锁
	globalMutex.Lock()
	defer globalMutex.Unlock()
	// 已初始化则返回当前连接池
	if currentConnection != nil {
		return currentConnection, nil
	}
	// 初始化多数据源
	var dy DynamicDataSource
	dyData, err := dy.New(sourceConfig)
	if err != nil {
		return nil, err
	}
	// 初始化数据连接池，返回
	pc := new(PooledConnection)
	pc.dataSource = dyData
	currentConnection = pc
	return pc, nil
}

// 创建SqlSession,初始化sqlSession
func NewSqlSession() *DefaultSqlSession {
	// 获取判断全局Connection是否已经初始化
	if currentConnection == nil {
		log.Logger.Error("DbConnection is not Initialized!")
		return nil
	}
	return &DefaultSqlSession{
		Connection: currentConnection,
	}
}

// 数据库配置
type Configuration struct {
	Driver           string // 驱动名
	ConnectionString string // 连接DSN
	ConnectionName   string // 连接名称
	ConnectionId     uint8  // 连接ID
	ConnectionState  bool   // 连接状态
	MaxConn          uint8  // 连接池最大连接数
	MaxIdleConn      uint8  // 连接池空闲初始连接数
	ConnMaxLifetime  int    // 连接生命周期
}

// 数据源定义
type DataSource interface {
	GetConnection() (*sql.DB, error)
}

var _ DataSource = (*SimpleDataSource)(nil)

//var _ DataSource = (*DynamicDataSource)(nil)
var _ SqlSession = (*DefaultSqlSession)(nil)

// 简单数据源
type SimpleDataSource struct {
	DbConfig         config.DataConfig
	DirectConnection *DirectConnection
}

// 直连数据连接
type DirectConnection struct {
	// 连接
	Connection *sql.DB
	Config     Configuration
}

// 初始化数据连接
func (conn *DirectConnection) InitConnection() error {
	db, err := sql.Open(conn.Config.Driver, conn.Config.ConnectionString)
	if err != nil {
		return err
	}
	maxConns := conn.Config.MaxConn
	maxIdle := conn.Config.MaxIdleConn
	maxLifeTime := conn.Config.ConnMaxLifetime
	if maxConns > 0 && maxIdle > 0 && maxLifeTime > 0 {
		db.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Minute)
		db.SetMaxIdleConns(int(maxIdle))
		db.SetMaxOpenConns(int(maxConns))
	} else {
		db.SetConnMaxLifetime(2 * time.Minute)
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
	}
	conn.Config.ConnectionState = true
	conn.Connection = db
	return nil
}

// 初始化数据源，同时初始化数据库连接
func (simpleDataSource *SimpleDataSource) InitDataSource() (err error) {
	dsn := simpleDataSource.DbConfig.ConnectionString
	driver := simpleDataSource.DbConfig.Driver
	maxConns := simpleDataSource.DbConfig.MaxOpenConnections
	maxIdleConns := simpleDataSource.DbConfig.MaxIdleConnections
	id := simpleDataSource.DbConfig.ConnId
	name := simpleDataSource.DbConfig.ConnName
	maxLifeTime := simpleDataSource.DbConfig.ConnMaxLifetime
	conn := new(DirectConnection)
	simpleDataSource.DirectConnection = conn
	simpleDataSource.DirectConnection.Config.Driver = driver
	simpleDataSource.DirectConnection.Config.ConnectionString = dsn
	simpleDataSource.DirectConnection.Config.MaxConn = maxConns
	simpleDataSource.DirectConnection.Config.MaxIdleConn = maxIdleConns
	simpleDataSource.DirectConnection.Config.ConnectionId = id
	simpleDataSource.DirectConnection.Config.ConnectionName = name
	simpleDataSource.DirectConnection.Config.ConnMaxLifetime = maxLifeTime
	err = simpleDataSource.DirectConnection.InitConnection()
	return
}

// 获取数据库连接
func (simpleDataSource *SimpleDataSource) GetConnection() (*sql.DB, error) {
	if simpleDataSource.DirectConnection.Config.ConnectionState {
		return simpleDataSource.DirectConnection.Connection, nil
	} else {
		return nil, errors.New("connection status error!")
	}
}

// 动态连接池，通过channel传入不同的数据源标识,可切换不同的数据源
type DynamicDataSource struct {
	DbConfigs   *config.DataSourceConfig
	DataSources map[string]*SimpleDataSource
	DbFlag      string
	dbMutex     sync.Mutex
}

func (dynamicDb *DynamicDataSource) CloseAll() {
	for k, v := range dynamicDb.DataSources {
		db, err := v.GetConnection()
		if err != nil {
			log.Logger.WithError(err).Error("close error")
			continue
		}
		err = db.Close()
		if err != nil {
			log.Logger.WithError(err).WithField("db", k).Error("close error")
		}
	}
}

// 初始化多数据源
func (dynamicDb *DynamicDataSource) New(sourceConfig *config.DataSourceConfig) (*DynamicDataSource, error) {
	dataSource := new(DynamicDataSource)
	dataSources := make(map[string]*SimpleDataSource, len(sourceConfig.Database))

	dataSource.DataSources = dataSources
	dataSource.DbConfigs = sourceConfig
	for key, value := range dataSource.DbConfigs.Database {
		simpleDataSource := new(SimpleDataSource)
		simpleDataSource.DbConfig = value
		err := simpleDataSource.InitDataSource()
		if err != nil {
			return nil, err
		}
		dataSource.DataSources[key] = simpleDataSource
	}
	return dataSource, nil
}

// 多个数据源可以通过配置的key获取相应的数据源，切换完之后如需再使用，需要将dbFlag设置回去
func (dynamicDb *DynamicDataSource) GetConnection(dbFlag string) (*sql.DB, error) {
	//dynamicDb.dbMutex.Lock()
	//defer dynamicDb.dbMutex.Unlock()
	if dbFlag == "" {
		dbFlag = "default"
	}
	db := dynamicDb.DataSources[dbFlag]
	if db == nil {
		return nil, errors.New(fmt.Sprint("there is no db =>", dynamicDb.DbFlag))
	}
	return db.GetConnection()
}

// 连接池的Connection
type PooledConnection struct {
	dataSource *DynamicDataSource
}

// 数据库操作绑定列定义struct可实现此接口,实现TableName方法，可自行定义操作
type RowBinder interface {
	GetTableName() string
}

// 获取不同的数据源连接池
func (pooledConnection *PooledConnection) GetConnection(dbFlag string) (*sql.DB, error) {
	return (*pooledConnection.dataSource).GetConnection(dbFlag)
}

func (pooledConnection *PooledConnection) Close() {
	pooledConnection.dataSource.CloseAll()
}

type IDataSourceFactory interface {
	GetDataSource() (*DataSource, error)
}

// SqlSession接口，为数据库通过操作提供统一接口，可通过实现SqlSession提供不同类型的数据库操作
// 这里主要提供了 select / insert / update / delete 接口，使用时直接操作即可
type SqlSession interface {
	// selectAll传入struct 结构体，返回结构体集合,结构体必须带相应的Tag
	// Tag内容有 DbType:数据库类型如VARCHAR/INT/TIMESTAMP等,ColumnName:数据库表列名，如t_table的id列，则定义为ColumnName:"id"
	// 可通过反射机制将结构体赋值，处理之后将结构体slice返回，如有异常抛出
	SelectAll(sql string, binder interface{}, params ...interface{}) ([]interface{}, error)
	// 查询单条数据，返回结构体
	SelectOne(sql string, binder interface{}, params ...interface{}) (interface{}, error)
	// 插入操作，这里定义 binder指针是为了操作后如有像MySQL数据库中获取LastInsertId，获取之后为主键赋值，定义Tag进行操作
	// 如PK:"LastId",当PK定义后则将LastInsertId赋值给PK主键,传入为指针，所以操作赋值要加锁
	// 返回值为插入影响行数
	Insert(sqlStr string, binder interface{}, params ...interface{}) (int64, interface{}, error)
	// Update处理 返回值为插入影响行数 如有异常抛出
	Update(Sql string, params ...interface{}) (int64, error)
	// 删除处理 返回值为插入影响行数 如有异常抛出
	Delete(sql string, params ...interface{}) (int64, error)
	// 获取统一的DbConnection
	GetConnection() (*sql.DB, error)
	// 关闭操作 如有异常抛出
	Close() error
	// Transaction操作, commit 提交自行实现
	Commit() error
	// Transaction操作, rollback 提交自行实现
	Rollback() error
}

// DefaultSqlSession为SqlSession默认实现
type DefaultSqlSession struct {
	// 连接池Connection
	Connection *PooledConnection
	// 连接标识，绑定不同的库
	DbFlag string
}

// 设置当前数据源，默认为default,修改完之后如需要再用，设置回去
func (sqlSession *DefaultSqlSession) SetDbFlag(dbFlag string) {
	if dbFlag != "" {
		sqlSession.DbFlag = dbFlag
	}
}

// 获取当前数据连接
func (sqlSession *DefaultSqlSession) GetConnection() (*sql.DB, error) {
	return sqlSession.Connection.GetConnection(sqlSession.DbFlag)
}

// 关闭数据库连接，这里不建议使用，close后则将所有的数据库连接全部回收
func (sqlSession *DefaultSqlSession) Close() error {
	conn, err := sqlSession.GetConnection()
	if err != nil {
		return err
	}
	return conn.Close()
}

// 未实现，需要开启Transaction
func (sqlSession *DefaultSqlSession) Commit() (err error) {
	return
}

// 未实现，需要开启Transaction
func (sqlSession *DefaultSqlSession) Rollback() (err error) {
	return
}

// 初始化默认的SqlSession,最好使用Prototype设计模式，隔离所有操作，避免多线程使用
func (sqlSession *DefaultSqlSession) New(pooledConnection *PooledConnection) *DefaultSqlSession {
	conn := new(DefaultSqlSession)
	conn.Connection = pooledConnection
	return conn
}

// 通过反射进行类型转换并赋值
// 将数据库字段类型进行枚举，这里clob blob 直接使用string来简单处理，后期添加以[]byte进行处理
// 传入数据库字段当前值为string类型，查询出后类型统一为sql.RowBytes即[]byte 难以处理，统一强转成string后进行处理
func (sqlSession *DefaultSqlSession) ConvertAssign(fieldValue *reflect.Value, src string, dbType string) (err error) {
	switch dbType {
	case "VARCHAR", "BLOB", "CLOB",
		"TINYBLOB", "TINYTEXT", "TEXT",
		"MEDIUMBLOB", "MEDIUMTEXT", "LONGBLOB", "LONGTEXT", "CHAR":
		fieldValue.SetString(src)
		break
	case "INT", "UINT", "TINYINT", "SMALLINT", "BIGINT":
		var tmp int64
		tmp, err := strconv.ParseInt(src, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetInt(tmp)
		break
	case "FLOAT", "DOUBLE", "DECIMAL":
		var tmp float64
		tmp, err = strconv.ParseFloat(src, 64)
		if err != nil {
			return
		}
		fieldValue.SetFloat(tmp)
		break
	case "TIMESTAMP", "DATETIME", "TIME", "DATE":
		var tmp time.Time
		// 此处必须为这个时间才可以进行转换
		tmp, err = time.Parse("2006-01-02 15:04:05", src)
		fieldValue.Set(reflect.ValueOf(tmp))
		break
	default:
		fieldValue.SetString(src)
		break
	}
	return
}

// 绑定Struct结构体，每个结构体统一将 数据库字段类型和数据库字段名称进行Tag处理，在这里通过反射进行统一处理
// 这里是以行为单位进行处理,统一数据转换为string后再进行下一步的操作,处理过程中可能为异常，直接向上抛,
func (sqlSession *DefaultSqlSession) BindingStruct(bind *interface{}, rowData map[string]string) (bindData interface{}, err error) {
	//bindType := reflect.TypeOf(*bind).Elem()
	bindType := reflect.TypeOf(*bind)
	bindValue := reflect.New(bindType).Elem()
	for i := 0; i < bindType.NumField(); i++ {
		dbType := bindType.Field(i).Tag.Get("DbType")
		columnName := bindType.Field(i).Tag.Get("ColumnName")
		propertyName := bindType.Field(i).Name
		if dbType == "" || columnName == "" {
			continue
		}
		propertyValue := rowData[columnName]
		fieldValue := bindValue.FieldByName(propertyName)
		if propertyValue == "NULL" {
			//fieldValue.SetString(propertyValue)
			continue
		}
		err := sqlSession.ConvertAssign(&fieldValue, propertyValue, dbType)
		if err != nil {
			return nil, err
		}
	}
	bindData = bindValue.Interface()
	return
}

// 可以理解为SelectAll的处理，如为SelectOne刚循环一次后结束即可
func (sqlSession *DefaultSqlSession) SelectAllData(rows *sql.Rows, bind interface{}) (vDatas []interface{}, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	// 无数据
	if len(columns) == 0 {
		return make([]interface{}, 0), nil
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	//var vDatas []interface{}
	for rows.Next() {
		vData := make(map[string]string)
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		if err != nil {
			return
		}
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				// 处理正常流程 先将列值赋值为string，然后再进行统一转换
				value = string(col)
			}
			vData[columns[i]] = value
		}
		bindData, err := sqlSession.BindingStruct(&bind, vData)
		if err != nil {
			return nil, err
		}
		vDatas = append(vDatas, bindData)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

// SelectOne的实现，binder是将数据转换为struct格式;sqlString为传入sql语句
func (sqlSession *DefaultSqlSession) SelectOne(sqlString string, binder interface{}, params ...interface{}) (row interface{}, err error) {
	rows, err := sqlSession.SelectAll(sqlString, binder, params...)
	if err != nil {
		return
	}
	if len(rows) > 0 {
		return rows[0], nil
	}
	return
}

// SelectOne的实现，binder是将数据转换为struct格式;sqlString为传入sql语句 params此处使用Prepared模式，直接传入即可
func (sqlSession *DefaultSqlSession) SelectAll(sqlInput string, binder interface{}, params ...interface{}) (dataRows []interface{}, err error) {
	connection, err := sqlSession.GetConnection()
	if err != nil {
		return
	}
	//connection.Ping()
	var rows *sql.Rows
	stmt, err := connection.Prepare(sqlInput)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err = stmt.Query(params...)
	if err != nil {
		return
	}
	defer rows.Close()
	return sqlSession.SelectAllData(rows, binder)
}

// 删除操作，返回影响行数
func (sqlSession *DefaultSqlSession) Delete(sqlStr string, params ...interface{}) (affectedRows int64, err error) {
	connection, err := sqlSession.GetConnection()
	if err != nil {
		return 0, err
	}
	result, err := connection.Exec(sqlStr, params...)
	if err != nil {
		return
	}
	affectedRows, err = result.RowsAffected()
	return
}

// 更新操作，返回影响行数
func (sqlSession *DefaultSqlSession) Update(sqlStr string, params ...interface{}) (affectedRows int64, err error) {
	connection, err := sqlSession.GetConnection()
	if err != nil {
		return 0, err
	}
	result, err := connection.Exec(sqlStr, params...)
	if err != nil {
		return
	}
	affectedRows, err = result.RowsAffected()
	return
}

// 插入操作，如需要要将Id返回，传入相应的struct结构体，返回影响行数和带有LastInsertId的interface
func (sqlSession *DefaultSqlSession) Insert(sqlStr string, binder interface{}, params ...interface{}) (affectiveRows int64, retValue interface{}, err error) {
	connection, err := sqlSession.GetConnection()
	if err != nil {
		return 0, nil, err
	}
	result, err := connection.Exec(sqlStr, params...)
	if err != nil {
		return
	}
	//result.LastInsertId()
	// 获取影响行数
	if affectiveRows, err = result.RowsAffected(); err != nil {
		return
	}
	// 不需要绑定
	if binder == nil {
		return
	}
	binderType := reflect.TypeOf(binder)
	binderValue := reflect.New(binderType).Elem()
	for i := 0; i < binderType.NumField(); i++ {
		// 如有找到 Tag里面有 PK:为"LastId" 说明需要将id以lastInsertId进行赋值
		pkName := binderType.Field(i).Tag.Get("PK")
		if pkName != "" && pkName == "LastId" {
			lastId, err := result.LastInsertId()
			if err != nil {
				return 0, nil, err
			}
			// 通过反射赋值后直接跳出循环
			name := binderType.Field(i).Name
			fieldValue := binderValue.FieldByName(name)
			fieldValue.SetInt(lastId)
			retValue = binderValue.Interface()
			break
		}
	}
	return
}
