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
	profileG := r.Group("/profile", user.ValidateJWT)
	playerG := r.Group("/player", user.ValidateJWT)

	profileG.GET("/", user.ValidateJWT, user.Profile)
	profileG.POST("/change-password", user.ValidateJWT, user.ChangePasswd)
	profileG.PUT("/update", user.ValidateJWT, user.UpdateProfile)

	playerG.POST("/listen", user.ValidateJWT, player.Listen)
	playerG.DELETE("/delete", user.ValidateJWT, player.RemoveSong)
	playerG.POST("/updatesong", user.ValidateJWT, player.SyncSong)
	playerG.POST("/search", user.ValidateJWT, player.Search)
	playerG.POST("/list", user.ValidateJWT, player.GetAll)
	playerG.POST("/addpl", user.ValidateJWT, player.CreatPlaylist)
	playerG.DELETE("/rmpl", user.ValidateJWT, player.RemovePlayList)
	playerG.POST("/likesong", user.ValidateJWT, player.LikeSong)
	r.Run(":1234")
}
