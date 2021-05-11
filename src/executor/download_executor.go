package executor

import (
	"awesomeProject/src/config"
	log "awesomeProject/src/logs"
	"awesomeProject/src/message"
	"awesomeProject/src/model"
	"awesomeProject/src/utils"
	"fmt"
	"os"
	"reflect"
	"time"
)

type DownloadExecutor struct {
}

// 处理下载
func (download *DownloadExecutor) Execute(info *model.SpiderInfo, msg message.IMessage) message.IMessage {
	// 如果不需要下载直接返回
	if info.NeedDownLoad != 1 {
		log.Logger.Info("此请求不需要下载,", info.Name)
		return msg
	}
	// 判断下载字段是否有值
	msgValue := reflect.ValueOf(msg).Elem()
	fieldName := msg.GetDownField()
	if fieldName == "" {
		log.Logger.Info("没有设置下载字段，不处理下载")
		return msg
	}
	field := msgValue.FieldByName(fieldName)
	if field.Kind() == reflect.Invalid || field.String() == "" {
		log.Logger.WithField("field", fieldName).Info("下载字段错误或下载字段为空，请重新设置")
		return msg
	}
	downloadConfig := config.GetConfigInstance().GetDownloadConfig()
	if downloadConfig == nil {
		log.Logger.Error("没有下载配置，消息不处理")
		return msg
	}
	date := time.Now().Format("2006010215")
	downloadPath := downloadConfig.DownloadPath
	suffix := downloadConfig.DownloadSuffix
	realPath := fmt.Sprintf("%s/%s", downloadPath, date)
	exists, err := utils.PathExists(realPath)
	if err != nil {
		log.Logger.WithError(err).Error("系统文件读取失败，直接返回")
		return msg
	}
	// 文件夹不存在则创建
	if !exists {
		os.Mkdir(realPath, 0744)
	}
	downImg := field.String()
	resultName, err := utils.DownloadFile(realPath, suffix, downImg)
	if err != nil {
		log.Logger.WithError(err).Error("下载文件失败")
		return msg
	}
	downloadUrl := fmt.Sprintf("%s%s/%s.%s", downloadConfig.PrefixUrl, date, resultName, suffix)
	log.Logger.Info("down file finished the file name is ", downloadUrl)
	field.SetString(downloadUrl)
	return msg
}
