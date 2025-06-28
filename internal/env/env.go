package env

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ROOT_DIR       string = getProjectRoot()
	ENVIRONMENT    string
	SERVER_ADDR    string
	DB_ADDR        string
	API_URL        string
	MAX_IDLE_CONNS int
	MAX_OPEN_CONNS int
	MAX_IDLE_TIME  string
)

func LoadEnvs() {
	println("LOADING ENV")
	envFile := filepath.Join(ROOT_DIR, ".env")

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}

	ENVIRONMENT = getString("ENVIRONMENT", "local")
	SERVER_ADDR = getString("SERVER_ADDR", ":8001")
	API_URL = getString("API_URL", "localhost:8001")
	DB_ADDR = getString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/gophersocial?sslmode=disable")
	MAX_IDLE_CONNS = GetInt("MAX_IDLE_CONNS", 5)
	MAX_OPEN_CONNS = GetInt("MAX_OPEN_CONNS", 10)
	MAX_IDLE_TIME = getString("MAX_IDLE_TIME", "5m")
}

func getString(key, fallback string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	return val
}
func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get root workdir:", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal("Faild to find root workdir")
		}
		dir = parent
	}
}
