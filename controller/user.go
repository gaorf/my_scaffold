package controller

import (
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"my_scaffold/dao"
	"my_scaffold/dto"
	"my_scaffold/middleware"
)

type UserController struct {
}


func ApiUserRegister(router *gin.RouterGroup){
	curd := UserController{}
	router.GET("/user/info", curd.Info)
}

func (u *UserController) Info(c *gin.Context) {
	infoInput := &dto.InfoUserInput{}
	if err:= infoInput.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	cc, _ := c.Get("claims")
	bb := cc.(*middleware.CustomClaims)
	fmt.Println(bb.Id)


	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	info, err := (&dao.User{}).Find(c, tx, infoInput.Id)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, info)
	return
}



