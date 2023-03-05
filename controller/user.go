package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pvg/domain"
	"pvg/helper"
	"strconv"
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

func (u *UserController) GetUserByUsername(c *gin.Context) {
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
		case err.Error() == helper.RecordNotFound:
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

func (u *UserController) Create(c *gin.Context) {
	var (
		err  error
		resp = helper.Response{
			Status: "Success",
		}
		pwd      = ""
		phone    = 0
		birthday = time.Time{}
		usr      = domain.Users{
			Username:  c.PostForm("username"),
			Firstname: c.PostForm("firstname"),
			Lastname:  c.PostForm("lastname"),
			Email:     c.PostForm("email"),
			Status:    helper.Inactive,
			CreatedAt: time.Now(),
		}
	)
	//c.Header("Content-Type", "application/x-www-form-urlencoded")

	ctx, cancel := context.WithTimeout(c.Request.Context(), u.Timeout)
	defer cancel()

	loc, _ := time.LoadLocation("Local")
	birthday, err = time.ParseInLocation(helper.BirthdayLayout, c.PostForm("birthday"), loc)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}
	fmt.Println("birthday nya:", birthday)

	err = c.ShouldBind(&usr)
	if err != nil {
		fmt.Println("err nya dari sini", err.Error())
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	//loc, _ := time.LoadLocation("Local")
	//birthday, err = time.ParseInLocation(helper.BirthdayLayout, c.PostForm("birthday"), loc)
	//fmt.Println("birthday nya:", birthday)
	//birthday, err = time.Parse(helper.BirthdayLayout, c.PostForm("birthday"))

	phone, err = strconv.Atoi(c.PostForm("phone"))
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	pwd, err = helper.HashPassword(c.PostForm("password"))
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}
	usr.Password = pwd
	usr.Phone = uint(phone)
	usr.Birthday = birthday

	err = u.UserServ.CreateUser(ctx, usr)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
