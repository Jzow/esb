package util

import (
	"errors"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSecret []byte //old gateway string value: Wz8AuYeknwoevs2qR2o9r6gYcY4ns4JO

const (
	MiniProgramPlatformID = 6
	WebPlatform           = "Web"
	NormalToken           = 0
	uidPidToken           = "UID_PID_TOKEN_STATUS"
)

const HoursOneDay = 24
const secondBefore = 5

type Claims struct {
	UserID     string
	PlatformID int
	jwt.RegisteredClaims
}

const (
	UidPidToken = "UID_PID_TOKEN_STATUS:"
)

func GetServerAuthorization() map[string]any {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
	authorization := BuildClaims("ticket", 10, 1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authorization)
	tokenString, _ := token.SignedString(jwtSecret)
	return map[string]any{"Authorization": tokenString}
}

func GetTokenKey(userID string, platformID int) string {
	return UidPidToken + userID + ":" + PlatformIDToName(platformID)
}
func GetAllPlatformTokenKey(userID string) []string {
	res := make([]string, len(PlatformID2Name))
	for k := range PlatformID2Name {
		res[k-1] = GetTokenKey(userID, k)
	}
	return res
}

type Config struct {
	batchSize       int
	continueOnError bool
	concurrentLimit int
}
type Option func(c *Config)

var PlatformName2ID = map[string]int{
	IOSPlatformStr:        IOSPlatformID,
	AndroidPlatformStr:    AndroidPlatformID,
	WindowsPlatformStr:    WindowsPlatformID,
	OSXPlatformStr:        OSXPlatformID,
	WebPlatformStr:        WebPlatformID,
	MiniWebPlatformStr:    MiniWebPlatformID,
	LinuxPlatformStr:      LinuxPlatformID,
	AndroidPadPlatformStr: AndroidPadPlatformID,
	IPadPlatformStr:       IPadPlatformID,
	AdminPlatformStr:      AdminPlatformID,
}

func BuildClaims(uid string, platformID int, ttl int64) Claims {
	now := time.Now()
	before := now.Add(-time.Second * time.Duration(secondBefore))
	return Claims{
		UserID:     uid,
		PlatformID: platformID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttl*HoursOneDay) * time.Hour)), // Expiration time
			IssuedAt:  jwt.NewNumericDate(before),                                              // Issuing time
			//NotBefore: jwt.NewNumericDate(before),                                              // Begin Effective time
		},
	}
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	claims, err := GetClaimFromToken(token, Secret(string(jwtSecret)))
	if err != nil {
		return nil, err
	}
	if claims != nil {
		return claims, nil
	}
	return nil, err
}
func Secret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}

func GetClaimFromToken(tokensString string, secretFunc jwt.Keyfunc) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secretFunc)
	if err == nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.New("claims unknown")
	}
	if ve, ok := err.(*jwt.ValidationError); ok {
		return nil, ve
	}

	return nil, errors.New("jwt parse error")
}

func PlatformIDToName(num int) string {
	return PlatformID2Name[num]
}

const (
	// Platform ID.
	IOSPlatformID        = 1
	AndroidPlatformID    = 2
	WindowsPlatformID    = 3
	OSXPlatformID        = 4
	WebPlatformID        = 5
	MiniWebPlatformID    = 6
	LinuxPlatformID      = 7
	AndroidPadPlatformID = 8
	IPadPlatformID       = 9
	AdminPlatformID      = 10

	// Platform string match to Platform ID.
	IOSPlatformStr        = "IOS"
	AndroidPlatformStr    = "Android"
	WindowsPlatformStr    = "Windows"
	OSXPlatformStr        = "OSX"
	WebPlatformStr        = "Web"
	MiniWebPlatformStr    = "MiniWeb"
	LinuxPlatformStr      = "Linux"
	AndroidPadPlatformStr = "APad"
	IPadPlatformStr       = "IPad"
	AdminPlatformStr      = "Admin"

	// terminal types.
	TerminalPC     = "PC"
	TerminalMobile = "Mobile"
	TerminalPad    = "Pad"
)

var PlatformID2Name = map[int]string{
	IOSPlatformID:        IOSPlatformStr,
	AndroidPlatformID:    AndroidPlatformStr,
	WindowsPlatformID:    WindowsPlatformStr,
	OSXPlatformID:        OSXPlatformStr,
	WebPlatformID:        WebPlatformStr,
	MiniWebPlatformID:    MiniWebPlatformStr,
	LinuxPlatformID:      LinuxPlatformStr,
	AndroidPadPlatformID: AndroidPadPlatformStr,
	IPadPlatformID:       IPadPlatformStr,
	AdminPlatformID:      AdminPlatformStr,
}
