package api

import (
	"net/http"

	"github.com/EDDYCJY/go-gin-example/middleware/service_auth"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	authservice "github.com/EDDYCJY/go-gin-example/service/service_auth"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, err.Error())
		return
	}

	resp, err := authservice.Login(req.Username, req.Password)
	if err != nil {
		appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, nil, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

func Logout(c *gin.Context) {
	appG := app.Gin{C: c}
	if err := authservice.Logout(service_auth.ExtractToken(c)); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
