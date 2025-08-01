package main

import (
	"fmt"

	"github.com/inlovewithgo/transit-backend/main/config"
)

func init() {
    config.InitDatabase()
}

func main() {

    fmt.Println("Transit Backend Service is starting...")

}