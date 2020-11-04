package main

//TODO : finish file

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func getPort() int {
	//check for arg //os.Args provides access to raw command-line arguments. Note that the first value in this slice is the path to the program, and os.Args[1:] holds the arguments to the program.
	if len(os.Args) != 2 {
		fmt.Printf("Vous devez utiliser le client ainsi : go run Client.go <portNumber>\n")
		os.Exit(1)
	} else {
		//l'arg doit etre int
		fmt.Printf("Vous avez indiqué le port : %v \n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("\"Vous devez utiliser le client ainsi : go run Client.go <portNumber>\n")
			os.Exit(1)
		} else {
			return portNumber
		}
	}
	return -1 //ne devrait jamais être atteint
}

func main() {
	port := getPort()
	fmt.Printf("Dialing TCP server sur port :", port)

	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf(portString)
	connection, err := net.Dial("tcp", portString)
	if err != nil {
		fmt.Printf("Connection echouée \n")
		os.Exit(1)
	} else {
		defer connection.Close()
		reader := bufio.NewReader(connection)
		fmt.Printf("Vous etes bien connectés \n")
		//test d'envois
		for i := 0; i < 10; i++ {

			io.WriteString(connection, fmt.Sprintf("Yoooo %d\n", i))
			resultString, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Lecture impossible")
				os.Exit(1)
			}
			resultString = strings.TrimSuffix(resultString, "\n")
			fmt.Printf("Le serveur à repondu : %v \n", resultString)
			time.Sleep(1000 * time.Millisecond)
		}

	}

}
