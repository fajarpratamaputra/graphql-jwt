package db

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
)

func ConnectGORM() *gorm.DB {
    DBMS     := "mysql"
    USER     := "root"
    PASS     := ""
    PROTOCOL := "tcp(127.0.0.1)"
    DBNAME   := "go_test"

    CONNECT := USER+":"+PASS+"@"+PROTOCOL+"/"+DBNAME
    db,err := gorm.Open(DBMS, CONNECT)

    if err != nil {
        panic(err.Error())
    }
    return db
}
