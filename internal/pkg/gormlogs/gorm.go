// see https://github.com/orandin/slog-gorm/tree/master

package gormlogs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type LogType string

const (
	ErrorLogType     LogType = "sql_error"
	SlowQueryLogType LogType = "slow_query"
	DefaultLogType   LogType = "default"

	SourceField    = "file"
	ErrorField     = "error"
	QueryField     = "query"
	DurationField  = "duration"
	SlowQueryField = "slow_query"
	RowsField      = "rows"
)

// New creates a new logger for gorm.io/gorm
func New(options ...Option) *logger {
	l := logger{
		ignoreRecordNotFoundError: true,
		errorField:                ErrorField,
		sourceField:               SourceField,

		// log levels
		logLevel: map[LogType]slog.Level{
			ErrorLogType:     slog.LevelError,
			SlowQueryLogType: slog.LevelWarn,
			DefaultLogType:   slog.LevelInfo,
		},
	}

	// Apply options
	for _, option := range options {
		option(&l)
	}

	if l.slogger == nil {
		// If no slogger is defined, use the default Logger
		l.slogger = slog.Default()
	}

	return &l
}

type logger struct {
	slogger                   *slog.Logger
	ignoreTrace               bool
	ignoreRecordNotFoundError bool
	traceAll                  bool
	slowThreshold             time.Duration
	logLevel                  map[LogType]slog.Level
	contextKeys               map[string]string

	sourceField string
	errorField  string
}

// LogMode log mode
func (l logger) LogMode(_ gormlogger.LogLevel) gormlogger.Interface {
	// log level is set by slog
	return l
}

// Info logs info
func (l logger) Info(ctx context.Context, msg string, args ...any) {
	l.log(l.slogger.InfoContext, ctx, msg, args...)
}

// Warn logs warn messages
func (l logger) Warn(ctx context.Context, msg string, args ...any) {
	l.log(l.slogger.WarnContext, ctx, msg, args...)
}

// Error logs error messages
func (l logger) Error(ctx context.Context, msg string, args ...any) {
	l.log(l.slogger.ErrorContext, ctx, msg, args...)
}

// log adds context attributes and logs a message with the given slog function
func (l logger) log(f func(ctx context.Context, msg string, args ...any), ctx context.Context, msg string, args ...any) {

	// Append context attributes
	args = l.appendContextAttributes(ctx, args)

	// Call slog
	f(ctx, msg, args...)
}

// Trace logs sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.ignoreTrace {
		return // Silent
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		sql, rows := fc()

		// Append context attributes
		attributes := l.appendContextAttributes(ctx, []any{
			slog.Any(l.errorField, err),
			slog.String(QueryField, sql),
			slog.Duration(DurationField, elapsed),
			slog.Int64(RowsField, rows),
			slog.String(l.sourceField, utils.FileWithLineNum()),
		})

		l.slogger.Log(ctx, l.logLevel[ErrorLogType], err.Error(), attributes...)

	case l.slowThreshold != 0 && elapsed > l.slowThreshold:
		sql, rows := fc()

		// Append context attributes
		attributes := l.appendContextAttributes(ctx, []any{
			slog.Bool(SlowQueryField, true),
			slog.String(QueryField, sql),
			slog.Duration(DurationField, elapsed),
			slog.Int64(RowsField, rows),
			slog.String(l.sourceField, utils.FileWithLineNum()),
		})
		l.slogger.Log(ctx, l.logLevel[SlowQueryLogType], fmt.Sprintf("slow sql query [%s >= %v]", elapsed, l.slowThreshold), attributes...)

	case l.traceAll:
		sql, rows := fc()

		// Append context attributes
		attributes := l.appendContextAttributes(ctx, []any{
			slog.String(QueryField, sql),
			slog.Duration(DurationField, elapsed),
			slog.Int64(RowsField, rows),
			slog.String(l.sourceField, utils.FileWithLineNum()),
		})

		l.slogger.Log(ctx, l.logLevel[DefaultLogType], fmt.Sprintf("SQL query executed [%s]", elapsed), attributes...)
	}
}

func (l logger) appendContextAttributes(ctx context.Context, args []any) []any {
	if args == nil {
		args = []any{}
	}
	for k, v := range l.contextKeys {
		if value := ctx.Value(v); value != nil {
			args = append(args, slog.Any(k, value))
		}
	}
	return args
}
