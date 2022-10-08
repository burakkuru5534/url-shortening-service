package main

import (
	"github.com/burakkuru5534/src/api"
	"github.com/burakkuru5534/src/auth"
	"github.com/burakkuru5534/src/helper"
	"github.com/burakkuru5534/src/service"
	_ "github.com/lib/pq"

	"log"
)

const (
	host      = "127.0.0.1"
	port      = 5432
	dbName    = "url-shortening-service"
	user      = "postgres"
	password  = "tayitkan"
	sslmode   = "disable"
	secretKey = "2GcQCe7SuKxbaA3NSMBy8ztBPbfDsXJ4"
	devMode   = false
)

func main() {

	helper.InitConf()

	conInfo := helper.PgConnectionInfo{
		Host:     host,
		Port:     port,
		Database: dbName,
		Username: user,
		Password: password,
		SSLMode:  sslmode,
	}

	helper.Conf.Auth = auth.NewAuth(secretKey, devMode)

	db, err := helper.NewPgSqlxDbHandle(conInfo, 10)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	// Create Appplication Service
	err = helper.InitApp(db)
	if err != nil {
		log.Fatal("Error initializing application: ", err)
	}

	service.StartHttpService(8080, api.HttpService())
}
