package http_log

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(context.Background(), "stack", util.GetUUID32(""))
		util.Log("Request URI: %v %v", ctx, c.Request.Method, c.Request.RequestURI)
		userId, _ := c.Get("user_id")
		util.Log("Request user_id: %v", ctx, userId)
		util.Log("Request Header: %v", ctx, maskedHeader(c.Request.Header))
		ticket, _ := c.Cookie("ticket")
		if ticket != "" {
			util.Log("Request Cookie ticket: %v", ctx, ticket)
		}
		if len(c.Request.Form) > 0 {
			util.Log("Request Form: %v", c.Request.Form)
		}
		if len(c.Request.PostForm) > 0 {
			util.Log("Request PostForm: %v", ctx, c.Request.PostForm)
		}
		body, err := io.ReadAll(c.Request.Body)
		util.Log("Request Body: %v", ctx, requestBodyForLog(c, body), err)

		if err == nil {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		// 将 gin.ResponseWriter 和一个缓冲区绑定
		wrapped := &responseWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = wrapped

		// 继续处理请求
		c.Next()

		// 打印响应参数
		util.Log("Response Status=%v,Body: %v", ctx, c.Writer.Status(), responseBodyForLog(c, wrapped.body.String()))
	}
}

func maskedHeader(header http.Header) http.Header {
	cloned := header.Clone()
	for _, key := range []string{"Authorization", "Cookie", "X-API-Token"} {
		if cloned.Get(key) != "" {
			cloned.Set(key, "***")
		}
	}
	return cloned
}

func requestBodyForLog(c *gin.Context, body []byte) any {
	if c.FullPath() == "/api/auth/login" {
		return "***"
	}
	return body
}

func responseBodyForLog(c *gin.Context, body string) string {
	if c.FullPath() == "/api/auth/login" {
		return "***"
	}
	return body
}
