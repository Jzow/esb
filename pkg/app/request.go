package app

import (
	"github.com/astaxie/beego/validation"
	"strings"

	"github.com/EDDYCJY/go-gin-example/pkg/logging"
)

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Infof(err.Key, err.Message)
	}

	return
}

func MarkErrorTip(errors []*validation.Error) string {
	if len(errors) == 0 {
		return "valid error"
	}
	var tip []string
	for _, err := range errors {
		tip = append(tip, err.Key+" "+err.Message)
	}
	return strings.Join(tip, ";")
}
