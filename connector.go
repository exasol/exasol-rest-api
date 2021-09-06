package main

import (
	"context"
)

type connector struct {
	config *config
}

func (c *connector) Connect(ctx context.Context) (*connection, error) {
	conn := &connection{
		config:   c.config,
		ctx:      ctx,
		isClosed: true,
	}
	err := conn.connect()
	if err != nil {
		return nil, err
	}

	err = conn.login(ctx)
	if err != nil {
		return nil, err
	}

	return conn, err
}
