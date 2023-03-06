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
			Status:  "Success",
			Message: "Please check your email for the activation codes",
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
		switch {
		case err.Error() == helper.DataExists:
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}
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

func (u *UserController) ActivateUser(c *gin.Context) {
	var (
		err  error
		resp = helper.Response{
			Status: "Success",
		}
		usrname        = c.Query("username")
		code           int
		activationCode = domain.ActivationCodes{}
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), u.Timeout)
	defer cancel()

	if c.PostForm("code") != "" {
		code, err = strconv.Atoi(c.PostForm("code"))
		if err != nil {
			resp.Status = "Failed"
			resp.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}
		activationCode.Code = uint(code)
	}

	err = c.ShouldBind(&code)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	err = u.UserServ.ActivateUser(ctx, usrname, code)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		switch {
		case err.Error() == helper.RecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, resp)
			return
		case err.Error() == helper.NoValidCode || err.Error() == helper.DataActive:
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		case err.Error() == helper.NoValid:
			c.AbortWithStatusJSON(http.StatusGone, resp)
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (u *UserController) RequestCodeActivation(c *gin.Context) {
	var (
		err  error
		resp = helper.Response{
			Status: "Success",
		}
		usrname        = c.Query("username")
		activationCode = domain.ActivationCodes{}
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), u.Timeout)
	defer cancel()

	activationCode, err = u.UserServ.RequestActivationCode(ctx, usrname)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		switch {
		case err.Error() == helper.RecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, resp)
			return
		case err.Error() == helper.DataActive:
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		case err.Error() == helper.StillValid:
			resp.Message += fmt.Sprintf(", your code (%v) still active until %v",
				activationCode.Code, activationCode.ExpiresAt.Format("2006-01-02 15:04:05"))
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}
	}

	resp.Message = fmt.Sprintf("This is your activation code: %v", activationCode.Code)
	c.JSON(http.StatusOK, resp)
}
