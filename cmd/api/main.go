// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"fmt"

	"github.com/Owbird/KNUST-AIM-API/internal/server"
)

func main() {

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
