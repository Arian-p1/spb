package main

import (
  "net/http"

	"github.com/Arian-p1/spb/database"
	"github.com/Arian-p1/spb/player"
	"github.com/Arian-p1/spb/user"
	"github.com/Arian-p1/spb/templates"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
    panic(err)
  }
	err = database.DatabaseConnection()
	if err != nil {
    panic(err)
  }
	err = database.Migration()
	if err != nil {
    panic(err)
  }

	engine := gin.Default()
  a := engine.HTMLRender
  engine.HTMLRender = &templates.HTMLTemplRenderer{FallbackHtmlRenderer: a}

	engine.ForwardedByClientIP = true
	engine.SetTrustedProxies([]string{"127.0.0.1"})
  engine.StaticFile("./static/output.css", "./static/output.css")
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", templates.Hello())
	})
	engine.POST("/register", user.Register)
	engine.POST("/login", user.Login)
	profileG := engine.Group("/profile", user.ValidateJWT)
	playerG := engine.Group("/player", user.ValidateJWT)

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
	engine.Run(":1234")
}
