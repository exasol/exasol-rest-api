package exasol_rest_api

import (
	"bytes"
	"compress/zlib"
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
	connProperties *connectionProperties
	websocket      *websocket.Conn
}

func (connection *websocketConnection) close() error {
	err := connection.send(context.Background(), &Command{Command: "disconnect"}, nil)
	connection.websocket.Close()
	connection.websocket = nil
	return err
}

func (connection *websocketConnection) executeQuery(query string) (string, error) {
	command := &SQLCommand{
		Command: Command{"execute"},
		SQLText: query,
		Attributes: Attributes{
			ResultSetMaxRows: connection.connProperties.ResultSetMaxRows,
		},
	}
	result, err := connection.asyncSend2(command)
	if err != nil {
		return "", err
	} else {
		return result, err
	}
}

func (connection *websocketConnection) login() error {
	hasCompression := connection.connProperties.Compression
	connection.connProperties.Compression = false
	loginCommand := &LoginCommand{
		Command:         Command{"login"},
		ProtocolVersion: connection.connProperties.ApiVersion,
	}
	loginResponse := &PublicKeyResponse{}
	err := connection.send(context.Background(), loginCommand, loginResponse)
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
		ClientName:     connection.connProperties.ClientName,
		DriverName:     fmt.Sprintf("exasol-driver-go %s", "v1.0.0"),
		ClientOs:       runtime.GOOS,
		ClientVersion:  connection.connProperties.ClientName,
		ClientRuntime:  runtime.Version(),
		Attributes: Attributes{
			CurrentSchema:      connection.connProperties.Schema,
			CompressionEnabled: hasCompression,
		},
	}

	if osUser, err := user.Current(); err != nil {
		authRequest.ClientOsUsername = osUser.Username
	}

	err = connection.send(context.Background(), authRequest, nil)
	if err != nil {
		return err
	}
	connection.connProperties.Compression = hasCompression

	return nil
}

func (connection *websocketConnection) getURIScheme() string {
	if connection.connProperties.Encryption {
		return "wss"
	} else {
		return "ws"
	}
}

func (connection *websocketConnection) connect() error {
	host := connection.connProperties.Host
	uri := fmt.Sprintf("%s:%d", host, connection.connProperties.Port)
	u := url.URL{
		Scheme: connection.getURIScheme(),
		Host:   uri,
	}
	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: !connection.connProperties.UseTLS}

	ws, _, err := dialer.DialContext(context.Background(), u.String(), nil)
	if err == nil {
		connection.websocket = ws
		connection.websocket.EnableWriteCompression(false)
	}
	return err
}

func (connection *websocketConnection) send(ctx context.Context, request, response interface{}) error {
	receiver, err := connection.asyncSend(request)
	if err != nil {
		return err
	}
	channel := make(chan error, 1)
	go func() { channel <- receiver(response) }()
	select {
	case <-ctx.Done():
		_, err := connection.asyncSend(&Command{Command: "abortQuery"})
		if err != nil {
			return fmt.Errorf("could not abort query %w", ctx.Err())
		}
		return ctx.Err()
	case err := <-channel:
		return err
	}
}

func (connection *websocketConnection) asyncSend(request interface{}) (func(interface{}) error, error) {
	message, err := json.Marshal(request)
	if err != nil {
		ErrorLogger.Printf("could not marshal request, %s", err)
		return nil, driver.ErrBadConn
	}

	messageType := websocket.TextMessage
	if connection.connProperties.Compression {
		var b bytes.Buffer
		w := zlib.NewWriter(&b)
		_, err = w.Write(message)
		if err != nil {
			return nil, err
		}
		w.Close()
		message = b.Bytes()
		messageType = websocket.BinaryMessage
	}

	err = connection.websocket.WriteMessage(messageType, message)
	if err != nil {
		ErrorLogger.Printf("could not send request, %s", err)
		return nil, driver.ErrBadConn
	}

	return func(response interface{}) error {

		_, message, err := connection.websocket.ReadMessage()
		if err != nil {
			ErrorLogger.Printf("could not receive data, %s", err)
			return driver.ErrBadConn
		}

		result := &BaseResponse{}
		if connection.connProperties.Compression {
			b := bytes.NewReader(message)
			r, err := zlib.NewReader(b)
			if err != nil {
				ErrorLogger.Printf("could not decode compressed data, %s", err)
				return driver.ErrBadConn
			}
			err = json.NewDecoder(r).Decode(result)
			if err != nil {
				ErrorLogger.Printf("could not decode data, %s", err)
				return driver.ErrBadConn
			}

		} else {
			err = json.Unmarshal(message, result)
			if err != nil {
				ErrorLogger.Printf("could not receive data, %s", err)
				return driver.ErrBadConn
			}
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

func (connection *websocketConnection) asyncSend2(request interface{}) (string, error) {
	message, err := json.Marshal(request)
	if err != nil {
		ErrorLogger.Printf("could not marshal request, %s", err)
		return "", driver.ErrBadConn
	}

	messageType := websocket.TextMessage
	if connection.connProperties.Compression {
		var b bytes.Buffer
		w := zlib.NewWriter(&b)
		_, err = w.Write(message)
		if err != nil {
			return "", err
		}
		w.Close()
		message = b.Bytes()
		messageType = websocket.BinaryMessage
	}

	err = connection.websocket.WriteMessage(messageType, message)
	if err != nil {
		ErrorLogger.Printf("could not send request, %s", err)
		return "", driver.ErrBadConn
	}

	return connection.getResult()
}

func (connection *websocketConnection) getResult() (string, error) {
	_, message, err := connection.websocket.ReadMessage()
	if err != nil {
		ErrorLogger.Printf("could not receive data, %s", err)
		return "", driver.ErrBadConn
	}

	result := &BaseResponse{}
	if connection.connProperties.Compression {
		b := bytes.NewReader(message)
		r, err := zlib.NewReader(b)
		if err != nil {
			ErrorLogger.Printf("could not decode compressed data, %s", err)
			return "", driver.ErrBadConn
		}
		err = json.NewDecoder(r).Decode(result)
		if err != nil {
			ErrorLogger.Printf("could not decode data, %s", err)
			return "", driver.ErrBadConn
		}

	} else {
		err = json.Unmarshal(message, result)
		if err != nil {
			ErrorLogger.Printf("could not receive data, %s", err)
			return "", driver.ErrBadConn
		}
	}
	marshal, err := json.Marshal(result)
	return string(marshal), nil
}
