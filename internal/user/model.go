package user

import "github.com/google/uuid"

type CreateUser struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	GitUrl     string `json:"git_url"`
	LikedinUrl string `json:"likedin_url"`
	Bio        string `json:"bio"`
}

type User struct {
	CreateUser
	ID uuid.UUID `json:"id"`
}

type ResponseUser struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	GitUrl     string    `json:"git_url"`
	LikedinUrl string    `json:"likedin_url"`
}