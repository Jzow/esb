package v1

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/gormType"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/base_service"
	"github.com/gin-gonic/gin"
)

func Paging[T models.DBEntity, F any](c *gin.Context, callback ...func(f *F, page *gormType.Page) (*gormType.Page, error)) {
	var (
		appG   = app.Gin{C: c}
		entity = new(T)
		form   = new(F)
	)
	errTip := appG.BindAuthAndValid(form, 1)
	if errTip != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, errTip)
		return
	}

	queryFields := setEntity(entity, form, "searchField", "defaultVal")
	optionalFields, statusMap := setOptionalEntity(entity, form, "optionalField")
	for _, field := range optionalFields {
		queryFields = append(queryFields, field)
	}

	srv := base_service.NewService(entity, form)
	err := srv.Paging(queryFields, statusMap)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_FAIL, nil)
		return
	}
	if len(callback) == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, srv.Page)
		return
	}

	result, cErr := callback[0](form, srv.Page)
	if cErr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, result)
}

func All[T models.DBEntity, F any](c *gin.Context, callback ...func(f *F, all []T) (any, error)) {
	var (
		appG   = app.Gin{C: c}
		entity = new(T)
		form   = new(F)
	)
	errTip := appG.BindAuthAndValid(form, 1)
	if errTip != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, errTip)
		return
	}

	queryFields := setEntity(entity, form, "searchField", "defaultVal")
	optionalFields, statusMap := setOptionalEntity(entity, form, "optionalField")
	for _, field := range optionalFields {
		queryFields = append(queryFields, field)
	}

	srv := base_service.NewService(entity, form)
	orders := getFieldName[F]("order")
	list, err := srv.GetAll(queryFields, orders, statusMap)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_FAIL, nil)
		return
	}
	if len(callback) == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, list)
		return
	}

	result, cErr := callback[0](form, list)
	if cErr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, result)
}

func Get[T models.DBEntity, F any](c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		entity = new(T)
		form   = new(F)
	)
	errTip := appG.BindAuthAndValid(form, 0)
	if errTip != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, errTip)
		return
	}
	reflect.ValueOf(entity).Elem().FieldByName(getPrimaryFieldName[T]()).Set(reflect.ValueOf(form).Elem().FieldByName("Id"))
	srv := base_service.NewService(entity)
	err := srv.Take(entity).Error
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, entity)
}

func Save[T models.DBEntity, F any](c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		entity = new(T)
		form   = new(F)
	)

	errTip := appG.BindAuthAndValid(form)
	if errTip != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, errTip)
		return
	}

	setEntity(entity, form, "addField", "$insert", "autoId", "autoPwd", "defaultVal")

	srv := base_service.NewService(entity)

	existResetFields := getFieldName[F]("existReset")
	if len(existResetFields) > 0 { //update
		dbItem, err := srv.Get(existResetFields)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_FAIL, nil)
			return
		}
		if dbItem != nil {
			statusField := reflect.ValueOf(dbItem).Elem().FieldByName("Status")
			if statusField.IsValid() && util.ToStr(statusField.Interface()) == "3" {
				resetEntity[T, F](dbItem, entity, "existReset")
			} else {
				appG.Response(http.StatusInternalServerError, e.ERROR_EXIST, nil, "已存在当前记录")
				return
			}
		}
	}

	existReset2Fields := getFieldName[F]("existReset2")
	if len(existReset2Fields) > 0 { //update
		dbItem, err := srv.Get(existReset2Fields)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_FAIL, nil)
			return
		}
		if dbItem != nil {
			statusField := reflect.ValueOf(dbItem).Elem().FieldByName("Status")
			if statusField.IsValid() && util.ToStr(statusField.Interface()) == "3" {
				resetEntity[T, F](dbItem, entity, "existReset2")
			} else {
				appG.Response(http.StatusInternalServerError, e.ERROR_EXIST, nil, "已存在当前记录")
				return
			}
		}
	}

	existUpdateFields := getFieldName[F]("existUpdate")
	if len(existUpdateFields) > 0 { //update
		dbItem, err := srv.Get(existUpdateFields)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_FAIL, nil)
			return
		}
		if dbItem != nil {
			resetEntity[T, F](dbItem, entity, "existUpdate")
		}
	}

	err := srv.Save(entity).Error
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, reflect.ValueOf(entity).Elem().FieldByName(getPrimaryFieldName[T]()).Interface())
}

func Edit[T models.DBEntity, F any](c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		entity = new(T)
		form   = new(F)
	)
	errTip := appG.BindAuthAndValid(form)
	if errTip != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, errTip)
		return
	}
	updateFields := setEntity(entity, form, "updateField", "$edit", "autoId", "autoPwd", "defaultVal")
	srv := base_service.NewService(entity)
	nums := srv.UpdateFieldById(updateFields)
	appG.Response(http.StatusOK, e.SUCCESS, nums)
}

