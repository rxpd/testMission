package DBModule

import (
	"avito/config"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
)

var db *sqlx.DB

func InitializeDB(config config.DBConfig) {
	fmt.Println()
	var err error
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=%s", config.Address, config.Port, config.Username, config.Password, config.DBName, config.SSLMode)

	if db, err = sqlx.Connect(config.Driver, connectionUrl); err != nil {
		log.Fatal(fmt.Sprintf("\nCannot connect to \"%s\" database\n", config.DBName))
	}

	fmt.Printf("\nSuccessful connection \"%s\" database\n", config.DBName)
}

func init() {
	InitializeDB(config.DBConf)
}

func GetDB() *sqlx.DB {
	return db
}
