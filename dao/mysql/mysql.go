package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var db *sqlx.DB

func Init() (err error) {
	var (
		host = viper.GetString("mysql.host")
		port = viper.GetInt("mysql.port")
		user = viper.GetString("mysql.user")
		password = viper.GetString("mysql.password")
		dbname = viper.GetString("mysql.dbname")
		maxOpenConnections = viper.GetInt("max_open_connections")
		maxIdleConnections = viper.GetInt("max_idle_connections")
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	return
}

func Close() {
	_ = db.Close()
}
