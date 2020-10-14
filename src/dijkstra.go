// push test from goland
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
	for _, elt := range graph {
		if elt.from == point {
			neighbors = append(neighbors, elt)
		}
	}
	return neighbors
}

func main() {
	graph, _ := fileToSlice()
	neighbors := getNeighbors("A", graph)
	fmt.Printf("%v \n ", neighbors)

}
