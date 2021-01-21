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

/*
Ce fichier a pour but de communiquer avec le server et de récupérer la reponse ainsi:
1. Connection au serveur
2. Extraction des datas de mon fichier graph
3. Preparation et envoi des données au serveur
4. Récupération des datas envoyées en retour par le serveur et écriture dans un fichier texte de sortie

	- DEBUG commentaires de debug
*/

//fonction pour vérifier la présence d'erreurs
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

//fonction de traitement et de vérification des arguments passés en entrée
func getArgs() (int, string, string) {
	// Vérifie si le nombre d'arguments n'est pas dans l'intervalle requis
	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Println("Erreur : l'usage de Client.go nécessite l'appel suivant : go run Client.go <graph.txt> <portNumber> <ip_adresse> et l'ip adresse est facultative")
		os.Exit(1)
	} else {
		//récupère le port et vérifie si une erreur est intervenue à la conversion
		portNumber, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("\"/ERROR/ Vous devez utiliser le client ainsi : go run Client.go <graph.txt> <portNumber> <ip_adresse> et l'ip adresse est facultative\"\n")
			os.Exit(1)
		} else { //Si pas d'erreurs
			//récupère le nom du fichier et vérifie que le fichier existe bien
			filename := os.Args[1]
			_, err := os.Stat(filename)
			if os.IsNotExist(err) {
				fmt.Printf("Erreur : le fichier %v n'existe pas, ou il fait référence à un dossier", filename)
				os.Exit(1) // pourquoi pas panic ?
			} else { // J'ai mon port et mon fichier
				if len(os.Args) == 3 { //alors l'ip a été ommise
					ip_adress := "127.0.0.1" //adresse locale
					return portNumber, ip_adress, filename
				} else { //j'ai 4 args => ip donnée
					ip_adress := os.Args[3]
					return portNumber, ip_adress, filename
				}
			}
		}
		// Ne devrait jamais retourner
	}
	return -1, "", ""
}

func main() {
	//On démarre deux timer
	start := time.Now()
	s := time.Now()
	port, ip_adress, filename := getArgs() //On récupère les args

	//Connection au serveur
	fmt.Printf("Dialing TCP server sur port : %d \n", port)
	portString := fmt.Sprintf("%s:%s", ip_adress, strconv.Itoa(port)) //formatage selon x.x.x.x:xxxx ex: 127.0.0.1:4000
	fmt.Printf(portString + "\n")                                     //retour à la ligne ?
	connection, err := net.Dial("tcp", portString)                    //on établie TCP
	if err != nil {                                                   //si erreur exit
		fmt.Printf("Connection echouée \n")
		os.Exit(1)
	} else { //Si ma connection marche
		defer connection.Close()              //on defer la fermeture pour etre sur de faire toutes les actions avant
		reader := bufio.NewReader(connection) //On met un reader sur la connection pour ecouter le serveur en retour
		fmt.Printf("Vous etes bien connectés \n")

		//on va envoyer le contenu de notre fichier
		fmt.Printf("Connection au serveur en : %s\n", time.Since(s)) //temps de connection
		s = time.Now()                                               //on relance un timer pour plus tard (pour calculer le parsing)

		file, err := os.Open(filename) //on ouvre le fichier graph donnée en arg
		checkError(err)
		defer file.Close() //defer close
		// on va parser notre fichier pour ajouter les lignes dans un slice
		scanner := bufio.NewScanner(file)
		for scanner.Scan() { //Ici on est déjà en ligne par ligne
			// On récupère la ligne du fichier et on l'envoie au serveur
			txt := scanner.Text()
			io.WriteString(connection, txt+"\n") ///Ici on envoie la ligne plus un retour à la ligne
			//fmt.Printf("Envoie de : %v \n", txt) DEBUG
		}
		//check si on a une erreur avec le scanner
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}
		fmt.Printf("Fichier parsé et envoyé en in : %s\n", time.Since(s))
		s = time.Now() //encore un timer pour la réponse
		//Après avoir tout envoyé on récupère la réponse du serveur
		outfile := fmt.Sprintf("out/%v", filepath.Base(filename)) // ou nomme le fichier de sortir en fonction de celui d'entrée
		content := ""
		for {
			resultString, err := reader.ReadString('#') //Là on attend la réponse du serveur (par le reader instancié plus tôt)

			if err != nil { //dès qu'on a une erreur on arrete de recevoir
				fmt.Printf("Fin de traitement du serveur \n")
				break //on sort de la boucle
			}

			//fmt.Printf("Réponse du serveur : %v \n ", resultString) DEBUG
			content += strings.TrimSuffix(resultString, "#") //on incremente content avec les résultats récupérés à chaque passage dans le for

		}
		fmt.Printf("Réception et traitement des données in : %s\n", time.Since(s))
		s = time.Now() //encore un timer pour la réponse
		f, err := os.OpenFile(outfile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //meme open file que dans graph generator
		if err != nil { //on affiche l'erreur si il y en a
			log.Println(err)
			panic(err) // si crash voir ici
		}
		defer f.Close()                                          //defer close pour l'ouverture du fichier de sortie
		err = f.Truncate(0)                                      // je tronque mais je lève une erreur (je tronque pour que le fichier soir bien vide)
		if _, err := f.WriteString(content + "\n"); err != nil { //si il y a une erreur durant l'écriture l'afficher, sinon l'écrire
			log.Println(err)
		}
		//quelques print pour synthétiser le déroulement du processus
		fmt.Printf("L'analyse de dijkstra est contenu dans : %v \n", outfile)
		fmt.Printf("Écriture des données in : %s\n", time.Since(s))
		fmt.Printf("Done in : %s\n", time.Since(start))

	}

}
