package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
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
func getArgs() (int, string, string) {
	// Vérifie qu'il y ai bien un argument
	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Println("Erreur : l'usage de Client.go nécessite l'appel suivant : go run Client.go <graph.txt> <portNumber> <ip_adresse> et l'ip adresse est facultative")
		os.Exit(1)
	} else {
		//récupère le nom du fichier et vérifie que le fichier existe bien
		portNumber, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("\"Vous devez utiliser le client ainsi : go run Client.go <graph.txt> <portNumber>\n")
			os.Exit(1)
		} else {
			filename := os.Args[1]
			_, err := os.Stat(filename)
			if os.IsNotExist(err) {
				fmt.Printf("Erreur : le fichier %v n'existe pas, ou il fait référence à un dossier", filename)
				os.Exit(1)
			} else { // J'ai mon port et mon fichier
				if len(os.Args) == 3 { //alors l'ip a été ommise
					ip_adress := "127.0.0.1"
					return portNumber, ip_adress, filename
				} else { //j'ai 4 args => ip donnée
					ip_adress := os.Args[3]
					return portNumber, ip_adress, filename
				}
				// Tout est ok, je retourne le nom du fichier pour la suite du script
			}
		}
		// Ne devrait jamais retourner
	}
	return -1, "", ""
}

func main() {
	start := time.Now()
	s := time.Now()
	port, ip_adress, filename := getArgs()

	//Connection au serveur
	fmt.Printf("Dialing TCP server sur port : %d \n", port)
	portString := fmt.Sprintf("%s:%s", ip_adress, strconv.Itoa(port))
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
		fmt.Printf("Connection au serveur en : %s\n", time.Since(s))
		s = time.Now()

		file, err := os.Open(filename)
		checkError(err)
		defer file.Close()
		// on va parser notre fichier pour ajouter les lignes dans un slice
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// On récupère la ligne du fichier et on l'envoie au serveur
			txt := scanner.Text()
			io.WriteString(connection, txt+"\n") ///Ici on a l'envoie des datas
			//fmt.Printf("Envoie de : %v \n", txt)
		}
		//check si on a une erreur avec le scanner
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}
		fmt.Printf("Fichier parsé et envoyé en in : %s\n", time.Since(s))
		s = time.Now()
		//Après avoir tout envoyé on récupère la réponse du serveur
		//outfile := fmt.Sprintf("out/out_%v.txt", time.Now().Unix()) // passer en GUID -> passer avec le nom d'entrée
		outfile := fmt.Sprintf("out/%v", filepath.Base(filename)) // passer avec le nom d'entrée
		content := ""
		for {
			resultString, err := reader.ReadString('\n') //Là on attend la réponse du serveur

			if err != nil {
				fmt.Printf("Fin de traitement du serveur \n")
				break
			}

			resultString = strings.TrimSuffix(resultString, "\n")
			//fmt.Printf("Réponse du serveur : %v \n ", resultString)
			//TODO stocker dans une var et écrire à la fin de la boucle ??
			content += resultString + "\n"

		}
		f, err := os.OpenFile(outfile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		err = f.Truncate(0)
		if _, err := f.WriteString(content + "\n"); err != nil {
			log.Println(err)
		}
		fmt.Printf("L'analyse de dijkstra est contenu dans : %v \n", outfile)
		fmt.Printf("Écriture, réception et traitement des données in : %s\n", time.Since(s))
		fmt.Printf("Done in : %s\n", time.Since(start))

	}

}
