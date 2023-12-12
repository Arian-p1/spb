package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseConnection() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", os.Getenv("DBHOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("DBPORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func Migration() error {
	return DB.AutoMigrate(&UserModel{}, &PlayList{}, &Song{})
}

func AddUser(um *UserModel) error {
	return DB.Create(um).Error
}

func UserExsist(email string) bool {
	var user UserModel
	DB.Find(&user, "email= ?", email).First(&user)
	if user.Email == email {
		return true
	} else {
		return false
	}
}

// returns uid if the password was correct
func PassCheck(email string, pass string) (uint, error) {
	var user UserModel
	DB.Find(&user, "email= ?", email).First(&user)
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pass)) == nil {
		return user.ID, nil
	} else {
		return 0, errors.New("dont brute force mother fucker")
	}
}

func UpdateProfile(uid uint, changes func(user *UserModel)) error {
	var user UserModel
	if err := DB.Find(&user, uid).First(&user).Error; err == nil {
		changes(&user)
		err = DB.Updates(user).Error
		return err
	} else {
		return err
	}
}

func ChangePasswd(uid uint, oldpasswd string, newpasswd string) error {
	var user UserModel
	if err := DB.Find(&user, uid).First(&user).Error; err == nil {
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldpasswd)) == nil {
			p, err := bcrypt.GenerateFromPassword([]byte(newpasswd), 8)
			if err == nil {
				user.PasswordHash = string(p)
				DB.Save(&user)
			}
			return err
		}
		return err
	} else {
		return err
	}
}

func FindSongById(songid uint) (Song, error) {
	var song Song
	if err := DB.Find(&song, songid).First(&song).Error; err != nil {
		return song, err
	}
	return song, nil
}

func DeleteSong(songid uint) error {
	var song Song
	if err := DB.Find(&song, songid).First(&song).Error; err == nil {
		DB.Delete(song)
		return err
	} else {
		return err
	}
}

func SyncSong(changes func(song *Song), songid ...uint) (string, error) {
	if len(songid) == 0 {
		var song Song
		changes(&song)
		song.ID = uint(uuid.New().ID())
		song.Path = "/nfs/" + song.Artist + "/" + song.PlayList + "/" + fmt.Sprint(song.ID) + ".mp3"
		return song.Path, DB.Create(song).Error
	} else {
		song, err := FindSongById(songid[1]) //TODO: make sure index is ok
		if err != nil {
			return "", err
		}
		changes(&song)
		return song.Path, DB.Updates(song).Error
	}
}
