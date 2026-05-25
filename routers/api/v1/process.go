package v1

import (
	"encoding/json"
	"net/http"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/service/process_service"
	"github.com/gin-gonic/gin"
)

type ProxyReq struct {
	Payload map[string]any `json:"payload"`
}

func proxyProcessAPI(c *gin.Context, method, apiURL string) {
	appG := app.Gin{C: c}
	payload := map[string]any{}

	if method == http.MethodGet {
		for k, vals := range c.Request.URL.Query() {
			if len(vals) > 0 {
				payload[k] = vals[0]
			}
		}
	} else {
		var req ProxyReq
		if err := c.ShouldBindJSON(&req); err != nil {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, err.Error())
			return
		}
		payload = req.Payload
		if payload == nil {
			payload = map[string]any{}
		}
	}

	resp, err := process_service.Default().Proxy(method, apiURL, payload)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil, err.Error())
		return
	}

	var out any
	if err := json.Unmarshal(resp, &out); err != nil {
		appG.Response(http.StatusOK, e.SUCCESS, string(resp))
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, out)
}

func SubmitProcess(c *gin.Context) {
	proxyProcessAPI(c, http.MethodPost, setting.GaiaOpenAPISetting.SubmitLeavePath)
}
func MyProcesses(c *gin.Context) {
	proxyProcessAPI(c, http.MethodGet, setting.GaiaOpenAPISetting.MyApplicationsPath)
}
func ProcessQuota(c *gin.Context) {
	proxyProcessAPI(c, http.MethodPost, setting.GaiaOpenAPISetting.LeaveQuotaPath)
}
func ProcessTypes(c *gin.Context) {
	proxyProcessAPI(c, http.MethodGet, setting.GaiaOpenAPISetting.LeaveTypesPath)
}
func ProcessHours(c *gin.Context) {
	proxyProcessAPI(c, http.MethodGet, setting.GaiaOpenAPISetting.LeaveHoursPath)
}
func ProcessExceptions(c *gin.Context) {
	proxyProcessAPI(c, http.MethodGet, setting.GaiaOpenAPISetting.ExceptionListPath)
}
