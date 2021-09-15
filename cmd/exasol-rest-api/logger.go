package exasol_rest_api

import (
	"log"
	"os"
)

var errorLogger = Logger(log.New(os.Stderr, "[exasol] ", log.LstdFlags|log.Lshortfile))

// Logger is used to log critical error messages.
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}
