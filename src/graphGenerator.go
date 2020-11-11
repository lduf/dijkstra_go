package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
			if size > (25*26)/2 {
				println("La variable size est supérieur au nombre de combinaison possibles. Nous limiterons la taille au nombre maximal de combinaison ")
				size = (25 * 26) / 2
			}
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
func randLetter(letters string) string {
	rand.Seed(time.Now().UnixNano())
	//liste des noeuds possibles
	return fmt.Sprintf("%c", letters[rand.Intn(len(letters))]) //je prends une lettre aléatoire
}

//Permet de générer un graph pour une taille donnée
func generateTie(size int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	neighb := make(map[string]string)
	for _, letter := range letters {
		neighb[string(letter)] = strings.Replace(letters, string(letter), "", -1)
	}
	disp_from := letters
	rand.Seed(time.Now().UnixNano())
	var from, to, toWrite string //noeud de départ -> noeud d'arrivé -> résultat de la fonction
	run := true
	draw := true

	for i := 0; i < size && run; i++ {
		alea := rand.Float64()
		if draw || alea < 0.3 { //30% du temps on change de point de départ sinon on garde le meme point pour faire un graph un peu plus fournis
			//println("Nouvelle lettre de départ")
			from = randLetter(disp_from)
			draw = false
		}
		//println("Tirage de to")
		to = randLetter(neighb[from]) // Je prends une lettre qui est encore disponible à partir de from, on a aussi forcément que son reverse est dispo car on les gère en meme temps
		// Si le nombre de points encore accéssible est plus grand que 1 pas de soucis, je supprime la lettre que je viens de prendre
		if len(neighb[from]) > 1 {
			neighb[from] = strings.Replace(neighb[from], to, "", -1) // retrait de la liste que je viens de prendre
			//	fmt.Printf("TO n'est pas le dernier de %v, retrait de la liste from : %v \n", from,neighb[from])
			if len(neighb[to]) > 1 {
				neighb[to] = strings.Replace(neighb[to], from, "", -1) // retrait de son reverse
			} else {
				disp_from = strings.Replace(disp_from, to, "", -1) // Je retire son reverse
			}
		} else { // Si c'était la denière lettre, le from n'est plus disponible car il ne mene nulle part, je le supprime
			if len(disp_from) > 1 { // Si ce n'est pas le dernier from
				disp_from = strings.Replace(disp_from, from, "", -1) // Je le retire
				if len(neighb[to]) > 1 {
					neighb[to] = strings.Replace(neighb[to], from, "", -1) // retrait de son reverse
				} else {
					disp_from = strings.Replace(disp_from, to, "", -1) // Je retire son reverse
				}
				draw = true //je précise que je dois tirer un nouveau from
			} else { // Si c'est le dernier from beh c'est la merde on kick
				println("KILL !")
				run = false
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
	s := time.Now()
	writeGraph(getArgs())
	fmt.Printf("Éxécution en  : %s\n", time.Since(s))
}
