package gaia_service

import (
	"encoding/json"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/models/gaia"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	tokenLock   sync.RWMutex
	cachedToken string
	tokenExpire time.Time
)

type Client struct{}

func buildURL(path string) string {
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

func (c *Client) authToken() (string, error) {
	tokenLock.RLock()
	if cachedToken != "" && time.Now().Before(tokenExpire) {
		t := cachedToken
		tokenLock.RUnlock()
		return t, nil
	}
	tokenLock.RUnlock()

	tokenLock.Lock()
	defer tokenLock.Unlock()
	if cachedToken != "" && time.Now().Before(tokenExpire) {
		return cachedToken, nil
	}
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_secret", setting.GaiaApiSetting.ClientSecret)
	values.Set("corp_id", setting.GaiaApiSetting.CorpID)

	var resp gaia.OAuthResponse
	authURL := setting.GaiaApiSetting.AuthPath
	if authURL == "" {
		authURL = setting.GaiaApiSetting.AuthURL
	}
	_, err := util.Post(buildURL(authURL), values.Encode(), &resp, map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if err != nil {
		return "", err
	}
	if !resp.Result || resp.Data == "" {
		return "", fmt.Errorf("gaia auth failed: %s", resp.Message)
	}
	cachedToken = resp.Data
	tokenExpire = time.Now().Add(setting.GaiaApiSetting.TokenTTL - time.Minute)
	return cachedToken, nil
}

func (c *Client) post(path string, req interface{}) (map[string]interface{}, error) {
	tk, err := c.authToken()
	if err != nil {
		return nil, err
	}
	headers := map[string]string{"Authorization": "Bearer " + tk}
	var out map[string]interface{}
	_, err = util.Post(buildURL(path), req, &out, headers)
	if err != nil {
		logging.Errorf("gaia post path=%s err=%s", path, err.Error())
		return nil, err
	}
	return out, nil
}
func (c *Client) Post(path string, req interface{}) (map[string]interface{}, error) {
	return c.post(path, req)
}

func (c *Client) get(path string, query interface{}) (map[string]interface{}, error) {
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
	_, err = util.Get(api, &out, map[string]string{"Authorization": "Bearer " + tk})
	return out, err
}
func (c *Client) Get(path string, query interface{}) (map[string]interface{}, error) {
	return c.get(path, query)
}

func (c *Client) PostLeaveSubmit(req gaia.LeaveSubmitRequest) (map[string]interface{}, error) {
	return c.post(setting.GaiaApiSetting.LeaveSubmitPath, req)
}

// LeaveSubmit keeps compatibility with older call sites.
func (c *Client) LeaveSubmit(req gaia.LeaveSubmitRequest) (map[string]interface{}, error) {
	return c.PostLeaveSubmit(req)
}
