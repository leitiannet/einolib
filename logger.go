package einolib

import (
	"log"
	"sync/atomic"
)

// 日志接口
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// 默认日志实现，使用标准库log
type defaultLogger struct{}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// 仅作 atomic.Value 的固定存储类型（*loggerHolder），本身不实现 Logger
type loggerHolder struct {
	impl Logger
}

// 日志代理，方便包内用logger.xxx
type loggerProxy struct{}

func (loggerProxy) Debugf(format string, args ...interface{}) { GetLogger().Debugf(format, args...) }
func (loggerProxy) Infof(format string, args ...interface{})  { GetLogger().Infof(format, args...) }
func (loggerProxy) Warnf(format string, args ...interface{})  { GetLogger().Warnf(format, args...) }
func (loggerProxy) Errorf(format string, args ...interface{}) { GetLogger().Errorf(format, args...) }

var (
	loggerValue       atomic.Value // *loggerHolder
	logger            loggerProxy  // 零值结构体，无数据竞争问题
	defaultLoggerImpl defaultLogger // 用于 init 注入与 GetLogger 异常回退
)

func init() {
	loggerValue.Store(&loggerHolder{impl: &defaultLoggerImpl})
}

// 设置自定义日志实现，传入nil则不生效
func SetLogger(impl Logger) {
	if impl == nil {
		return
	}
	loggerValue.Store(&loggerHolder{impl: impl})
}

// 获取当前日志实现
func GetLogger() Logger {
	v := loggerValue.Load()
	if h, ok := v.(*loggerHolder); ok && h != nil && h.impl != nil {
		return h.impl
	}
	return &defaultLoggerImpl
}
