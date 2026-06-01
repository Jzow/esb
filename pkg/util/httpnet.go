package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func trimBodyForErr(body []byte) string {
	s := strings.TrimSpace(string(body))
	if len(s) > 256 {
		return s[:256]
	}
	return s
}

func validBody(param interface{}) bool {
	if param == nil {
		return true
	}
	switch param.(type) {
	case string, []byte:
		return true
	}
	switch reflect.ValueOf(param).Kind() {
	case reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		return true
	}
	return false
}
func validModel(param interface{}) bool {
	return param == nil || reflect.ValueOf(param).Kind() == reflect.Ptr || reflect.ValueOf(param).Kind() == reflect.Slice
}

// time.Duration类型为超时时长，第一个Map[string]string或map[string]interface{}为headers，其它类型为响应model
func parseParam(params ...interface{}) (int, map[string]string, time.Duration) {
	timeout := 15 * time.Second
	var headers map[string]string
	modelIndex := -1
	for i, param := range params {
		switch param.(type) {
		case time.Duration:
			timeout = param.(time.Duration)
			continue
		case map[string]string:
			if headers == nil || len(headers) == 0 {
				headers = param.(map[string]string)
				continue
			}
		case map[string]interface{}:
			if headers == nil || len(headers) == 0 {
				headers = make(map[string]string)
				for key, val := range param.(map[string]interface{}) {
					headers[key] = fmt.Sprintf("%v", val)
				}
				continue
			}
		default:
			modelIndex = i
		}
	}
	if modelIndex == -1 {
		return -1, headers, timeout
	}
	return modelIndex, headers, timeout
}

func Request(method, url string, body interface{}, model interface{}, headers map[string]string, timeout time.Duration) ([]byte, error) {
	if method == "" {
		return nil, errors.New("empty method")
	}
	if url == "" {
		return nil, errors.New("empty url")
	}
	if !validBody(body) {
		return nil, errors.New("not support the body param type")
	}
	if !validModel(model) {
		return nil, errors.New("not support the model param type")
	}

	client := &http.Client{
		Timeout: timeout,
	}
	var err error
	var reader io.Reader
	switch body.(type) {
	case string:
		reader = strings.NewReader(body.(string))
		break
	case []byte:
		reader = bytes.NewReader(body.([]byte))
		break
	}
	if body != nil {
		switch reflect.ValueOf(body).Kind() {
		case reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("json.Marshal request error,%s", err.Error()))
			}
			reader = bytes.NewReader(bodyBytes)
		}
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("NewRequest method error,%s", err.Error()))
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	//set headers
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("request %s", err.Error()))
	}
	resBytes, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("read response content error,%s", err.Error()))
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("http status=%d, body=%s", res.StatusCode, trimBodyForErr(resBytes)))
	}
	if model != nil {
		err = json.Unmarshal(resBytes, model)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("json.Unmarshal response error,%s, body=%s", err.Error(), trimBodyForErr(resBytes)))
		}
	}
	return resBytes, nil
}

func Get(url string, params ...interface{}) ([]byte, error) {
	modelIndex, headers, timeout := parseParam(params...)
	if modelIndex == -1 {
		return Request("GET", url, nil, nil, headers, timeout)
	}
	return Request("GET", url, nil, params[modelIndex], headers, timeout)
}
func Post(url string, body interface{}, params ...interface{}) ([]byte, error) {
	modelIndex, headers, timeout := parseParam(params...)
	if modelIndex == -1 {
		return Request("POST", url, body, -1, headers, timeout)
	}
	return Request("POST", url, body, params[modelIndex], headers, timeout)
}
func Put(url string, body interface{}, params ...interface{}) ([]byte, error) {
	modelIndex, headers, timeout := parseParam(params...)
	if modelIndex == -1 {
		return Request("PUT", url, body, nil, headers, timeout)
	}
	return Request("PUT", url, body, params[modelIndex], headers, timeout)
}
func Delete(url string, params ...interface{}) ([]byte, error) {
	modelIndex, headers, timeout := parseParam(params...)
	if modelIndex == -1 {
		return Request("DELETE", url, nil, nil, headers, timeout)
	}
	return Request("DELETE", url, nil, params[modelIndex], headers, timeout)
}
