package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	entity "user-management/internal/user-management/domain/entities"
	"user-management/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mockIPServer(t *testing.T, responseBody string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	}))
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	mockRepo.On("Create", mock.Anything).Return(int64(1), nil)

	ipServer := mockIPServer(t, `{"city":"TestCity","country":"TC"}`, 200)
	defer ipServer.Close()

	mockClient := mocks.NewIPInfoClient(t)
	mockClient.On("GetInfo", "1.1.1.1").Return(map[string]interface{}{"city": "Test"}, nil)

	svc := &userService{repo: mockRepo, ipInfoClient: mockClient}

	user := entity.User{Name: "Aren", Email: "aren@example.com"}
	result, err := svc.RegisterUser(user, "1.1.1.1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result["usedId"])
}

func TestRegisterUser_IPInfoFails(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	mockRepo.On("Create", mock.Anything).Return(int64(2), nil)

	mockClient := mocks.NewIPInfoClient(t)
	mockClient.On("GetInfo", "8.8.8.8").Return(nil, errors.New("ipinfo down"))

	svc := &userService{repo: mockRepo, ipInfoClient: mockClient}

	user := entity.User{Name: "Test", Email: "t@x.com"}
	result, err := svc.RegisterUser(user, "8.8.8.8")

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result["usedId"])
}

func TestRegisterUser_RepoFails(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	mockRepo.On("Create", mock.Anything).Return(int64(0), errors.New("db error"))

	mockClient := mocks.NewIPInfoClient(t)
	mockClient.On("GetInfo", mock.Anything).Return(map[string]interface{}{"city": "Test"}, nil)
	svc := &userService{repo: mockRepo, ipInfoClient: mockClient}

	user := entity.User{Name: "Fail", Email: "f@x.com"}
	_, err := svc.RegisterUser(user, "9.9.9.9")
	assert.Error(t, err)
}

func TestListUsers_Success(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	users := []entity.User{{ID: 1, Name: "Test"}}
	mockRepo.On("GetAll").Return(users, nil)

	svc := &userService{repo: mockRepo}
	out, err := svc.ListUsers()
	assert.NoError(t, err)
	assert.Equal(t, users, out)
}

func TestListUsers_Error(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	mockRepo.On("GetAll").Return(nil, errors.New("db fail"))

	svc := &userService{repo: mockRepo}
	_, err := svc.ListUsers()
	assert.Error(t, err)
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	u := entity.User{ID: 2, Name: "A"}
	mockRepo.On("GetByID", int64(2)).Return(u, nil)

	svc := &userService{repo: mockRepo}
	out, err := svc.GetUserByID(2)
	assert.NoError(t, err)
	assert.Equal(t, u, out)
}

func TestGetUserByID_Error(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	mockRepo.On("GetByID", int64(9)).Return(entity.User{}, errors.New("not found"))

	svc := &userService{repo: mockRepo}
	_, err := svc.GetUserByID(9)
	assert.Error(t, err)
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	u := entity.User{ID: 3, Name: "U"}
	mockRepo.On("Update", u).Return(nil)

	svc := &userService{repo: mockRepo}
	err := svc.UpdateUser(u)
	assert.NoError(t, err)
}

func TestUpdateUser_Error(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	u := entity.User{ID: 3, Name: "Bad"}
	mockRepo.On("Update", u).Return(errors.New("update fail"))

	svc := &userService{repo: mockRepo}
	err := svc.UpdateUser(u)
	assert.Error(t, err)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	mockRepo.On("Delete", int64(4)).Return(nil)

	svc := &userService{repo: mockRepo}
	err := svc.DeleteUser(4)
	assert.NoError(t, err)
}

func TestDeleteUser_Error(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	mockRepo.On("Delete", int64(7)).Return(errors.New("delete fail"))

	svc := &userService{repo: mockRepo}
	err := svc.DeleteUser(7)
	assert.Error(t, err)
}
