package db

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// ID       string `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	IsAdmin  bool   `gorm:"default:false" json:"isAdmin"`
}

func (u *User) CreateAdmin() error {
	user := User{
		Email:    os.Getenv("ADMIN_USER"),
		Password: os.Getenv("ADMIN_PASSWORD"),
		IsAdmin:  true,
	}

	// Hash the password
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error hashing password: %w", err)
	}

	user.Password = string(password)

	// create the user in the db
	if err := DBConn.Create(&user).Error; err != nil {
		return fmt.Errorf("Error creating admin: %w", err)
	}

	return nil
}

func (u *User) LoginAsAdmin(email, password string) (*User, error) {
	if err := DBConn.Where("email = ? AND is_admin = ?", email, true).First(&u).Error; err != nil {
		return nil, fmt.Errorf("Error finding user: %w", err)
	}

	// compare the password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("Invalid password: %w", err)
	}

	return u, nil
}
