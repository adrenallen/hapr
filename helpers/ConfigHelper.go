package helpers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Jeffail/gabs"
)

var configFilePath string
var configContainer *gabs.Container

func SetConfigFilePath(path string) {
	configFilePath = path
}

func GetConfigValue(key string) (string, error) {
	configContainer, err := parseConfig()
	if err != nil {
		return "", err
	}

	if configContainer.ExistsP(key) {
		return fmt.Sprintf("%v", configContainer.Path(key).Data()), nil
	}

	return "", errors.New("Config key not found")
}

func parseConfig() (*gabs.Container, error) {
	if configContainer != nil {
		return configContainer, nil
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	parsed, err := gabs.ParseJSON(b)
	if err != nil {
		return nil, err
	}

	configContainer = parsed
	return configContainer, nil
}
