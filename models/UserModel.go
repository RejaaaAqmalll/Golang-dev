package models

import "github.com/jinzhu/gorm"

// struct  tabel users
type Users struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique;" json:"email"`
	Password string `json:"password"`
	Code     string `json:"code_verify" gorm:"unique"`
}
