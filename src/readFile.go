//On attend ici un argument qui est le lien vers notre fichier a lire
package main

import (
	"fmt"
	"os"
	"sort"
	//	"io/ioutil"
	"bufio"
	"strconv"
	"strings"
)

//Cette fonction permet de vérifier l'état d'une erreur
// -> Si erreur on panic ^^
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getArgs() string {
	// Vérifie qu'il y ai bien un argument
	if len(os.Args) != 2 {
		fmt.Println("Erreur : l'usage de readFile.go nécessite un argument précisant le fichier à traiter")
		os.Exit(1)
	} else {
		//récupère le nom du fichier et vérifie que le fichier existe bien
		filename := os.Args[1]
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			fmt.Printf("Erreur : le fichier %v n'existe pas, ou il fait référence à un dossier", filename)
			os.Exit(1)
		} else {
			// Tout est ok, je retourne le nom du fichier pour la suite du script
			return filename
		}
		// Ne devrait jamais retourner
	}
	return ""
}

type elementGraph struct {
	from   string
	to     string
	weight int
}

//Cette fonction permet de mettre les caractères d'une liste en majuscule
func listToUpper(list []string) {
	for key, elt := range list {
		list[key] = strings.ToUpper(elt)
	}
}

//Cette fonction premet de retirer les valeurs dupliquées dans un slice
func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//On va parser notre fichier pour récupérer les datas de notre graph
// On profite pour créer une liste triée de nos différents noeuds (on prendre minuscule = majuscule (a=A))
func fileToSlice() ([]elementGraph, []string) {
	filename := getArgs()
	//Ici on a le nom du fichier (qui existe forcément car vérifier avec le getArgs()
	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()
	// on va parser notre fichier pour ajouter les lignes dans un slice
	scanner := bufio.NewScanner(file)
	var slice []elementGraph
	var noeuds []string
	for scanner.Scan() {
		// On récupère la ligne du fichier et on la p-split avec l'espace pour le mettre ensuite dans notre slice général (exemple A B 1) est contenu dans splitted[i]
		splitted := strings.Split(scanner.Text(), " ")
		if splitted[2] != "." {
			noeuds = append(noeuds, splitted[0], splitted[1])
			// Je convertis mon poids en entier pcq il était stocké comme un int
			weight, _ := strconv.Atoi(splitted[2])
			// J'ajoute à mon slice un elementGraph
			slice = append(slice, elementGraph{splitted[0], splitted[1], weight})
		}
	}
	//check si on a une erreur avec le scanner
	if err := scanner.Err(); err != nil {
		os.Exit(1)
	}
	//notre liste de noeuds est constituée mais il peut y avoir des doubles, et la liste est non triée
	listToUpper(noeuds)
	noeuds = unique(noeuds)
	sort.Strings(noeuds)

	//voilà mon slice
	return slice, noeuds
}

/*
func main() {
	graph, noeuds := fileToSlice()
	fmt.Printf("%v \n %v \n ",graph, noeuds)

}*/
