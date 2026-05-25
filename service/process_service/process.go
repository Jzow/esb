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
	form := url.Values{}
	form.Set("grant_type", cfg.GrantType)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("corp_id", cfg.CorpID)

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(cfg.AuthURL, "/"), strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("auth http status=%d, body=%s", resp.StatusCode, compactBody(body))
	}
	var authResp AuthResponse
	if err = json.Unmarshal(body, &authResp); err != nil {
		return "", fmt.Errorf("parse auth response error: status=%d, body=%s", resp.StatusCode, compactBody(body))
	}
	if !authResp.Result || authResp.Data == "" {
		return "", fmt.Errorf("auth failed: status=%d, body=%s", resp.StatusCode, compactBody(body))
	}
	s.token = authResp.Data
	s.expire = time.Now().Add(time.Duration(cfg.TokenTTLSeconds) * time.Second)
	return s.token, nil
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

	req, err := http.NewRequest(strings.ToUpper(method), targetURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
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
