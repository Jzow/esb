package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetCurrentTimeStamp() uint {
	intNum, _ := strconv.Atoi(fmt.Sprintf("%d", time.Now().Unix()))
	return uint(intNum)
}

func GetCurrentChinaTime() time.Time {
	region, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(region)
}

func GetTimeStampByStrFormat(timeStr string) uint {
	var (
		localTime time.Time
		err       error
	)
	for _, dateFormat := range []string{"2006-01-02 15:04:05", "2006-1-2 15:04:05", "2006/01/02 15:04:05", "2006/1/2 15:04:05", "2006-01-02", "2006-1-2", "2006/01/02", "2006/1/2", "2006-01-02 15:04"} {
		localTime, err = time.ParseInLocation(dateFormat, timeStr, time.Local)
		if err == nil {
			break
		}
	}
	if err != nil {
		return 0
	}
	intNum, err := strconv.Atoi(fmt.Sprintf("%d", localTime.Unix()))
	if err != nil {
		return 0
	}
	return uint(intNum)
}

func GetStructFieldTag(obj interface{}, keyTag, valueTag string) map[string]string {
	tagMap := make(map[string]string)
	elems := reflect.TypeOf(obj).Elem()
	for i := 0; i < elems.NumField(); i++ {
		field := elems.Field(i)
		tagMap[field.Tag.Get(keyTag)] = field.Tag.Get(valueTag)
	}
	return tagMap
}

func ObjToTimeStamp(obj interface{}) uint {
	return GetTimeStampByStrFormat(ObjToStr(obj))
}
func ObjToUint32(obj interface{}) uint32 {
	val, _ := strconv.ParseInt(ToStr(obj), 10, 32)
	return uint32(val)
}
func ObjToUint(obj interface{}) uint {
	val, _ := strconv.ParseInt(ToStr(obj), 10, 32)
	return uint(val)
}
func ObjToInt64(obj interface{}) int64 {
	val, err := strconv.ParseInt(ToStr(obj), 10, 64)
	if err != nil {
		Log("[ObjToInt64]obj=%v", obj, err)
	}
	return val
}
func ObjToUint64(obj interface{}) uint64 {
	val, err := strconv.ParseUint(ToStr(obj), 10, 64)
	if err != nil {
		Log("[ObjToUint64]obj=%v", obj, err)
	}
	return val
}
func ObjToInt(obj interface{}) int {
	val, err := strconv.Atoi(ToStr(obj))
	if err != nil {
		Log("[ObjToInt]obj=%v", obj, err)
	}
	return val
}
func ToStr(obj interface{}) string {
	if obj == nil {
		return ""
	}
	switch obj.(type) {
	case string:
		return obj.(string)
	case float64, float32:
		return fmt.Sprintf("%.0f", obj)
	default:

		return fmt.Sprintf("%v", obj)
	}
}
func ObjToStr(obj interface{}) string {
	if obj == nil {
		return ""
	}
	switch obj.(type) {
	case string:
		return obj.(string)
	case float64, float32:
		return fmt.Sprintf("%f", obj)
	default:

		return fmt.Sprintf("%v", obj)
	}
}
func ObjToStrByArray(objs ...interface{}) string {
	if objs != nil && len(objs) > 0 {
		for _, obj := range objs {
			str := ObjToStr(obj)
			if str != "" {
				return str
			}
		}
	}
	return ""
}

func ContainInt(array []int, element int) bool {
	for _, item := range array {
		if element == item {
			return true
		}
	}
	return false
}

func ContainStr(array []string, element string) bool {
	for _, item := range array {
		if element == item {
			return true
		}
	}
	return false
}

func UnionInt(A, B []int) []int {
	result := make([]int, 0)
	// 去重
	flagMap := make(map[int]bool, 0)
	A = append(A, B...)
	for _, a := range A {
		if _, ok := flagMap[a]; ok {
			continue
		}
		flagMap[a] = true
		result = append(result, a)
	}
	return result
}

func IntersectInt(A, B []int) []int {
	if len(A) < 1 || len(B) < 1 {
		return []int{}
	}
	result := make([]int, 0)
	// 去重
	flagMap := make(map[int]bool, 0)
	for _, a := range A {
		if _, ok := flagMap[a]; ok {
			continue
		}
		flagMap[a] = true
		for _, b := range B {
			if b == a {
				result = append(result, a)
				break
			}
		}
	}
	return result
}

func ExceptInt(A, B []int) []int {
	if len(A) < 1 || len(B) < 1 {
		return A
	}
	result := make([]int, 0)
	// 去重
	flagMap := make(map[int]bool, 0)
	for _, a := range A {
		if _, ok := flagMap[a]; ok {
			continue
		}
		flagMap[a] = true
		flag := true
		for _, b := range B {
			if b == a {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, a)
		}
	}
	return result
}

