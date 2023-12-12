package user

import (
	"errors"
	"net/http"

	"github.com/Arian-p1/spb/database"
	"github.com/Arian-p1/spb/helper"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Context = *gin.Context

type SignReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c Context) {
	var reqj SignReq
	if err := c.ShouldBindBodyWith(&reqj, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := helper.EmailValidator(reqj.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if database.UserExsist(reqj.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user exist"})
		return
	}
	p, _ := bcrypt.GenerateFromPassword([]byte(reqj.Password), 8)
	uid := uint(uuid.New().ID())
	token, err := database.GenerateJWT(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please try Again"})
		return
	}
	u := database.UserModel{
		ID:           uid,
		Username:     "",
		Email:        reqj.Email,
		Token:        token,
		Bio:          "",
		PasswordHash: string(p),
	}
	if err := database.AddUser(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"err": "fuuuuuuck"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": "added", "token": token})
}

func Login(c Context) {
	var reqj SignReq
	if err := c.ShouldBindBodyWith(&reqj, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if database.UserExsist(reqj.Email) {
		uid, err := database.PassCheck(reqj.Email, reqj.Password)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Status": err.Error()})
			return
		}
		token, err := database.GenerateJWT(uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Status": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Status": "Ok", "Token": token})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("user doesent exist")})
		return
	}
}

func Profile(c Context) {
	c.JSON(http.StatusOK, gin.H{"yaay": "yooo"})
}

func UpdateProfile(c Context) {
	uid, err := IdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
		return
	}
	var ureq UserUpdate
	if err := c.ShouldBindBodyWith(&ureq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = database.UpdateProfile(uid, func(user *database.UserModel) {
		if ureq.Username != "" {
			user.Username = ureq.Username
		}
		if ureq.Email != "" {
			user.Email = ureq.Email
		}
		if ureq.Bio != "" {
			user.Bio = ureq.Bio
		}
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "updated"})
}

type changepasswd struct {
	Old string `json:"oldpasswd"`
	New string `json:"newpasswd"`
}

func ChangePasswd(c Context) {
	var change changepasswd
	err := c.ShouldBindBodyWith(&change, binding.JSON)
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, gin.H{"err": err.Error()})
	}
	id, err := IdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, gin.H{"err": err.Error()})
		return
	}
	err = database.ChangePasswd(id, change.Old, change.New)
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusPreconditionFailed, gin.H{"status": "password changed"})
}
