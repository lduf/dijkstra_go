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
func randLetter(alphabet []int) int {
	rand.Seed(time.Now().UnixNano())
	//liste des noeuds possibles
	return alphabet[rand.Intn(len(alphabet))] //je prends une lettre aléatoire
}

func remove(slice []int, s int) []int {
	end := append(slice[:s], slice[s+1:]...)
	return end
}
func remove_element(slice []int, elt int) []int {
	i := 0
	for slice[i] != elt {
		i++
	}

	return remove(slice, i)
}

//Permet de générer un graph pour une taille donnée
func generateTie(size int) string {
	var alphabet []int
	for i := 0; i < size; i++ {
		alphabet = append(alphabet, i)
	}
	//fmt.Printf("Alphabet %d \n", alphabet)
	neighb := make(map[int][]int) //contient les lettres possibles de matcher
	/*
		[1] = [2],[3],[...],[n]
		[2] = [1],[3],[...],[n]
		[n] = [1], ... [n-1]

	*/
	//boucle d'initialisation de neighb (voisins disponibles)
	//fmt.Printf("Debug génération neighb\n")
	for _, letter := range alphabet {
		neighb[letter] = make([]int, len(alphabet))
		copy(neighb[letter], alphabet)
		//fmt.Printf("neighb[%d] :  %d \n",letter, neighb[letter])
		remove(neighb[letter], letter)
		neighb[letter] = neighb[letter][:len(neighb[letter])-1]
	}
	//fmt.Printf("neighb %d \n", neighb)
	disp_from := alphabet
	//rand.Seed(time.Now().UnixNano())
	var from, to int
	var toWrite string //noeud de départ -> noeud d'arrivé -> résultat de la fonction
	run := true
	draw := true

	for i := 0; i < size && run; i++ {
		alea := rand.Float64()
		if draw || alea < 0.3 { //30% du temps on change de point de départ sinon on garde le meme point pour faire un graph un peu plus fournis
			//println("Nouvelle lettre de départ")
			from = randLetter(disp_from)
			draw = false
		}
		//println(from)
		//println("Tirage de to")
		to = randLetter(neighb[from]) // Je prends une lettre qui est encore disponible à partir de from, on a aussi forcément que son reverse est dispo car on les gère en meme temps
		// Si le nombre de points encore accéssible est plus grand que 1 pas de soucis, je supprime la lettre que je viens de prendre
		if len(neighb[from]) > 1 {
			remove_element(neighb[from], to) // retrait de la liste que je viens de prendre
			//	fmt.Printf("TO n'est pas le dernier de %v, retrait de la liste from : %v \n", from,neighb[from])
			if len(neighb[to]) > 1 {
				remove_element(neighb[to], from) // retrait de son reverse
			} else {
				remove_element(disp_from, to) // Je retire son reverse
			}
		} else { // Si c'était la denière lettre, le from n'est plus disponible car il ne mene nulle part, je le supprime
			if len(disp_from) > 1 { // Si ce n'est pas le dernier from
				remove_element(disp_from, from) // Je le retire
				if len(neighb[to]) > 1 {
					remove_element(neighb[to], from) // retrait de son reverse
				} else {
					remove_element(disp_from, to) // Je retire son reverse
				}
				draw = true //je précise que je dois tirer un nouveau from
			} else { // Si c'est le dernier from beh c'est la merde on kick
				println("KILL !")
				run = false
			}
		}
		weight := randWeight()
		//Combinaison du graph générée
		toWrite += fmt.Sprintf("%d %d %d\n", from, to, weight)
		//À faire dans l'autre sens (from -> to avec le poids 1 mais to -> from avec un poid 2 possible)
		alea = rand.Float64()
		if alea < 0.75 { //75% du temps on garde la meme poids
			toWrite += fmt.Sprintf("%d %d %d\n", to, from, weight)
		} else {
			toWrite += fmt.Sprintf("%d %d %d\n", to, from, randWeight()) //nouveau poid tiré au sort
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
