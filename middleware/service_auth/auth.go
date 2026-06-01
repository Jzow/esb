package service_auth

import (
	"net/http"
	"strings"

	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	authservice "github.com/EDDYCJY/go-gin-example/service/service_auth"
	"github.com/gin-gonic/gin"
)

func Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtractToken(c)
		principal, err := authservice.Authenticate(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": e.UNAUTHORIZED, "msg": "invalid or expired token"})
			c.Abort()
			return
		}
		c.Set("auth_type", principal.Type)
		c.Set("auth_token", principal.Token)
		if principal.Username != "" {
			c.Set("username", principal.Username)
		}
		c.Next()
	}
}

func CheckOpenAPIURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		openAPI, ok, err := setting.FindOpenAPIByURL(c.Param("url"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": e.INVALID_PARAMS, "msg": "invalid openapi url"})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"code": e.INVALID_PARAMS, "msg": "unknown openapi url"})
			c.Abort()
			return
		}
		token := ExtractToken(c)
		principal, err := authservice.AuthenticateForApp(token, openAPI.Name)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": e.UNAUTHORIZED, "msg": "invalid token or openapi permission"})
			c.Abort()
			return
		}
		c.Set("auth_type", principal.Type)
		c.Set("auth_token", principal.Token)
		c.Set("openapi", openAPI.Name)
		if principal.Username != "" {
			c.Set("username", principal.Username)
		}
		c.Next()
	}
}

func ExtractToken(c *gin.Context) string {
	if token := strings.TrimSpace(c.GetHeader("X-API-Token")); token != "" {
		return token
	}
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	if auth == "" {
		return ""
	}
	parts := strings.Fields(auth)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return auth
}
