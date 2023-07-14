package command

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CmdParams struct {
	host   string
	port   int
	dbPath string
	dbName string
}

func readParams(args []string) CmdParams {
	env_host := os.Getenv("VAULT_AUTO_UNSEAL_HOST")
	if env_host == "" {
		env_host = "127.0.0.1"
	}
	env_port := os.Getenv("VAULT_AUTO_UNSEAL_PORT")
	if env_port == "" {
		env_port = "8200"
	}
	env_db_path := os.Getenv("VAULT_AUTO_UNSEAL_DB_PATH")
	if env_db_path == "" {
		env_db_path = "."
	}
	env_db_name := os.Getenv("VAULT_AUTO_UNSEAL_DB_NAME")
	if env_db_name == "" {
		env_db_name = "vault-auto-unseal.db"
	}
	int_env_port, err := strconv.Atoi(env_port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return CmdParams{host: env_host, port: int_env_port, dbPath: env_db_path, dbName: env_db_name}
}

func startServer(p CmdParams, r *chi.Mux) int {
	fmt.Printf("Starting server on %s:%d\n", p.host, p.port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", p.host, p.port), r); err != nil {
		fmt.Printf("Unable to start server, received error: %s", err)
		return 1
	}
	return 0
}

func Run(args []string) int {
	p := readParams(args)
	setDBConf(p.dbPath, p.dbName)
	r := createRouter()
	exitCode := startServer(p, r)

	return exitCode
}
