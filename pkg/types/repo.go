package types

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	yaml "gopkg.in/yaml.v2"
)

var ConfigDefaultLocation = filepath.Join(os.Getenv("HOME"), ".git-reposync")

type Config struct {
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	BaseDirectory string    `json:"baseDir"`
	LastSync      time.Time `json:"lastSync"`
}

func GetUserConfig() (*Config, error) {
	result := &Config{}
	content, err := ioutil.ReadFile(ConfigDefaultLocation)
	if err != nil {
		if _, isNotFound := err.(*os.PathError); isNotFound {
			return result, nil
		}
		return nil, err
	}
	err = yaml.Unmarshal(content, result)
	return result, err
}

func SaveUserConfig(config *Config) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ConfigDefaultLocation, content, 0600)
}
