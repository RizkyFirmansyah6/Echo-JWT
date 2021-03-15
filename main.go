package main

import (
	"EchoAPI/API"
	_ "fmt"
	_ "strings"
)

type M map[string]interface{}

func main() {
	API.StartServer("")
}
