package player

import (
	"net/http"
	"strconv"

	"github.com/Arian-p1/spb/database"
	"github.com/Arian-p1/spb/user"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type syncPlaylist struct {
	Name   string          `json:"name"`
	Privet bool            `json:"privet"`
	Songs  []database.Song `json:"songs"`
}

func CreatPlaylist(c user.Context) {
	var pl syncPlaylist
	err := c.ShouldBindBodyWith(pl, binding.JSON)
	if err != nil {
		c.JSON(http.StatusConflict, err.Error())
	}
	id := uint(uuid.New().ID())
	if err != nil {
		c.JSON(http.StatusConflict, err.Error())
	}
	userid, err := user.IdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusConflict, err.Error())
	}
	err = database.AddPlaylist(database.PlayList{ID: id, UserID: userid, Name: pl.Name, Privet: pl.Privet, Songs: pl.Songs})
	if err != nil {
		c.JSON(http.StatusConflict, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func RemovePlayList(c user.Context) {
	id, err := strconv.ParseUint(c.Query("PLID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	err = database.RemovePlaylist(uint(id))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "removed")
}

func UpdatePlayList(c user.Context) {
	id, err := strconv.ParseUint(c.Query("PLID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	var npl syncPlaylist
	err = c.ShouldBindBodyWith(npl, binding.JSON)
	if err != nil {
		c.JSON(http.StatusConflict, err.Error())
	}
	err = database.UpdatePlayList(uint(id), func(pl *database.PlayList) {
		pl.Name = npl.Name
		pl.Privet = npl.Privet
		pl.Songs = npl.Songs
	})
	c.JSON(http.StatusOK, "updated")
}

func LikeSong(c user.Context) {
	sid, err := strconv.ParseUint(c.Query("song_id"), 10, 64)
	if err != nil {
		c.Status(400)
		return
	}
	uid, err := user.IdFromJWT(c)
	if err != nil {
		c.Status(400)
		return
	}
	err = database.AddSongPlayList(uid, uint(sid))
	if err != nil {
		c.Status(400)
		return
	}
	c.Status(200)
}

// func UnLike(c user.Context) {
// 	sid, err := strconv.ParseUint(c.Query("song_id"), 10, 64)
// 	if err != nil {
// 		c.Status(400)
// 		return
// 	}
// 	uid, err := user.IdFromJWT(c)
// 	if err != nil {
// 		c.Status(400)
// 		return
// 	}
// 	err = database.Re(uid, uint(sid))
// 	if err != nil {
// 		c.Status(400)
// 		return
// 	}
// 	c.Status(200)
// }
