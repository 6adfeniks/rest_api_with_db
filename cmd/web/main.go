/*
	Package provides REST API backend for simple app CRUD using MySQL as database
*/
package main

import (
	"github.com/6adfeniks/rest_api_with_db/cmd/web/app"
	"github.com/6adfeniks/rest_api_with_db/internal/config"
	"log"
)

// main starts the app
func main() {
	cfg, err := config.NewConfig("../../configs/config2.yml")
	if err != nil {
		log.Fatal(err)
	}

	a := app.App{}
	a.Initialize(cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname)

	a.Run(cfg.Server.Port)
}
