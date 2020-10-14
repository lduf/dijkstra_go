package main

import (
	"fmt"
)

// On récupère la liste des noeuds pour la parser
//pour chaque noeud -> je récupère la liste des voisins avec lesquels j'ai un lien
func getNeighbors(point string, graph []elementGraph) []elementGraph {
	//on travaille sur le point "point" appartenant la liste "graph"

	//on parcourt notre slice graph et on regarde pour chaque élément si le point de départ est bien le point "point"
	var neighbors []elementGraph
	for _, elt := range graph { //elt est un élément du slice
		if elt.from == point {
			neighbors = append(neighbors, elt) //ajout de l'elt au slice
		}
	}
	return neighbors
}

//Cette fonction permet de récupérer tous les voisins de tous les noeuds
// La fonction retourne un map on peut donc appeler la liste des noeuds visins facilement
func getAllNeighbors(graph []elementGraph, noeuds []string) map[string][]elementGraph {
	allNeighbors := make(map[string][]elementGraph)
	for _, noeud := range noeuds { // parcours la liste des noeuds qui existe
		allNeighbors[noeud] = getNeighbors(noeud, graph) // Ajout de la liste des voisins au map
	}
	return allNeighbors
}

func main() {
	//var allNeighbors map[string][]elementGraph
	fmt.Printf("%v \n", getAllNeighbors(fileToSlice()))
}
