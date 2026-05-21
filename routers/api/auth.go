package api

type appAuthentication struct {
	AppKey    string `valid:"Required; MaxSize(255)"`
	AppSecret string `valid:"Required; MaxSize(255)"`
}

type UserAuthentication struct {
	UserName string `valid:"Required; MaxSize(255)"`
	Password string `valid:"Required; MaxSize(255)"`
}

// @Tags auth
// @Summary Get access token
// @Produce  json
// @Param appKey query string true "appKey"
// @Param appSecret query string true "appSecret"
// @Success 200 {object} app.Response "success http status"
// @Failure 500 {object} app.Response "failure http status"
// @Router /api/ticket/getaccesstoken [get]
//func GetAccessToken(c *gin.Context) {
//	appG := app.Gin{C: c}
//	valid := validation.Validation{}
//
//	appAuth := appAuthentication{c.Query("appKey"), c.Query("appSecret")}
//	ok, _ := valid.Valid(&appAuth)
//
//	if !ok {
//		app.MarkErrors(valid.Errors)
//		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
//		return
//	}
//
//	authService := auth_service.Auth{Username: appAuth.AppKey, Password: util.EncodeSha256(appAuth.AppSecret)}
//	userId, err := authService.Check()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, app.Response{OperateCode: e.ERROR_AUTH_CHECK_TOKEN_FAIL, Msg: "failed to check"})
//		return
//	}
//
//	if userId == "" {
//		c.JSON(http.StatusUnauthorized, app.Response{OperateCode: e.UNAUTHORIZED, Msg: "error authentication"})
//		return
//	}
//
//	token, err := util.GenerateToken()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, app.Response{OperateCode: e.ERROR_AUTH_TOKEN, Msg: "failed to generate the token"})
//		return
//	}
//
//	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
//		"token": token,
//	})
//}

// @Tags auth
// @Summary Login
// @Accept x-www-form-urlencoded
// @Produce  json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {object} app.Response "success http status"
// @Failure 500 {object} app.Response "failure http status"
// @Router /api/ticket/login [post]
//func Login(c *gin.Context) {
//	appG := app.Gin{C: c}
//	valid := validation.Validation{}
//
//	userAuth := UserAuthentication{c.PostForm("username"), c.PostForm("password")}
//	ok, _ := valid.Valid(&userAuth)
//
//	if !ok {
//		app.MarkErrors(valid.Errors)
//		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
//		return
//	}
//
//	authService := auth_service.Auth{Username: userAuth.UserName, Password: util.EncodeSha256(userAuth.Password)}
//	userId, err := authService.Check()
//	if err != nil {
//		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TICKET_FAIL, nil)
//		return
//	}
//
//	if userId == "" {
//		appG.Response(http.StatusUnauthorized, e.ERROR_NOT_EXIST_TTICKET, nil)
//		return
//	}
//
//	ticket, err := authService.CreateTicket()
//	if err != nil {
//		appG.Response(http.StatusInternalServerError, e.ERROR_CREATE_TTICKET, nil)
//		return
//	}
//
//	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
//		"userid":   userId,
//		"username": userAuth.UserName,
//		"ticket":   ticket,
//	})
//}

// @Tags auth
// @Summary Logout
// @Produce  json
// @Param ticket header string true "ticket"
// @Success 200 {object} app.Response "success http status"
// @Failure 500 {object} app.Response "failure http status"
// @Router /api/ticket/logout [get]
//func Logout(c *gin.Context) {
//	appG := app.Gin{C: c}
//
//	ticket := c.GetHeader("ticket")
//	if ticket != "" {
//		bytes, err := gredis.Get(ticket)
//		if err == nil {
//			isOK, err2 := gredis.Delete(ticket)
//			if err2 != nil {
//				logging.Infof("[auth]failed to delete redis ticket value,ticket=%s, error:%s", ticket, err2.Error())
//			} else if !isOK {
//				logging.Infof("[auth]failed to delete redis ticket value,ticket=%s", ticket)
//			} else {
//				logging.Infof("[auth]%s logout,ticket=%s", bytes, ticket)
//				appG.Response(http.StatusOK, e.SUCCESS, nil)
//				return
//			}
//		} else {
//			logging.Infof("[auth]failed to get redis ticket value, error:%s", err.Error())
//		}
//	} else {
//		logging.Infof("[auth]failed to logout,ticket is empty")
//	}
//
//	appG.Response(http.StatusInternalServerError, e.ERROR, nil)
//}
