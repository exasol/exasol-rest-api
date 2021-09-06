package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type DSNConfig struct {
	host          string
	port          int
	user          string
	password      string
	autocommit    *bool
	encryption    *bool
	compression   *bool
	clientName    string
	clientVersion string
	fetchSize     int
	useTLS        *bool
}

func NewConfig(user, password string) *DSNConfig {
	return &DSNConfig{
		host:     "localhost",
		port:     8563,
		user:     user,
		password: password,
	}
}

func (c *DSNConfig) Compression(enabled bool) *DSNConfig {
	c.compression = &enabled
	return c
}
func (c *DSNConfig) Encryption(enabled bool) *DSNConfig {
	c.encryption = &enabled
	return c
}
func (c *DSNConfig) Autocommit(enabled bool) *DSNConfig {
	c.autocommit = &enabled
	return c
}
func (c *DSNConfig) UseTLS(enabled bool) *DSNConfig {
	c.useTLS = &enabled
	return c
}
func (c *DSNConfig) FetchSize(size int) *DSNConfig {
	c.fetchSize = size
	return c
}
func (c *DSNConfig) ClientName(name string) *DSNConfig {
	c.clientName = name
	return c
}
func (c *DSNConfig) ClientVersion(version string) *DSNConfig {
	c.clientVersion = version
	return c
}
func (c *DSNConfig) Host(host string) *DSNConfig {
	c.host = host
	return c
}
func (c *DSNConfig) Port(port int) *DSNConfig {
	c.port = port
	return c
}

func (c *DSNConfig) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("exa:%s:%d;user=%s;password=%s;", c.host, c.port, c.user, c.password))
	if c.autocommit != nil {
		sb.WriteString(fmt.Sprintf("autocommit=%d;", boolToInt(*c.autocommit)))
	}
	if c.compression != nil {
		sb.WriteString(fmt.Sprintf("compression=%d;", boolToInt(*c.compression)))
	}
	if c.encryption != nil {
		sb.WriteString(fmt.Sprintf("encryption=%d;", boolToInt(*c.encryption)))
	}
	if c.useTLS != nil {
		sb.WriteString(fmt.Sprintf("usetls=%d;", boolToInt(*c.useTLS)))
	}
	if c.fetchSize != 0 {
		sb.WriteString(fmt.Sprintf("fetchsize=%d;", c.fetchSize))
	}
	if c.clientName != "" {
		sb.WriteString(fmt.Sprintf("clientname=%s;", c.clientName))
	}
	if c.clientVersion != "" {
		sb.WriteString(fmt.Sprintf("clientversion=%s;", c.clientVersion))
	}
	return strings.TrimRight(sb.String(), ";")
}

func parseDSN(dsn string) (*config, error) {
	if !strings.HasPrefix(dsn, "exa:") {
		return nil, fmt.Errorf("invalid connection string, must start with 'exa:'")
	}

	splitDsn := splitIntoConnectionStringAndParameters(dsn)
	host, port, err := extractHostAndPort(splitDsn[0])
	if err != nil {
		return nil, err
	}

	if len(splitDsn) < 2 {
		return getDefaultConfig(host, port), nil
	} else {
		return getConfigWithParameters(host, port, splitDsn[1])
	}
}

func splitIntoConnectionStringAndParameters(dsn string) []string {
	cleanDsn := strings.Replace(dsn, "exa:", "", 1)
	return strings.SplitN(cleanDsn, ";", 2)
}

func extractHostAndPort(connectionString string) (string, int, error) {
	hostPort := strings.Split(connectionString, ":")
	if len(hostPort) != 2 {
		return "", 0, fmt.Errorf("invalid host or port, expected format: <host>:<port>")
	}
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid `port` value, numeric port expected")
	}
	return hostPort[0], port, nil
}

func getDefaultConfig(host string, port int) *config {
	return &config{
		Host:        host,
		Port:        port,
		ApiVersion:  2,
		Autocommit:  true,
		Encryption:  true,
		Compression: false,
		UseTLS:      true,
		ClientName:  "Go client",
		Params:      map[string]string{},
		FetchSize:   128 * 1024,
	}
}

func getConfigWithParameters(host string, port int, parametersString string) (*config, error) {
	config := getDefaultConfig(host, port)
	parameters := extractParameters(parametersString)
	for _, parameter := range parameters {
		parameter = strings.TrimRight(parameter, ";")
		keyValuePair := strings.SplitN(parameter, "=", 2)
		if len(keyValuePair) != 2 {
			return nil, fmt.Errorf("invalid parameter %s, expected format <parameter>=<value>", parameter)
		}
		key := keyValuePair[0]
		value := keyValuePair[1]

		switch key {
		case "password":
			config.Password = unescape(value, ";")
		case "user":
			config.User = unescape(value, ";")
		case "autocommit":
			config.Autocommit = value == "1"
		case "encryption":
			config.Encryption = value == "1"
		case "usetls":
			config.UseTLS = value == "1"
		case "compression":
			config.Compression = value == "1"
		case "clientname":
			config.ClientName = value
		case "clientversion":
			config.ClientVersion = value
		case "schema":
			config.Schema = value
		case "fetchsize":
			value, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid `fetchsize` value, numeric expected")
			}
			config.FetchSize = value
		case "resultsetmaxrows":
			value, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid `resultsetmaxrows` value, numeric expected")
			}
			config.ResultSetMaxRows = value
		default:
			config.Params[key] = unescape(value, ";")
		}
	}
	return config, nil
}

func extractParameters(parametersString string) []string {
	reg := regexp.MustCompile(`[\w];`)
	return splitAfter(parametersString, reg)
}

func unescape(s, char string) string {
	return strings.ReplaceAll(s, `\`+char, char)
}

func splitAfter(s string, re *regexp.Regexp) []string {
	var (
		r []string
		p int
	)
	is := re.FindAllStringIndex(s, -1)
	if is == nil {
		return append(r, s)
	}
	for _, i := range is {
		r = append(r, s[p:i[1]])
		p = i[1]
	}
	return append(r, s[p:])
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
