package service

import (
	"user-management/internal/user-management/domain"
	entity "user-management/internal/user-management/domain/entities"

	"github.com/rs/zerolog/log"
)

type IUserService interface {
	RegisterUser(user entity.User, ip string) (map[string]interface{}, error)
	ListUsers() ([]entity.User, error)
	GetUserByID(id int64) (entity.User, error)
	UpdateUser(user entity.User) error
	DeleteUser(id int64) error
}

type userService struct {
	repo         domain.IUserRepository
	ipInfoClient IPInfoClient
}

func NewUserService(r domain.IUserRepository, ipInfoClient IPInfoClient) IUserService {
	return &userService{repo: r, ipInfoClient: ipInfoClient}
}

func (s *userService) RegisterUser(user entity.User, ip string) (map[string]interface{}, error) {
	ipInfo, err := s.ipInfoClient.GetInfo(ip)
	if err != nil {
		log.Error().Msgf("RegisterUser error getting Geo API")
		ipInfo = map[string]interface{}{
			"Couldn't call the geo API": "true",
		}
	}
	id, err := s.repo.Create(user)
	if err != nil {
		return map[string]interface{}{}, err
	}
	ipInfo["usedId"] = id
	return ipInfo, nil
}

func (s *userService) ListUsers() ([]entity.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetUserByID(id int64) (entity.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) UpdateUser(user entity.User) error {
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id int64) error {
	return s.repo.Delete(id)
}
