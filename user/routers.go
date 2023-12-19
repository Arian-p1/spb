package user

import (
	"net/http"

	"github.com/Arian-p1/spb/database"
	"github.com/Arian-p1/spb/helper"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	u := database.UserModel{
		Username:     "",
		Email:        reqj.Email,
		Bio:          "",
		PasswordHash: string(p),
	}
	if err := database.AddUser(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"err": "fuuuuuuck"})
		return
	}

	err := database.AddPlaylist(database.PlayList{UserID: u.ID, Name: "Liked", Privet: true, Songs: []database.Song{}})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := database.GenerateJWT(u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please try Again"})
		return
	}
	c.JSON(201, gin.H{"token": token})
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
			c.JSON(http.StatusOK, gin.H{"Status": "email or password is wrong"})
			return
		}
		token, err := database.GenerateJWT(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is wrong"})
		return
	}
}

func Profile(c Context) {
	id, err := IdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	data := database.GetUser(id)
	c.JSON(http.StatusOK, gin.H{"id": data.ID, "email": data.Email, "username": data.Username, "Bio": data.Bio})
}

func UpdateProfile(c Context) {
	uid, err := IdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		if ureq.Email != "" && helper.EmailValidator(ureq.Email) == nil {
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
	c.Status(http.StatusOK)
}

type changepasswd struct {
	Old string `json:"oldpasswd"`
	New string `json:"newpasswd"`
}

func ChangePasswd(c Context) {
	var change changepasswd
	err := c.ShouldBindBodyWith(&change, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}
	id, err := IdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	err = database.ChangePasswd(id, change.Old, change.New)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if err == nil {
		c.Status(200)
	}
}
