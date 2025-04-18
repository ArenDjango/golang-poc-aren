package repository

import (
	"context"
	"user-management/internal/user-management/domain"
	entity "user-management/internal/user-management/domain/entities"
	"user-management/internal/user-management/infrastructure/model"

	"github.com/uptrace/bun"
)

type userRepo struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) domain.IUserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user entity.User) (int64, error) {
	u := entity.FromEntity(user)
	_, err := r.db.NewInsert().Model(&u).Exec(context.Background())
	return u.ID, err
}

func (r *userRepo) GetAll() ([]entity.User, error) {
	var users []model.User
	err := r.db.NewSelect().Model(&users).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	var result []entity.User
	for _, u := range users {
		result = append(result, entity.ToEntity(u))
	}
	return result, nil
}

func (r *userRepo) GetByID(id int64) (entity.User, error) {
	var user model.User
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return entity.User{}, err
	}
	return entity.ToEntity(user), nil
}

func (r *userRepo) Update(user entity.User) error {
	u := entity.FromEntity(user)
	_, err := r.db.NewUpdate().Model(&u).Where("id = ?", u.ID).Exec(context.Background())
	return err
}

func (r *userRepo) Delete(id int64) error {
	_, err := r.db.NewDelete().Model(&model.User{}).Where("id = ?", id).Exec(context.Background())
	return err
}
