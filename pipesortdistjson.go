package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Team struct {
	Name   string `json:"Name"`
	Tokens int    `json:"Tokens"`
}

type Mensaje struct {
	Numero int
}

var min int
var cont int
var chCont chan int
var chTeam chan Team
var numero int
var addressLocal string
var addressRemoto string
var teamStart Team
var teamsArr []Team
var cantTeams int


func main() {

	//lectura por consola del host origen
	brIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del host local: ")
	puertoLocal, _ := brIn.ReadString('\n')
	puertoLocal = strings.TrimSpace(puertoLocal)
	addressLocal = fmt.Sprintf("localhost:%s", puertoLocal)

	//lectura por consola del host destino
	brIn = bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del host remoto: ")
	puertoRemoto, _ := brIn.ReadString('\n')
	puertoRemoto = strings.TrimSpace(puertoRemoto)
	addressRemoto = fmt.Sprintf("localhost:%s", puertoRemoto)

	//lectura de nro de mensajes a recibir
	brIn = bufio.NewReader(os.Stdin)
	fmt.Print("Ingreese el numero de mensajes a recibir: ")
	numstr, _ := brIn.ReadString('\n')
	numstr = strings.TrimSpace(numstr)
	numero, _ = strconv.Atoi(numstr)
	cantTeams=0
	//creamos canal
	// chTeam = make(chan Team, 1)
	// chTeam <- player

	//habilitar el modo escucha (servidor) nodo local
	// ln, _ := net.Listen("tcp", addressLocal)
	// defer ln.Close()

	//manejo de concurrencia para aceptar conexion de clientes

	ln := escuchar()
	defer ln.Close()

	for cantTeams < 2 {
		conn, _ := ln.Accept()
		go manejador(conn)

	}
	game()

}

func escuchar() net.Listener {
	conn, _ := net.Listen("tcp", addressLocal)
	return conn
}

func manejador(conn net.Conn) {

	defer conn.Close()
	br := bufio.NewReader(conn)
	msgJson, _ := br.ReadString('\n')

	//deserializar
	json.Unmarshal([]byte(msgJson), &teamStart)
	fmt.Println("Mensaje recibido: ")
	fmt.Println(teamStart)

	teamsArr = append(teamsArr, teamStart)
	cantTeams+=1
	for _, t := range teamsArr {
		fmt.Println(t.Name)
		fmt.Println(t.Tokens)
	}

}

func enviar(num int) {
	conn, _ := net.Dial("tcp", addressRemoto)
	defer conn.Close()
	// fmt.Fprintf(conn, "%d\n", num)

	teamStart.Tokens = num

	//Serializar

	arrBytesMsg, _ := json.MarshalIndent(teamStart, "", " ")
	jsonStrMsg := string(arrBytesMsg)

	fmt.Println("Mensaje enviado: ")
	fmt.Println(jsonStrMsg)

	fmt.Fprintf(conn, jsonStrMsg)
}

type Player struct {
	Team   *Team
	Name   string
	InGame bool
}

func game() {
	// teams := []*Team{
	// 	{Name: "Team 1", Tokens: 1},
	// 	{Name: "Team 2", Tokens: 1},
	// 	{Name: "Team 3", Tokens: 1},
	// }
	var teams []*Team
	for _, t := range teamsArr {
		teams = append(teams, &Team{Name: t.Name, Tokens: t.Tokens})
	}

	players := make([]*Player, 0)
	for ind, team := range teams {
		players = append(players, &Player{Team: team, Name: "Player " + strconv.Itoa(ind+1), InGame: true})
	}

	// for _, p := range players {
	// 	fmt.Println(p)
	// }

	rand.Seed(time.Now().UnixNano())

	playingBoard := createPlayingBoard(len(teams))
	currentPlayerIndex := 0

	var wg sync.WaitGroup
	moveCh := make(chan string)

	wg.Add(len(players))

	for i := 0; i < len(players); i++ {
		go func(playerIndex int) {
			defer wg.Done()
			player := players[playerIndex]

			for player.InGame {
				if isGameFinished(players) {
					break
				}
				if currentPlayerIndex != playerIndex {
					continue
				}
				// if !players[currentPlayerIndex].InGame{break}
				move := <-moveCh
				if move == "Bucket" {
					fmt.Printf("%s from %s reached another team's Bucket!\n", player.Name, player.Team.Name)
					otherTeamIndex := getRandomTeamIndex(player.Team, teams)
					otherTeam := teams[otherTeamIndex]
					otherPlayer := players[otherTeamIndex]
					tokenObtained := obtainToken(player.Team, otherTeam)
					if tokenObtained {
						fmt.Printf("%s obtained a token from %s! Tokens: %d\n", player.Name, otherTeam.Name, player.Team.Tokens)
						fmt.Printf("%s lost a token! Tokens: %d\n", otherTeam.Name, otherTeam.Tokens)
						if otherTeam.Tokens == 0 {
							fmt.Printf("%s has no more tokens. They are out of the game.\n", otherTeam.Name)
							otherPlayer.InGame = false
							// deletePlayer(teams,otherTeamIndex)
							// isGameFinished(players)
							// player.InGame = false
						}
					}
					if player.Team.Tokens == 0 {
						player.InGame = false
						// deletePlayer(teams,currentPlayerIndex)
						fmt.Printf("%s from %s has no more tokens. They are out of the game.\n", player.Name, player.Team.Name)
					}
				} else if move == "Player" {

					otherPlayerIndex := getRandomPlayerIndex(playerIndex, players)
					otherPlayer := players[otherPlayerIndex]

					fmt.Printf("%s from %s encountered a player from %s!\n", player.Name, player.Team.Name, otherPlayer.Team.Name)
					rpsResult := playRPS()

					if rpsResult == "Win" {
						fmt.Printf("%s from %s won in Rock, Paper, Scissors!\n", player.Name, player.Team.Name)
						currentPlayerIndex = otherPlayerIndex
					} else {
						fmt.Printf("%s from %s lost in Rock, Paper, Scissors!\n", player.Name, player.Team.Name)
						// player.InGame = false
						// fmt.Printf("%s from %s is out of the game.\n", player.Team.Name, player.Team.Name)
					}
				} else {
					fmt.Printf("%s from %s made a jump!\n", player.Name, player.Team.Name)
					currentPlayerIndex = getNextPlayerIndex(currentPlayerIndex, len(players), players)
				}
			}

		}(i)
		for !isGameFinished(players) {
			currentPlayer := players[currentPlayerIndex]
			if currentPlayer.InGame {
				fmt.Printf("%s from %s turn!\n", currentPlayer.Name, currentPlayer.Team.Name)
				fmt.Println("Jump into each hoop to move across the board.")

				moveCh <- makeJump(playingBoard, currentPlayerIndex)
			}
			time.Sleep(2 * time.Second)
		}

		close(moveCh)
		wg.Wait()

		fmt.Println("Game over!")
		for _, team := range teams {
			if team.Tokens != 0 {
				fmt.Printf("%s is the winner, they collected %d tokens.\n", team.Name, team.Tokens)
			}
		}
	}
}

