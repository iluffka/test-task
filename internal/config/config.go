package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	Port		string		`json:"port"`
	Start		time.Time	`json:"start"`
	StorageName	string		`json:"storage_name"`
	URLPattern	string		`json:"url_pattern"`
	path		string
	file		*string
}

func New(configPath string, fileName *string) *Config {
	return &Config{
		path: configPath,
		file: fileName,
	}
}

func(c *Config) Load() {
	b, err := ioutil.ReadFile(c.getConfigName(c.path, c.file))
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(b, &c); err != nil {
		log.Fatal(err)
	}
	c.Start = time.Now().UTC()
}
