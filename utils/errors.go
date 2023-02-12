package utils

import (
	"log"
)

func FailsOnError(err error, message string) {
	if err != nil {
		log.Println(message)
		log.Fatal(err)
	}
}
