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

func RemoveSong(c user.Context) {
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
	err = database.RemoveSong(uint(id))
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
	errChan := make(chan error)

	go func() {
		file, err := c.FormFile("file")
		if err != nil {
			errChan <- err
			return
		}

		if err := c.SaveUploadedFile(file, path); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	err := <-errChan
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your file has been successfully uploaded."})
	return nil
}

func GetAll(c user.Context) {
	t := c.Query("type")
	f := c.Query("from")
	uid := uint(0)
	if f == "me" {
		var err error
		uid, err = user.IdFromJWT(c)
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

func Search(c user.Context) {
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
