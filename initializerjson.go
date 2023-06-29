package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Player struct {
	Name      string
	Tokens    int
	positionX int
	positionY int
	index     int
	inGame    bool
	delay     int
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


	fmt.Print("Ingrese su nombre: ")
	name, _ := brIn.ReadString('\n')
	name = strings.TrimSpace(name)

	var player Player
	player.Name = name
	player.Tokens = 2
	player.inGame=true
	player.positionX=1
	enviarJson(player)
}
func enviarJson(player Player) {
	conn, _ := net.Dial("tcp", addressRemoto)
	defer conn.Close()

	//Serializar
	arrBytesMsg, _ := json.Marshal(player)
	jsonStrMsg := string(arrBytesMsg)

	fmt.Println("Mensaje enviado: ")
	// fmt.Println(jsonStrMsg)
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