func UnionStr(A, B []string) []string {
	result := make([]string, 0)
	// 去重
	flagMap := make(map[string]bool, 0)
	A = append(A, B...)
	for _, a := range A {
		if _, ok := flagMap[a]; ok {
			continue
		}
		flagMap[a] = true
		result = append(result, a)
	}
	return result
}

func IntersectStr(A, B []string) []string {
	if len(A) < 1 || len(B) < 1 {
		return []string{}
	}
	result := make([]string, 0)
	// 去重
	flagMap := make(map[string]bool, 0)
	for _, a := range A {
		if _, ok := flagMap[a]; ok {
			continue
		}
		flagMap[a] = true
		for _, b := range B {
			if b == a {
				result = append(result, a)
				break
			}
		}
	}
	return result
}

func ExceptStr(A, B []string) []string {
	if len(A) < 1 || len(B) < 1 {
		return A
	}
	result := make([]string, 0)
	// 去重
	flagMap := make(map[string]bool, 0)
	for _, a := range A {
		if _, ok := flagMap[a]; ok {
			continue
		}
		flagMap[a] = true
		flag := true
		for _, b := range B {
			if b == a {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, a)
		}
	}
	return result
}

func Distinct[T any](arr []T) []T {
	var newArr []T
	for i := 0; i < len(arr); i++ {
		var exist bool
		for j := 0; j < len(newArr); j++ {
			if ToStr(newArr[j]) == ToStr(arr[i]) {
				exist = true
				break
			}
		}
		if !exist {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

func UnmarshalObj[T any](data any) (T, error) {
	var obj T
	var reader io.Reader
	switch data.(type) {
	case string:
		reader = strings.NewReader(data.(string))
		break
	case []byte:
		reader = bytes.NewBuffer(data.([]byte))
		break
	case json.RawMessage:
		reader = bytes.NewBuffer(data.(json.RawMessage))
		break
	default:
		return obj, errors.New("error data format")
	}
	dec := json.NewDecoder(reader)
	dec.UseNumber()

	err := dec.Decode(&obj)
	return obj, err
}

func Unmarshal[T any](data any) (*T, error) {
	var reader io.Reader
	switch data.(type) {
	case string:
		reader = strings.NewReader(data.(string))
		break
	case []byte:
		reader = bytes.NewBuffer(data.([]byte))
		break
	case json.RawMessage:
		reader = bytes.NewBuffer(data.(json.RawMessage))
		break
	default:
		return nil, errors.New("error data format")
	}
	dec := json.NewDecoder(reader)
	dec.UseNumber()
	var obj T
	err := dec.Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func Marshal(data any) []byte {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logging.Errorf("[UtilMarshal]failed to marshal,data:%+v", data)
	}
	return dataBytes
}

func To[T any](data any) (*T, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logging.Errorf("[UtilMarshal]failed to marshal,data:%+v", data)
		return nil, err
	}
	return Unmarshal[T](dataBytes)
}

func ToArrayString[T any](array []T, sep string) string {
	var strArray []string
	for _, item := range array {
		strArray = append(strArray, fmt.Sprintf("%v", item))
	}
	return strings.Join(strArray, sep)
}

func ToStringArray[T any](array []T, sep string) []string {
	var strArray []string
	for _, item := range array {
		strArray = append(strArray, fmt.Sprintf("%v", item))
	}
	return strArray
}

func Log(params ...any) {
	var (
		ctx    any
		err    any
		format string
		others []any
	)
	for i := 0; i < len(params); i++ {
		switch params[i].(type) {
		case string:
			if format == "" {
				format = params[i].(string)
			} else {
				others = append(others, params[i])
			}
		case context.Context:
			ctx = params[i]
		case error:
			err = params[i]
		case []byte:
			others = append(others, string(params[i].([]byte)))
		case nil:
			others = append(others, params[i])
		default:
			switch reflect.TypeOf(params[i]).Kind() {
			case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map, reflect.Pointer:
				if params[i] != nil {
					others = append(others, string(Marshal(params[i])))
				} else {
					others = append(others, "nil")
				}
			default:
				others = append(others, params[i])
			}
		}
	}
	length := len(regexp.MustCompile(`%v`).FindAllString(format, -1))
	otherLength := len(others)
	if length < otherLength {
		for i := 0; i < otherLength-length; i++ {
			format += " %v"
		}
	}
	if ctx != nil {
		stack := ctx.(context.Context).Value("stack")
		if stack != nil {
			format = fmt.Sprintf("[%v]%v", stack, format)
		}
	}
	if err != nil {
		format += fmt.Sprintf(",error:%v", err)
		logging.Errorf(format, others...)
	} else {
		logging.Infof(format, others...)
	}
}
