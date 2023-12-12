package database

type UserModel struct {
	ID           uint   `gorm:"primary_key"`
	Username     string `gorm:"column:username"`
	Email        string `gorm:"column:email;unique_index"`
	Token        string `gorm:"column:token;unique_index"`
	Bio          string `gorm:"column:bio;size:1024"`
	PasswordHash string `gorm:"column:password;not null"`
}

type Song struct {
	ID         uint   `gorm:"primary_key"`
	UploadedBy uint   `gorm:"uploaded_by"`
	Name       string `gorm:"column:name;not null"`
	Artist     string `gorm:"column:artist"`
	PlayList   string `gorm:"column:playlist"`
	Path       string `gorm:"column:path"`
}

type PlayList struct {
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"column:name"`
	Privet bool   `gorm:"column:privet"`
	Songs  []Song `gorm:"foreignKey:id"`
}

// for liked Songs we make a default playlist named Liked Songs
