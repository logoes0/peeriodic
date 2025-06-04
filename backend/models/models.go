package models

import (
	"firebase.google.com/go/v4/auth"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthService struct {
	DB       *gorm.DB
	FireAuth *auth.Client
}
