package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func getPortS() int {
	//check for arg //os.Args provides access to raw command-line arguments. Note that the first value in this slice is the path to the program, and os.Args[1:] holds the arguments to the program.
	if len(os.Args) != 2 {
		fmt.Printf("Vous devez utiliser le server ainsi : go run Server.go <portNumber>\n")
		os.Exit(1)
	} else {
		//l'arg doit etre int
		fmt.Printf("Vous avez indiqué le port :\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("\"Vous devez utiliser le server ainsi : go run Server.go <portNumber>\n")
			os.Exit(1)
		} else {
			return portNumber
		}
	}
	return -1 //ne devrait jamais être atteint
}

func main() {
	port := getPortS()
	fmt.Printf("Creation d'un server TCP local sur le port :\n", port)
	//creation du portString avec le bon format pour écouter
	portString := fmt.Sprintf(":%s", strconv.Itoa(port))

	ecoute, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Printf("L'instance ecoute n'a pas pu être crée\n")
		panic(err)
	}
	//si nous sommes ici il n'y a pas d'erreur et panic ne s'est pas executée

	ct := 1

	for {
		fmt.Printf("Acceptation de la prochaine connection\n")
		connection, errc := ecoute.Accept()

		if errc != nil {
			fmt.Printf("Erreur lors de l'acceptation de la prochaine connection")
			panic(errc)
		}
		//si nous sommes ici il n'y a pas d'erreur et panic ne s'est pas executée

		go handleConnection(connection, ct)
		ct += 1
	}
}

func handleConnection(connect net.Conn, ct int) {

	defer connect.Close()
	connectReader := bufio.NewReader(connect)

	for {
		inputLine, err := connectReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error but no panic")
			fmt.Printf("Error :\n", err.Error())
			break
		}

		inputLine = strings.TrimSuffix(inputLine, "\n")
		fmt.Printf("RCV %d %v \n", ct, inputLine)
		splitLine := strings.Split(inputLine, " ")
		returnedString := splitLine[len(splitLine)-1]
		fmt.Printf("SND %d %v \n", ct, returnedString)
		io.WriteString(connect, fmt.Sprintf("%s\n", returnedString))
	}
}
