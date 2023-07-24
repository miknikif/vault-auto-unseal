package common

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
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
	ENV_HOSTNAME               = "HOSTNAME"
	ENV_PREFIX                 = "VAULT_AUTO_UNSEAL"
	ENV_HOST                   = "HOST"
	ENV_PORT                   = "PORT"
	ENV_TLS_CA_CRT_PATH        = "CA_CRT_PATH"
	ENV_TLS_CLIENT_CA_CRT_PATH = "CLIENT_CA_CRT_PATH"
	ENV_TLS_CRT_PATH           = "TLS_CRT_PATH"
	ENV_TLS_KEY_PATH           = "TLS_KEY_PATH"
	ENV_DB_PATH                = "DB_PATH"
	ENV_DB_NAME                = "DB_NAME"
	ENV_LOG_FORMAT             = "LOG_FORMAT"
	ENV_LOG_LEVEL              = "LOG_LEVEL"
	ENV_PRODUCTION             = "PRODUCTION"
)

// Log specific configuration provided during startup
type LogConfig struct {
	LogFormat string
	LogLevel  string
}

// App params provided during startup
// This struct is inialized only once, during startup
// It's should remain unchanged
type Params struct {
	Host         string
	Port         int
	DBPath       string
	DBName       string
	IsProduction bool
	LogConfig    *LogConfig
}

// TLS conf
type TLSConfig struct {
	TLSConfig *tls.Config
	BundleCrt string
	CACrt     string
	TLSCrt    string
	TLSKey    string
}

// App Config Struct
type Config struct {
	Lock   sync.RWMutex
	Args   *Params
	TLS    *TLSConfig
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

// Configure TLS
func (c *Config) configureTLS() error {
	tlsEnabled := false
	hostname := readEnv(ENV_HOSTNAME, "vau-server")
	tlsCAPath := readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_TLS_CA_CRT_PATH), "")
	tlsClientCAPath := readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_TLS_CLIENT_CA_CRT_PATH), "")
	tlsCRTPath := readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_TLS_CRT_PATH), "")
	tlsKeyPath := readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_TLS_KEY_PATH), "")

	tlsConfig := &TLSConfig{}

	t := &tls.Config{
		ServerName:               hostname,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		ClientAuth:               tls.NoClientCert,
	}

	// Load custom CA to the store
	if fileExists(tlsCAPath) {
		crt, err := ioutil.ReadFile(tlsCAPath)
		if err != nil {
			return err
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(crt)
		t.RootCAs = certPool
		tlsConfig.CACrt = tlsCAPath
	}

	// Load custom client CA to the store
	if fileExists(tlsClientCAPath) {
		crt, err := ioutil.ReadFile(tlsClientCAPath)
		if err != nil {
			return err
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(crt)
		t.ClientCAs = certPool
		t.ClientAuth = tls.RequireAndVerifyClientCert
	}

	// Load server certificate
	if fileExists(tlsCRTPath) && fileExists(tlsKeyPath) {
		tlsEnabled = true
		bundleCrt := tlsCRTPath

		if fileExists(tlsCAPath) {
			f, err := os.CreateTemp(os.TempDir(), "vau-crt")
			if err != nil {
				return err
			}

			caCrt, err := ioutil.ReadFile(tlsCAPath)
			if err != nil {
				return err
			}

			tlsCrt, err := ioutil.ReadFile(tlsCRTPath)
			if err != nil {
				return err
			}

			_, err = f.Write(tlsCrt)
			if err != nil {
				return err
			}
			_, err = f.Write(caCrt)
			if err != nil {
				return err
			}
			bundleCrt = f.Name()
		}

		crt, err := tls.LoadX509KeyPair(bundleCrt, tlsKeyPath)
		if err != nil {
			return err
		}

		t.Certificates = []tls.Certificate{crt}
		tlsConfig.BundleCrt = bundleCrt
		tlsConfig.TLSCrt = tlsCRTPath
		tlsConfig.TLSKey = tlsKeyPath
	}

	if tlsEnabled {
		tlsConfig.TLSConfig = t
		c.TLS = tlsConfig
	}

	return nil
}

// Read all supported ENV Vars
func (c *Config) readEnv() error {
	log := &LogConfig{
		LogFormat: readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_LOG_FORMAT), "standard"),
		LogLevel:  readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_LOG_LEVEL), "info"),
	}

	cp := &Params{
		LogConfig:    log,
		Host:         readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_HOST), "0.0.0.0"),
		Port:         readEnvInt(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_PORT), 8200),
		DBPath:       readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_DB_PATH), "."),
		DBName:       readEnv(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_DB_NAME), "vaseal.db"),
		IsProduction: readEnvBool(fmt.Sprintf("%s_%s", ENV_PREFIX, ENV_PRODUCTION), true),
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.Args = cp
	err := c.configureLogging()
	if err != nil {
		return err
	}
	err = c.configureTLS()
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
