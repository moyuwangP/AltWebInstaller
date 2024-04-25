package util

import (
	"log"
	"os"
)

var errLogger *log.Logger
var debuLogger *log.Logger

func init() {
	errLogger = log.New(os.Stderr, "[ERROR]", 0)
	debuLogger = log.New(os.Stdout, "[DEBU]", 0)
}

func LogError(err string) {
	errLogger.Println(err)
}

func LogDebug(err string) {
	debuLogger.Println(err)
}

func LogErrorf(err string, arg ...string) {
	errLogger.Printf(err, arg)
}
