package domain

import entity "user-management/internal/user-management/domain/entities"

type IUserRepository interface {
	Create(user entity.User) (int64, error)
	GetAll() ([]entity.User, error)
	GetByID(id int64) (entity.User, error)
	Update(user entity.User) error
	Delete(id int64) error
}
