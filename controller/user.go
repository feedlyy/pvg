package controller

import (
	"context"
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

	err = c.ShouldBind(&usr)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

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

func (u *UserController) Update(c *gin.Context) {
	var (
		err  error
		resp = helper.Response{
			Status: "Success",
		}
		pwd      = ""
		phone    = 0
		id       = c.Param("id")
		birthday = time.Time{}
		usr      = domain.Users{
			Firstname: c.PostForm("firstname"),
			Lastname:  c.PostForm("lastname"),
		}
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), u.Timeout)
	defer cancel()

	usrID, err := strconv.Atoi(id)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}
	usr.ID = uint(usrID)

	switch {
	case c.PostForm("birthday") != "":
		loc, _ := time.LoadLocation("Local")
		birthday, err = time.ParseInLocation(helper.BirthdayLayout, c.PostForm("birthday"), loc)
		if err != nil {
			resp.Status = "Failed"
			resp.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}
		usr.Birthday = birthday
	case c.PostForm("phone") != "":
		phone, err = strconv.Atoi(c.PostForm("phone"))
		if err != nil {
			resp.Status = "Failed"
			resp.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}
		usr.Phone = uint(phone)
	case c.PostForm("password") != "":
		pwd, err = helper.HashPassword(c.PostForm("password"))
		if err != nil {
			resp.Status = "Failed"
			resp.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}
		usr.Password = pwd
	default:
	}

	err = u.UserServ.UpdateUser(ctx, usr)
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

	c.JSON(http.StatusOK, resp)
}

func (u *UserController) DeleteUser(c *gin.Context) {
	var (
		err  error
		resp = helper.Response{
			Status: "Success",
		}
		id = c.Param("id")
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), u.Timeout)
	defer cancel()

	usrID, err := strconv.Atoi(id)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	err = u.UserServ.DeleteUser(ctx, usrID)
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

	c.JSON(http.StatusOK, resp)
}
