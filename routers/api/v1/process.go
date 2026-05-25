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
	Payload map[string]any `json:"payload" binding:"required"`
}

func proxyProcessAPI(c *gin.Context, apiURL string) {
	var req ProxyReq
	appG := app.Gin{C: c}
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, err.Error())
		return
	}

	resp, err := process_service.Default().Proxy(apiURL, req.Payload)
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

func SubmitProcess(c *gin.Context) { proxyProcessAPI(c, setting.GaiaOpenAPISetting.SubmitLeavePath) }
func MyProcesses(c *gin.Context)   { proxyProcessAPI(c, setting.GaiaOpenAPISetting.MyApplicationsPath) }
func ProcessQuota(c *gin.Context)  { proxyProcessAPI(c, setting.GaiaOpenAPISetting.LeaveQuotaPath) }
func ProcessTypes(c *gin.Context)  { proxyProcessAPI(c, setting.GaiaOpenAPISetting.LeaveTypesPath) }
func ProcessHours(c *gin.Context)  { proxyProcessAPI(c, setting.GaiaOpenAPISetting.LeaveHoursPath) }
func ProcessExceptions(c *gin.Context) {
	proxyProcessAPI(c, setting.GaiaOpenAPISetting.ExceptionListPath)
}
