package app

import (
	"context"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xuri/excelize/v2"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type SuccessResponse struct {
	Code int         `json:"code" example:"200"`
	Msg  string      `json:"msg" example:"ok"`
	Data interface{} `json:"data,omitempty"`
}

func (g *Gin) IsLocalDebug() bool {
	return util.ContainStr([]string{"localhost", "127.0.0.1", "0.0.0.0", "::1"}, strings.Split(g.C.Request.Host, ":")[0])
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}, msg ...string) {
	var tip string
	if len(msg) == 0 {
		tip = e.GetMsg(errCode)
	} else {
		tip = msg[0]
	}
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  tip,
		Data: data,
	})
	return
}

func (g *Gin) QueryAuth(key, authField string) string {
	if authField == "tenant_id" || authField == "project_id" || authField == "user_id" {
		channel, _ := g.C.Get("channel")
		if channel == "user" || (channel == "tenant" && authField == "tenant_id") {
			return g.C.GetString(authField)
		}
	}
	return g.C.Query(key)
}

func (g *Gin) ParamUint(key string) uint64 {
	val := g.C.Param(key)
	valInt, _ := strconv.ParseUint(val, 10, 64)
	return valInt
}

func (g *Gin) Query(authField string) string {
	return g.QueryAuth(authField, authField)
}

// BindAndValid binds and validates data
func (g *Gin) BindAndValid(form interface{}) string {
	err := g.C.BindJSON(form)
	if err != nil {
		return "failed to parse struct param"
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return "failed to valid param"
	}
	if !check {
		return MarkErrorTip(valid.Errors)
	}

	return ""
}

// BindAuthAndValid binds and validates data , reset auth tag field value if auth field exist
func (g *Gin) BindAuthAndValid(form interface{}, bindType ...int) string {
	var err error
	if len(bindType) == 0 {
		err = g.C.ShouldBindBodyWith(form, binding.JSON)
	} else if bindType[0] == 1 {
		err = g.C.BindQuery(form)
	}
	if err != nil {
		return "failed to parse struct param"
	}

	elems := reflect.TypeOf(form).Elem()
	for i := 0; i < elems.NumField(); i++ {
		field := elems.Field(i)

		if g.IsLocalDebug() {
			localFieldValue := field.Tag.Get("localField")
			if localFieldValue != "" {
				piField := reflect.ValueOf(form).Elem().FieldByName(field.Name)
				if piField.CanSet() {
					switch piField.Kind() {
					case reflect.String:
						piField.SetString(localFieldValue)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						piField.SetInt(util.ObjToInt64(localFieldValue))
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						piField.SetUint(util.ObjToUint64(localFieldValue))
					}
				}
			}
		}

		//reset query param
		queryField := field.Tag.Get("queryField")
		if queryField != "" {
			queryFieldValue := g.C.Query(queryField)
			piField := reflect.ValueOf(form).Elem().FieldByName(field.Name)
			if piField.CanSet() {
				switch piField.Kind() {
				case reflect.String:
					piField.SetString(queryFieldValue)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					piField.SetInt(util.ObjToInt64(queryFieldValue))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					piField.SetUint(util.ObjToUint64(queryFieldValue))
				}
			}
		}

		//reset path param
		paramField := field.Tag.Get("paramField")
		if paramField != "" {
			paramFieldValue := g.C.Param(paramField)
			piField := reflect.ValueOf(form).Elem().FieldByName(field.Name)
			if piField.CanSet() {
				switch piField.Kind() {
				case reflect.String:
					piField.SetString(paramFieldValue)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					piField.SetInt(util.ObjToInt64(paramFieldValue))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					piField.SetUint(util.ObjToUint64(paramFieldValue))
				}
			}
		}

		//reset auth param
		authField := field.Tag.Get("authField")
		if authField != "" {
			authFieldValue, exist := g.C.Get(authField)
			if exist {
				piField := reflect.ValueOf(form).Elem().FieldByName(field.Name)
				if piField.CanSet() {
					switch piField.Kind() {
					case reflect.String:
						fieldValue := util.ObjToStr(authFieldValue)
						if authField == "tenant_id" || authField == "project_id" || authField == "user_id" {
							channel, _ := g.C.Get("channel")
							if channel == "user" || (channel == "tenant" && authField == "tenant_id") {
								piField.SetString(fieldValue)
								break
							}
							break
						}
						piField.SetString(fieldValue)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						piField.SetInt(util.ObjToInt64(authFieldValue))
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						piField.SetUint(util.ObjToUint64(authFieldValue))
					}
				}
			}
		}

	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return "failed to valid param"
	}
	if !check {
		return MarkErrorTip(valid.Errors)
	}

	return ""
}

func (g *Gin) ExportExcel(ctx context.Context, filename string, column []any, data [][]any) {
	xlsx := excelize.NewFile()
	defer xlsx.Close()
	sheetName := "Sheet1"
	xlsx.SetSheetRow(sheetName, "A1", &column)
	for i := 0; i < len(data); i++ {
		xlsx.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+2), &data[i])
	}
	err := xlsx.Write(g.C.Writer)
	util.Log("[ExportExcel]filename=%v,count=%v,go version=%v", ctx, filename, len(data), runtime.Version(), err)
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR, nil, "failed to export")
	} else {
		g.C.Writer.Header().Set("Content-Type", "application/octet-stream")
		g.C.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	}
}
