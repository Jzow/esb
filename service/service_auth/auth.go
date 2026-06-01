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
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("empty token")
	}
	if isAppToken(token) {
		return &Principal{Type: "app", Token: token}, nil
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
	if token == "" || isAppToken(token) {
		return nil
	}
	if util.Rdb == nil {
		return errors.New("redis is not initialized")
	}
	return util.Rdb.Del(context.Background(), userTokenKey(token)).Err()
}

func isAppToken(token string) bool {
	for _, configured := range strings.Split(setting.ServiceAuthSetting.ApiTokens, ",") {
		configured = strings.TrimSpace(configured)
		if configured == "" {
			continue
		}
		if matchSecret(token, configured) {
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
