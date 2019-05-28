package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

const (
	ConfigPath    = "configs/"
	defaultConfig = "development.json"
)

type Config struct {
	Port        string    `json:"port"`
	Start       time.Time `json:"start"`
	StorageName string    `json:"storage_name"`
	URLPattern  string    `json:"url_pattern"`
	path        string
	file        string
}

func New(configPath, fileName string) *Config {
	return &Config{
		path: configPath,
		file: fileName,
	}
}

func (c *Config) Load() {
	b, err := ioutil.ReadFile(c.getConfigName(c.path, c.file))
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(b, &c); err != nil {
		log.Fatal(err)
	}
	c.Start = time.Now().UTC()
}

func (c *Config) getConfigName(path, name string) string {
	if name == "" {
		return fmt.Sprintf("%s%s", path, defaultConfig)
	}
	return fmt.Sprintf("%s%s", path, name)
}
