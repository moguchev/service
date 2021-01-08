package middleware

import "github.com/sirupsen/logrus"

// Middleware - represent the data-struct for middleware
type Middleware struct {
	log *logrus.Logger
}

// InitMiddleware - initialize the middleware
func InitMiddleware(l *logrus.Logger) *Middleware {
	return &Middleware{
		log: l,
	}
}
