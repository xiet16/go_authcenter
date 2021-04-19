package dao

import (
	"github.com/xiet16/authcenter/common/lib"
	"log"
)

type User struct {
	ID int `gorm:"primary_key" json:"id"`
	Name string ` json:"user_name" gorm:"column:user_name"`
	Password string `json:"user_pwd" gorm:"column:user_pwd"`
	CreateTime int `json:"create_time" gorm:"column:create_time"`
	UpdateTime int `json:"update_time" gorm:"column:update_time"`
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

func (u *User)Create(user *User) error {
    tx,err :=lib.GetGormPool("default")
	if err!=nil {
		log.Println(err)
		return err
	}

	tx.Create(user)
	return nil
}

func(u *User)Update(user *User) error {
	tx ,err := lib.GetGormPool("default")
	if err!=nil {
		log.Println(err)
		return err
	}

	tx.Model(&user).Updates(user)
	return err
}