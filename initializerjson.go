package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Team struct {
	Name   string `json:"name"`
	Tokens int    `json:"tokens"`
}

type Mensaje struct {
	Numero int
}

var addressRemoto string
var mensaje Mensaje

func main() {
	brIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del host remoto: ")
	puertoRemoto, _ := brIn.ReadString('\n')
	puertoRemoto = strings.TrimSpace(puertoRemoto)
	addressRemoto = fmt.Sprintf("localhost:%s", puertoRemoto)

	// enviar(6)
	// enviar(3)
	// enviar(1)
	// enviar(5)
	// enviarMensaje(6)
	// enviarMensaje(3)
	// enviarMensaje(1)
	// enviarMensaje(5)

	fmt.Print("Ingrese su nombre: ")
	name, _ := brIn.ReadString('\n')
	name = strings.TrimSpace(name)

	var team Team
	team.Name = name
	team.Tokens = 3

	enviarJson(team)
}
func enviarJson(team Team) {
	conn, _ := net.Dial("tcp", addressRemoto)
	defer conn.Close()

	//Serializar
	arrBytesMsg, _ := json.Marshal(team)
	jsonStrMsg := string(arrBytesMsg)

	fmt.Println("Mensaje enviado: ")
	fmt.Println(jsonStrMsg)
	fmt.Fprintf(conn, jsonStrMsg)

}
func enviarMensaje(num int) {
	conn, _ := net.Dial("tcp", addressRemoto)
	defer conn.Close()

	mensaje.Numero = num

	//Serializar
	arrBytesMsg, _ := json.Marshal(mensaje)
	jsonStrMsg := string(arrBytesMsg)

	fmt.Println("Mensaje enviado: ")
	fmt.Println(jsonStrMsg)
	fmt.Fprintf(conn, jsonStrMsg)

}
