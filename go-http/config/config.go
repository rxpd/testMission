package config

import (
	_ "github.com/jackc/pgx/log/logrusadapter"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type DBConfig struct {
	Address, Port, Username, Password, DBName, SSLMode, Driver string
}

var DBConf = DBConfig{}

var Logging bool
var ParseIntervalInSeconds int
var ManualCheckCooldownInMinutes int

//var MainViewerObjectsPerPage int
//var NormalizedObjectsPerPage int

func init() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal(err)
	}
	DBConf = DBConfig{
		Address:  os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		Username: os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		DBName:   os.Getenv("PG_DBNAME"),
		SSLMode:  os.Getenv("PG_SSLMODE"),
		Driver:   os.Getenv("PG_DRIVER"),
	}

	Logging, err = strconv.ParseBool(os.Getenv("LOGGING"))
	if err != nil {
		log.Fatal(err)
	}



	ManualCheckCooldownInMinutes, err = strconv.Atoi(os.Getenv("MANUAL_CHECK_COOLDOWN_IN_MINUTES"))
	if err != nil {
		log.Fatal(err)
	}
	ParseIntervalInSeconds, err = strconv.Atoi(os.Getenv("PARSE_INTERVAL_IN_SECONDS"))
	if err != nil {
		log.Fatal(err)
	}
	//ParseIntervalInSeconds = time.Duration(configInterval) * time.Second
	//fmt.Println(ParseIntervalInSeconds)
}
