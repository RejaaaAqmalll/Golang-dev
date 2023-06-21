package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// struct  tabel users
type Users struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique;" json:"email"`
	Password string `json:"password"`
	Code     string `json:"code_verify" gorm:"unique"`
}

type PersonalInfo struct {
	gorm.Model
	FullName       string    `json:"full_name" gorm:"type:varchar(255)"`
	Email          string    `json:"email" gorm:"type:varchar(255);unique"`
	Phone          string    `json:"phone" gorm:"type:varchar(13)"`
	PlaceOfBirth   string    `json:"place_of_birth" gorm:"type:varchar(255)"`
	DateOfBirth    time.Time `json:"date_of_birth" gorm:"type:Date"`
	MaritalStatus  string    `json:"marital_status" gorm:"type:enum('Married', 'Single', 'Widower', 'Widow')"`
	Religion       string    `json:"religion" gorm:"type:enum('Islam', 'Christian', 'Hindu', 'Buddha', 'Catholic', 'Konghucu')"`
	Gender         string    `json:"gender" gorm:"type:enum('Laki-laki', 'Perempuan')"`
	IdentityId     int       `json:"identity_id" gorm:"type:int"`
	IdentityType   Identity  `json:"identity_name" gorm:"foreignkey:IdentityId"`
	IdentityNumber int       `json:"identity_number" gorm:"type:varchar(16)"`
	Address        string    `json:"address" gorm:"type:varchar(255)"`
	Code           string    `json:"code" gorm:"type:varchar(255)"`
	Password       string    `json:"password" gorm:"type:varchar(255)"`
}

type Identity struct {
	gorm.Model
	IdentityName string `json:"identity_name"`
}

type JobLevel struct {
	gorm.Model
	JobLevelName string `json:"job_level_name" gorm:"type:varchar(255)"`
}

type JobPosition struct {
	gorm.Model
	JobPositionName string `json:"job_position_name" gorm:"type:varchar(255)"`
}

type Division struct {
	gorm.Model
	DivisionName string `json:"division_name" gorm:"type:varchar(255)"`
}
