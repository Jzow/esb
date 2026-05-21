package models

import "gorm.io/gorm"

type Auth struct {
	UserId   string `gorm:"primary_key" json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (auth *Auth) TableName() string {
	return "user"
}

// CheckAuth checks if authentication information exists
func CheckAuth(username, password string) (string, error) {
	var auth Auth
	err := UDB.Select("user_id").Where("status<3").Where(Auth{UserName: username, Password: password}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	return auth.UserId, nil
}
