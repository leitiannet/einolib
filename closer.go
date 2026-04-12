package einolib

import (
	"context"
	"io"
	"sync"
)

// 说明：必须先 WithCloser 再 AddCloser，并在结束时 CloseCloser。否则资源无法被正确清理

// 资源清理收集器（通过 context 传递）
type closerCollector struct {
	mu      sync.Mutex
	closers []io.Closer
	once    sync.Once
}

func (cc *closerCollector) Close() {
	cc.once.Do(func() {
		cc.mu.Lock()
		defer cc.mu.Unlock()
		// 后注册的资源先关闭(LIFO)
		for i := len(cc.closers) - 1; i >= 0; i-- {
			if err := cc.closers[i].Close(); err != nil {
				logger.Warnf("close resource failed: %v", err)
			}
		}
		cc.closers = nil
	})
}

type closerCollectorKey struct{}

// 为 context 附加资源清理收集器
func WithCloser(ctx context.Context) context.Context {
	if ctx == nil {
		logger.Warnf("context is nil, using context.Background()")
		ctx = context.Background()
	}
	if _, ok := ctx.Value(closerCollectorKey{}).(*closerCollector); ok {
		return ctx
	}
	return context.WithValue(ctx, closerCollectorKey{}, &closerCollector{})
}

// 将需要清理的资源注册到 context 中的收集器
func AddCloser(ctx context.Context, closers ...io.Closer) {
	if ctx == nil {
		logger.Warnf("context is nil, closers not registered")
		return
	}
	cc, ok := ctx.Value(closerCollectorKey{}).(*closerCollector)
	if !ok || cc == nil {
		logger.Warnf("context has no closer collector, closers not registered")
		return
	}
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for _, c := range closers {
		if c != nil {
			cc.closers = append(cc.closers, c)
		}
	}
}

// 关闭 context 中收集器已注册的所有资源
func CloseCloser(ctx context.Context) {
	if ctx == nil {
		logger.Warnf("context is nil, nothing to close")
		return
	}
	if cc, ok := ctx.Value(closerCollectorKey{}).(*closerCollector); ok && cc != nil {
		cc.Close()
	}
}
