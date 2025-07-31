package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	WithRequestID(ctx context.Context) Logger
}

type zapLogger struct {
	l *zap.Logger
}

func NewLogger() Logger {
	mode := os.Getenv("GIN_MODE")
	isProd := mode == "release" || mode == "production"

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderCfg.TimeKey = "time"
	encoderCfg.MessageKey = "msg"
	encoderCfg.CallerKey = "caller"

	var encoder zapcore.Encoder
	if isProd {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	level := zapcore.InfoLevel
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		if parsed, err := zapcore.ParseLevel(v); err == nil {
			level = parsed
		}
	}

	logFile := "./logs/app.log"
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})

	var ws zapcore.WriteSyncer
	if isProd {
		ws = zapcore.NewMultiWriteSyncer(os.Stdout, fileWriter)
	} else {
		ws = os.Stdout
	}

	core := zapcore.NewCore(encoder, ws, level)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &zapLogger{l: logger}
}

func (z *zapLogger) Debug(msg string, fields ...zap.Field) { z.l.Debug(msg, fields...) }
func (z *zapLogger) Info(msg string, fields ...zap.Field)  { z.l.Info(msg, fields...) }
func (z *zapLogger) Warn(msg string, fields ...zap.Field)  { z.l.Warn(msg, fields...) }
func (z *zapLogger) Error(msg string, fields ...zap.Field) { z.l.Error(msg, fields...) }
func (z *zapLogger) Fatal(msg string, fields ...zap.Field) { z.l.Fatal(msg, fields...) }
func (z *zapLogger) With(fields ...zap.Field) Logger       { return &zapLogger{l: z.l.With(fields...)} }

func (z *zapLogger) WithRequestID(ctx context.Context) Logger {
	fields := []zap.Field{}
	if rid, ok := ctx.Value("request_id").(string); ok && rid != "" {
		fields = append(fields, zap.String("request_id", rid))
	}
	return z.With(fields...)
}
