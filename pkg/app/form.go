package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) string {
	err := c.BindJSON(form)
	if err != nil {
		return "failed to parse struct param"
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return "failed to valid param"
	}
	if !check {
		return MarkErrorTip(valid.Errors)
	}

	return ""
}
