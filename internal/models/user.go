package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	HashPassword string    `json:"password"`
	IsAdmin      bool      `json:"isAdmin"`
	CreatedAt    time.Time `json:"created_at"`
}

type Post struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdPost"`
	UpdatedAt   time.Time `json:"updatedPost"`
}

//type UsersPosts struct {
//	UserP *User         `json:"user"`
//	Posts map[int]*Post `json:"posts"`
//}

type UsersPosts struct {
	UserP *User         `json:"user"`
	Posts map[int]*Post `json:"posts"`
}

//Posts []*Post `json:"posts"`

// POSTS

func UserConstructor(uName string, uEmail string, uHashPassword string) (*User, error) {
	ecrypt, err := bcrypt.GenerateFromPassword([]byte(uHashPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Name:         uName,
		Email:        uEmail,
		HashPassword: string(ecrypt),
		CreatedAt:    time.Now().UTC(),
	}, nil
}
