package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Address string     `json:"address"`
	Mysql   Mysql      `json:"mysql"`
	Routers []FileTree `json:"routers"`
}

type Mysql struct {
	Address string `json:"address"`
	User    string `json:"user"`
	Key     string `json:"key"`
	DB      string `json:"db"`
}

type FileTree struct {
	ServerPath string `json:"server_path"`
	LocalPath  string `json:"local_path"`
	Forbidden  []string `json:"forbidden"`
}

func LoadConfig(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
