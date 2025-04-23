package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	ROOT_DIR    string = getProjectRoot()
	ENVIRONMENT string
	SERVER_ADDR string
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
}

func getString(key, fallback string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	return val
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
