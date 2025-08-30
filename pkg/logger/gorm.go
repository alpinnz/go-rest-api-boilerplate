package logger

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logLevel logger.LogLevel
}

func NewGormLogger(level logger.LogLevel) logger.Interface {
	return &GormLogger{logLevel: level}
}

func (gl *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	gl.logLevel = level
	return gl
}

func (gl *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if gl.logLevel >= logger.Info {
		Log.Info(msg, "args", data)
	}
}

func (gl *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if gl.logLevel >= logger.Warn {
		Log.Warn(msg, "args", data)
	}
}

func (gl *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if gl.logLevel >= logger.Error {
		Log.Error(msg, "args", data)
	}
}

func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if gl.logLevel <= logger.Silent {
		return // skip semua log
	}

	sql, rows := fc()
	elapsed := time.Since(begin)

	if err != nil {
		if gl.logLevel >= logger.Error {
			Log.Error("SQL execution failed",
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
				"error", err,
			)
		}
	} else {
		if gl.logLevel >= logger.Info {
			Log.Debug("SQL executed",
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
			)
		}
	}
}
