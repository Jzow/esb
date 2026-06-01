package gaia_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/EDDYCJY/go-gin-example/models/gaia"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/go-redis/redis/v8"
)

var (
	tokenLock sync.Mutex
)

type Client struct{}

func gaiaCommonHeaders() map[string]string {
	h := map[string]string{}
	if ua := strings.TrimSpace(setting.GaiaApiSetting.UserAgent); ua != "" {
		h["User-Agent"] = ua
	}
	if o := strings.TrimSpace(setting.GaiaApiSetting.Origin); o != "" {
		h["Origin"] = o
	}
	if r := strings.TrimSpace(setting.GaiaApiSetting.Referer); r != "" {
		h["Referer"] = r
	}
	return h
}

func authHeaderValue(token string) string {
	prefix := strings.TrimSpace(setting.GaiaApiSetting.TokenPrefix)
	if prefix == "" {
		prefix = "Bearer"
	}
	return prefix + " " + token
}

func buildURL(path string) string {
	path = resolvePath(path)
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	base := strings.TrimRight(setting.GaiaApiSetting.BaseUrl, "/")
	if base == "" {
		return path
	}
	if strings.HasPrefix(path, "/") {
		return base + path
	}
	return base + "/" + path
}

func resolvePath(path string) string {
	tenant := strings.TrimSpace(setting.GaiaApiSetting.CorpID)
	path = strings.ReplaceAll(path, "{tenant}", url.PathEscape(tenant))
	return path
}

func (c *Client) authToken() (string, error) {
	tokenLock.Lock()
	defer tokenLock.Unlock()
	if util.Rdb == nil {
		return "", errors.New("redis is not initialized")
	}
	if token, ok := redisAuthToken(); ok {
		return token, nil
	}

	values := url.Values{}
	grantType := strings.TrimSpace(setting.GaiaApiSetting.GrantType)
	if grantType == "" {
		grantType = "client_credentials"
	}
	values.Set("grant_type", grantType)
	values.Set("client_secret", setting.GaiaApiSetting.ClientSecret)
	values.Set("corp_id", setting.GaiaApiSetting.CorpID)

	var resp gaia.OAuthResponse
	authURL := setting.GaiaApiSetting.AuthPath
	if authURL == "" {
		authURL = setting.GaiaApiSetting.AuthURL
	}
	headers := gaiaCommonHeaders()
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	_, err := util.Post(buildURL(authURL), values.Encode(), &resp, headers)
	if err != nil {
		logging.Errorf("gaia auth url=%s err=%s", buildURL(authURL), err.Error())
		return "", err
	}
	if !resp.Result || resp.Data == "" {
		return "", fmt.Errorf("gaia auth failed: %s", resp.Message)
	}
	if err := cacheAuthToken(resp.Data, normalizedTokenTTL()); err != nil {
		return "", err
	}
	return resp.Data, nil
}

func redisAuthToken() (string, bool) {
	if util.Rdb == nil {
		return "", false
	}
	token, err := util.Rdb.Get(context.Background(), gaiaTokenKey()).Result()
	if err == redis.Nil {
		return "", false
	}
	if err != nil {
		logging.Errorf("gaia token redis get key=%s err=%s", gaiaTokenKey(), err.Error())
		return "", false
	}
	return token, token != ""
}

func cacheAuthToken(token string, ttl time.Duration) error {
	if util.Rdb == nil {
		return errors.New("redis is not initialized")
	}
	if err := util.Rdb.Set(context.Background(), gaiaTokenKey(), token, ttl).Err(); err != nil {
		logging.Errorf("gaia token redis set key=%s err=%s", gaiaTokenKey(), err.Error())
		return err
	}
	return nil
}

func normalizedTokenTTL() time.Duration {
	ttl := setting.GaiaApiSetting.TokenTTL
	if ttl <= 0 {
		return 2 * time.Hour
	}
	return ttl
}

func gaiaTokenKey() string {
	return "esb:openapi:token:gaiastandard"
}

func normalizeResponse(out map[string]interface{}) (interface{}, error) {
	if code, ok := numericCode(out["code"]); ok && code != 200 {
		return nil, fmt.Errorf("gaia api failed: code=%d, message=%s, reason=%s", code, stringValue(out["message"]), stringValue(out["reason"]))
	}
	if result, ok := out["result"].(bool); ok && !result {
		return nil, fmt.Errorf("gaia api failed: result=false, message=%s, reason=%s", stringValue(out["message"]), stringValue(out["reason"]))
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

func (c *Client) post(path string, req interface{}) (interface{}, error) {
	if strings.TrimSpace(path) == "" || strings.Contains(path, "<") {
		return nil, errors.New("empty or invalid gaia api path")
	}
	tk, err := c.authToken()
	if err != nil {
		return nil, err
	}
	headers := gaiaCommonHeaders()
	headers["Authorization"] = authHeaderValue(tk)
	var out map[string]interface{}
	_, err = util.Post(buildURL(path), req, &out, headers)
	if err != nil {
		logging.Errorf("gaia post path=%s err=%s", path, err.Error())
		return nil, err
	}
	return normalizeResponse(out)
}
func (c *Client) Post(path string, req interface{}) (interface{}, error) {
	return c.post(path, req)
}

func (c *Client) get(path string, query interface{}) (interface{}, error) {
	if strings.TrimSpace(path) == "" || strings.Contains(path, "<") {
		return nil, errors.New("empty or invalid gaia api path")
	}
	tk, err := c.authToken()
	if err != nil {
		return nil, err
	}
	qBytes, _ := json.Marshal(query)
	var qMap map[string]interface{}
	_ = json.Unmarshal(qBytes, &qMap)
	values := url.Values{}
	for k, v := range qMap {
		if fmt.Sprintf("%v", v) != "" {
			values.Set(k, fmt.Sprintf("%v", v))
		}
	}
	api := buildURL(path)
	if values.Encode() != "" {
		api += "?" + values.Encode()
	}
	var out map[string]interface{}
	headers := gaiaCommonHeaders()
	headers["Authorization"] = authHeaderValue(tk)
	_, err = util.Get(api, &out, headers)
	if err != nil {
		logging.Errorf("gaia get url=%s err=%s", api, err.Error())
		return nil, err
	}
	return normalizeResponse(out)
}
func (c *Client) Get(path string, query interface{}) (interface{}, error) {
	return c.get(path, query)
}

func (c *Client) PostLeaveSubmit(req gaia.LeaveSubmitRequest) (interface{}, error) {
	return c.post(setting.GaiaApiSetting.LeaveSubmitPath, req)
}

// LeaveSubmit keeps compatibility with older call sites.
func (c *Client) LeaveSubmit(req gaia.LeaveSubmitRequest) (interface{}, error) {
	return c.PostLeaveSubmit(req)
}
