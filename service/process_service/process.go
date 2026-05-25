package process_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
)

type AuthResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
	Data    string `json:"data"`
	Code    int    `json:"code"`
}

type ProcessService struct {
	client *http.Client
	mu     sync.Mutex
	token  string
	expire time.Time
}

var defaultProcessService = &ProcessService{client: &http.Client{Timeout: 20 * time.Second}}

func Default() *ProcessService { return defaultProcessService }

func (s *ProcessService) getToken() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.token != "" && time.Now().Before(s.expire.Add(-2*time.Minute)) {
		return s.token, nil
	}
	cfg := setting.GaiaOpenAPISetting

	// 方案1：按文档/示例使用 form-urlencoded body
	form := url.Values{}
	form.Set("grant_type", cfg.GrantType)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("corp_id", cfg.CorpID)
	bodyReq, err := http.NewRequest(http.MethodPost, strings.TrimRight(cfg.AuthURL, "/"), strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	bodyReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	bodyReq.Header.Set("User-Agent", "Mozilla/5.0 ESB/1.0")
	bodyReq.Header.Set("Accept", "application/json")
	bodyResp, bodyBytes, bodyErr := s.doAuthRequest(bodyReq)
	if bodyErr == nil {
		s.token = bodyResp.Data
		util.Log("[GaiaAuth] token acquired by body mode, token=%v, ttl=%v", maskToken(s.token), cfg.TokenTTLSeconds)
		s.expire = time.Now().Add(time.Duration(cfg.TokenTTLSeconds) * time.Second)
		return s.token, nil
	}

	// 方案2：兼容你在 Postman 成功的 query-string 方式
	query := url.Values{}
	query.Set("grant_type", cfg.GrantType)
	query.Set("client_secret", cfg.ClientSecret)
	query.Set("corp_id", cfg.CorpID)
	authURL := strings.TrimRight(cfg.AuthURL, "/") + "?" + query.Encode()
	queryReq, err := http.NewRequest(http.MethodPost, authURL, nil)
	if err != nil {
		return "", err
	}
	queryReq.Header.Set("User-Agent", "Mozilla/5.0 ESB/1.0")
	queryReq.Header.Set("Accept", "application/json")
	queryResp, _, queryErr := s.doAuthRequest(queryReq)
	if queryErr == nil {
		s.token = queryResp.Data
		util.Log("[GaiaAuth] token acquired by query mode, token=%v, ttl=%v", maskToken(s.token), cfg.TokenTTLSeconds)
		s.expire = time.Now().Add(time.Duration(cfg.TokenTTLSeconds) * time.Second)
		return s.token, nil
	}

	return "", fmt.Errorf("auth failed with both body and query modes; body_mode_err=%v; body_mode_resp=%s; query_mode_err=%v", bodyErr, compactBody(bodyBytes), queryErr)
}

func (s *ProcessService) doAuthRequest(req *http.Request) (AuthResponse, []byte, error) {
	util.Log("[GaiaAuth] request method=%v,url=%v", req.Method, req.URL.String())
	resp, err := s.client.Do(req)
	if err != nil {
		return AuthResponse{}, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AuthResponse{}, nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return AuthResponse{}, body, fmt.Errorf("auth http status=%d, body=%s (可能被Gaia网关/WAF拦截)", resp.StatusCode, compactBody(body))
	}
	var authResp AuthResponse
	if err = json.Unmarshal(body, &authResp); err != nil {
		return AuthResponse{}, body, fmt.Errorf("parse auth response error: status=%d, body=%s", resp.StatusCode, compactBody(body))
	}
	if !authResp.Result || authResp.Data == "" {
		return AuthResponse{}, body, fmt.Errorf("auth failed: status=%d, body=%s", resp.StatusCode, compactBody(body))
	}
	return authResp, body, nil
}

func (s *ProcessService) Proxy(method, apiURL string, payload map[string]any) ([]byte, error) {
	token, err := s.getToken()
	if err != nil {
		return nil, err
	}
	targetURL := strings.TrimSpace(apiURL)
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = strings.TrimRight(setting.GaiaOpenAPISetting.BaseURL, "/") + "/" + strings.TrimLeft(apiURL, "/")
	}

	var body io.Reader
	if strings.EqualFold(method, http.MethodGet) {
		u, err := url.Parse(targetURL)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		for k, v := range payload {
			q.Set(k, fmt.Sprintf("%v", v))
		}
		u.RawQuery = q.Encode()
		targetURL = u.String()
	} else {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bodyBytes)
	}

	callOnce := func(authHeader string) (int, []byte, error) {
		req, err := http.NewRequest(strings.ToUpper(method), targetURL, body)
		if err != nil {
			return 0, nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("User-Agent", "PostmanRuntime/7.51.1")
		req.Header.Set("Authorization", authHeader)
		util.Log("[GaiaProxy] request method=%v,url=%v,auth_mode=%v,headers=%v", strings.ToUpper(method), targetURL, authHeader[:minInt(len(authHeader), 10)], dumpHeaders(req.Header))
		resp, err := s.client.Do(req)
		if err != nil {
			return 0, nil, err
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, nil, err
		}
		util.Log("[GaiaProxy] response status=%v, body=%v", resp.StatusCode, compactBody(respBody))
		return resp.StatusCode, respBody, nil
	}

	status, respBody, err := callOnce("Bearer " + token)
	if err != nil {
		return nil, err
	}
	if status >= 200 && status < 300 {
		return respBody, nil
	}
	// 兼容部分网关/接口仅接受原始 token 头
	status2, respBody2, err2 := callOnce(token)
	if err2 == nil && status2 >= 200 && status2 < 300 {
		return respBody2, nil
	}
	if err2 != nil {
		return nil, fmt.Errorf("gaia api failed status=%d, body=%s; retry_raw_token_err=%v", status, compactBody(respBody), err2)
	}
	return nil, fmt.Errorf("gaia api http status=%d, body=%s; retry_raw_token_status=%d, body=%s", status, compactBody(respBody), status2, compactBody(respBody2))
}

func compactBody(body []byte) string {
	text := strings.TrimSpace(string(body))
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\r", "")
	if len(text) > 300 {
		return text[:300] + "..."
	}
	return text
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maskToken(token string) string {
	if len(token) <= 16 {
		return token
	}
	return token[:8] + "..." + token[len(token)-8:]
}

func dumpHeaders(h http.Header) string {
	m := map[string]string{}
	for k, v := range h {
		if len(v) == 0 {
			continue
		}
		if strings.EqualFold(k, "Authorization") {
			m[k] = maskToken(v[0])
		} else {
			m[k] = v[0]
		}
	}
	b, _ := json.Marshal(m)
	return string(b)
}
