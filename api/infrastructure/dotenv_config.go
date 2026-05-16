package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func EnvLoad(name string) error {
	err := godotenv.Load(preffixDir(name))
	return err
}

func EnvGet(k string) (string, error) {
	value, exists := os.LookupEnv(k)
	if !exists {
		msg := fmt.Sprintf("Key(%s) not in environment", k)
		return "", errors.New(msg)
	}
	return value, nil
}

// EnvGetFileName check command "ENV" variable and return `.env` file name.
//
// Return:
//   - ".env" if "ENV=production"
//   - ".env.flag_name" if "ENV" is "test" or "development".
func EnvGetFileName() (string, error) {
	ENV_FLAGS := []string{"test", "development", "production"}
	envFlag := os.Getenv("ENV")

	if !slices.Contains(ENV_FLAGS, envFlag) {
		flagStr := strings.Join(ENV_FLAGS[:], ",")
		msg := fmt.Sprintf("Please provide one of the valid ENV: %s", flagStr)
		return "", errors.New(msg)
	}

	if envFlag == "production" {
		return ".env", nil
	}

	fileName := fmt.Sprintf(".env.%s", envFlag)
	return fileName, nil
}

func EnvGetAsBool(k string) (bool, error) {
	vStr, err := EnvGet(k)
	if err != nil {
		return false, err
	}
	vBool, err := strconv.ParseBool(vStr)
	return vBool, err
}

// SEE:: https://github.com/joho/godotenv/issues/126#issuecomment-1474645022
// Dir returns the absolute path of the given environment file (envFile) in the Go module's
// root directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the envFile to the directory containing 'go.mod'.
// It panics if it fails to find the 'go.mod' file.
func preffixDir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}
