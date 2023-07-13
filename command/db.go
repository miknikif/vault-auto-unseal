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

type InternalKey struct {
	id      int
	name    string
	aesKey  string
	version int
	//rsaKey *rsa.PrivateKey
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

var lock = &sync.Mutex{}
var dbconf *DBConfig

func setDBConf(path string, name string) {
	lock.Lock()
	defer lock.Unlock()
	dbconf = &DBConfig{path: path, name: name}
}

func getDBConf() *DBConfig {
	if dbconf == nil {
		setDBConf(".", "vault-auto-unseal.db")
	}

	return dbconf
}

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

func getDBConn(conf *DBConfig) (*sql.DB, error) {
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

	return db, nil
}

func addAESKeyToDB(key InternalKey) error {
	db, err := getDBConn(getDBConf())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO aes_keys (name, version, key) VALUES (?, ?, ?)")
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

func getAESKeyFromDB(name string, version int) (InternalKey, error) {
	db, err := getDBConn(getDBConf())
	if err != nil {
		return InternalKey{}, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT * FROM aes_keys WHERE name = ? AND version = ? LIMIT 1", name, version)
	ik := InternalKey{}
	if err = row.Scan(&ik.id, &ik.name, &ik.version, &ik.aesKey); err != nil {
		return InternalKey{}, err
	}
	return ik, nil
}

func dbGetKey(name string) (Key, error) {
	versions := []int{1}
	sik := map[int]InternalKey{}
	for i, _ := range versions {
		ik, err := getAESKeyFromDB(name, 1)
		if err != nil {
			return Key{}, err
		}
		sik[i] = ik
	}
	k := Key{Name: name, LatestVersion: versions[len(versions)-1], Type: KEY_TYPE_AES256_GCM96, keys: sik}
	return k, nil
}