func Delete[T models.DBEntity, F any](c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   = new(F)
		entity = new(T)
	)
	errTip := appG.BindAuthAndValid(form, 1)
	if errTip != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, errTip)
		return
	}
	entityElem := reflect.ValueOf(entity).Elem()
	formElem := reflect.ValueOf(form).Elem()
	entityElem.FieldByName(getPrimaryFieldName[T]()).Set(formElem.FieldByName("Id"))

	srv := base_service.NewService(entity)
	if err := srv.Delete(entity).Error; err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func Disable[T models.DBEntity, F any](c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   = new(F)
		entity = new(T)
	)
	errTip := appG.BindAuthAndValid(form, 1)
	if errTip != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, errTip)
		return
	}
	entityElem := reflect.ValueOf(entity).Elem()
	formElem := reflect.ValueOf(form).Elem()

	entityElem.FieldByName(getPrimaryFieldName[T]()).Set(formElem.FieldByName("Id"))
	fields := []string{"Status"}
	if entityElem.FieldByName("ModifyTime").IsValid() {
		entityElem.FieldByName("ModifyTime").Set(reflect.ValueOf(util.GetCurrentTimeStamp()))
		fields = append(fields, "ModifyTime")
	}
	entityElem.FieldByName("Status").SetInt(3)

	srv := base_service.NewService(entity)
	nums := srv.UpdateFieldById(fields)
	appG.Response(http.StatusOK, e.SUCCESS, nums)
}

func setEntity[T models.DBEntity, F any](entity *T, form *F, targetField string, tags ...string) []string {

	var attachType int
	var autoId, autoPwd, defaultField string
	for _, tag := range tags {
		switch tag {
		case "autoId":
			autoId = tag
		case "autoPwd":
			autoPwd = tag
		case "defaultVal":
			defaultField = tag
		case "$insert":
			attachType = 1
		case "$update":
			attachType = 2
		case "$edit":
			attachType = 3
		}
	}

	var (
		fields          = getFieldName[F](targetField)
		autoIdFieldMap  = getFieldMap[F](autoId)
		autoPwdFieldMap = getFieldMap[F](autoPwd)
		defaultFieldMap = getFieldMap[F](defaultField)
	)

	entityElem := reflect.ValueOf(entity).Elem()
	formElem := reflect.ValueOf(form).Elem()

	for _, field := range fields {
		if idType, ok := autoIdFieldMap[field]; ok {
			if idType == "uid" {
				if field == "Id" {
					field = getPrimaryFieldName[T]()
				}
				entityElem.FieldByName(field).SetString(util.GetUUID32(""))
			}
			if idType == "sf" {
				entityElem.FieldByName(field).SetInt(util.GetSnowFlakeId())
			}
			if idType == "sha256" {
				entityElem.FieldByName(field).SetString(util.Sha256(formElem.FieldByName(field).String()))
			}
			continue
		}
		if pwd, ok := autoPwdFieldMap[field]; ok {
			entityElem.FieldByName(field).SetString(util.Sha256(pwd.(string)))
			continue
		}
		if defaultVal, ok := defaultFieldMap[field]; ok {
			entityElem.FieldByName(field).Set(getReflectValue(formElem.FieldByName(field).Interface(), defaultVal))
			continue
		}
		entityElem.FieldByName(field).Set(formElem.FieldByName(field))
	}

	if attachType == 1 {
		if entityElem.FieldByName("CreateUserId").IsValid() {
			entityElem.FieldByName("CreateUserId").Set(formElem.FieldByName("CreateUserId"))
		}
		if entityElem.FieldByName("CreateTime").IsValid() {
			entityElem.FieldByName("CreateTime").Set(reflect.ValueOf(util.GetCurrentTimeStamp()))
		}
	}
	if attachType == 2 || attachType == 3 {
		if entityElem.FieldByName("ModifyTime").IsValid() {
			entityElem.FieldByName("ModifyTime").Set(reflect.ValueOf(util.GetCurrentTimeStamp()))
			fields = append(fields, "ModifyTime")
		}
	}
	if attachType == 3 {
		entityElem.FieldByName(getPrimaryFieldName[T]()).Set(formElem.FieldByName("Id"))
	}

	return fields
}

func resetEntity[T models.DBEntity, F any](dbEntity *T, entity *T, fieldTag string) []string {
	tagFieldNames := getFieldName[F](fieldTag)
	dbEntityElem := reflect.ValueOf(dbEntity).Elem()
	entityElem := reflect.ValueOf(entity).Elem()
	for _, fieldName := range tagFieldNames {
		if fieldName == "existUpdate" { //保存更新
			if util.ContainStr([]string{"CreateTime", "create_user_id"}, fieldName) {
				continue
			}
		}
		entityElem.FieldByName(fieldName).Set(dbEntityElem.FieldByName(fieldName))
	}
	if fieldTag != "existUpdate" {
		if entityElem.FieldByName("Status").IsValid() {
			entityElem.FieldByName("Status").Set(reflect.ValueOf(1))
		}
		if entityElem.FieldByName("ModifyTime").IsValid() {
			entityElem.FieldByName("ModifyTime").Set(reflect.ValueOf(uint(0)))
		}
		if entityElem.FieldByName("ModifyUserId").IsValid() {
			entityElem.FieldByName("ModifyUserId").Set(reflect.ValueOf(""))
		}
		id := getPrimaryFieldName[T]()
		if entityElem.FieldByName(id).IsValid() {
			entityElem.FieldByName(id).Set(dbEntityElem.FieldByName(id))
		}
	}
	return tagFieldNames
}

