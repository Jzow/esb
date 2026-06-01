package routers

import (
	"net/http"

	"github.com/EDDYCJY/go-gin-example/middleware/auth"
	"github.com/EDDYCJY/go-gin-example/middleware/http_log"
	serviceauth "github.com/EDDYCJY/go-gin-example/middleware/service_auth"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/routers/api"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
		//测试
		/*		apiv1.GET("/test", v1.GetTests)
				apiv1.GET("/test/:id", v1.GetTest)
				apiv1.POST("/test", v1.AddTest)
				apiv1.PUT("/test/:id", v1.EditTest)
				apiv1.DELETE("/test/:id", v1.DeleteTest)*/
	}

	authRouter := r.Group("/api/auth")
	authRouter.Use(http_log.LogMiddleware())
	{
		authRouter.POST("/login", api.Login)
		authRouter.POST("/logout", serviceauth.Check(), api.Logout)
	}

	esb := r.Group("/api/esb/v1")
	esb.Use(http_log.LogMiddleware(), serviceauth.CheckOpenAPIURL())
	{
		esb.GET("/*url", api.ProxyOpenAPI)
		esb.POST("/*url", api.ProxyOpenAPI)
		esb.PUT("/*url", api.ProxyOpenAPI)
		esb.DELETE("/*url", api.ProxyOpenAPI)
	}
	return r
}
