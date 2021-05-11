package dao

import (
	"awesomeProject/src/db"
	log "awesomeProject/src/logs"
	"awesomeProject/src/message"
	"awesomeProject/src/model"
	"awesomeProject/src/module"
	"awesomeProject/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	movieCategoryBinderKey = "movie_category_binder"
)

type MovieDao struct {
	sqlSession *db.DefaultSqlSession
}

func NewMovieDao() *MovieDao {
	md := &MovieDao{
		sqlSession: db.NewSqlSession(),
	}
	// 先切库
	md.sqlSession.SetDbFlag("vod")
	return md
}

// 查询ids
func (movieDao *MovieDao) GetCategoryIdsBySourceType(sourceType string) string {
	for i := 0; i < 3; i++ {
		ret, err := module.GetCacheProviderInstance().HGet(movieCategoryBinderKey, utils.Md5WithString(sourceType))
		if err == redis.Nil {
			// 查询为空 判断redis中是否已经存过了
			size, err := module.GetCacheProviderInstance().HLen(movieCategoryBinderKey)
			if err != nil {
				log.Logger.WithError(err).Error("Redis获取长度失败")
				return ""
			}
			if size == 0 {
				// redis中没有，数据加载后再去查询 重试三遍之后处理不成功则退出
				movieDao.PreloadCategory()
				continue
			}
		} else if err != nil {
			log.Logger.WithError(err).Error("Redis查询出错!")
			return ""
		} else {
			// 查询到了
			binder := model.VodCategoryBinder{}
			err = json.Unmarshal([]byte(ret), &binder)
			if err != nil {
				log.Logger.WithError(err).Error("binder unmarshal failed!")
				return ""
			}
			return binder.CategoryIds
		}
	}
	return ""
}

// 初始化数据
func (movieDao *MovieDao) PreloadCategory() {
	binders, err := movieDao.GetVodCategoryBinders()
	if err != nil {
		log.Logger.WithError(err).Error("查询绑定数据失败")
	} else {
		if binders != nil {
			key := "movie_category_binder"
			for _, binder := range binders {
				data, err := json.Marshal(binder)
				if err != nil {
					log.Logger.WithError(err).Error("Json marshal failed!")
					continue
				}
				err = module.GetCacheProviderInstance().HSet(key, utils.Md5WithString(binder.SourceType), data)
				if err != nil {
					log.Logger.WithError(err).Error("HSet cache error!!")
				}
			}
		}
	}
}

// 获取绑定信息
func (movieDao *MovieDao) GetVodCategoryBinders() ([]model.VodCategoryBinder, error) {
	vb := model.VodCategoryBinder{}
	querySql := "select * from t_vod_category_binder"
	rows, err := movieDao.sqlSession.SelectAll(querySql, vb)
	log.Logger.WithField("sql", querySql).Info("Query vod category")
	if err != nil {
		return nil, err
	}
	if len(rows) > 0 {
		result := make([]model.VodCategoryBinder, len(rows))
		for i := range rows {
			result[i] = rows[i].(model.VodCategoryBinder)
		}
		return result, nil
	}
	return nil, nil
}

// 通过名称查询视频信息
func (movieDao *MovieDao) SelectMoviesByName(name string) ([]message.MovieMessage, error) {
	msg := message.MovieMessage{}
	querySql := "select * from t_vod where name = ?"
	log.Logger.WithFields(logrus.Fields{"sql": querySql, "params": name}).Info("")
	rows, err := movieDao.sqlSession.SelectAll(querySql, msg, name)
	if err != nil {
		return nil, err
	}
	msgs := make([]message.MovieMessage, len(rows))
	for i := range rows {
		msgs[i] = rows[i].(message.MovieMessage)
	}
	return msgs, nil
}

// 获取所有视频播放内容
func (movieDao *MovieDao) SelectMovieContents(vodId int64) ([]message.PlayContent, error) {
	contentSql := "select * from t_vod_content where vod_id = ?"
	log.Logger.WithFields(logrus.Fields{"sql": contentSql, "params": vodId}).Info("Query movie contents")
	contentBinder := message.PlayContent{}
	rows, err := movieDao.sqlSession.SelectAll(contentSql, contentBinder, vodId)
	if err != nil {
		return nil, err
	}
	contents := make([]message.PlayContent, len(rows))
	for _, row := range rows {
		contents = append(contents, row.(message.PlayContent))
	}
	return contents, nil
}

// 更新视频信息 只更新最新播放集数和是否完结
func (movieDao *MovieDao) UpdateMovie(id int64, playNum int, isEnd int8) (bool, error) {
	updateSql := "update t_vod set new_episode = ?,is_end = ? where id = ?"
	log.Logger.WithFields(logrus.Fields{
		"sql":     updateSql,
		"id":      id,
		"playNum": playNum,
		"isEnd":   isEnd}).Info("Update movie start")
	affectRow, err := movieDao.sqlSession.Update(updateSql, playNum, isEnd, id)
	if err != nil {
		return false, err
	}
	return affectRow > 0, nil
}

