package entity

import "user-management/internal/user-management/infrastructure/model"

type User struct {
	ID    int64  `json:"id" bun:",pk,autoincrement"`
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

func ToEntity(u model.User) User {
	return User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func FromEntity(e User) model.User {
	return model.User{
		ID:    e.ID,
		Name:  e.Name,
		Email: e.Email,
	}
}