func getPointValue(val any) reflect.Value {
	switch val.(type) {
	case *int:
		return reflect.ValueOf(*val.(*int))
	case *int8:
		return reflect.ValueOf(*val.(*int8))
	case *int16:
		return reflect.ValueOf(*val.(*int16))
	case *int32:
		return reflect.ValueOf(*val.(*int16))
	case *int64:
		return reflect.ValueOf(*val.(*int64))
	case *uint:
		return reflect.ValueOf(*val.(*uint))
	case *uint8:
		return reflect.ValueOf(*val.(*uint8))
	case *uint16:
		return reflect.ValueOf(*val.(*uint16))
	case *uint32:
		return reflect.ValueOf(*val.(*uint16))
	case *uint64:
		return reflect.ValueOf(*val.(*uint64))
	case *string:
		return reflect.ValueOf(*val.(*string))
	case *bool:
		return reflect.ValueOf(*val.(*bool))
	case *struct{}:
		return reflect.ValueOf(*val.(*struct{}))
	}
	return reflect.ValueOf(val)
}

func getReflectValue(source any, val any) reflect.Value {
	switch source.(type) {
	case int, int32:
		return reflect.ValueOf(util.ObjToInt(val))
	case int64:
		return reflect.ValueOf(util.ObjToInt64(val))
	case uint, uint32:
		return reflect.ValueOf(util.ObjToUint32(val))
	case uint64:
		return reflect.ValueOf(util.ObjToUint64(val))
	case string:
		return reflect.ValueOf(util.ObjToStr(val))
	case bool:
		return reflect.ValueOf(val.(bool))
	case struct{}:
		return reflect.ValueOf(val.(struct{}))
	}
	return reflect.ValueOf(val)
}

func setOptionalEntity[T models.DBEntity, F any](entity *T, form *F, targetField string) ([]string, map[string]any) {
	var (
		fields         = getFieldName[F](targetField)
		statusMap      map[string]any
		optionalFields []string
	)

	entityElem := reflect.ValueOf(entity).Elem()
	formElem := reflect.ValueOf(form).Elem()

	for _, field := range fields {
		formField := formElem.FieldByName(field)
		if formField.Kind() == reflect.Pointer {
			if !formField.IsNil() {
				entityElem.FieldByName(field).Set(getPointValue(formField.Interface()))
				optionalFields = append(optionalFields, field)
			} else if field == "Status" {
				sField, _ := reflect.TypeOf(form).Elem().FieldByName("Status")
				tagVal := sField.Tag.Get("nil")
				if tagVal != "" {
					var intVals []int
					for _, s := range strings.Split(tagVal, ",") {
						intVals = append(intVals, util.ObjToInt(s))
					}
					statusMap = map[string]any{"status": intVals}
				}
			}
			continue
		}
		entityElem.FieldByName(field).Set(formField)
		optionalFields = append(optionalFields, field)
	}
	return optionalFields, statusMap
}

func getPrimaryFieldName[T models.DBEntity]() string {
	ele := reflect.TypeOf(new(T)).Elem()
	for i := 0; i < ele.NumField(); i++ {
		tag := ele.Field(i).Tag.Get("gorm")
		if strings.Contains(tag, "primary") {
			return ele.Field(i).Name
		}
	}
	return ""
}

func getFieldName[F any](tagName string) []string {
	var updateFields []string
	if tagName == "" {
		return updateFields
	}
	ele := reflect.TypeOf(new(F)).Elem()
	for i := 0; i < ele.NumField(); i++ {
		tag := ele.Field(i).Tag.Get(tagName)
		if tag == "true" {
			updateFields = append(updateFields, ele.Field(i).Name)
			continue
		}
		if tag != "" {
			updateFields = append(updateFields, tag)
		}
	}
	return updateFields
}

func getFieldMap[F any](tagName string) map[string]interface{} {
	fieldMap := make(map[string]interface{})
	if tagName == "" {
		return fieldMap
	}
	ele := reflect.TypeOf(new(F)).Elem()
	for i := 0; i < ele.NumField(); i++ {
		tag := ele.Field(i).Tag.Get(tagName)
		if tag != "" {
			fieldMap[ele.Field(i).Name] = tag
		}
	}
	return fieldMap
}
