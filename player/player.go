package player

import (
	"net/http"
	"os"
	"strconv"

	"github.com/Arian-p1/spb/database"
	"github.com/Arian-p1/spb/user"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type songSync struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	PlayList string `json:"playlist"`
	File     bool   `json:"file"`
}

func Listen(c user.Context) {
	id, err := strconv.ParseUint(c.Query("SongID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	song, err := database.FindSongById(uint(id))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	file, err := os.ReadFile(song.Path)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=file-name.txt")
	c.Data(http.StatusOK, "application/octet-stream", file)
}

func DeleteSong(c user.Context) {
	id, err := strconv.ParseUint(c.Query("SongID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	song, err := database.FindSongById(uint(id))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	err = os.Remove(song.Path)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	err = database.DeleteSong(1)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
}

func SyncSong(c user.Context) {
	var nsong songSync
	err := c.ShouldBindBodyWith(&nsong, binding.JSON)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	if nsong.ID == 0 {
		path, err := database.SyncSong(func(song *database.Song) {
			song.Name = nsong.Name
			song.Artist = nsong.Artist
			song.PlayList = nsong.PlayList
		})
		if nsong.File {
			if err = HandleFile(c, path); err == nil {
				c.JSON(http.StatusOK, gin.H{"song": "uploaded"})
			}
		}
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
	} else {
		path, err := database.SyncSong(func(song *database.Song) {
			song.Name = nsong.Name
			song.Artist = nsong.Artist
			song.PlayList = nsong.PlayList
		}, nsong.ID)
		if nsong.File {
			if err = HandleFile(c, path); err == nil {
				c.JSON(http.StatusOK, gin.H{"song": "uploaded"})
			}
		}
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
	}
}

func HandleFile(c user.Context, path string) error {
	file, _ := c.FormFile("file")
	if err := c.SaveUploadedFile(file, path); err != nil {
		return err
	}
	c.JSON(http.StatusOK, gin.H{"song": "uploaded"})
	return nil
}
