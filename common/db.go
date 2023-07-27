package common

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// "github.com/miknikif/vault-auto-unseal/policies"
	// "github.com/miknikif/vault-auto-unseal/tokens"
)

const (
	INIT_DB_RES_ERROR   = 0
	INIT_DB_RES_EXISTED = 1
	INIT_DB_RES_CREATED = 2
)

// Create db file if it doesn't exist
func CreateDBIfNotExists(c *Config) (int, error) {
	c.Logger.Debug(fmt.Sprintf("Checking if DB at the following path %s/%s exists", c.Args.DBPath, c.Args.DBName))
	_, err := os.Stat(fmt.Sprintf("%s/%s", c.Args.DBPath, c.Args.DBName))

	if err == nil {
		return INIT_DB_RES_EXISTED, nil
	}

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return INIT_DB_RES_ERROR, err
	}

	c.Logger.Warn(fmt.Sprintf("Creating new sqlite3 DB at the following path %s/%s because it doesn't exist", c.Args.DBPath, c.Args.DBName))
	err = os.MkdirAll(c.Args.DBPath, 0755)
	if err != nil {
		return INIT_DB_RES_ERROR, err
	}
	_, err = os.Create(fmt.Sprintf("%s/%s", c.Args.DBPath, c.Args.DBName))
	if err != nil {
		return INIT_DB_RES_ERROR, err
	}
	c.Logger.Debug(fmt.Sprintf("Created new sqlite3 DB at the following path %s/%s because it doesn't exist", c.Args.DBPath, c.Args.DBName))

	return INIT_DB_RES_CREATED, nil
}

// Init DB Connection
func (c *Config) initDB() error {
	res, err := CreateDBIfNotExists(c)
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
	c.DBStatus = res
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
