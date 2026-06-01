package service_auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/go-redis/redis/v8"
)

const userTokenPrefix = "esb:auth:user_token:"

type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	Username  string `json:"username"`
}

type Principal struct {
	Type     string
	Username string
	Token    string
	Apps     []string
}

func Login(username, password string) (*LoginResponse, error) {
	if username == "" || password == "" {
		return nil, errors.New("empty username or password")
	}
	if !matchSecret(username, setting.ServiceAuthSetting.Username) ||
		!matchSecret(password, setting.ServiceAuthSetting.Password) {
		return nil, errors.New("invalid username or password")
	}
	if util.Rdb == nil {
		return nil, errors.New("redis is not initialized")
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	ttl := time.Duration(setting.ServiceAuthSetting.TokenTTLSeconds) * time.Second
	if err := util.Rdb.Set(context.Background(), userTokenKey(token), username, ttl).Err(); err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: setting.ServiceAuthSetting.TokenTTLSeconds,
		Username:  username,
	}, nil
}

func Authenticate(token string) (*Principal, error) {
	return AuthenticateForApp(token, "")
}

func AuthenticateForApp(token, appName string) (*Principal, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("empty token")
	}
	if apps, ok := appTokenApps(token); ok {
		if appName != "" && !allowApp(apps, appName) {
			return nil, errors.New("token is not allowed to access openapi app")
		}
		return &Principal{Type: "app", Token: token, Apps: apps}, nil
	}
	if util.Rdb == nil {
		return nil, errors.New("redis is not initialized")
	}

	username, err := util.Rdb.Get(context.Background(), userTokenKey(token)).Result()
	if err == redis.Nil {
		return nil, errors.New("invalid or expired token")
	}
	if err != nil {
		return nil, err
	}
	return &Principal{Type: "user", Username: username, Token: token}, nil
}

func Logout(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil
	}
	if _, ok := appTokenApps(token); ok {
		return nil
	}
	if util.Rdb == nil {
		return errors.New("redis is not initialized")
	}
	return util.Rdb.Del(context.Background(), userTokenKey(token)).Err()
}

func appTokenApps(token string) ([]string, bool) {
	for _, configured := range strings.Split(setting.ServiceAuthSetting.ApiTokens, ",") {
		configured = strings.TrimSpace(configured)
		if configured == "" {
			continue
		}
		configuredToken, apps := parseAppToken(configured)
		if matchSecret(token, configuredToken) {
			return apps, true
		}
	}
	return nil, false
}

func parseAppToken(configured string) (string, []string) {
	if parts := strings.SplitN(configured, "=", 2); len(parts) == 2 {
		return strings.TrimSpace(parts[1]), splitApps(parts[0])
	}
	if parts := strings.SplitN(configured, ":", 2); len(parts) == 2 {
		return strings.TrimSpace(parts[0]), splitApps(parts[1])
	}
	return configured, []string{"*"}
}

func splitApps(apps string) []string {
	var out []string
	for _, app := range strings.FieldsFunc(apps, func(r rune) bool {
		return r == '|' || r == ';' || r == ','
	}) {
		app = strings.TrimSpace(app)
		if app != "" {
			out = append(out, app)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func allowApp(apps []string, appName string) bool {
	for _, app := range apps {
		if app == "*" || strings.EqualFold(app, appName) {
			return true
		}
	}
	return false
}

func matchSecret(actual, expected string) bool {
	return subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) == 1
}

func userTokenKey(token string) string {
	return userTokenPrefix + token
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token failed: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
