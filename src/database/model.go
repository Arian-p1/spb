package database

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Username     string `gorm:"column:username"`
	Email        string `gorm:"column:email;unique_index"`
	Bio          string `gorm:"column:bio;size:2024"`
	PasswordHash string `gorm:"column:password;not null"`
}

type Song struct {
	gorm.Model
	UserID   uint   `gorm:"column:userid;not null"`
	Name     string `gorm:"column:name;not null"`
	Artist   string `gorm:"column:artist"`
	PlayList string `gorm:"column:playlist"`
	ObjName  string `gorm:"column:objname"`
}

type PlayList struct {
	gorm.Model
	UserID uint   `gorm:"column:userid;not null"`
	Name   string `gorm:"column:name"`
	Privet bool   `gorm:"column:privet"`
	Songs  []Song `gorm:"foreignKey:ID"`
}
