package main

import (
	"bufio"
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

	// Creamos un reador para poder leer de input del teclado.
	reader := bufio.NewReader(os.Stdin)

	clientLog(reader, server)

	go recibirMensajesDelServer(server)

	// Por ahora, un loop infinito de mensajes entre server y clientes.
	for {

		fmt.Print("Escriba su mensaje: ")
		input, _ := reader.ReadString('\n')

		// Le mando el mensaje al Servidor.
		server.Write([]byte(input))

	}

}

func clientLog(reader *bufio.Reader, server net.Conn) {
	fmt.Print("Por favor ingrese su nombre de usuario: ")
	input, _ := reader.ReadString('\n')

	// Le mando el mensaje al Servidor.
	server.Write([]byte(input))
}

func recibirMensajesDelServer(server net.Conn) {
	for { //Todo el tiempo vamos a tener que estar escuchando por nuevos mensajes del servidor.
		message, _ := bufio.NewReader(server).ReadString('.')

		fmt.Print("\n\n" + message)

		fmt.Print("\n\nEscriba su mensaje: ")
	}
}
