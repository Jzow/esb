package util

import (
	"crypto/rand"
	"github.com/google/uuid"
	"math/big"
	"strings"
)

func GetRandomString(length int, kind string) string {
	passwd := make([]rune, length)
	var codeModel []rune
	switch kind {
	case "num":
		codeModel = []rune("0123456789")
	case "char":
		codeModel = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case "mix":
		codeModel = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case "advance":
		codeModel = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+=-!@#$%*[]")
	default:
		codeModel = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	}
	for i := range passwd {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(codeModel))))
		passwd[i] = codeModel[int(index.Int64())]
	}
	return string(passwd)
}

func GetUUID36() string {
	return uuid.New().String()
}

func GetUUID32(replaceStr string) string {
	return strings.ReplaceAll(uuid.New().String(), "-", replaceStr)
}
