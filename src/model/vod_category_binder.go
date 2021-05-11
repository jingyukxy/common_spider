package model

// category与vod源的绑定关系
type VodCategoryBinder struct {
	Id          int64  `json:"id" DbType:"BIGINT" ColumnName:"id"`
	SpiderId    int64  `json:"spider_id" DbType:"BIGINT" ColumnName:"spider_id"`
	SourceType  string `json:"source_type" DbType:"VARCHAR" ColumnName:"source_type"`
	CategoryIds string `json:"category_ids" DbType:"VARCHAR" ColumnName:"category_ids"`
}
