package logger

import (
	"dapp_cess_client/conf"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

var (
	OutPutLogger *zap.Logger
)

func InitLogger() {
	if len(conf.ClientConf.BoardInfo.BoardPath) == 0 {
		conf.ClientConf.BoardInfo.BoardPath = conf.Board_Path_D
	}
	_, err := os.Stat(conf.ClientConf.BoardInfo.BoardPath)
	if err != nil {
		err = os.MkdirAll(conf.ClientConf.BoardInfo.BoardPath, os.ModePerm)
		if err != nil {
			fmt.Printf("\x1b[%dm[err]\x1b[0m Create '%v' file output.log error\n", 41, conf.ClientConf.BoardInfo.BoardPath)
			os.Exit(conf.Exit_ConfErr)
		}
	}
	initOutPutLogger()
}

// output log
func initOutPutLogger() {
	outputlogpath := filepath.Join(conf.ClientConf.BoardInfo.BoardPath, "output.log")
	hook := lumberjack.Logger{
		Filename:   outputlogpath,
		MaxSize:    10,
		MaxAge:     360,
		MaxBackups: 0,
		LocalTime:  true,
		Compress:   true,
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
		TimeKey:    "time",
		//CallerKey:    "file",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime:  formatEncodeTime,
		//EncodeCaller: zapcore.ShortCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)
	caller := zap.AddCaller()
	development := zap.Development()
	OutPutLogger = zap.New(core, caller, development)
}

func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}
