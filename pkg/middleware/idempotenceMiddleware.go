package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type IdempotenceMiddleware struct {
	*redis.Redis
}

func NewIdempotenceMiddleware() *IdempotenceMiddleware {
	return &IdempotenceMiddleware{}
}

func (m *IdempotenceMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}
