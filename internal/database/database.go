package database

import (
	"fmt"

	"github.com/topboyasante/go-snip/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
)

var (
	DB *gorm.DB
)

func ConnectToDB() {
	var err error
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.ENV.DBUser, config.ENV.DBPassword, config.ENV.DBHost, config.ENV.DBPort, config.ENV.DBName)


	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the DB", err)
	}

	fmt.Println("CONNECTED TO DB")

}
