package middleware

import (
	"bank-service-app/pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestLogger(log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("query", req.URL.RawQuery),
				zap.Int("status", res.Status),
				zap.String("ip", c.RealIP()),
				zap.String("user-agent", req.UserAgent()),
				zap.Duration("latency", time.Since(start)),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			fields = append(fields, zap.String("request_id", id))

			n := res.Status
			switch {
			case n >= 500:
				log.ErrorWithContext("Server error", err, fields...)
			case n >= 400:
				log.WarnWithContext("Client error", fields...)
			case n >= 300:
				log.InfoWithContext("Redirect", fields...)
			default:
				log.InfoWithContext("Success", fields...)
			}

			return nil
		}
	}
}
