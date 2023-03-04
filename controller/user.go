package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"pvg/domain"
	"pvg/helper"
	"time"
)

type UserController struct {
	UserServ domain.UserService
	Timeout  time.Duration
}

func NewUserController(u domain.UserService, timeout time.Duration) UserController {
	return UserController{
		UserServ: u,
		Timeout:  timeout,
	}
}

func (u *UserController) GetUsers(c *gin.Context) {
	var (
		res  []domain.Users
		err  error
		resp = helper.Response{
			Message: "Success",
		}
	)
	ctx, cancel := context.WithTimeout(c.Request.Context(), u.Timeout)
	defer cancel()

	res, err = u.UserServ.GetAllUser(ctx)
	if err != nil {
		resp.Message = "Failed"
		resp.Error = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = res
	c.JSON(http.StatusOK, resp)
}
