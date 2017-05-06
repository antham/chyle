package chyle

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "CHYLE - ", log.Ldate|log.Ltime)
}
