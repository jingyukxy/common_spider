package dao

import (
	"awesomeProject/src/db"
	log "awesomeProject/src/logs"
	"awesomeProject/src/model"
)

type SpiderDao struct {
	sqlSession *db.DefaultSqlSession
}

func NewSpiderDao() *SpiderDao {
	return &SpiderDao{
		sqlSession: db.NewSqlSession(),
	}
}

func (spiderDao *SpiderDao) GetSpiderInfoByCode(code int) (spiderInfo model.SpiderInfo, err error) {
	info := model.SpiderInfo{}
	row, err := spiderDao.sqlSession.SelectOne("select * from t_spider where source_code = ? and status =1", info, code)
	if err != nil {
		log.Logger.WithError(err).Error("get spider info error!!")
		return
	}
	return row.(model.SpiderInfo), nil
}
