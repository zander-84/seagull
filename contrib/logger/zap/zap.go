package zap

import (
	"context"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/tool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

type zLog struct {
	engine  *zap.Logger
	conf    Conf
	writers []io.Writer
}

func NewZapLog(conf Conf, writers []io.Writer) (log contract.Logger, cancel func(ctx context.Context) error, err error) {
	z := new(zLog)
	z.init(conf, writers)

	return z, z.Stop, z.start()
}

func (z *zLog) init(conf Conf, writers []io.Writer) {
	z.conf = conf.SetDefault()
	z.writers = writers
}
func (z *zLog) start() error {

	newCore := make([]zapcore.Core, 0)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(i time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Now().Format(tool.TimeFormat))
	}
	encoderConfig.MessageKey = ""

	logLevel := zap.DebugLevel
	switch z.conf.Level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.InfoLevel
	}

	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	})

	//____ 控制台输出
	if z.conf.ConsoleHook.Enable {
		console := zapcore.Lock(os.Stdout)
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		newCore = append(newCore,
			zapcore.NewCore(consoleEncoder, console, priority),
		)
	}

	//____ 文件写入
	if z.conf.FileHook.Enable {
		if err := openOrCreateWithAction(z.conf.FileHook.Path, z.conf.Name+".log", func(f *os.File) {
			//log.SetOutput(f)
			//log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
			//log.Println("starting log...")
		}); err != nil {
			return err
		}

		fileHook := lumberjack.Logger{
			Filename:   getPrefixPath(z.conf.FileHook.Path, z.conf.Name+".log"), // 日志文件路径
			MaxSize:    z.conf.FileHook.MaxSize,                                 // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: z.conf.FileHook.MaxBackups,                              // 日志文件最多保存多少个备份
			MaxAge:     z.conf.FileHook.MaxAge,                                  // 文件最多保存多少天
			Compress:   false,                                                   // 是否压缩
		}
		fileWriter := zapcore.AddSync(&fileHook)
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		newCore = append(newCore,
			zapcore.NewCore(jsonEncoder, fileWriter, priority),
		)
	}

	if len(z.writers) > 0 {
		for _, writer := range z.writers {
			encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
			newCore = append(newCore, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), priority))
		}
	}

	if z.conf.AddCaller {
		z.engine = zap.New(zapcore.NewTee(newCore...), zap.AddCaller(), zap.AddCallerSkip(2))
	} else {
		z.engine = zap.New(zapcore.NewTee(newCore...))
	}

	return nil
}

func (z *zLog) Stop(ctx context.Context) error {
	return tool.ExitWithContext(ctx, z.engine.Sync)
}

func (z *zLog) Engine() *zap.Logger {
	return z.engine
}
