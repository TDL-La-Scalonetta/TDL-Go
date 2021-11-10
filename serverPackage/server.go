package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"os"
)

type client struct {
	nombre string
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
		newClient, err := listener.Accept()
		if err != nil {
			fmt.Println("Error conectando:", err.Error())
			return
		}

		clientLog(newClient, clients)

		fmt.Println("El nombre del ultimo cliente conectado es ", clients.Front().Value.(client).nombre)

		go handleConnection(newClient) // En esta parte manejamos los mensajes entre servidor y cliente.
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

	reader := bufio.NewReader(os.Stdin) // Usamos la variable reader para leer del teclado.

	fmt.Print("Mensaje a mandarle al cliente: ")
	input, _ := reader.ReadString('\n')

	client.Write([]byte(input)) // Aca le respondemos al cliente.

	// Repetimos el proceso hasta que el cliente se vaya.
	handleConnection(client)
}

func clientLog(clientSocket net.Conn, clients *list.List) {
	buffer, err := bufio.NewReader(clientSocket).ReadBytes('\n') // Leo el mensaje del cliente, funcion bloqueante.

	// Cierro las conexiones cuando el cliente se va.
	if err != nil {
		fmt.Println("Se fue el cliente.")
		clientSocket.Close()
		return
	}

	newClient := client{
		nombre: string(buffer[:len(buffer)-1]),
	}

	clients.PushBack(newClient)
}
