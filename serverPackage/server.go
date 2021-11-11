package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"os"
)

type Client struct {
	nombre          string
	socket          net.Conn
	productoElegido string
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

		eleccionDeProducto(clients)

		go reenviarMensajesDeClientes(newClient, clients)
	}

}

func reenviarMensajesDeClientes(client Client, clients *list.List) {

	for { // Constantemente estaremos escuchando mensajes de los clientes.

		buffer, err := bufio.NewReader(client.socket).ReadBytes('\n') // Leo el mensaje del cliente, funcion bloqueante.

		// Cierro las conexiones cuando el cliente se va.
		if err != nil {
			fmt.Println("Se fue el cliente.")
			client.socket.Close()
			return
		}

		mensaje := client.nombre + string(buffer) + "." // Esto del punto lo usamos para señalar donde termina el mensaje. Es temporal.

		fmt.Println(mensaje)

		for c := clients.Front(); c != nil; c = c.Next() {
			c.Value.(Client).socket.Write([]byte(mensaje))
		}
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

func eleccionDeProducto(clients *list.List) {
	buffer, err := bufio.NewReader(clients.Back().Value.(Client).socket).ReadBytes('\n')

	if err != nil {
		fmt.Println("Se fue el cliente.")
		clients.Back().Value.(Client).socket.Close()
		return
	}

	auxClient := clients.Back().Value.(Client)

	auxClient.productoElegido = string(buffer)

	clients.Back().Value = auxClient

}
