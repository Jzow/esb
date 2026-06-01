package openapi_proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	mu sync.Mutex
}

type oauthResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
	Data    string `json:"data"`
	Code    int    `json:"code"`
}

type Response struct {
	StatusCode int
	Data       interface{}
	RawBody    []byte
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Proxy(method, proxyPath, rawQuery string, body []byte, contentType string) (*Response, error) {
	openAPI, ok, err := setting.FindOpenAPIByURL(proxyPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("unknown openapi url: %s", proxyPath)
	}
	targetURL, err := buildTargetURL(openAPI, proxyPath, rawQuery)
	if err != nil {
		return nil, err
	}
	headers := commonHeaders(openAPI)
	if contentType != "" {
		headers["Content-Type"] = contentType
	}
	token, err := c.authToken(openAPI)
	if err != nil {
		return nil, err
	}
	if token != "" {
		headers["Authorization"] = authHeaderValue(openAPI, token)
	}

	req, err := http.NewRequest(method, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	httpClient := &http.Client{Timeout: 15 * time.Second}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("openapi http status=%d, body=%s", res.StatusCode, trimBodyForErr(resBytes))
	}

	data, err := normalizeResponse(resBytes)
	if err != nil {
		return nil, err
	}
	return &Response{StatusCode: res.StatusCode, Data: data, RawBody: resBytes}, nil
}

func (c *Client) authToken(openAPI *setting.OpenAPI) (string, error) {
	if strings.TrimSpace(openAPI.AuthPath) == "" && strings.TrimSpace(openAPI.AuthURL) == "" {
		return "", nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	ttl := normalizedTokenTTL(openAPI)
	if util.Rdb == nil {
		return "", errors.New("redis is not initialized")
	}
	if token, ok := c.redisAuthToken(openAPI); ok {
		return token, nil
	}

	values := url.Values{}
	grantType := strings.TrimSpace(openAPI.GrantType)
	if grantType == "" {
		grantType = "client_credentials"
	}
	values.Set("grant_type", grantType)
	values.Set("client_secret", openAPI.ClientSecret)
	values.Set("corp_id", openAPI.CorpID)

	authURL := openAPI.AuthPath
	if authURL == "" {
		authURL = openAPI.AuthURL
	}
	targetAuthURL, err := buildAbsoluteURL(openAPI, authURL)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, targetAuthURL, strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	for k, v := range commonHeaders(openAPI) {
		req.Header.Set(k, v)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: 15 * time.Second}
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("openapi auth status=%d, body=%s", res.StatusCode, trimBodyForErr(resBytes))
	}

	var authResp oauthResponse
	if err := json.Unmarshal(resBytes, &authResp); err != nil {
		return "", fmt.Errorf("openapi auth json error=%s, body=%s", err.Error(), trimBodyForErr(resBytes))
	}
	if !authResp.Result || authResp.Data == "" {
		return "", fmt.Errorf("openapi auth failed: code=%d, message=%s", authResp.Code, authResp.Message)
	}
	if err := c.cacheAuthToken(openAPI, authResp.Data, ttl); err != nil {
		return "", err
	}
	return authResp.Data, nil
}

func (c *Client) redisAuthToken(openAPI *setting.OpenAPI) (string, bool) {
	if util.Rdb == nil {
		return "", false
	}
	ctx := context.Background()
	key := openAPITokenKey(openAPI.Name)
	token, err := util.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false
	}
	if err != nil {
		logging.Errorf("openapi token redis get key=%s err=%s", key, err.Error())
		return "", false
	}
	return token, token != ""
}

func (c *Client) cacheAuthToken(openAPI *setting.OpenAPI, token string, ttl time.Duration) error {
	if util.Rdb == nil {
		return errors.New("redis is not initialized")
	}
	key := openAPITokenKey(openAPI.Name)
	if err := util.Rdb.Set(context.Background(), key, token, ttl).Err(); err != nil {
		logging.Errorf("openapi token redis set key=%s err=%s", key, err.Error())
		return err
	}
	return nil
}

func normalizedTokenTTL(openAPI *setting.OpenAPI) time.Duration {
	ttl := openAPI.TokenTTL
	if ttl <= 0 {
		ttl = 2 * time.Hour
	}
	return ttl
}

func openAPITokenKey(appName string) string {
	return "esb:openapi:token:" + strings.TrimSpace(appName)
}

