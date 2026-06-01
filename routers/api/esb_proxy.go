package api

import (
	"io"
	"net/http"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/service/openapi_proxy"
	"github.com/gin-gonic/gin"
)

var proxyClient = openapi_proxy.NewClient()

func ProxyOpenAPI(c *gin.Context) {
	appG := app.Gin{C: c}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, err.Error())
		return
	}

	resp, err := proxyClient.Proxy(
		c.Request.Method,
		c.Param("url"),
		c.Request.URL.RawQuery,
		body,
		c.GetHeader("Content-Type"),
	)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, resp.Data)
}
