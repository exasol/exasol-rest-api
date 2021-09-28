package exasol_rest_api

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
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
	connProperties *ApplicationProperties
	websocket      *websocket.Conn
}

func (connection *websocketConnection) connect() error {
	uri := fmt.Sprintf("%s:%d", connection.connProperties.ExasolHost, connection.connProperties.ExasolPort)
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

func (connection *websocketConnection) close() {
	err := connection.send(&command{Command: "disconnect"}, nil)
	connection.websocket.Close()
	connection.websocket = nil
	if err != nil {
		errorLogger.Printf("error closing a connection: %s", err)
	}
}

func (connection *websocketConnection) getURIScheme() string {
	if connection.connProperties.Encryption {
		return "wss"
	} else {
		return "ws"
	}
}

func (connection *websocketConnection) executeQuery(query string) ([]byte, error) {
	command := &sqlCommand{
		command: command{"execute"},
		SQLText: query,
		Attributes: attributes{
			ResultSetMaxRows: 1000,
		},
	}
	return connection.sendRequestWithStringResponse(command)
}

func (connection *websocketConnection) login() error {
	loginCommand := &loginCommand{
		command:         command{"login"},
		ProtocolVersion: connection.connProperties.ExasolWebsocketAPIVersion,
	}
	loginResponse := &publicKeyResponse{}
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
	password := []byte(connection.connProperties.ExasolPassword)
	encPass, err := rsa.EncryptPKCS1v15(rand.Reader, &pubKey, password)
	if err != nil {
		errorLogger.Printf("password encryption error: %s", err)
		return err
	}
	b64Pass := base64.StdEncoding.EncodeToString(encPass)

	authRequest := authCommand{
		Username:       connection.connProperties.ExasolUser,
		Password:       b64Pass,
		UseCompression: false,
		ClientName:     "Exasol REST API",
		ClientOs:       runtime.GOOS,
		ClientRuntime:  runtime.Version(),
	}

	if osUser, err := user.Current(); err != nil {
		authRequest.ClientOsUsername = osUser.Username
	}

	return connection.send(authRequest, nil)
}

func (connection *websocketConnection) send(request, response interface{}) error {
	receiver, err := connection.sendRequestWithInterfaceResponse(request)
	if err != nil {
		return err
	}
	return receiver(response)
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

	return func(responseType interface{}) error {
		_, message, err := connection.websocket.ReadMessage()
		if err != nil {
			return err
		}
		result := &baseResponse{}
		err = json.Unmarshal(message, result)
		if err != nil {
			return err
		}
		if result.Status != "ok" {
			return fmt.Errorf("[%s] %s", result.Exception.SQLCode, result.Exception.Text)
		}
		if responseType == nil {
			return nil
		}
		return json.Unmarshal(result.ResponseData, responseType)
	}, nil
}

func (connection *websocketConnection) sendRequestWithStringResponse(request interface{}) ([]byte, error) {
	requestJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	messageType := websocket.TextMessage
	err = connection.websocket.WriteMessage(messageType, requestJson)
	if err != nil {
		return nil, err
	}
	_, message, err := connection.websocket.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}
