package main

import (
	"fmt"
	"net"
	"os"
)

const (
	connHost = "localhost" // Por ahora en entorno local.
	connPort = "8080"      // Puerto al que se van a conectar los clientes.
	connType = "tcp"
)

func main() {

	server, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error conectando:", err.Error())
		os.Exit(1)
	}

}
