package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	connHost = "localhost" // Por ahora el programa funciona en el entorno local, mas adelante se analizara si se extiende a que funcione entre diferentes equipos.
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

		go handleConnection(client) // En esta parte manejamos los mensajes entre servidor y cliente.
	}

}

func handleConnection(client net.Conn) {
	buffer, err := bufio.NewReader(client).ReadBytes('\n') // Leo el mensaje del cliente, funcion bloqueante.

	// Cierro las conexiones cuando el cliente se va.
	if err != nil {
		fmt.Println("Se fue el cliente.")
		client.Close()
		return
	}

	fmt.Println("Mensaje del Cliente:", string(buffer[:len(buffer)-1]))

	// Le respondemos al cliente.
	respuesta := []byte("Gracias por hablarle al server")
	client.Write(respuesta)

	// Repetimos el proceso hasta que el cliente se vaya.
	handleConnection(client)
}
