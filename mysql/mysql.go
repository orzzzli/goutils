package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

var GlobalMysqlPool *sqlx.DB

func NewMysqlPool(user string, pass string, host string, port string, database string,maxOpen int, maxIdle int, maxLife int) error {
	driverStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",user,pass,host,port,database)
	pool, _ := sqlx.Open("mysql",driverStr)
	pool.SetMaxOpenConns(maxOpen)
	pool.SetMaxIdleConns(maxIdle)
	pool.SetConnMaxLifetime(time.Second * time.Duration(maxLife))
	err := pool.Ping()
	if err != nil {
		return err
	}
	GlobalMysqlPool = pool
	return nil
}