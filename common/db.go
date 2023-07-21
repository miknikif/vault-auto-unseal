package common

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Init DB Connection
func (c *Config) initDB() error {
	err := CreateDBIfNotExists(c)
	if err != nil {
		return err
	}
	db, err := gorm.Open("sqlite3", fmt.Sprintf("%s/%s", c.Args.DBPath, c.Args.DBName))
	if err != nil {
		return err
	}
	db.DB().SetMaxIdleConns(10)
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.DB = db

	return nil
}

// Get DB object
func GetDB() (*gorm.DB, error) {
	c, err := GetConfig()
	if err != nil {
		return nil, err
	}
	return c.DB, nil
}
