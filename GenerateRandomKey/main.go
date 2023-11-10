package main

import (
	"encoding/hex"
	"log"

	"github.com/gorilla/securecookie"
)

func main() {
	log.Println()
	log.Println(hex.EncodeToString(securecookie.GenerateRandomKey(64)))
	log.Println(hex.EncodeToString(securecookie.GenerateRandomKey(32)))
	log.Println()
}