func getNextPlayerIndex(currentPlayerIndex, numPlayers int, players []*Player) int {
	nextPlayerIndex := (currentPlayerIndex + 1) % numPlayers
	for !players[nextPlayerIndex].InGame {
		nextPlayerIndex = (nextPlayerIndex + 1) % numPlayers
	}
	return nextPlayerIndex
}

func createPlayingBoard(numTeams int) []string {
	playingBoard := make([]string, numTeams+1)
	for i := 0; i < numTeams; i++ {
		playingBoard[i] = "Hoop"
	}
	playingBoard[numTeams] = "Bucket"
	return playingBoard
}

func makeJump(playingBoard []string, currentPlayerIndex int) string {
	jumpOptions := []string{"Hoop", "Hoop", "Hoop", "Bucket", "Player"}
	move := jumpOptions[rand.Intn(len(jumpOptions))]

	if move == "Player" && playingBoard[currentPlayerIndex] != "Hoop" {
		move = "Hoop"
	}

	return move
}

func getRandomPlayerIndex(currentPlayerIndex int, players []*Player) int {
	otherPlayerIndex := rand.Intn(len(players))
	for otherPlayerIndex == currentPlayerIndex || !players[otherPlayerIndex].InGame {
		otherPlayerIndex = rand.Intn(len(players))
	}
	return otherPlayerIndex
}

func getRandomTeamIndex(currentTeam *Team, teams []*Team) int {
	otherTeamIndex := rand.Intn(len(teams))
	for otherTeamIndex == getTeamIndex(currentTeam, teams) || teams[otherTeamIndex].Tokens == 0 {
		otherTeamIndex = rand.Intn(len(teams))
	}
	return otherTeamIndex
}

func deletePlayer(teams []*Team, index int) {
	result := append(teams[:index], teams[index+1:]...)
	teams = result
}

func getTeamIndex(team *Team, teams []*Team) int {
	for i, t := range teams {
		if t == team {
			return i
		}
	}
	return -1
}

func playRPS() string {
	rpsOptions := []string{"Rock", "Paper", "Scissors"}
	playerChoice := rpsOptions[rand.Intn(len(rpsOptions))]
	otherPlayerChoice := rpsOptions[rand.Intn(len(rpsOptions))]

	fmt.Printf("You chose %s. The other player chose %s.\n", playerChoice, otherPlayerChoice)

	switch {
	case playerChoice == otherPlayerChoice:
		return "Tie"
	case playerChoice == "Rock" && otherPlayerChoice == "Scissors":
		return "Win"
	case playerChoice == "Paper" && otherPlayerChoice == "Rock":
		return "Win"
	case playerChoice == "Scissors" && otherPlayerChoice == "Paper":
		return "Win"
	default:
		return "Loss"
	}
}

func obtainToken(currentTeam *Team, otherTeam *Team) bool {
	if otherTeam.Tokens > 0 {
		currentTeam.Tokens++
		otherTeam.Tokens--
		return true
	}
	return false
}

func isGameFinished(players []*Player) bool {
	numPlayers := len(players)
	numPlayersLost := 0

	for _, player := range players {

		if player.InGame == false {
			numPlayersLost++
		}
	}
	return numPlayersLost == numPlayers
}
