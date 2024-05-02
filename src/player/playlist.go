package player

import (
	"net/http"
	"strconv"

	"github.com/Arian-p1/spb/src/database"
	"github.com/Arian-p1/spb/src/security"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type syncPlaylist struct {
	Name   string          `json:"name"`
	Privet bool            `json:"privet"`
	Songs  []database.Song `json:"songs"`
}

func CreatPlaylist(c *gin.Context) {
	var pl syncPlaylist
	err := c.ShouldBindBodyWith(pl, binding.JSON)
	if err != nil {
		c.JSON(http.StatusConflict, err.Error())
	}

	userid, err := security.IdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusConflict, err.Error())
	}
	err = database.AddPlaylist(database.PlayList{UserID: userid, Name: pl.Name, Privet: pl.Privet, Songs: pl.Songs})
	if err != nil {
		c.JSON(http.StatusConflict, err.Error())
	}
	c.Status(http.StatusOK)
}

func RemovePlayList(c *gin.Context) {
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

func UpdatePlayList(c *gin.Context) {
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

	if err != nil {
		c.JSON(http.StatusExpectationFailed, "failed to update playlist")
	}

	c.JSON(http.StatusOK, "updated")
}

func LikeSong(c *gin.Context) {
	sid, err := strconv.ParseUint(c.Query("song_id"), 10, 64)
	if err != nil {
		c.Status(400)
		return
	}
	uid, err := security.IdFromJWT(c)
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

func UnLike(c *gin.Context) {
	sid, err := strconv.ParseUint(c.Query("song_id"), 10, 64)
	if err != nil {
		c.Status(400)
		return
	}
	uid, err := security.IdFromJWT(c)
	if err != nil {
		c.Status(400)
		return
	}
	err = database.RemoveSongPlaylist(uid, uint(sid))
	if err != nil {
		c.Status(400)
		return
	}
	c.Status(200)
}
