package player

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Arian-p1/spb/src/database"
	"github.com/Arian-p1/spb/src/objectstorage"
	"github.com/Arian-p1/spb/src/security"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type songSync struct {
	ID       int   `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	PlayList string `json:"playlist"`
	File     bool   `json:"file"`
}

func Listen(c *gin.Context) {
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

	ctx := context.Background()
	slink, err := objectstorage.Get(ctx, song.ObjName)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	if ctx.Err() != nil {
		c.JSON(http.StatusConflict, gin.H{"err": ctx.Err})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=file-name.txt")
	c.Header("application/octet-stream", slink.Path)
	c.Status(http.StatusOK)
}

func RemoveSong(c *gin.Context) {
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
	ctx := context.Background()
	err = objectstorage.Remove(ctx, song.ObjName)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	if ctx.Err() != nil {
		c.JSON(http.StatusConflict, gin.H{"err": ctx.Err})
		return
	}
	err = database.RemoveSong(uint(id))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
}

func SyncSong(c *gin.Context) {
	var ssong songSync
	err := c.ShouldBindBodyWith(&ssong, binding.JSON)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		return
	}
	ctx := context.Background()
	if ssong.ID == -1 {
		obj, err := database.SyncSong(false, func(song *database.Song) {
			song.Name = ssong.Name
			song.Artist = ssong.Artist
			song.PlayList = ssong.PlayList
		})
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
		fheader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
		file, err := fheader.Open()
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
		err = objectstorage.Upload(ctx, obj, file, *fheader)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
    if ctx.Err() != nil {
      c.JSON(http.StatusConflict, gin.H{"err": ctx.Err})
      return
    }
	} else {
		err = database.RemoveSong(uint(ssong.ID))
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
		obj, err := database.SyncSong(true, func(song *database.Song) {
		song.ID = uint(ssong.ID)
			song.Name = ssong.Name
			song.Artist = ssong.Artist
			song.PlayList = ssong.PlayList
		})
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
    if fheader, err := c.FormFile("file"); err == nil {
      file, err := fheader.Open()
      if err != nil {
        c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
        return
      }
      err = objectstorage.Upload(ctx, obj, file, *fheader)
      if err != nil {
        c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
        return
      }
      if ctx.Err() != nil {
        c.JSON(http.StatusConflict, gin.H{"err": ctx.Err})
        return
      }
    }
	}
}

func GetAll(c *gin.Context) {
	t := c.Query("type")
	f := c.Query("from")
	uid := uint(0)
	if f == "me" {
		var err error
		uid, err = security.IdFromJWT(c)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
			return
		}
	}
	switch t {
	case "song":
		if songs, err := database.FindAllSongs(uid); err == nil {
			c.JSON(http.StatusOK, songs)
		} else {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		}
	case "playlist":
		if playlists, err := database.FindAllPL(uid); err == nil {
			c.JSON(http.StatusOK, playlists)
		} else {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		}
	default:
		c.JSON(http.StatusConflict, gin.H{"err": "shit"})

	}
}

func Search(c *gin.Context) {
	t := c.Query("type")
	n := c.Query("name")
	switch t {
	case "song":
		var songs []database.Song
		if err := database.Search(songs, n); err == nil {
			c.JSON(http.StatusOK, songs)
		} else {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		}
	case "playlist":
		var playlists []database.PlayList
		if err := database.Search(playlists, n); err == nil {
			c.JSON(http.StatusOK, playlists)
		} else {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		}
	default:
		c.JSON(http.StatusConflict, gin.H{"err": "shit"})

	}
}
