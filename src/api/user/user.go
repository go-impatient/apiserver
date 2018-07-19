package user

import (
	"github.com/moocss/apiserver/src/model"
	"github.com/moocss/apiserver/src/service"
	"github.com/gin-gonic/gin"
	"github.com/moocss/apiserver/src/pkg/errno"
	"github.com/moocss/apiserver/src/util"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	"strconv"
)

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Username string `json:"username"`
}

type ListRequest struct {
	Username string `json:"username"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

type ListResponse struct {
	TotalCount uint64            `json:"totalCount"`
	UserList   []*model.UserResult `json:"userList"`
}

type SwaggerListResponse struct {
	TotalCount uint64           `json:"totalCount"`
	UserList   []model.UserResult `json:"userList"`
}

// @Summary Get an user by the user identifier
// @Description Get an user by username
// @Tags user
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Success 200 {object} model.UserModel "{"code":0,"message":"OK","data":{"username":"kong","password":"$2a$10$E0kwtmtLZbwW/bDQ8qI8e.eHPqhQOW9tvjwpyo/p05f/f4Qvr3OmS"}}"
// @Router /user/{username} [get]
func Get(c *gin.Context) {
	username := c.Param("username")
	// Get the user by the `username` from the database.
	user :=  service.User.GetUserByName(username)
	if user != nil {
		util.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	util.SendResponse(c, nil, user)
}


// @Summary Add new user to the database
// @Description Add a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body user.CreateRequest true "Create a new user"
// @Success 200 {object} user.CreateResponse "{"code":0,"message":"OK","data":{"username":"kong"}}"
// @Router /user [post]
func Create(c *gin.Context) {
	log.Info("User Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		util.SendResponse(c, errno.ErrBind, nil)
		return
	}

	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	// Validate the data.
	if err := service.User.Validate(&u); err != nil {
		util.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// Encrypt the user password.
	if err := service.User.Encrypt(&u); err != nil {
		util.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}
	// Insert the user to the database.
	//if err := service.User.CreateUser(&u); err != nil {
	//	util.SendResponse(c, errno.ErrDatabase, nil)
	//	return
	//}

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information.
	util.SendResponse(c, nil, rsp)
}


// @Summary Delete an user by the user identifier
// @Description Delete user by ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path uint64 true "The user's database id index num"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /user/{id} [delete]
func Delete(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	if err := service.User.DeleteUser(uint64(userId)); err != nil {
		util.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	util.SendResponse(c, nil, nil)
}


// @Summary Update a user info by the user identifier
// @Description Update a user by ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path uint64 true "The user's database id index num"
// @Param user body model.UserModel true "The user info"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /user/{id} [put]
func Update(c *gin.Context) {
	log.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	userId, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var u model.UserModel
	if err := c.Bind(&u); err != nil {
		util.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// We update the record based on the user id.
	u.ID = uint64(userId)

	// Validate the data.
	if err := service.User.Validate(&u); err != nil {
		util.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// Encrypt the user password.
	if err := service.User.Encrypt(&u); err != nil {
		util.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	// Save changed fields.
	if err :=  service.User.UpdateUser(&u); err != nil {
		util.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	util.SendResponse(c, nil, nil)
}
