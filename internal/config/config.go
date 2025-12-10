package config

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-chi/jwtauth/v5"
)

var (
	once      sync.Once
	appConfig *AppConfig
)

type AppConfig struct {
	DB   *sql.DB
	Auth *jwtauth.JWTAuth
}

// NewAppConfig initializes and returns a singleton instance of AppConfig.
// It ensures that the configuration is loaded only once using sync.Once.
// This function sets up the JWT authentication client and the database client.
// If there is an error initializing the database client, it will panic.
// Wherever you need any config variables, use this function call directly as it's a singleton.

func NewAppConfig() *AppConfig {
	once.Do(func() {
		appConfig = &AppConfig{
			Auth: newJWTAuthClient(),
		}
		db, err := newDBClient()
		if err != nil {
			panic(err)
		}
		appConfig.DB = db
	})

	return appConfig
}

func newJWTAuthClient() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(Envs.JWT_SECRET_KEY), nil)

}

func newDBClient() (*sql.DB, error) {
	db, err := sql.Open(Envs.DB_DRIVER, Envs.DB_URL)
	if err != nil {
		fmt.Println("error connecting to database", err)
		return nil, err
	}
	db.SetConnMaxLifetime(time.Duration(Envs.DB_MAX_CONN_TIME_SEC) * time.Second)
	db.SetMaxOpenConns(Envs.DB_MAX_OPEN_CONN)
	db.SetMaxIdleConns(Envs.DB_MAX_IDLE_CONN)

	if err := db.Ping(); err != nil {
		fmt.Println("error pinging database", err)
		return nil, err
	}

	fmt.Println("connected to database")
	return db, nil

}
