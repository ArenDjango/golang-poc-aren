package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user-management/internal/user-management/domain"
	entity "user-management/internal/user-management/domain/entities"
)

type IUserService interface {
	RegisterUser(user entity.User, ip string) (map[string]interface{}, error)
	ListUsers() ([]entity.User, error)
	GetUserByID(id int64) (entity.User, error)
	UpdateUser(user entity.User) error
	DeleteUser(id int64) error
}

type userService struct {
	repo            domain.IUserRepository
	userGeoApiToken string
}

func NewUserService(r domain.IUserRepository, userGeoApiToken string) IUserService {
	return &userService{repo: r, userGeoApiToken: userGeoApiToken}
}

func (s *userService) RegisterUser(user entity.User, ip string) (map[string]interface{}, error) {
	ipInfo, err := s.getIPInfo(ip)
	if err == nil {
		fmt.Print("ipinfo", ipInfo)
		// user.Metadata["location"] = ipInfo
	}
	id, err := s.repo.Create(user)
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

func (s *userService) getIPInfo(ip string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s?token=%s", ip, s.userGeoApiToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	fmt.Print(resp)
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
