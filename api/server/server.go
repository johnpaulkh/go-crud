package server

import (
	"context"
	"fmt"

	"johnpaulkh/go-crud/api/config"
	"johnpaulkh/go-crud/api/handler"
	"johnpaulkh/go-crud/api/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Initialize(config config.Configuration) {

	// Create a new instance of the logger. You can have any number of instances.
	var log = logrus.New()

	log.WithFields(logrus.Fields{
		"mongo_url":   config.Database.Url,
		"server_port": config.Server.Port,
		"db_name":     config.Database.DbName,
		"collection":  config.Database.Collection,
	}).Info("Configuration informations")

	logrus.Infof("Application Name %s is starting....", config.App.Name)

	router := gin.Default()

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Database.Url))
	repository := repository.NewUserRepository(client, &config)
	userHandler := handler.NewUserHandler(client, &config, repository)

	router.GET("/api/v1/users/:id", userHandler.GetUser)
	router.GET("/api/v1/users", userHandler.ListUser)
	router.POST("/api/v1/users", userHandler.CreateUser)
	router.PUT("/api/v1/users/:id", userHandler.UpdateUser)

	// PORT environment variable was defined.
	formattedUrl := fmt.Sprintf(": %s", config.Server.Port)

	router.Run(formattedUrl)
}
