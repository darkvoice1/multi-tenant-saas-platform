package middleware

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

type bodyCaptureWriter struct {
	gin.ResponseWriter
	body   bytes.Buffer
	status int
}

func (w *bodyCaptureWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func (w *bodyCaptureWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func (w *bodyCaptureWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func ErrorCodeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		orig := c.Writer
		capture := &bodyCaptureWriter{ResponseWriter: orig}
		c.Writer = capture

		c.Next()

		status := capture.status
		if status == 0 {
			status = c.Writer.Status()
			if status == 0 {
				status = 200
			}
		}

		body := capture.body.Bytes()
		contentType := orig.Header().Get("Content-Type")
		if status >= 400 && strings.Contains(contentType, "application/json") && len(body) > 0 {
			var obj map[string]interface{}
			if err := json.Unmarshal(body, &obj); err == nil {
				if _, ok := obj["code"]; !ok {
					obj["code"] = defaultErrorCode(status)
				}
				if _, ok := obj["error"]; ok {
					body, _ = json.Marshal(obj)
				}
			}
		}

		orig.WriteHeader(status)
		_, _ = orig.Write(body)
	}
}

func defaultErrorCode(status int) string {
	switch status {
	case 400:
		return "BAD_REQUEST"
	case 401:
		return "UNAUTHORIZED"
	case 403:
		return "FORBIDDEN"
	case 404:
		return "NOT_FOUND"
	case 409:
		return "CONFLICT"
	case 413:
		return "PAYLOAD_TOO_LARGE"
	case 429:
		return "RATE_LIMITED"
	default:
		return "INTERNAL_ERROR"
	}
}
