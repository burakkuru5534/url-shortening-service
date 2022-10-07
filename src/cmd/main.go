package main

import (
	"github.com/burakkuru5534/src/api"
	"github.com/burakkuru5534/src/auth"
	"github.com/burakkuru5534/src/helper"
	"github.com/burakkuru5534/src/service"
	_ "github.com/lib/pq"

	"log"
)

func main() {

	helper.InitConf()

	conInfo := helper.PgConnectionInfo{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "url-shortening-service",
		Username: "postgres",
		Password: "tayitkan",
		SSLMode:  "disable",
	}

	helper.Conf.Auth = auth.NewAuth("2GcQCe7SuKxbaA3NSMBy8ztBPbfDsXJ4", false)

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
