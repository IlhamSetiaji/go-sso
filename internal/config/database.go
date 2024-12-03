package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func NewDatabase() *gorm.DB {
	once.Do(func() {
		config := viper.New()

		config.SetConfigName("config")
		config.SetConfigType("json")
		config.AddConfigPath("./")
		err := config.ReadInConfig()

		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
		driver := config.GetString("database.driver")
		username := config.GetString("database.username")
		password := config.GetString("database.password")
		host := config.GetString("database.host")
		port := config.GetInt("database.port")
		database := config.GetString("database.name")
		idleConnection := config.GetInt("database.pool.idle")
		maxConnection := config.GetInt("database.pool.max")
		maxLifeTimeConnection := config.GetInt("database.pool.lifetime")

		if idleConnection <= 0 {
			idleConnection = 10 // Default to 10 idle connections
		}
		if maxConnection <= 0 {
			maxConnection = 100 // Default to 100 max connections
		}
		if maxLifeTimeConnection <= 0 {
			maxLifeTimeConnection = 14400 // Default to 4 hours (in seconds)
		}

		var dsn string
		var db *gorm.DB

		switch driver {
		case "mysql":
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		case "postgres":
			dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, username, database, password)
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		default:
			log.Fatalf("unsupported database driver: %s", driver)
		}

		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		connection, err := db.DB()
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		connection.SetMaxIdleConns(idleConnection)
		connection.SetMaxOpenConns(maxConnection)
		connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

		dbInstance = db
	})

	return dbInstance
}
