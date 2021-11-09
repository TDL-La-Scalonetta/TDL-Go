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
		fmt.Println("Error escuchando:", err.Error())
		os.Exit(1)
	}
	defer listener.Close() //Siempre se cerrara el Listener.

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Error conectando:", err.Error())
			return
		}

		fmt.Println("Client " + client.RemoteAddr().String() + " connected.")
	}

}
