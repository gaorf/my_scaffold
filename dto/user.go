package dto

import (
	"github.com/gin-gonic/gin"
	"my_scaffold/public"
)

type InfoUserInput struct {
	Id int64 `form:"id" json:"id" comment:"ID" validate:"required"`
}

func (params *InfoUserInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
