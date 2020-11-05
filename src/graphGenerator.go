package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//récupère les arguments fournis au lancement de go
func getArgs() (int, string) {
	// Vérifie qu'il y ai bien un argument
	if len(os.Args) != 3 {
		fmt.Println("Erreur : l'usage de graphGenerator.go nécessite l'appel suivant : go run graphGenerator.go <size> <graph.txt>")
		os.Exit(1)
	} else {
		//récupère le nom du fichier et vérifie que le fichier existe bien
		size, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Vous devez utiliser le générateur ainsi : go run graphGenerator.go <size>\n")
			os.Exit(1)
		} else {
			filename := os.Args[2]
			// Tout est ok, je retourne le nom du fichier pour la suite du script
			return size, filename
		}
		// Ne devrait jamais retourner
	}
	return -1, ""
}

//Génère un poids entre 1 et 16
func randWeight() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(15) + 1
}

//Génère une lettre aléatoire
func randLetter() string {
	rand.Seed(time.Now().UnixNano())
	letters := "AZERTYUIOPQSDFGHJKLMWXCVBN"                    //liste des noeuds possibles
	return fmt.Sprintf("%c", letters[rand.Intn(len(letters))]) //je prends une lettre aléatoire
}

//Permet de générer un graph pour une taille donnée
func generateTie(size int) string {
	rand.Seed(time.Now().UnixNano())
	var from, to, toWrite string //noeud d'arrivé
	combinaison := make(map[string]map[string]int)
	for i := 0; i < size; i++ {
		alea := rand.Float64()
		if i == 0 || alea < 0.3 { //30% du temps on change de point de départ sinon on garde le meme point pour faire un graph un peu plus fournis
			//	println("Nouvelle lettre de départ")
			from = randLetter()
		}

		//je prends une lettre d'arrivé qui doit être différente de celle de départ et la combinaison from-to ne doit pas exister
		for {
			//println("Tirage de to")
			to = randLetter()
			if _, ok := combinaison[from][to]; !ok && to != from {
				combinaison[from] = map[string]int{to: 1} //on bloque dans les deux sens car dans la suite de la fonction on va faire dans le sens to -> from
				combinaison[to] = map[string]int{from: 1}
				//combinaison[from][to]=1 //j'ajoute ma combinaison au tableau (usage d'un map fait bcp moins d'appel en mémoire pour vérifier l'existance (je crois ^^)
				break
			}
		}
		weight := randWeight()

		//Combinaison du graph générée
		toWrite += fmt.Sprintf("%v %v %d\n", from, to, weight)

		//À faire dans l'autre sens (from -> to avec le poids 1 mais to -> from avec un poid 2 possible)
		alea = rand.Float64()
		if alea < 0.75 { //75% du temps on garde la meme poids
			toWrite += fmt.Sprintf("%v %v %d\n", to, from, weight)
		} else {
			toWrite += fmt.Sprintf("%v %v %d\n", to, from, randWeight()) //nouveau poid tiré au sort
		}
	}
	toWrite += ". . ."
	return toWrite
}

func writeGraph(size int, path string) {
	fmt.Printf("Création du fichier %v et génaration d'un graph de taille %d \n", path, size)
	f, err := os.OpenFile(path,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(generateTie(size)); err != nil {
		log.Println(err)
	}

}
func main() {
	writeGraph(getArgs())
}
