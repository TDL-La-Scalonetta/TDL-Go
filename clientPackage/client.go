package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	connHost = "localhost" // Por ahora, el programa va a funcionar en un entorno local.
	connPort = "8080"      // Puerto al que se van a conectar los clientes.
	connType = "tcp"       // Protocolo de comunicacion.
)

func main() {

	server, err := net.Dial(connType, connHost+":"+connPort) // Me conecto al servidor.
	if err != nil {
		fmt.Println("Error conectando:", err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin) // Creamos un reador para poder leer el input del teclado.

	clientLog(reader, server) // El usuario ingresa sus datos.

	eleccionDeProducto(reader, server)

	go actualizarSubasta(server) // Goroutine utilizada para actualizar los datos de la subasta provistos por el servidor.

	ejecutarAccionDeSubasta(reader, server) // Ofertar o retirarse.

}

func clientLog(reader *bufio.Reader, server net.Conn) {
	fmt.Print("\n\nBienvenido a Scalonetta, el mejor sitio para las subastas! \n\n")
	fmt.Print("Por favor ingrese su nombre de usuario: ")
	input, _ := reader.ReadString('\n')

	server.Write([]byte(input)) // Le mando el mensaje al Servidor.
}

func actualizarSubasta(server net.Conn) {
	for { //Constantemente estamos escuchando por nuevas actualizaciones de la subasta.
		message, _ := bufio.NewReader(server).ReadString('.')

		fmt.Print("\n\n" + message)

	}
}

func ejecutarAccionDeSubasta(reader *bufio.Reader, server net.Conn) {

	for {
		fmt.Print("\n\nElija que desea hacer con el producto: \n\n")
		fmt.Print("\n\nA): Hacer una oferta. \n")
		fmt.Print("\n\nB): Retirarse. \n")

		input, _ := reader.ReadString('\n')

		if strings.Contains(input, "A") { // Si se elige la opcion de ofertar.
			fmt.Print("\n\nPor favor ingrese el monto a ofertar: ")
			input, _ = reader.ReadString('\n')
			server.Write([]byte(input)) // Le paso el monto de la oferta.
		} else {
			fmt.Print("\n\nGracias por participar de la subasta! \n\n")
			server.Write([]byte("Abandona subasta.\n")) // Le aviso que el cliente se retira de la subasta.
			return
		}
	}
}

func eleccionDeProducto(reader *bufio.Reader, server net.Conn) {

	fmt.Print("\nPor favor elija el producto que desea obtener: \n")

	fmt.Print("\nA) Bicicleta.\n")
	fmt.Print("\nB) Auto.\n")
	fmt.Print("\nC) Avión.\n\n")

	fmt.Print("Escriba su opcion (A, B o C): ")
	input, _ := reader.ReadString('\n')

	server.Write([]byte(input)) // Le mando la opcion elegida al servidor.

	fmt.Print("\nUsted ha elegido la opcion ", input)

	salaDeEspera(server)

}

func salaDeEspera(server net.Conn) {

	fmt.Print("\nPor favor espere a que ingrese una persona mas para poder comenzar con la subasta.\n\n")

	message, _ := bufio.NewReader(server).ReadString('.') // Esperamos a que llegue un cliente mas que desea el mismo producto.

	fmt.Print("\n\n" + message)
}
