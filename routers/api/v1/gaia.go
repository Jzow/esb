package v1

import (
	"github.com/EDDYCJY/go-gin-example/models/gaia"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/service/gaia_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var gaiaClient = &gaia_service.Client{}

func LeaveSubmit(c *gin.Context) {
	appG := app.Gin{C: c}
	var req gaia.LeaveSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, err.Error())
		return
	}
	resp, err := gaiaClient.PostLeaveSubmit(req)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

func MyApplications(c *gin.Context) {
	commonGet(c, setting.GaiaApiSetting.MyApplicationsPath, &gaia.MyApplicationsRequest{})
}
func LeaveQuota(c *gin.Context) {
	appG := app.Gin{C: c}
	var req gaia.LeaveQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, err.Error())
		return
	}
	resp, err := gaiaClient.Post(setting.GaiaApiSetting.LeaveQuotaPath, req)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
func LeaveTypes(c *gin.Context) {
	commonGet(c, setting.GaiaApiSetting.LeaveTypesPath, &gaia.LeaveTypesRequest{})
}
func LeaveHours(c *gin.Context) {
	appG := app.Gin{C: c}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, err.Error())
		return
	}
	resp, err := gaiaClient.Post(setting.GaiaApiSetting.LeaveHoursPath, req)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
func ExceptionList(c *gin.Context) {
	commonGet(c, setting.GaiaApiSetting.ExceptionListPath, &gaia.ExceptionListRequest{})
}

func commonGet(c *gin.Context, path string, model interface{}) {
	appG := app.Gin{C: c}
	if err := c.ShouldBindQuery(model); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, err.Error())
		return
	}
	resp, err := gaiaClient.Get(path, model)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
