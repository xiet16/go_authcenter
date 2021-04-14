package dao

import (
	"github.com/xiet16/authcenter/lib"
	"log"
)

type User struct {
	ID int `gorm:"primary_key" json:"id"`
	Name string ` json:"user_name" gorm:"column:user_name"`
	Password string `json:"user_pwd" gorm:"column:user_pwd"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) GetUserIDByPwd(search *User) (*User,error) {
	tx ,err:= lib.GetGormPool("default")
	if err!=nil {
		log.Println(err)
		return nil ,err
	}

	out := &User{}
	err = tx.Where(search).Find(out).Error
	if err != nil {
		return nil,err
	}
	return out,nil
}