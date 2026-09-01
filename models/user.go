package models

import (
	"time"
)

type User struct{
	ID uint `json:"ID" gorm:"primary_key"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password" gorm"not null"`
	Role AuthRole `json:"role" gorm:"type:varchar(20);not null;default: user"`
	CreatedAt time.Time 
	UpdatedAt time.Time
}