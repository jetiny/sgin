package base

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/jetiny/sgin/common"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"xorm.io/core"
)

var gxormLogger *xorm.SimpleLogger

func initLogger() error {
	logFilePath := gEnvLogDir.String()
	logFileName := gEnvLogFileName.String()
	err := os.MkdirAll(logFilePath, 0777)
	if err != nil {
		return err
	}
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	// 创建文件
	logFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	logFile.Close()

	// 写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	writer := io.MultiWriter(os.Stdout, src)

	logger := common.Logger

	//设置日志级别
	level := logrus.Level(gEnvLogLevel.Int())
	logger.SetLevel(level)
	//设置输出
	logger.Out = writer
	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",
		// rotatelogs.ForceNewFile(),
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(time.Duration(gEnvLogExpired.Int())*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(time.Duration(gEnvLogCutDays.Int())*time.Hour),
	)
	if err != nil {
		return err
	}
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
		logrus.TraceLevel: logWriter,
	}
	logger.AddHook(lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
	gxormLogger = xorm.NewSimpleLogger3(writer, xorm.DEFAULT_LOG_PREFIX, xorm.DEFAULT_LOG_FLAG, xormLogLevel(level))
	gxormLogger.ShowSQL(true)
	gin.DefaultWriter = writer
	return nil
}

func xormLogLevel(level logrus.Level) core.LogLevel {
	switch level {
	case logrus.DebugLevel:
		return core.LOG_DEBUG
	case logrus.ErrorLevel:
		return core.LOG_ERR
	case logrus.FatalLevel:
		return core.LOG_ERR
	case logrus.InfoLevel:
		return core.LOG_INFO
	case logrus.WarnLevel:
		return core.LOG_WARNING
	case logrus.TraceLevel:
		return core.LOG_DEBUG
	default:
		return core.LOG_DEBUG
	}
}
