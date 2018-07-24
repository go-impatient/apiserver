package service

import (
	"sync"
	"github.com/moocss/apiserver/src/model"
	"github.com/moocss/apiserver/src/pkg/auth"
	validator "gopkg.in/go-playground/validator.v9"
)

// User service
var User = &userService{
	mutex: &sync.Mutex{},
}

type userService struct {
	mutex *sync.Mutex
}

const (
	pageSize = 20
)

func (srv *userService) CreateUser(user *model.UserModel) error {
	srv.mutex.Lock()
	defer  srv.mutex.Unlock()

	tx := DB.Self.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (srv *userService) DeleteUser(id uint64) error  {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	user := model.UserModel{}
	user.ID = id

	tx := DB.Self.Begin()
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (srv *userService) UpdateUser(user *model.UserModel) error  {
	srv.mutex.Lock()
	defer  srv.mutex.Unlock()

	tx := DB.Self.Begin()
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()

		return  err
	}
	tx.Commit()

	return nil
}

func (srv *userService) GetUser(id uint64) *model.UserModel  {
	u := &model.UserModel{}

	if err := DB.Self.First(&u, id).Error; err != nil {
		return nil
	}

	return u
}

func (srv *userService) GetUserByName(username string) *model.UserModel  {
	u := &model.UserModel{}
	
	if err := DB.Self.Where("`username` = ?", username).First(&u).Error; err != nil {
		return nil
	}
	return u
}

func GetUserList(username string, page int) ([]*model.UserModel, uint64, error) {
	offset := (page - 1) * pageSize
	var count uint64

	users := make([]*model.UserModel, 0)

	if err := DB.Self.Model(&model.UserModel{}).
		Select("`id`, `username`, `created_at`").
		Where("`username` LIKE ?", username).
		Order("`id` DESC").
		Count(&count).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {

		return users, count, err
	}

	return users, count, nil
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (srv *userService) Compare(u *model.UserModel, pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}

// Encrypt the user password.
func (srv *userService) Encrypt(u *model.UserModel) (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}

// Validate the fields.
func (srv *userService) Validate(u *model.UserModel) error {
	validate := validator.New()
	return validate.Struct(u)
}
