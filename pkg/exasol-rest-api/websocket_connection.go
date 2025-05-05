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
	"math/big"
	"net/url"
	"os/user"
	"runtime"
	"strconv"

	exaerror "github.com/exasol/error-reporting-go"
	"github.com/gorilla/websocket"
)

type websocketConnection struct {
	connProperties *ApplicationProperties
	websocket      *websocket.Conn
}

// [impl->dsn~communicate-with-database~1]
func (connection *websocketConnection) connect() error {
	uri := fmt.Sprintf("%s:%d", connection.connProperties.ExasolHost, connection.connProperties.ExasolPort)
	exaURL := url.URL{
		Scheme: connection.getURIScheme(),
		Host:   uri,
	}
	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: getInsecureConnection(connection)}

	websocketConnection, _, err := dialer.DialContext(context.Background(), exaURL.String(), nil)
	if err == nil {
		connection.websocket = websocketConnection
		connection.websocket.EnableWriteCompression(false)
		return nil
	} else {
		return exaerror.New("E-ERA-14").
			Message("error while establishing a websockets connection: {{error|uq}}").
			Parameter("error", err.Error())
	}
}

func getInsecureConnection(connection *websocketConnection) bool {
	return connection.connProperties.UseTLS != 1
}

func (connection *websocketConnection) close() {
	err := connection.send(&command{Command: "disconnect"}, nil)
	if err != nil {
		errorLogger.Printf("error closing a websockets connection: %s", err)
	}
	err = connection.websocket.Close()
	if err != nil {
		errorLogger.Printf("error closing websocket: %s", err)
	}
	connection.websocket = nil
}

func (connection *websocketConnection) getURIScheme() string {
	if connection.connProperties.Encryption == 1 {
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
			Autocommit:       true,
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
		return exaerror.New("E-ERA-15").
			Message("error while sending a login command via websockets connection: {{error|uq}}").
			Parameter("error", err.Error())
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
		return exaerror.New("F-ERA-21").
			Message("password encryption error during login via websockets connection: {{error|uq}}").
			Parameter("error", err.Error())
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

func (connection *websocketConnection) sendRequestWithInterfaceResponse(request interface{}) (func(interface{}) error,
	error) {
	message, err := connection.sendRequestWithStringResponse(request)
	if err != nil {
		return nil, err
	}

	return func(responseType interface{}) error {
		result := &webSocketsBaseResponse{}
		err = json.Unmarshal(message, result)

		if err != nil {
			return exaerror.New("F-ERA-27").
				Message("error converting JSON message from websockets into response struct: {{error|uq}}").
				Parameter("error", err.Error())
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
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, exaerror.New("F-ERA-24").
			Message("cannot convert request into JSON format: {{error|uq}}").
			Parameter("error", err.Error())
	}

	messageType := websocket.TextMessage
	err = connection.websocket.WriteMessage(messageType, requestJSON)
	if err != nil {
		return nil, exaerror.New("F-ERA-25").
			Message("error writing a message via websocket connection: {{error|uq}}").
			Parameter("error", err.Error())
	}

	_, message, err := connection.websocket.ReadMessage()
	if err != nil {
		return nil, exaerror.New("F-ERA-26").
			Message("error reading a message from websocket: {{error|uq}}").
			Parameter("error", err.Error())
	}

	return message, nil
}
