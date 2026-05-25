package routers

import (
	"github.com/EDDYCJY/go-gin-example/middleware/auth"
	"github.com/EDDYCJY/go-gin-example/middleware/http_log"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	v1 "github.com/EDDYCJY/go-gin-example/routers/api/v1"
	"github.com/gin-gonic/gin"

	_ "github.com/EDDYCJY/go-gin-example/docs"
	"github.com/EDDYCJY/go-gin-example/pkg/export"
	"github.com/EDDYCJY/go-gin-example/pkg/qrcode"
	"github.com/EDDYCJY/go-gin-example/pkg/upload"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFS("/api/test/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/api/test/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/api/test/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	if setting.ServerSetting.RunMode == "debug" {
		r.StaticFS("/api/test/logs", http.Dir("runtime/logs"))
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	apiv1 := r.Group("/api/test/v1")
	apiv1.Use(auth.Check(), http_log.LogMiddleware())
	{
		// 其他内部受保护接口

		//测试
		/*		apiv1.GET("/test", v1.GetTests)
				apiv1.GET("/test/:id", v1.GetTest)
				apiv1.POST("/test", v1.AddTest)
				apiv1.PUT("/test/:id", v1.EditTest)
				apiv1.DELETE("/test/:id", v1.DeleteTest)*/

	}

	processApi := r.Group("/api/test/v1/process")
	processApi.Use(http_log.LogMiddleware())
	{
		processApi.POST("/submit", v1.SubmitProcess)
		processApi.GET("/my-applications", v1.MyProcesses)
		processApi.POST("/quota", v1.ProcessQuota)
		processApi.GET("/types", v1.ProcessTypes)
		processApi.POST("/hours", v1.ProcessHours)
		processApi.GET("/exceptions", v1.ProcessExceptions)
	}
	return r
}
