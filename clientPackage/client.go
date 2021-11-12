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

	eleccionDeProducto(reader, server)

	enviarAccionDeCliente(reader, server)

	go actualizarSubasta(server)

}

func clientLog(reader *bufio.Reader, server net.Conn) {
	fmt.Print("\n\nBienvenido a Scalonetta, el mejor sitio para las subastas! \n\n")
	fmt.Print("Por favor ingrese su nombre de usuario: ")
	input, _ := reader.ReadString('\n')

	// Le mando el mensaje al Servidor.
	server.Write([]byte(input))
}

func actualizarSubasta(server net.Conn) {
	for { //Todo el tiempo vamos a tener que estar escuchando por nuevos mensajes del servidor.
		message, _ := bufio.NewReader(server).ReadString('.')

		fmt.Print("\n\n" + message)

	}
}

func enviarAccionDeCliente(reader *bufio.Reader, server net.Conn) {

	for {
		fmt.Print("\n\nElija que desea hacer con el producto: \n\n")
		fmt.Print("\n\nA): Hacer una oferta. \n")
		fmt.Print("\n\nB): Retirarse. \n")

		input, _ := reader.ReadString('\n')

		// Le mando el mensaje al Servidor.
		server.Write([]byte(input))
	}
}

func eleccionDeProducto(reader *bufio.Reader, server net.Conn) {

	fmt.Print("\nPor favor elija el producto que desea obtener: \n")

	fmt.Print("\nA) Bicicleta.\n")
	fmt.Print("\nB) Auto.\n")
	fmt.Print("\nC) Avión.\n\n")

	fmt.Print("Escriba su opcion (A, B o C): ")
	input, _ := reader.ReadString('\n')

	// Le mando el mensaje al Servidor.
	server.Write([]byte(input))

	fmt.Print("\nUsted ha elegido la opcion ", input)
	fmt.Print("\nPor favor espere a que ingrese una persona mas para poder comenzar con la subasta.\n\n")

	// Esperamos a que llegue un cliente mas que desea el mismo producto.

	message, _ := bufio.NewReader(server).ReadString('.')

	fmt.Print("\n\n" + message)

}
