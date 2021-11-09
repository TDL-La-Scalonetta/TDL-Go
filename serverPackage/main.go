package main

import (
	"fmt"
	"net"
	"os"
)

const (
	connHost = "localhost" // Por ahora el programa funciona en el entorno local.
	connPort = "8080"
	connType = "tcp" // Protocolo de transmision de mensajes.
)

func main() {
	listener, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close() //Siempre se cerrara el Listener.

}
