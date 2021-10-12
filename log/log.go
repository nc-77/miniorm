package log

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	infoLog  = log.New(os.Stdout, "\033[34m[info]\033[0m ", log.LstdFlags)
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags)
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	InfoLevel  = iota // InfoLevel prints all logs on stdout
	ErrorLevel        // ErrorLevel prints only error logs on stdout
	Disabled          // Disabled does not print any logs
)

func SetLevel(level int) {
	if level >= ErrorLevel {
		infoLog.SetOutput(ioutil.Discard)
	}
	if level >= Disabled {
		errorLog.SetOutput(ioutil.Discard)
	}
}
