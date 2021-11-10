package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"os"
)

type Client struct {
	nombre string
	socket net.Conn
}

const (
	connHost = "localhost" // Por ahora el programa funciona en el entorno local, mas adelante se analizara si se extiende a que funcione entre diferentes equipos.
	connPort = "8080"
	connType = "tcp" // Protocolo de transmision de mensajes.
)

func main() {

	clients := list.New()

	listener, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error escuchando:", err.Error())
		os.Exit(1)
	}
	defer listener.Close() //Siempre se cerrara el Listener.

	for {
		newClientSocket, err := listener.Accept()
		if err != nil {
			fmt.Println("Error conectando:", err.Error())
			return
		}

		clientLog(newClientSocket, clients)

		newClient := clients.Back().Value.(Client)

		go recibirMensajesDeClientes(newClient) // En esta parte manejamos los mensajes entre servidor y cliente.
	}

}

func recibirMensajesDeClientes(client Client) {

	for { // Constantemente estaremos escuchando mensajes de los clientes.

		buffer, err := bufio.NewReader(client.socket).ReadBytes('\n') // Leo el mensaje del cliente, funcion bloqueante.

		// Cierro las conexiones cuando el cliente se va.
		if err != nil {
			fmt.Println("Se fue el cliente.")
			client.socket.Close()
			return
		}

		mensaje := client.nombre + string(buffer)

		fmt.Println(mensaje)
	}

}

func clientLog(clientSocket net.Conn, clients *list.List) {
	buffer, err := bufio.NewReader(clientSocket).ReadBytes('\n') // Leo el mensaje del cliente, funcion bloqueante.

	// Cierro las conexiones cuando el cliente se va.
	if err != nil {
		fmt.Println("Se fue el cliente.")
		clientSocket.Close()
		return
	}

	newClient := Client{
		nombre: string(buffer),
		socket: clientSocket,
	}

	clients.PushBack(newClient)
}
