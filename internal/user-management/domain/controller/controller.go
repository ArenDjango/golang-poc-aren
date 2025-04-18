package controller

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	entity "user-management/internal/user-management/domain/entities"
	"user-management/internal/user-management/domain/service"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type controller struct {
	userService service.IUserService
	validator   *validator.Validate
}

func NewController(userService service.IUserService) *controller {
	return &controller{
		userService: userService,
		validator:   validator.New(),
	}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Register a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      entity.User  true  "User info"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {string}  string  "Invalid request"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /users [post]
func (c *controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u entity.User
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.validator.Struct(u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUserInfo, err := c.userService.RegisterUser(u, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUserInfo)
}

// GetUsers godoc
// @Summary      List users
// @Description  Get a list of all users
// @Tags         users
// @Produce      json
// @Success      200  {array}   entity.User
// @Failure      500  {string}  string  "Internal server error"
// @Router       /users [get]
func (c *controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieve a single user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  entity.User
// @Failure      400  {string}  string  "Invalid user ID"
// @Failure      404  {string}  string  "User not found"
// @Router       /users/{id} [get]
func (c *controller) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.userService.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update an existing user's information
// @Tags         users
// @Accept       json
// @Param        id    path      int          true  "User ID"
// @Param        user  body      entity.User  true  "User data"
// @Success      200   {string}  string       "OK"
// @Failure      400   {string}  string       "Invalid input"
// @Failure      500   {string}  string       "Internal server error"
// @Router       /users/{id} [put]
func (c *controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.validator.Struct(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = id

	if err := c.userService.UpdateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Remove a user by ID
// @Tags         users
// @Param        id   path      int  true  "User ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {string}  string  "Invalid user ID"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /users/{id} [delete]
func (c *controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := c.userService.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
