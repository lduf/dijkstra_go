//On attend ici un argument qui est le lien vers notre fichier a lire
package main

import(
		"fmt"
		"os"
	//	"io/ioutil"
		"strconv"
		"bufio"
		"strings"
      )

//Cette fonction permet de vérifier l'état d'une erreur 
// -> Si erreur on panic ^^
func checkError(err error){
	if err != nil{
		panic(err)
	}
}

func getArgs() string{
	// Vérifie qu'il y ai bien un argument
	if len(os.Args) != 2 {
		fmt.Println("Erreur : l'usage de readFile.go nécessite un argument précisant le fichier à traiter")
			os.Exit(1)
	} else{
		//récupère le nom du fichier et vérifie que le fichier existe bien
filename :=  os.Args[1]
		  _, err := os.Stat(filename)
		  if os.IsNotExist(err) {
			  fmt.Printf("Erreur : le fichier %v n'existe pas, ou il fait référence à un dossier", filename)
				  os.Exit(1)
		  } else{
			  // Tout est ok, je retourne le nom du fichier pour la suite du script
			  return filename
		  }
	  // Ne devrait jamais retourner
	}
	return ""
}
type elementGraph struct{
	from string
	to string
	weight int
}
func fileToSlice() []elementGraph{
	filename := getArgs()
	//Ici on a le nom du fichier (qui existe forcément car vérifier avec le getArgs() 
	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()
// on va parser notre fichier pour ajouter les lignes dans un slice
	scanner := bufio.NewScanner(file)
	var slice []elementGraph
	for scanner.Scan() {
// On récupère la ligne du fichier et on la p-split avec l'espace pour le mettre ensuite dans notre slice général (exemple A B 1) est contenu dans splitted[i]
		splitted := strings.Split(scanner.Text(), " ")
		if splitted[2] != "."{
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
	//voilà mon slice 
	return slice
}
//func fileToMap() map[string]map[string]{
//	maped := make(map[string]map[string])
//	return maped
//}

func main() {
	fmt.Printf("%v",fileToSlice())
//	fmt.Printf("%v",fileToMap())

}
