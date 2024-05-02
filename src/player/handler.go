package player

import "github.com/gin-gonic/gin"

func Handler(rg *gin.RouterGroup) {
	rg.POST("/listen", Listen)
	rg.DELETE("/delete", RemoveSong)
	rg.POST("/updatesong", SyncSong)
	rg.POST("/search", Search)
	rg.POST("/list", GetAll)
	rg.POST("/addpl", CreatPlaylist)
	rg.DELETE("/rmpl", RemovePlayList)
	rg.POST("/likesong", LikeSong)
}
