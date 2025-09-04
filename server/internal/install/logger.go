package install

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"bbs-go/internal/pkg/config"

	"github.com/pkg/errors"

	"bbs-go/internal/pkg/simple/common/strs"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	conf := config.Instance.Logger
	var writer io.Writer = os.Stdout
	if strs.IsNotBlank(conf.Filename) {
		writer = io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   conf.Filename,   // 日志文件的位置
				MaxSize:    conf.MaxSize,    // 文件最大尺寸（以MB为单位）
				MaxAge:     conf.MaxAge,     // 保留旧文件的最大天数
				MaxBackups: conf.MaxBackups, // 保留的最大旧文件数量
				LocalTime:  true,            // 使用本地时间创建时间戳
				Compress:   true,            // 是否压缩/归档旧文件
			},
		)
	}
	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
		AddSource:   true,
	})
	slog.SetDefault(slog.New(handler))
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	// other cases

	case slog.KindAny:
		switch v := a.Value.Any().(type) {
		case error:
			a.Value = fmtErr(v)
		}
	}

	return a
}

// fmtErr returns a slog.GroupValue with keys "msg" and "trace". If the error
// does not implement interface { StackTrace() errors.StackTrace }, the "trace"
// key is omitted.
func fmtErr(err error) slog.Value {
	var groupValues []slog.Attr

	groupValues = append(groupValues, slog.String("msg", err.Error()))

	type StackTracer interface {
		StackTrace() errors.StackTrace
	}

	// Find the trace to the location of the first errors.New,
	// errors.Wrap, or errors.WithStack call.
	var st StackTracer
	for err := err; err != nil; err = errors.Unwrap(err) {
		if x, ok := err.(StackTracer); ok {
			st = x
		}
	}

	if st != nil {
		groupValues = append(groupValues,
			slog.Any("trace", traceLines(st.StackTrace())),
		)
	}

	return slog.GroupValue(groupValues...)
}

func traceLines(frames errors.StackTrace) []string {
	traceLines := make([]string, len(frames))

	// Iterate in reverse to skip uninteresting, consecutive runtime frames at
	// the bottom of the trace.
	var skipped int
	skipping := true
	for i := len(frames) - 1; i >= 0; i-- {
		// Adapted from errors.Frame.MarshalText(), but avoiding repeated
		// calls to FuncForPC and FileLine.
		pc := uintptr(frames[i]) - 1
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			traceLines[i] = "unknown"
			skipping = false
			continue
		}

		name := fn.Name()

		if skipping && strings.HasPrefix(name, "runtime.") {
			skipped++
			continue
		} else {
			skipping = false
		}

		filename, lineNr := fn.FileLine(pc)

		traceLines[i] = fmt.Sprintf("%s %s:%d", name, filename, lineNr)
	}

	return traceLines[:len(traceLines)-skipped]
}
