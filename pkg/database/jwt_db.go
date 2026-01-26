package database

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "gojwt/pkg/config"
)

func NewMySQL() (*sql.DB, error) {
    cfg := config.Get()
    
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", 
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
    )
    
    db, err := sql.Open("mysql", dsn)
    if (err != nil) {
        return nil, err
    }
    
    err = db.Ping()
    if err != nil {
    	return nil, err
    }
    
    return db, nil
}