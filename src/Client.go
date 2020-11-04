package main

//TODO : finish file

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//Il faut ici que je me connecte à mon serveur
//que j'extraie les datas de mon fichier pour les envoyer et prépare l'envoie des données  au serveur
// ensuite j'envoie les datas au serveur
// il traite les datas et me renvoie un graph que j'écris dans un out.txt

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
func getArgs() (int, string) {
	// Vérifie qu'il y ai bien un argument
	if len(os.Args) != 3 {
		fmt.Println("Erreur : l'usage de Client.go nécessite l'appel suivant : go run Client.go <portNumber> <graph.txt>")
		os.Exit(1)
	} else {
		//récupère le nom du fichier et vérifie que le fichier existe bien
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("\"Vous devez utiliser le client ainsi : go run Client.go <portNumber>\n")
			os.Exit(1)
		} else {
			filename := os.Args[2]
			_, err := os.Stat(filename)
			if os.IsNotExist(err) {
				fmt.Printf("Erreur : le fichier %v n'existe pas, ou il fait référence à un dossier", filename)
				os.Exit(1)
			} else {
				// Tout est ok, je retourne le nom du fichier pour la suite du script
				return portNumber, filename
			}
		}
		// Ne devrait jamais retourner
	}
	return -1, ""
}

/*
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
}*/

func main() {
	port, filename := getArgs()

	//Connection au serveur
	fmt.Printf("Dialing TCP server sur port : %d \n", port)
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf(portString + "\n")
	connection, err := net.Dial("tcp", portString)
	if err != nil {
		fmt.Printf("Connection echouée \n")
		os.Exit(1)
	} else { //Si ma connection marche
		defer connection.Close()
		reader := bufio.NewReader(connection)
		fmt.Printf("Vous etes bien connectés \n")

		//on va envoyer le contenu de notre fichier

		file, err := os.Open(filename)
		checkError(err)
		defer file.Close()
		// on va parser notre fichier pour ajouter les lignes dans un slice
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// On récupère la ligne du fichier et on l'envoie au serveur
			txt := scanner.Text()
			io.WriteString(connection, txt+"\n") ///Ici on a l'envoie des datas
			fmt.Printf("Envoie de : %v \n", txt)
		}
		//check si on a une erreur avec le scanner
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}
		//Après avoir tout envoyé on récupère la réponse du serveur
		outfile := fmt.Sprintf("out/out_%v.txt", time.Now().Unix())
		for {
			resultString, err := reader.ReadString('\n') //Là on attend la réponse du serveur

			if err != nil {
				fmt.Printf("Le serveur ne renvoie aucune donnée \n")
				break
			}

			resultString = strings.TrimSuffix(resultString, "\n")
			fmt.Printf("Réponse du serveur : %v \n ", resultString)

			f, err := os.OpenFile(outfile,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()
			if _, err := f.WriteString(resultString + "\n"); err != nil {
				log.Println(err)
			}
		}
		fmt.Printf("L'analyse de dijkstra est contenu dans : %v \n", outfile)

	}

}
