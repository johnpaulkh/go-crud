package main

import (
	"johnpaulkh/go-crud/api/config"
	"johnpaulkh/go-crud/api/server"
	_ "johnpaulkh/go-crud/docs"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title			User APIs
// @version		1.0
// @description	User APIs.
// @termsOfService	http://swagger.io/terms/
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:8080
// @BasePath		/api/v1
// @schemes		http
func main() {

	config := readConfiguration(read())

	server.Initialize(config)
}

func readConfiguration(conf config.Configuration) config.Configuration {

	mongoUri := os.Getenv("MONGODB_URL")
	port := os.Getenv("SERVER_PORT")
	dbName := os.Getenv("DB_NAME")
	collection := os.Getenv("COLLECTION")
	appName := os.Getenv("APP_NAME")

	if mongoUri != "" || port != "" || dbName != "" || collection != "" || appName != "" {
		return config.Configuration{
			App: config.Application{
				Name: appName},
			Database: config.DatabaseSetting{
				Url:        mongoUri,
				DbName:     dbName,
				Collection: collection},
			Server: config.ServerSettings{
				Port: port},
		}
	}

	// return config.yml variable
	return config.Configuration{
		App: config.Application{
			Name: conf.App.Name},
		Database: config.DatabaseSetting{
			Url:        conf.Database.Url,
			DbName:     conf.Database.DbName,
			Collection: conf.Database.Collection},
		Server: config.ServerSettings{
			Port: conf.Server.Port},
	}
}

func read() config.Configuration {
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigFile("config.yml")

	var config config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		logrus.Errorf("Unable to decode into struct, %v", err)
	}

	logrus.Warnf("Config with variables %v", config)

	return config
}