func buildTargetURL(openAPI *setting.OpenAPI, proxyPath, rawQuery string) (string, error) {
	if strings.TrimSpace(proxyPath) == "" || proxyPath == "/" {
		return "", errors.New("empty proxy url")
	}
	path := strings.TrimLeft(proxyPath, "/")
	path = strings.ReplaceAll(path, "{tenant}", url.PathEscape(openAPI.CorpID))
	if !isAbsoluteURL(path) {
		return "", errors.New("proxy url must be a full http or https url")
	}
	if err := validateTargetURL(openAPI, path); err != nil {
		return "", err
	}
	if rawQuery != "" {
		if strings.Contains(path, "?") {
			path += "&" + rawQuery
		} else {
			path += "?" + rawQuery
		}
	}
	return path, nil
}

func buildAbsoluteURL(openAPI *setting.OpenAPI, path string) (string, error) {
	if isAbsoluteURL(path) {
		return path, nil
	}
	if strings.TrimSpace(openAPI.BaseUrl) == "" {
		return "", fmt.Errorf("empty base url for openapi app: %s", openAPI.Name)
	}
	if strings.HasPrefix(path, "/") {
		return openAPI.BaseUrl + path, nil
	}
	return openAPI.BaseUrl + "/" + path, nil
}

func isAbsoluteURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func validateTargetURL(openAPI *setting.OpenAPI, target string) error {
	if strings.TrimSpace(openAPI.BaseUrl) == "" {
		return nil
	}
	targetURL, err := url.Parse(target)
	if err != nil {
		return err
	}
	baseURL, err := url.Parse(openAPI.BaseUrl)
	if err != nil {
		return err
	}
	if !strings.EqualFold(targetURL.Scheme, baseURL.Scheme) || !strings.EqualFold(targetURL.Host, baseURL.Host) {
		return fmt.Errorf("target url host %s is not allowed for openapi app %s", targetURL.Host, openAPI.Name)
	}
	return nil
}

func commonHeaders(openAPI *setting.OpenAPI) map[string]string {
	headers := map[string]string{}
	if openAPI.UserAgent != "" {
		headers["User-Agent"] = openAPI.UserAgent
	}
	if openAPI.Origin != "" {
		headers["Origin"] = openAPI.Origin
	}
	if openAPI.Referer != "" {
		headers["Referer"] = openAPI.Referer
	}
	for k, v := range fixedHeaders(openAPI) {
		headers[k] = v
	}
	return headers
}

func fixedHeaders(openAPI *setting.OpenAPI) map[string]string {
	headers := map[string]string{}
	for _, item := range strings.Split(openAPI.FixedHeaders, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := replaceVariables(openAPI, strings.TrimSpace(parts[1]))
		if key != "" {
			headers[key] = value
		}
	}
	return headers
}

func replaceVariables(openAPI *setting.OpenAPI, value string) string {
	value = strings.ReplaceAll(value, "{CorpID}", openAPI.CorpID)
	value = strings.ReplaceAll(value, "{corpId}", openAPI.CorpID)
	value = strings.ReplaceAll(value, "{tenant}", openAPI.CorpID)
	return value
}

func authHeaderValue(openAPI *setting.OpenAPI, token string) string {
	prefix := strings.TrimSpace(openAPI.TokenPrefix)
	if prefix == "" {
		prefix = "Bearer"
	}
	return prefix + " " + token
}

func normalizeResponse(body []byte) (interface{}, error) {
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, nil
	}
	var out map[string]interface{}
	if err := json.Unmarshal(body, &out); err != nil {
		return string(body), nil
	}
	if code, ok := numericCode(out["code"]); ok && code != 200 {
		return nil, fmt.Errorf("openapi failed: code=%d, message=%s, reason=%s", code, stringValue(out["message"]), stringValue(out["reason"]))
	}
	if result, ok := out["result"].(bool); ok && !result {
		return nil, fmt.Errorf("openapi failed: result=false, message=%s, reason=%s", stringValue(out["message"]), stringValue(out["reason"]))
	}
	if details, ok := out["details"]; ok {
		return details, nil
	}
	if data, ok := out["data"]; ok {
		return data, nil
	}
	return out, nil
}

func numericCode(v interface{}) (int, bool) {
	switch code := v.(type) {
	case int:
		return code, true
	case int64:
		return int(code), true
	case float64:
		return int(code), true
	case json.Number:
		n, err := code.Int64()
		return int(n), err == nil
	default:
		return 0, false
	}
}

func stringValue(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func trimBodyForErr(body []byte) string {
	s := strings.TrimSpace(string(body))
	if len(s) > 512 {
		return s[:512]
	}
	return s
}
