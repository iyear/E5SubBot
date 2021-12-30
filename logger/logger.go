package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// 日志切割设置
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/latest.log", // 日志文件位置
		MaxSize:    1,                  // 日志文件最大大小(MB)
		MaxBackups: 5,                  // 保留旧文件最大数量
		MaxAge:     30,                 // 保留旧文件最长天数
		Compress:   true,               // 是否压缩旧文件
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}

// 编码器
func getEncoder() zapcore.Encoder {
	// 使用默认的JSON编码
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// Init 初始化Logger
func Init() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	zap.ReplaceGlobals(zap.New(core, zap.AddCaller()))
}
