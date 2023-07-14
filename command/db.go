package command

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync"
)

type KeyType string

type DBConfig struct {
	path string
	name string
}

type DBConn struct {
	*DBConfig
	db   *sql.DB
	lock *sync.Mutex
}

type InternalKey struct {
	id      int
	name    string
	aesKey  string
	version int
}

type Key struct {
	keys                 map[int]InternalKey
	AllowPlaintextBackup bool        `json:"allow_plaintext_backup"`
	AutoRotatePeriod     int         `json:"auto_rotate_period"`
	DeletionAllowed      bool        `json:"deletion_allowed"`
	Derived              bool        `json:"derived"`
	Exportable           bool        `json:"exportable"`
	ImportedKey          bool        `json:"imported_key"`
	Keys                 map[int]int `json:"keys"`
	LatestVersion        int         `json:"latest_version"`
	MinAvailableVersion  int         `json:"min_available_version"`
	MinDecryptionVersion int         `json:"min_decryption_version"`
	MinEncryptionVersion int         `json:"min_encryption_version"`
	Name                 string      `json:"name"`
	SupportsDecryption   bool        `json:"supports_decryption"`
	SupportsDerivation   bool        `json:"supports_derivation"`
	SupportsEncryption   bool        `json:"supports_encryption"`
	SupportsSigning      bool        `json:"supports_signing"`
	Type                 KeyType     `json:"type"`
}

const (
	KEY_TYPE_AES256_GCM96 KeyType = "aes256-gcm96"
	KEY_TYPE_ED25519      KeyType = "ed25519"
	KEY_TYPE_RSA_4096     KeyType = "rsa-4096"
)

var ldbconf = &sync.Mutex{}
var dbconf *DBConfig
var ldbconn = &sync.Mutex{}
var dbconn *DBConn

func setDBConf(path string, name string) {
	ldbconf.Lock()
	defer ldbconf.Unlock()
	dbconf = &DBConfig{path: path, name: name}
}

// Returns a singleton of DBConfig object
func getDBConf() *DBConfig {
	if dbconf == nil {
		setDBConf(".", "vault-auto-unseal.db")
	}
	ldbconf.Lock()
	defer ldbconf.Unlock()

	return dbconf
}

// Function executed when sqlite DB doesn't exists
// This small utility doesn't support any migrations on the DB yet
func initDB(conf *DBConfig) error {
	fmt.Printf("Creating new sqlite3 DB at the %s/%s\n", conf.path, conf.name)
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s", conf.path, conf.name))
	if err != nil {
		return err
	}
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS aes_keys (id INTEGER PRIMARY KEY, name TEXT, version INTEGER, key TEXT)")
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	db.Close()
	return nil
}

// Helper function to get a DB connection
func createDBConn(conf *DBConfig) (*DBConn, error) {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", conf.path, conf.name)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := initDB(conf)
			if err != nil {
				return nil, err
			}
		}
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s", conf.path, conf.name))
	if err != nil {
		return nil, err
	}

	dbconn := &DBConn{DBConfig: conf, db: db, lock: &sync.Mutex{}}

	return dbconn, nil
}

func getDBConn() (*DBConn, error) {
	ldbconn.Lock()
	defer ldbconn.Unlock()

	if dbconn == nil {
		conn, err := createDBConn(getDBConf())
		if err != nil {
			return nil, err
		}
		dbconn = conn
	}

	return dbconn, nil
}

// Add new AES key to DB.
func addAESKeyToDB(key InternalKey) error {
	dbconn, err := getDBConn()
	if err != nil {
		return err
	}
	dbconn.lock.Lock()
	defer dbconn.lock.Unlock()

	stmt, err := dbconn.db.Prepare("INSERT INTO aes_keys (name, version, key) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(key.name, key.version, key.aesKey)
	if err != nil {
		return err
	}

	return nil
}

// Get AES key from DB. Queried by key name and key version
func getAESKeyFromDB(name string, version int) (InternalKey, error) {
	dbconn, err := getDBConn()
	if err != nil {
		return InternalKey{}, err
	}
	dbconn.lock.Lock()
	defer dbconn.lock.Unlock()
	row := dbconn.db.QueryRow("SELECT * FROM aes_keys WHERE name = ? AND version = ? LIMIT 1", name, version)
	ik := InternalKey{}
	if err = row.Scan(&ik.id, &ik.name, &ik.version, &ik.aesKey); err != nil {
		return InternalKey{}, err
	}
	return ik, nil
}

// Delete AES key from DB. Queried by key name and key version
func delAESKeyFromDB(key InternalKey) error {
	dbconn, err := getDBConn()
	if err != nil {
		return err
	}
	dbconn.lock.Lock()
	defer dbconn.lock.Unlock()
	fmt.Println(key.id)
	_, err = dbconn.db.Exec("DELETE FROM aes_keys WHERE id = ?", key.id)
	if err != nil {
		return err
	}
	return nil
}

// Getting Key by quering db. Currently we're checking only a single key version, version 1
func dbGetKey(name string) (Key, error) {
	versions := []int{1}
	sik := map[int]InternalKey{}
	for i := range versions {
		ik, err := getAESKeyFromDB(name, 1)
		if err != nil {
			return Key{}, err
		}
		sik[i+1] = ik
	}
	k := Key{Name: name, LatestVersion: versions[len(versions)-1], Type: KEY_TYPE_AES256_GCM96, keys: sik}
	fmt.Println(k.keys)
	return k, nil
}
