package main

import (
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type config struct {
	User             string
	Password         string
	Host             string
	Port             int
	Params           map[string]string // Connection parameters
	ApiVersion       int
	ClientName       string
	ClientVersion    string
	Schema           string
	Autocommit       bool
	FetchSize        int
	Compression      bool
	ResultSetMaxRows int
	Timeout          time.Time
	Encryption       bool
	UseTLS           bool
}

func getPropertiesFromFile() *config {
	f, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	defaultConfig := getDefaultConfig(cfg.Host, cfg.Port)
	defaultConfig.User = cfg.User
	defaultConfig.Password = cfg.Password
	defaultConfig.Port = cfg.Port
	defaultConfig.Params = cfg.Params
	defaultConfig.ApiVersion = cfg.ApiVersion
	defaultConfig.ClientName = cfg.ClientName
	defaultConfig.ClientVersion = cfg.ClientVersion
	defaultConfig.Schema = cfg.Schema
	defaultConfig.FetchSize = cfg.FetchSize
	defaultConfig.Compression = cfg.Compression
	defaultConfig.ResultSetMaxRows = cfg.ResultSetMaxRows
	defaultConfig.Timeout = cfg.Timeout
	defaultConfig.Encryption = cfg.Encryption
	defaultConfig.UseTLS = cfg.UseTLS
	return defaultConfig
}

func getDefaultConfig(host string, port int) *config {
	return &config{
		Host:        host,
		Port:        port,
		ApiVersion:  2,
		Autocommit:  true,
		Encryption:  true,
		Compression: false,
		UseTLS:      false,
		ClientName:  "Go client",
		Params:      map[string]string{},
		FetchSize:   128 * 1024,
	}
}