// 通过名称查询演职人员信息
func (movieDao *MovieDao) SelectActorByName(name string) (message.Actor, error) {
	querySql := "select * from t_vod_staff where name = ?"
	log.Logger.WithFields(logrus.Fields{
		"sql":  querySql,
		"name": name,
	}).Info("Query actor name")
	actor := message.Actor{}
	row, err := movieDao.sqlSession.SelectOne(querySql, actor, name)
	if err != nil {
		return actor, err
	}
	if row != nil {
		return row.(message.Actor), nil
	} else {
		return message.Actor{}, nil
	}
}

// 更新演职人员信息
func (movieDao *MovieDao) SaveStaff(actor message.Actor) (int64, error) {
	insertSql := "insert t_vod_staff (name, staff_type, gender, image, height, country, description, update_time) values(?,?,?,?,?,?,?,?)"
	name := actor.Name
	log.Logger.WithFields(logrus.Fields{"sql": insertSql, "name": name}).Info("Insert staff start")
	staffType := actor.StaffType
	gender := 1
	image := actor.Img
	height := 0
	country := actor.Area
	description := actor.Description
	updateTime := time.Now().Format("2006-01-02 15:04:05")
	_, binder, err := movieDao.sqlSession.Insert(insertSql, actor, name, staffType, gender, image, height, country, description, updateTime)
	if err != nil {
		return 0, err
	}
	staffId := binder.(message.Actor).Id
	return staffId, err
}

// 更新视频播放信息
func (movieDao *MovieDao) SaveContent(id int64, content message.PlayContent) error {
	insertSql := "insert into t_vod_content(title, play_type, play_url, vod_id) values(?,?,?,?)"
	log.Logger.WithFields(logrus.Fields{"sql": insertSql, "name": content.PlayUrl}).Info("Insert staff start")
	title := content.PlayName
	playUrl := content.PlayUrl
	playNum := content.PlayNum
	playType := content.PlayType
	if title == "" {
		title = fmt.Sprintf("第%v集", playNum)
	}
	_, _, err := movieDao.sqlSession.Insert(insertSql, content, title, playType, playUrl, id)
	return err
}

// 每存储一个演职人员，还需要把关联信息存储
func (movieDao *MovieDao) SaveVodStaffAssociate(vodId int64, staffId int64) error {
	insertSql := "insert into t_vod_staff_associate  (vod_id, staff_id) values (?,?)"
	log.Logger.WithFields(logrus.Fields{"sql": insertSql, "vodId": vodId}).Info("Insert vod staff start")
	_, _, err := movieDao.sqlSession.Insert(insertSql, nil, vodId, staffId)
	return err
}

// 保存movie信息
func (movieDao *MovieDao) SaveMovie(spiderInfo *model.SpiderInfo, msg message.MovieMessage) error {
	// 先查询是视频已经存在
	movies, err := movieDao.SelectMoviesByName(msg.Name)
	if err != nil {
		return err
	}
	// 数据库中已经有相应的视频
	if len(movies) > 0 {
		// 如果有多个默认取第一个
		movie := movies[0]
		contents, err := movieDao.SelectMovieContents(movie.Id)
		if err != nil {
			return err
		}
		playNums := len(msg.PlayUrls)
		currentNums := len(contents)
		// 如视频地址大于现在的播放地址，则需要更新，默认是排序的, 否则不更新
		if playNums > currentNums {
			// 先保存内容
			contents := msg.PlayUrls[currentNums:]
			for _, content := range contents {
				err := movieDao.SaveContent(movie.Id, content)
				if err != nil {
					// 保存出错，记录一下
					return err
				}
			}
			// 再更新视频信息
			isEnd := 0
			ret, err := strconv.ParseInt(movie.TotalEpisodes, 0, 64)
			// 转换出错
			if err != nil {
				return err
			}
			if playNums >= int(ret) {
				isEnd = 1
			}
			_, err = movieDao.UpdateMovie(movie.Id, playNums, int8(isEnd))
			// 出错还需要处理
		}
	} else {
		// 如还未有数据，则需要增加数据
		id, err := movieDao.InsertMovie(msg)
		// 新增失败
		if err != nil {
			return err
		}
		err = movieDao.saveVodCategory(msg.MovieType, id, spiderInfo.AppId)
		if err != nil {
			return err
		}
		// 新增播放内容
		for _, content := range msg.PlayUrls {
			movieDao.SaveContent(id, content)
		}
		// 保存演职人员信息
		movieDao.SaveStaffByType(id, msg.Director, 1)
		// 保存演员
		for _, actor := range msg.Actors {
			err := movieDao.SaveStaffByType(id, actor, 2)
			if err != nil {
				//保存出错打印日志继续处理
				log.Logger.WithField("msg", msg).WithError(err).Error("保存角色信息失败")
				continue
			}
		}
	}
	return nil
}

