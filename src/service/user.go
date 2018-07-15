package model

import (
	"github.com/moocss/apiserver/src/service"
	"github.com/moocss/apiserver/src/schema"
)

type User struct {
	db *service.Database
}

func (m *User) TableName() string {
	return "tb_users"
}

func (m *User) Create(user *schema.UserModel) error {
	return m.db.Self.Create(&user).Error
}

func (m *User) DeleteUser(id uint64) error  {
	user := schema.UserModel{}
	user.ID = id
	return m.db.Self.Delete(&user).Error
}

func (m *User) Update(user *schema.UserModel) error  {
	return m.db.Self.Save(user).Error
}

func (m *User) Find(id uint64) (*schema.UserModel, error)  {

}