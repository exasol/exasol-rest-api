package exasol_rest_api

import (
	"context"
)

type connector struct {
	connProperties *connectionProperties
}

func (c *connector) Connect(ctx context.Context) (*connection, error) {
	conn := &connection{
		config:   c.connProperties,
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
