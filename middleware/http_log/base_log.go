package http_log

import (
	"context"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/gin-gonic/gin"
)

func BaseLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(context.Background(), "stack", util.GetUUID32(""))
		util.Log("Request URI: %v %v", ctx, c.Request.Method, c.Request.RequestURI)
		userId, _ := c.Get("user_id")
		util.Log("Request user_id: %v", ctx, userId)
		util.Log("Request Header: %v", ctx, c.Request.Header)
		ticket, _ := c.Cookie("ticket")
		if ticket != "" {
			util.Log("Request Cookie ticket: %v", ctx, ticket)
		}
		c.Next()
	}
}
