//On attend ici un argument qui est le lien vers notre fichier a lire
package main

import(
	"fmt"
	"os"
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

func main() {
	filename := getArgs()
	fmt.Printf("%v", filename)
}
