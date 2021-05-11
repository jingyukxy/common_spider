package logs

import (
	"awesomeProject/src/config"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var Logger *logrus.Logger

// 默认字段hook
type DefaultFieldHook struct {
}

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, "/bin")
	if index != -1 {
		return path[:index]
	}
	return path
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["appId"] = 1
	entry.Data["appName"] = "spider"
	return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func NewLfs(logfile string, logFilePath string, maxCount int) (lfsHook *lfshook.LfsHook, err error) {
	baseLogPath := path.Join(logFilePath, logfile)
	writer, err := rotatelogs.New(
		baseLogPath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithRotationCount(maxCount),
	)
	if err != nil {
		Logger.Error("config file log error %v", errors.WithStack(err))
		panic(err)
	}
	errWriter, err := rotatelogs.New(
		baseLogPath+".wf_%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath+".wf"),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithRotationCount(maxCount),
	)
	if err != nil {
		Logger.Error("err config file log error %v", errors.WithStack(err))
		panic(err)
	}
	lfsHook = lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.ErrorLevel: errWriter,
		logrus.InfoLevel:  writer,
		logrus.FatalLevel: errWriter,
		logrus.PanicLevel: errWriter,
		logrus.TraceLevel: writer,
		logrus.WarnLevel:  writer,
	}, Logger.Formatter)
	return
}

func InitLogger(config *config.LogConfig) {
	Logger = logrus.New()
	Logger.ReportCaller = true
	Logger.Out = os.Stdout
	//Logger.Formatter = &logrus.TextFormatter{
	//	FullTimestamp:   true,
	//	TimestampFormat: "2006-01-02 15:04:05.000000000",
	//}
	Logger.Formatter = &nested.Formatter{
		HideKeys:        true,
		NoColors:        true,
		TimestampFormat: time.RFC3339Nano,
	}
	//Logger.AddHook(&DefaultFieldHook{})
	filePath := GetAppPath() + "/logs"
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(filePath, 0755)
	}
	fileName := config.FileName
	lfsHook, err := NewLfs(fileName, filePath, config.MaxCount)
	if err != nil {
		Logger.Panic("%v", err)
	}
	Logger.AddHook(lfsHook)
}
