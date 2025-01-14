package infrastructure

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/maronfranc/poc-golang-ddd/util"
)

func EnvLoad(name string) error {
	// TODO: https://github.com/joho/godotenv/issues/126#issuecomment-1474645022
	return godotenv.Load(name)
}

func EnvGet(k string) (string, error) {
	value, exists := os.LookupEnv(k)
	if !exists {
		msg := fmt.Sprintf("Key(%s) not in environment", k)
		return "", errors.New(msg)
	}
	return value, nil
}

// EnvGetFileName check command "-env" flag and return .env file name.
//
// Return:
//   - ".env" if "-env=production"
//   - ".env.flag_name" if "-env" is "test" or "development".
func EnvGetFileName() (string, error) {
	ENV_FLAGS := []string{"test", "development", "production"}
	envFlag := flag.String(
		"env",
		"flag_error",
		"-env flags: ['test','development','production']",
	)
	flag.Parse()

	if !util.ArrayContains(ENV_FLAGS, *envFlag) {
		flagStr := strings.Join(ENV_FLAGS[:], ",")
		msg := fmt.Sprintf("Please provide one of the valid -env flag: %s", flagStr)
		return "", errors.New(msg)
	}

	if *envFlag == "production" {
		return ".env", nil
	}

	fileName := fmt.Sprintf(".env.%s", *envFlag)
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
