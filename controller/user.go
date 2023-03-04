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
			Status: "Success",
		}
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), u.Timeout)
	defer cancel()

	res, err = u.UserServ.GetAllUser(ctx)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

func (u *UserController) GetUser(c *gin.Context) {
	var (
		res  domain.Users
		err  error
		resp = helper.Response{
			Status: "Success",
		}
		usrname = c.Query("username")
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), u.Timeout)
	defer cancel()

	err = helper.ValidateGet(usrname)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	res, err = u.UserServ.GetUser(ctx, usrname)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		switch {
		case err.Error() == "record not found":
			c.AbortWithStatusJSON(http.StatusNotFound, resp)
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}
	}

	resp.Data = res
	c.JSON(http.StatusOK, resp)
}
