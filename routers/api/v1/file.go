package v1

import (
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Tags file
// @Summary FileUpload
// @Produce  json
// @Param formData file true "File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/wiki/v1/file/upload [post]

func FileUpload(c *gin.Context) {
	appG := app.Gin{C: c}
	file, info, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if info == nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	name := upload.GetImageName(info.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + name

	if !upload.CheckImageExt(name) || !upload.CheckImageSize(file) {
		appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	if err := c.SaveUploadedFile(info, src); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"url":      upload.GetImageFullUrl(name),
		"save_url": savePath + name,
	})
}

// @Tags file
// @Summary FileDownload
// @Produce  json
// @Param formData file true "File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/wiki/v1/file/download [get]

func FileDownload(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
