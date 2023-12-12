package main

import (
	"github.com/Arian-p1/spb/database"
	"github.com/Arian-p1/spb/helper"
	"github.com/Arian-p1/spb/player"
	"github.com/Arian-p1/spb/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	helper.Panic(err)

	err = database.DatabaseConnection()
	helper.Panic(err)

	err = database.Migration()
	helper.Panic(err)

	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Ok Response",
		})
	})
	r.POST("/register", user.Register)
	r.POST("/login", user.Login)

	r.GET("/profile", user.ValidateJWT, user.Profile)
	r.POST("/profile/change-password", user.ValidateJWT, user.ChangePasswd)
	r.PUT("/profile/update", user.ValidateJWT, user.UpdateProfile)

	r.POST("/player/listen", user.ValidateJWT, player.Listen)
	r.POST("/player/delete", user.ValidateJWT, player.DeleteSong)
	r.POST("/player/update", user.ValidateJWT, player.SyncSong)
	r.Run(":1234")
}