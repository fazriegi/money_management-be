package config

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var DB *sqlx.DB

func NewDatabase(viper *viper.Viper) {
	host := viper.GetString("db.host")
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	name := viper.GetString("db.name")
	port := viper.GetInt32("db.port")

	dbDriver := "mysql"
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		username,
		password,
		host,
		port,
		name,
	)

	conn, err := sqlx.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	if err = conn.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	log.Println("connected to database successfully")

	DB = conn
}

func GetDatabase() *sqlx.DB {
	if DB == nil {
		log.Fatal("database connection is not initialized")
	}
	return DB
}
