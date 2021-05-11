package model

import "time"

// spider info
type SpiderInfo struct {
	Id            int64     `DbType:"BIGINT" ColumnName:"id" json:"id"`                          // id
	AppId         int64     `DbType:"BIGINT" ColumnName:"app_id" json:"app_id"`                  // app id
	Name          string    `DbType:"VARCHAR" ColumnName:"name" json:"name"`                     // 名称
	Description   string    `DbType:"VARCHAR" ColumnName:"description" json:"description"`       // 描述
	AppType       int64     `DbType:"BIGINT" ColumnName:"app_type" json:"app_type"`              // app类型
	SourceIp      string    `DbType:"VARCHAR" ColumnName:"source_ip" json:"source_ip"`           // 源ip
	QueueName     string    `DbType:"VARCHAR" ColumnName:"queue_name" json:"queue_name"`         // 队列名称
	Rule          string    `DbType:"VARCHAR" ColumnName:"rule" json:"rule"`                     // rule规则 ajax/html
	Status        int8      `DbType:"TINYINT" ColumnName:"status" json:"status"`                 // 状态
	NeedDownLoad  int8      `DbType:"TINYINT" ColumnName:"need_download" json:"need_download"`   // 是否需要下载
	NeedTranslate int8      `DbType:"TINYINT" ColumnName:"need_translate" json:"need_translate"` // 是否需要翻译
	NeedTag       int8      `DbType:"TINYINT" ColumnName:"need_tag" json:"need_tag"`             // 是否需要打标签
	NeedProxy     int8      `DbType:"TINYINT" ColumnName:"need_proxy" json:"need_proxy"`         // 是否需要代理
	CreateTime    time.Time `DbType:"DATETIME" ColumnName:"create_time" json:"create_time"`      // 创建时间
}
