package auth_service

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"time"
)

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (string, error) {
	return models.CheckAuth(a.Username, a.Password)
}

func (a *Auth) CreateTicket() (string, error) {
	ticket := util.GetUUID32("")
	err := util.SetRedisCache(ticket, a.Username, ticket, 86400*time.Second)
	if err != nil {
		logging.Infof("[ticket]failed to set %s's redis ticket, error: %s", a.Username, err.Error())
		return "", err
	}
	return ticket, nil
}