// 存储演职人员信息
func (movieDao *MovieDao) SaveStaffByType(vodId int64, actor message.Actor, staffType int8) error {
	if actor.Name == "" {
		return nil
	}
	director, err := movieDao.SelectActorByName(actor.Name)
	// 查询出错
	if err != nil {
		return err
	}
	// 导演不存在存储
	if director.Id == 0 {
		actor.StaffType = 1
		staffId, err := movieDao.SaveStaff(actor)
		if err != nil {
			return err
		}
		// 处理完之后还需要将演职人员关系表信息处理一下
		return movieDao.SaveVodStaffAssociate(vodId, staffId)
	}
	// 已经存在需要打印日志
	return nil
}

// 存储关联表
func (movieDao *MovieDao) saveVodCategory(sourceType string, vodId int64, appId int64) error {
	ids := movieDao.GetCategoryIdsBySourceType(sourceType)
	if ids != "" {
		idsLimit := strings.Split(ids, ",")
		for _, id := range idsLimit {
			insertSql := "insert into t_vod_category (vod_id, category_id,app_id) values(?,?,?)"
			rid, err := strconv.ParseInt(id, 0, 64)
			_, _, err = movieDao.sqlSession.Insert(insertSql, nil, vodId, rid, appId)
			if err != nil {
				log.Logger.WithError(err).Error("存储绑定数据失败,返回")
				return err
			}
		}
	} else {
		log.Logger.Error("此数据没有绑定来源")
	}
	return nil
}

// 先保存视频，然后保存演职人员信息和视频播放地址
func (movieDao *MovieDao) InsertMovie(movieMessage message.MovieMessage) (id int64, err error) {
	// 存储之前看一下是否绑定类别，没有绑定直接返回
	ids := movieDao.GetCategoryIdsBySourceType(movieMessage.MovieType)
	if ids == "" {
		return 0, errors.New("没有绑定数据")
	}
	movieInsert := "INSERT INTO `t_vod` (`name`,`category_id`,`title`,`keywords`,`image`,`area`,`language`,`show_time`," +
		"`total_episode`,`new_episode`,`is_end`,`show_year`,`add_time`,`hit_day_count`,`hit_week_count`," +
		"`hit_month_count`,`stars`,`status`,`editor`,`letter`,`score`,`douban_score`,`is_film`,`vod_length`," +
		"`description`,`update_time`)VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	name := movieMessage.Name
	title := movieMessage.Title
	keywords := movieMessage.Keywords
	image := movieMessage.Img
	area := movieMessage.Area
	upTime := movieMessage.UpdateTime
	var showTime time.Time
	if upTime != "" {
		if len(upTime) == 4 {
			showTime, err = time.Parse("2006", upTime)
			if err != nil {
				log.Logger.WithField("showtime", movieMessage.UpdateTime).Info("update time is error!")
			}
		} else if len(upTime) >= 10 {
			showTime, err = time.Parse("2006-01-02", upTime)
			if err != nil {
				log.Logger.WithField("showtime", movieMessage.UpdateTime).Info("update time is error!")
			}
		} else {
			showTime = time.Now()
		}
	}

	language := movieMessage.Language
	totalEpisode := utils.If(movieMessage.TotalEpisodes != "", movieMessage.TotalEpisodes, "1").(string)
	newEpisodes := len(movieMessage.PlayUrls)
	isEnd := 0
	if totalEpisode == string(newEpisodes) {
		isEnd = 1
	}
	hitDayCount := 0
	hitWeekCount := 0
	hitMonthCount := 0
	stars := 0
	// 状态统一为可用
	status := 1
	showYear := utils.If(movieMessage.Year != "", movieMessage.Year, "2019").(string)
	editor := "spider_editor"
	score := utils.If(movieMessage.Score != "", movieMessage.Score, "0").(string)
	doubanScore := utils.If(movieMessage.DoubanScore != "", movieMessage.Score, "0").(string)
	description := movieMessage.Description
	updateTime := time.Now().Format("2006-01-02 15:04:05")
	addTime := time.Now().Format("2006-01-02 15:04:05")
	// 需要加判断是否是电影
	isFilm := 1
	categoryId := 1
	// 需要判断首字母
	letter := utils.GetFirstLetterOfPinYin(name)
	// 电视长度需要加判断
	vodLength := 0
	//vodLength := movieMessage.l
	_, retValue, err := movieDao.sqlSession.Insert(movieInsert, movieMessage,
		name, categoryId, title, keywords, image, area, language, showTime, totalEpisode,
		newEpisodes, isEnd, showYear, addTime, hitDayCount, hitWeekCount, hitMonthCount,
		stars, status, editor, letter, score, doubanScore, isFilm, vodLength, description, updateTime)
	if err != nil {
		return 0, err
	}
	retMsg := retValue.(message.MovieMessage)
	return retMsg.Id, err
}
