package handler

import (
	"fmt"

	"johnpaulkh/go-crud/api/config"
	"johnpaulkh/go-crud/api/model"
	"johnpaulkh/go-crud/api/repository"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type userHandler struct {
	client *mongo.Client
	config *config.Configuration
	repo   repository.UserRepository
}

type UserHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Get(*gin.Context)
	List(*gin.Context)
}

func NewUserHandler(client *mongo.Client, config *config.Configuration, repo repository.UserRepository) UserHandler {
	return &userHandler{
		client: client,
		config: config,
		repo:   repo,
	}
}

// CreateUser		godoc
//
//	@Summary		Create user
//	@Description	Save user data in Db.
//	@Param			user	body	model.User	true	"Create user"
//	@Produce		application/json
//	@Tags			users
//	@Success		200	{object}	model.User{}
//	@Router			/api/v1/users [post]
func (app *userHandler) Create(c *gin.Context) {
	var request model.User

	err := c.BindJSON(&request)
	if err != nil {
		fmt.Printf("error : %v", err)
	}

	fromDB, err := app.repo.Create(request, c)
	if err != nil {
		logrus.Error("Error during create in Handler", err)
	}

	c.JSON(200, fromDB)
}

// GetUserByID godoc
//
//	@Summary		Get User by Id
//	@Description	Get User by Id
//	@Tags			users
//	@Produce		json
//	@Router			/api/v1/users/{userId} [get]
func (app *userHandler) Get(c *gin.Context) {
	id := c.Param("id")

	user, err := app.repo.Get(id, c)
	if err != nil {
		logrus.Error("Error during get in Handler", err)
	}

	c.JSON(200, user)
}

// UpdateUser		godoc
//
//	@Summary		User user
//	@Description	Save user data in Db by Id.
//	@Param			user	body	model.User	true	"Update user"
//	@Produce		application/json
//	@Tags			users
//	@Success		200	{object}	model.User{}
//	@Router			/api/v1/users/{userId} [put]
func (app *userHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var request model.User
	err := c.ShouldBindJSON(&request)
	if err != nil {
		logrus.Error("Error during reading from JSON", err)
	}

	user, err := app.repo.Update(id, request, c)
	if err != nil {
		logrus.Error("Error during put in handler", err)
	}

	c.JSON(200, user)
}

// GetUsersPage godoc
//
//	@Summary		Get Users with Page
//	@Description	Get Users with Page
//	@Tags			users
//	@Produce		json
//	@Router			/api/v1/users [get]
func (app *userHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		logrus.Error("Error during reading from param page", err)
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		logrus.Error("Error during reading from param size", err)
	}

	result, err := app.repo.List(page, size, c)
	if err != nil {
		logrus.Error("Error during listing in handler ", err)
	}

	c.JSON(200, result)
}
