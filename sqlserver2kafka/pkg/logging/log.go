package logging

import (
	log2 "log"
	"time"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"

)
//var log *zap.SugaredLogger
var log *zap.Logger
type Field = zapcore.Field
// 时间格式
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
}
// zap日志配置信息
func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// 初始化日志
type level struct {
	RunMode string `yaml:"runMode"`
}
func Init() () {
	//loglevel := setting.ServerSetting.RunMode
	// 读取配置文件
	data,err := ioutil.ReadFile("sqlserver2kafka/conf/app.yaml")
	if err != nil {
		log2.Panic("parse conf/app.yaml err:",err)
	}
	// 初始化结构体，解析配置文件
	loglevel := level{}
	yaml.Unmarshal(data,&loglevel)
	// 根据配置文件赋值
	var level zapcore.Level
	switch loglevel.RunMode {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(newEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		level,
	)
	// 要改json格式的话，把NewConsoleEncoder 改成NewJSONEncoder
	// 改成文件输出的话修改zapcore.AddSync(os.Stdout))
	// 连接 https://blog.csdn.net/NUCEMLS/article/details/86534444
	//logger := zap.New(core, zap.AddCaller(),zap.AddCallerSkip(1)) // AddCallerSkip(1) 跳过封装函数的调用
	//log = logger.Sugar() // 用这个，可以使用%s 格式化输出，但是再用zap.string会打印出多余的内容
	// 不用格式化输出，这个效率更高
	log = zap.New(core, zap.AddCaller(),zap.AddCallerSkip(1)) // AddCallerSkip(1) 跳过封装函数的调用
}
func Debug(msg string, fields ...Field) {
	log.Debug(msg, fields ...)
}
func Info(msg string, fields ...Field)  {
	log.Info(msg, fields ...) // fiels ... 被打散传入
}
func Warn(msg string, fields ...Field)  {
	log.Warn(msg, fields ...)
}
func Error(msg string, fields ...Field)  {
	log.Error(msg, fields ...)
}
func DPanic(msg string, fields ...Field)  {
	log.DPanic(msg, fields ...)
}
func Panic(msg string, fields ...Field)  {
	log.Panic(msg, fields ...)
}
func Fatal(msg string, fields ...Field)  {
	log.Fatal(msg, fields ...)
}

// 使用printf 格式化
//func Debug(args ...interface{})  {
//	log.Debug(args)
//}
//func Debugf(template string, args ...interface{})  {
//	log.Debugf(template, args...)
//}
//func Info(args ...interface{}) {
//	log.Info(args...)
//}
//
//func Infof(template string, args ...interface{}) {
//	log.Infof(template, args...)
//}
//
//func Warn(args ...interface{}) {
//	log.Warn(args...)
//}
//
//func Warnf(template string, args ...interface{}) {
//	log.Warnf(template, args...)
//}
//
//func Error(args ... interface{}) {
//	log.Error(args...)
//}
//
//func Errorf(template string, args ...interface{}) {
//	log.Errorf(template, args...)
//}
//
//func Panic(args ...interface{}) {
//	log.Panic(args...)
//}
//
//func Panicf(template string, args ...interface{}) {
//	log.Panicf(template, args...)
//}