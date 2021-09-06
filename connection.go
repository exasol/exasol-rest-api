package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"os/user"
	"runtime"
	"strconv"

	"github.com/gorilla/websocket"
)

type connection struct {
	config    *config
	websocket *websocket.Conn
	ctx       context.Context
	isClosed  bool
}

func (c *connection) Close() error {
	c.isClosed = true
	err := c.send(context.Background(), &Command{Command: "disconnect"}, nil)
	c.websocket.Close()
	c.websocket = nil
	return err
}

func (c *connection) simpleExec(query string) (string, error) {
	command := &SQLCommand{
		Command: Command{"execute"},
		SQLText: query,
		Attributes: Attributes{
			ResultSetMaxRows: c.config.ResultSetMaxRows,
		},
	}
	result, err := c.asyncSend2(command)
	if err != nil {
		return "", err
	}
	return result, err
}

func (c *connection) login(ctx context.Context) error {
	hasCompression := c.config.Compression
	c.config.Compression = false
	loginCommand := &LoginCommand{
		Command:         Command{"login"},
		ProtocolVersion: c.config.ApiVersion,
	}
	loginResponse := &PublicKeyResponse{}
	err := c.send(ctx, loginCommand, loginResponse)
	if err != nil {
		return err
	}

	pubKeyMod, _ := hex.DecodeString(loginResponse.PublicKeyModulus)
	var modulus big.Int
	modulus.SetBytes(pubKeyMod)

	pubKeyExp, _ := strconv.ParseUint(loginResponse.PublicKeyExponent, 16, 32)

	pubKey := rsa.PublicKey{
		N: &modulus,
		E: int(pubKeyExp),
	}
	password := []byte(c.config.Password)
	encPass, err := rsa.EncryptPKCS1v15(rand.Reader, &pubKey, password)
	if err != nil {
		errorLogger.Printf("password encryption error: %s", err)
		return driver.ErrBadConn
	}
	b64Pass := base64.StdEncoding.EncodeToString(encPass)

	authRequest := AuthCommand{
		Username:       c.config.User,
		Password:       b64Pass,
		UseCompression: false,
		ClientName:     c.config.ClientName,
		DriverName:     fmt.Sprintf("exasol-driver-go %s", "v1.0.0"),
		ClientOs:       runtime.GOOS,
		ClientVersion:  c.config.ClientName,
		ClientRuntime:  runtime.Version(),
		Attributes: Attributes{
			Autocommit:         c.config.Autocommit,
			CurrentSchema:      c.config.Schema,
			CompressionEnabled: hasCompression,
		},
	}

	if osUser, err := user.Current(); err != nil {
		authRequest.ClientOsUsername = osUser.Username
	}

	authResponse := &AuthResponse{}
	err = c.send(ctx, authRequest, authResponse)
	if err != nil {
		return err
	}
	c.isClosed = false
	c.config.Compression = hasCompression

	return nil
}
