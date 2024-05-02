package src

import (
	"github.com/Arian-p1/spb/src/player"
	"github.com/Arian-p1/spb/src/security"
	"github.com/Arian-p1/spb/src/user"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	engine := gin.Default()
  	engine.POST("/register", user.Register)
	engine.POST("/login", user.Login)
	profileGroup := engine.Group("/profile", security.ValidateJWT)
	playerGroup := engine.Group("/player", security.ValidateJWT)
	user.Handler(profileGroup)
	player.Handler(playerGroup)
	return engine
}
