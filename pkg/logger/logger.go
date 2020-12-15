package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ctxlog struct{}

// WithLogger put logger to context
func WithLogger(ctx context.Context, l *logrus.Logger) context.Context {
	return context.WithValue(ctx, ctxlog{}, l)
}

var DefaultLogger = logrus.New()

// GetLogger get logger from context, or DefaultLogger if not exists
func GetLogger(ctx context.Context) *logrus.Logger {
	l, ok := ctx.Value(ctxlog{}).(*logrus.Logger)
	if !ok {
		l = DefaultLogger
	}
	return l
}
