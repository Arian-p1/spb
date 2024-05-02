package user

import "github.com/gin-gonic/gin"

func Handler(rg *gin.RouterGroup) {
  rg.GET("/", Profile)
	rg.POST("/change-password", ChangePasswd)
	rg.PUT("/update", UpdateProfile)
}
