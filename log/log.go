package log

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
)

const (
	// LogFlag 控制日志的前缀
	logFlag = log.LstdFlags | log.Lmicroseconds | log.Llongfile
)

var (
	errLogger  *log.Logger
	infoLogger *log.Logger
)

func init() {
	errLogger = log.New(os.Stderr, "", logFlag)
	infoLogger = log.New(os.Stdout, "", logFlag)
}

// Fatalf 打印错误日志并退出
func Fatalf(ctx context.Context, format string, v ...interface{}) {
	errLogger.Output(2, fmt.Sprintf(buildFormat(ctx, "FATAL", format), v...))
	os.Exit(1)
}

// Errorf 打印错误日志
func Errorf(ctx context.Context, format string, v ...interface{}) {
	errLogger.Output(2, fmt.Sprintf(buildFormat(ctx, "ERROR", format), v...))
}

// Warnf 打印警告日志
func Warnf(ctx context.Context, format string, v ...interface{}) {
	infoLogger.Output(2, fmt.Sprintf(buildFormat(ctx, "WARN", format), v...))
}

// Infof 打印错误日志
func Infof(ctx context.Context, format string, v ...interface{}) {
	infoLogger.Output(2, fmt.Sprintf(buildFormat(ctx, "INFO", format), v...))
}

func buildFormat(ctx context.Context, level, format string) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s ", level))

	buf.WriteString(format)

	return buf.String()
}
