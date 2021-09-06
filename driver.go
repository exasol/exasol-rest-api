package main

import (
	"context"
	"time"
)

type ExasolDriver struct{}

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

func (e ExasolDriver) Open(dsn string) (*connection, error) {
	config, err := parseDSN(dsn)
	if err != nil {
		return nil, err
	}
	exasolConnector := &connector{
		config: config,
	}
	return exasolConnector.Connect(context.Background())
}
