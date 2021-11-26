package controller

import (
	"errors"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"my_scaffold/dao"
	"my_scaffold/dto"
	"my_scaffold/middleware"
	"strings"
	"time"
)

type ApiController struct {
}

func ApiRegister(router *gin.RouterGroup) {
	curd := ApiController{}
	router.POST("/login", curd.Login)
	router.GET("/loginout", curd.LoginOut)
}

func ApiLoginRegister(router *gin.RouterGroup) {
	curd := ApiController{}
	router.GET("/user/listpage", curd.ListPage)
	router.POST("/user/add", curd.AddUser)
	router.POST("/user/edit", curd.EditUser)
	router.POST("/user/remove", curd.RemoveUser)
	router.POST("/user/batchremove", curd.RemoveUser)
}

func (demo *ApiController) Login(c *gin.Context) {
	api := &dto.LoginInput{}
	if err := api.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	maps := fmt.Sprintf(" name = '%v' and password='%v' ", api.Username, api.Password)
	users ,err := (&dao.User{}).FindUserByName(c, tx, maps)

	if err != nil {
		middleware.ResponseError(c, 2002, errors.New("账号或密码错误"))
		return
	}

	token,err := GenerateToken(users)
	if err != nil {
		middleware.ResponseError(c, 2002, errors.New(err.Error()))
		return
	}

	type UserLogin struct {
		Token string
	}
	data := UserLogin{
		Token: token,
	}
	middleware.ResponseSuccess(c, data)
	return
}

func (demo *ApiController) LoginOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	session.Save()
	return
}

func (demo *ApiController) ListPage(c *gin.Context) {
	listInput := &dto.ListPageInput{}
	if err := listInput.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	if listInput.PageSize == 0 {
		listInput.PageSize = 10
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	userList, total, err := (&dao.User{}).PageList(c, tx, listInput)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	m := &dao.ListPageOutput{
		List:  userList,
		Total: total,
	}
	middleware.ResponseSuccess(c, m)
	return
}

func (demo *ApiController) AddUser(c *gin.Context) {
	addInput := &dto.AddUserInput{}
	if err := addInput.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	user := &dao.User{
		Name:  addInput.Name,
		Sex:   addInput.Sex,
		Age:   addInput.Age,
		Birth: addInput.Birth,
		Addr:  addInput.Addr,
	}
	if err := user.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

func (demo *ApiController) EditUser(c *gin.Context) {
	editInput := &dto.EditUserInput{}
	if err := editInput.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	user, err := (&dao.User{}).Find(c, tx, int64(editInput.Id))
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	user.Name = editInput.Name
	user.Sex = editInput.Sex
	user.Age = editInput.Age
	user.Birth = editInput.Birth
	user.Addr = editInput.Addr
	if err := user.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

func (demo *ApiController) RemoveUser(c *gin.Context) {
	removeInput := &dto.RemoveUserInput{}
	if err := removeInput.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	if err := (&dao.User{}).Del(c, tx, strings.Split(removeInput.IDS, ",")); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}


func GenerateToken(user dao.User) (string, error){
	j := &middleware.JWT{
		[]byte("newtrekWang"),
	}

	claims := middleware.CustomClaims{
		user.Id,
		user.Name,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 签名过期时间
			Issuer:    "newtrekWang",                    // 签名颁发者
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}