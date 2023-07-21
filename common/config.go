package common

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	loghelper "github.com/miknikif/vault-auto-unseal/helper/logging"
)

// All currently used ENV VARS
// Used in the following form: <ENV_PREFIX>_<ENV_VAR>
const (
	ENV_PREFIX     = "VAULT_AUTO_UNSEAL"
	ENV_HOST       = "HOST"
	ENV_PORT       = "PORT"
	ENV_TLS_CA_CRT = "CA_CRT"
	ENV_TLS_CRT    = "TLS_CRT"
	ENV_TLS_KEY    = "TLS_KEY"
	ENV_DB_PATH    = "DB_PATH"
	ENV_DB_NAME    = "DB_NAME"
	ENV_LOG_FORMAT = "LOG_FORMAT"
	ENV_LOG_LEVEL  = "LOG_LEVEL"
)

// TLS Specific configuration provided during startup
// TODO: watch for the certificate changes to restart the app/listener
type TLSConfig struct {
	Enabled bool
	Proto   string // Can be http or https. We're setting this automatically based on Enabled field
	CACrt   string
	TLSCrt  string
	TLSKey  string
}

// Log specific configuration provided during startup
type LogConfig struct {
	LogFormat string
	LogLevel  string
}

// App params provided during startup
// This struct is inialized only once, during startup
// It's should remain unchanged
type Params struct {
	Host      string
	Port      int
	TLS       *TLSConfig
	DBPath    string
	DBName    string
	LogConfig *LogConfig
}

// App Config Struct
type Config struct {
	Lock   sync.RWMutex
	Args   *Params
	Logger hclog.Logger
	DB     *gorm.DB
}

var c *Config

// Configure logging. We're using logger which is used in the Vault project
func (c *Config) configureLogging() error {
	// Parse all the log related config
	logLevel, err := loghelper.ParseLogLevel(c.Args.LogConfig.LogLevel)
	if err != nil {
		return err
	}

	logFormat, err := loghelper.ParseLogFormat(c.Args.LogConfig.LogFormat)
	if err != nil {
		return err
	}

	logCfg := &loghelper.LogConfig{
		LogLevel:  logLevel,
		LogFormat: logFormat,
	}

	c.Logger, err = loghelper.Setup(logCfg, os.Stdout)
	if err != nil {
		return err
	}

	return nil
}

// Read all supported ENV Vars
func (c *Config) readEnv() error {
	tlsEnabled := false
	tlsProto := "http"
	tlsKey := readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_TLS_KEY), "")

	if tlsKey != "" {
		tlsEnabled = true
		tlsProto = "https"
	}

	t := &TLSConfig{
		Enabled: tlsEnabled,
		Proto:   tlsProto,
		CACrt:   readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_TLS_CA_CRT), ""),
		TLSCrt:  readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_TLS_CRT), ""),
		TLSKey:  tlsKey,
	}

	log := &LogConfig{
		LogFormat: readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_LOG_FORMAT), "standard"),
		LogLevel:  readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_LOG_LEVEL), "info"),
	}

	cp := &Params{
		TLS:       t,
		LogConfig: log,
		Host:      readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_HOST), "0.0.0.0"),
		Port:      readEnvInt(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_PORT), 8200),
		DBPath:    readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_DB_PATH), "."),
		DBName:    readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_DB_NAME), "vaseal.db"),
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.Args = cp
	err := c.configureLogging()
	if err != nil {
		return err
	}

	return nil
}

// Create db file if it doesn't exist
func CreateDBIfNotExists(c *Config) error {
	c.Logger.Debug(fmt.Sprintf("Checking if DB at the following path %s/%s exists", c.Args.DBPath, c.Args.DBName))
	_, err := os.Stat(fmt.Sprintf("%s/%s", c.Args.DBPath, c.Args.DBName))

	if err == nil {
		return nil
	}

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	c.Logger.Warn(fmt.Sprintf("Creating new sqlite3 DB at the following path %s/%s because it doesn't exist", c.Args.DBPath, c.Args.DBName))
	err = os.MkdirAll(c.Args.DBPath, 0755)
	if err != nil {
		return err
	}
	_, err = os.Create(fmt.Sprintf("%s/%s", c.Args.DBPath, c.Args.DBName))
	if err != nil {
		return err
	}
	c.Logger.Debug(fmt.Sprintf("Created new sqlite3 DB at the following path %s/%s because it doesn't exist", c.Args.DBPath, c.Args.DBName))

	return nil
}

// Get DB object
func (c *Config) getDB() *gorm.DB {
	return c.DB
}

// Init default Config struct
func DefaultConfig() (*Config, error) {
	c := &Config{}
	// Read ENV Vars to the Config struct
	err := c.readEnv()
	if err != nil {
		return nil, err
	}

	// Init DB
	err = c.initDB()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Return global config
func GetConfig() (*Config, error) {
	if c == nil {
		conf, err := DefaultConfig()
		if err != nil {
			return nil, err
		}

		c = conf
	}

	return c, nil
}

// Get logger interface
func GetLogger() (hclog.Logger, error) {
	var logger hclog.Logger
	c, err := GetConfig()
	if err != nil {
		return logger, err
	}
	logger = c.Logger
	return logger, nil
}
