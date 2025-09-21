package mysql

import (
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/AhmadNizar/cata-dtc/internal/entity"
)

// NewMysqlRepository initiate mysql database client connection
func NewMysqlRepository(option *entity.MysqlDBConnOption, logLevel string) *gorm.DB {
	client, err := gorm.Open(mysql.Open(option.URL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Reference: https://www.alexedwards.net/blog/configuring-sqldb
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	maxIdleConn, _ := strconv.Atoi(option.MaxIdleConn)
	sqlDB, _ := client.DB()
	sqlDB.SetMaxIdleConns(maxIdleConn)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	maxOpenConn, _ := strconv.Atoi(option.MaxOpenConn)
	sqlDB.SetMaxOpenConns(maxOpenConn)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	connMaxLifeTimeInMinutes, _ := strconv.Atoi(option.MaxLifetimeInMinute)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifeTimeInMinutes) * time.Minute)

	if logLevel == "info" {
		client = client.Debug()
	}

	return client
}