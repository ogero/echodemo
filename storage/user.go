package storage

import (
	"bitbucket.org/ogero/echodemo/storage/types"
	"github.com/elithrar/simple-scrypt"
	"github.com/jinzhu/gorm"
	"strings"
)

type User struct {
	gorm.Model
	Email    string     `form:"email" gorm:"type:varchar(100);unique_index;NOT NULL"`
	Password string     `form:"password" gorm:"NOT NULL"`
	Role     types.Role `form:"role" gorm:"NOT NULL"`
}

func (o *User) CheckPasswordMatch(password string) bool {
	if scrypt.CompareHashAndPassword([]byte(o.Password), []byte(password)) != nil {
		return false
	}
	return true
}

func (o *User) BeforeSave() (err error) {
	// TODO: DOS attack possible when setting scrypt abusive pass on form and then logging in
	if len(o.Password) > 0 && strings.Count(o.Password, "$") != 4 {
		var pass []byte
		if pass, err = scrypt.GenerateFromPassword([]byte(o.Password), scrypt.DefaultParams); err == nil {
			o.Password = string(pass)
		}
	}
	return
}
