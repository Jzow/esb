package auth

import (
	"context"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	userCenter_service "github.com/EDDYCJY/go-gin-example/service/usercenter_service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtTicket, err := c.Cookie("token")
		util.Log("[auth]cookie token=%v", jwtTicket, err)
		tenantId, err := c.Cookie("tenant_id")
		util.Log("[auth]cookie tenantId=%v", tenantId, err)
		if jwtTicket != "" {
			ClientCheck(c)
			return
		}
		ServerCheck(c)
	}
}

func ClientCheck(c *gin.Context) {
	jwtTicket, _ := c.Cookie("token")
	ctx := context.WithValue(context.Background(), "stack", util.GetUUID32(""))
	tenantId, _ := c.Cookie("tenant_id")
	util.Log("[ClientCheck]cookie token=%v,cookie tenantId=%v", jwtTicket, tenantId)

	if jwtTicket == "" || tenantId == "" {
		if setting.ServerSetting.RunMode == "debug" { //debug mode
			util.Log("[ClientCheck]authorize by debug mode", ctx)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"code": e.UNAUTHORIZED, "msg": "illegal token"})
		c.Abort()
		return
	}
	claims, err := util.ParseToken(jwtTicket)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			util.Log("[ClientCheck]the jwt token is expired,jwtToken=%v, error: %v", ctx, jwtTicket, err)
		default:
			util.Log("[ClientCheck]failed to parse the jwt token, jwtToken=%v, error: %v", ctx, jwtTicket, err)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"code": e.UNAUTHORIZED, "msg": "invalid token"})
		c.Abort()
		return
	}

	resp := userCenter_service.CheckSignature(jwtTicket, tenantId, claims.UserID)
	if resp == nil || resp.Code != 200 {
		util.Log("[ClientCheck]tenantId=%v, UserID=%v jwtTicket=%v, not match tenant user,illegal visit", ctx, tenantId, claims.UserID, jwtTicket)
		c.JSON(http.StatusForbidden, gin.H{"code": e.FORBIDDEN, "msg": "illegal visit"})
		c.Abort()
		return
	}
	userResp := userCenter_service.GetTenantUsers(jwtTicket, tenantId, claims.UserID)
	if resp == nil || resp.Code != 200 {
		logging.Infof("[GetTenantUserInfo]tenantId=%v, UserID=%v, not match tenant user,illegal visit", tenantId, claims.UserID)
		c.JSON(http.StatusForbidden, gin.H{"code": e.FORBIDDEN, "msg": "illegal visit"})
		c.Abort()
		return
	}
	c.Set("user_id", userResp.Data[0].UserId)
	c.Set("tenant_id", tenantId)
	projectId := c.GetHeader("X-Project-Id")
	c.Set("project_id", projectId)
	c.Set("channel", "user")
	util.Log("[ClientCheck]user pass, userid=%v, tenantId=%v, projectId=%v, jwtTicket=%v", ctx, claims.UserID, tenantId, projectId, jwtTicket)
	c.Next()
}

func ServerCheck(c *gin.Context) {
	authParts := c.GetHeader("Authorization")
	ctx := context.WithValue(context.Background(), "stack", util.GetUUID32(""))
	/*if authParts == "" {
		if setting.ServerSetting.RunMode == "debug" {
			util.Log("[serviceAuth]authorize by debug mode", ctx)
			c.Next()
			return
		}
		util.Log("[serviceAuth]empty header Authorization or error format", ctx)
		c.JSON(http.StatusUnauthorized, gin.H{"code": e.UNAUTHORIZED, "msg": "illegal Authorization"})
		c.Abort()
		return
	}*/
	jwtToken := authParts
	claims, err := util.ParseToken(jwtToken)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			util.Log("[serviceAuth]the jwt token is expired,jwtToken=%v", ctx, jwtToken, err)
		default:
			util.Log("[serviceAuth]failed to parse the jwt token, jwtToken=%v", ctx, jwtToken, err)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"code": e.UNAUTHORIZED, "msg": "invalid Authorization"})
		c.Abort()
		return
	}

	if claims.PlatformID == 10 {
		util.Log("[serviceAuth]pass,identity=%v, jwtToken=%v", ctx, claims.UserID, jwtToken)
		c.Next()
	} else {
		// old client check
		tenantId := c.GetHeader("X-Tenant-Id")
		resp := userCenter_service.CheckSignature(jwtToken, tenantId, claims.UserID)
		if resp == nil || resp.Code != 200 {
			logging.Infof("[serviceAuth-oldCheck]tenantId=%v, UserID=%v, not match tenant user,illegal visit", tenantId, claims.UserID)
			c.JSON(http.StatusForbidden, gin.H{"code": e.FORBIDDEN, "msg": "illegal visit"})
			c.Abort()
			return
		}
		userResp := userCenter_service.GetTenantUsers(jwtToken, tenantId, claims.UserID)
		if resp == nil || resp.Code != 200 {
			logging.Infof("[GetTenantUserInfo]tenantId=%v, UserID=%v, not match tenant user,illegal visit", tenantId, claims.UserID)
			c.JSON(http.StatusForbidden, gin.H{"code": e.FORBIDDEN, "msg": "illegal visit"})
			c.Abort()
			return
		}
		c.Set("user_id", userResp.Data[0].UserId)
		c.Set("tenant_id", tenantId)
		projectId := c.GetHeader("X-Project-Id")
		c.Set("project_id", projectId)
		c.Set("channel", "user")
		logging.Infof("[serviceAuth-oldCheck]user pass, userid=%s, tenantId=%s, projectId=%s, jwtToken=%s", claims.UserID, tenantId, projectId, jwtToken)
		c.Next()
	}
}
