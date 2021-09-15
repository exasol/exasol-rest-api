package exasol_rest_api

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math/big"
	"net/url"
	"os/user"
	"runtime"
	"strconv"
)

type websocketConnection struct {
	connProperties *ConnectionProperties
	websocket      *websocket.Conn
}

func (connection *websocketConnection) close() error {
	err := connection.send(&Command{Command: "disconnect"}, nil)
	connection.websocket.Close()
	connection.websocket = nil
	return err
}

func (connection *websocketConnection) connect() error {
	uri := fmt.Sprintf("%s:%d", connection.connProperties.Host, connection.connProperties.Port)
	exaURL := url.URL{
		Scheme: connection.getURIScheme(),
		Host:   uri,
	}
	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: !connection.connProperties.UseTLS}
	websocketConnection, _, err := dialer.DialContext(context.Background(), exaURL.String(), nil)
	if err == nil {
		connection.websocket = websocketConnection
		connection.websocket.EnableWriteCompression(false)
	}
	return err
}

func (connection *websocketConnection) getURIScheme() string {
	if connection.connProperties.Encryption {
		return "wss"
	} else {
		return "ws"
	}
}

func (connection *websocketConnection) executeQuery(query string) (string, error) {
	command := &SQLCommand{
		Command: Command{"execute"},
		SQLText: query,
		Attributes: Attributes{
			ResultSetMaxRows: 1000,
		},
	}
	result, err := connection.sendRequestWithStringResponse(command)
	if err != nil {
		return "", err
	} else {
		return result, err
	}
}

func (connection *websocketConnection) login() error {
	loginCommand := &LoginCommand{
		Command:         Command{"login"},
		ProtocolVersion: connection.connProperties.ApiVersion,
	}
	loginResponse := &PublicKeyResponse{}
	err := connection.send(loginCommand, loginResponse)
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
	password := []byte(connection.connProperties.Password)
	encPass, err := rsa.EncryptPKCS1v15(rand.Reader, &pubKey, password)
	if err != nil {
		ErrorLogger.Printf("password encryption error: %s", err)
		return driver.ErrBadConn
	}
	b64Pass := base64.StdEncoding.EncodeToString(encPass)

	authRequest := AuthCommand{
		Username:       connection.connProperties.User,
		Password:       b64Pass,
		UseCompression: false,
		ClientName:     "Exasol REST API",
		ClientOs:       runtime.GOOS,
		ClientRuntime:  runtime.Version(),
	}

	if osUser, err := user.Current(); err != nil {
		authRequest.ClientOsUsername = osUser.Username
	}

	err = connection.send(authRequest, nil)
	if err != nil {
		return err
	}

	return nil
}

func (connection *websocketConnection) send(request, response interface{}) error {
	receiver, err := connection.sendRequestWithInterfaceResponse(request)
	if err != nil {
		return err
	}
	err = receiver(response)
	if err != nil {
		return err
	}
	return nil
}

func (connection *websocketConnection) sendRequestWithInterfaceResponse(request interface{}) (func(interface{}) error, error) {
	requestAsJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	messageType := websocket.TextMessage
	err = connection.websocket.WriteMessage(messageType, requestAsJson)
	if err != nil {
		return nil, err
	}

	return func(response interface{}) error {
		_, message, err := connection.websocket.ReadMessage()
		if err != nil {
			return err
		}
		result := &BaseResponse{}
		err = json.Unmarshal(message, result)
		if err != nil {
			return err
		}
		if result.Status != "ok" {
			return fmt.Errorf("[%s] %s", result.Exception.SQLCode, result.Exception.Text)
		}
		if response == nil {
			return nil
		}
		return json.Unmarshal(result.ResponseData, response)
	}, nil
}

func (connection *websocketConnection) sendRequestWithStringResponse(request interface{}) (string, error) {
	requestJson, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	messageType := websocket.TextMessage
	err = connection.websocket.WriteMessage(messageType, requestJson)
	if err != nil {
		return "", err
	}
	_, message, err := connection.websocket.ReadMessage()
	if err != nil {
		return "", err
	}
	return string(message), nil
}
