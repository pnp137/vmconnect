package system

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"linksupply.io/vmconnect/database"
)

type DbConfig struct {
	Host     string `json:"db_host"`
	Port     int    `json:"db_port"`
	Username string `json:"db_user"`
	Password string `json:"db_password"`
	Name     string `json:"db_name"`
}

type DataSource struct {
	Db *gorm.DB
}

var dataSource *DataSource = nil

func NewDataSource(config *DbConfig, logEnabled bool) *DataSource {
	if dataSource != nil {
		return dataSource
	}

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Name)
	// Add customization to DB log level
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %s", err))
	}

	sqlDB, err := connection.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %s", err))
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("Successfully connected to the Database")

	// Migrate DB
	database.Migrate(connection)
	log.Printf("DB migrations completed successfully")

	dataSource = &DataSource{Db: connection}
	return dataSource
}
