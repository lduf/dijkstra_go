package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// On récupère la liste des noeuds pour la parser
//pour chaque noeud -> je récupère la liste des voisins avec lesquels j'ai un lien pour un noeud donné
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

//permet d'obtenir le poid le plus petit pour un graph donné
func getMin(graphPart []chemin) chemin {
	minKey := 0
	minValue := graphPart[0].weight
	for k, elt := range graphPart {
		if elt.weight < minValue {
			minValue = elt.weight
			minKey = k
		}
	}
	return graphPart[minKey]
}

func getMinDijk(dijksTAB map[string][]chemin, deadPoints map[string]int) (string, int) {
	min := -1
	minPoint := ""
	minKey := 0
	for point, i := range dijksTAB { //(i pour iterration) -> c'est un slice de chemin
		////fmt.Printf("##### DEBUG min:   Je suis au point %v \n", point)
		if _, ok := deadPoints[point]; !ok {
			////fmt.Printf("##### DEBUG min: %v n'est pas un point mort \n", point)
			for k, chm := range i {
				if min != -1 && chm.weight < min {

					min = chm.weight
					minPoint = point
					minKey = k
				}
				if min == -1 {
					min = chm.weight
					minPoint = point
					minKey = k
				}
				////fmt.Printf("##### DEBUG min: Pour l'instant le chemin le plus court est : %v depuis le point %v \n", dijksTAB[minPoint][minKey], minPoint)
			}
		}
	}
	return minPoint, minKey
}

type chemin struct {
	from   string
	weight int
}

func reverse(s []string) []string { // Permet de reverse un slice de string
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

/*
func getDijkstra1to1(from string, to string) ([]string, int) {
	ways, distances := getDijkstra(from)
	return ways[to], distances[to]
}*/

func getDijkstra(from string, wg *sync.WaitGroup) (map[string][]string, map[string]int) {
	defer wg.Done()
	//initialisation des variables
	ways := make(map[string][]string) //va contenir tous les chemins

	//var way []string //contient le nom des noeuds par lesquels passer
	distances := make(map[string]int) //distance totale parcourue
	// À la main on utilise un tableau à 2 entrées : A | B | C | ... | Z (nom des noeuds). On itère à chaque tour pour trouver la distance la plus courte.
	dijksTAB := make(map[string][]chemin)
	deadPoints := make(map[string]int) //nom des noeuds par lesquels on ne peut pas repasser

	//Étape 1 : on créé notre tableau dans lequel on appliquera l'algo
	graph, noeuds := fileToSlice()              //On récupère ici le graph sur lequel on va travailler et la liste des noeuds présents dans le graph
	neighbors := getAllNeighbors(graph, noeuds) //voisins de tous les noeuds

	dijksTAB[from] = append(dijksTAB[from], chemin{from, 0})
	//deadPoints[from] = 0

	//on récupère la distance la plus courte
	for i := 0; i < len(noeuds); i += 1 {
		//fmt.Printf("\n \n \n### DEBUG : Nouveau tour \n")
		pt, k := getMinDijk(dijksTAB, deadPoints) //je récupère le point et la clé contenant la distance la plus courte
		smallestWay := dijksTAB[pt][k]            // Ici j'ai le chemin le plus court
		//fmt.Printf("### DEBUG : J'ai récupéré le chemin le plus court : lettre %v poid : %d \n", pt, smallestWay.weight)
		deadPoints[pt] = i // J'ajoute au point mort le point contenu dans le chemin le plus court
		//fmt.Printf("### DEBUG : Ajout de %v aux points morts avecc l'indice : %d \n", pt, i)
		//fmt.Printf("### DEBUG : Parcourt des voisins de lettre %v  \n", pt)
		for _, direction := range neighbors[pt] { //pour tous les voisins du point du chemin le plus court
			//fmt.Printf("### DEBUG :  %v à pour voisin %v \n", pt, direction.to)
			if _, ok := deadPoints[direction.to]; !ok { // si la direction du point vers lequel je vais n'est pas dans la liste des points morts alors
				//fmt.Printf("### DEBUG :  %v n'est pas un point mort \n", direction.to)
				dijksTAB[direction.to] = append(dijksTAB[direction.to], chemin{direction.from, direction.weight + smallestWay.weight}) //De mon point, il existe un nouveau chemin vers le point direction.to de poid total := plus petit
				//fmt.Printf("### DEBUG :  Je peux aller en %v depuis %v. Cette route à un poid de %d menant à un poid total de %d \n", direction.to, direction.from, direction.weight, direction.weight+smallestWay.weight)
			}
		}
		//fmt.Printf("### DEBUG : voici à quoi ressemble mon Dijsktra %v \n \n \n", dijksTAB)
	}
	//À ce niveau, on a récupéré le tableau dijkstra pour à partir de la lettre from.
	// Il ne reste plus qu'à remonter le tableau pour retourner à from !
	for _, noeud := range noeuds {
		ways[noeud] = append(ways[noeud], noeud) //Ajout du noeud d'arrivé
		n := noeud
		for getMin(dijksTAB[n]).from != from {
			ways[noeud] = append(ways[noeud], getMin(dijksTAB[n]).from)
			n = getMin(dijksTAB[n]).from
		}
		ways[noeud] = append(ways[noeud], from) //Ajout du noeud de départ
		distances[noeud] = getMin(dijksTAB[noeud]).weight
		//Je reverse mon way
		ways[noeud] = reverse(ways[noeud])
	}

	return ways, distances
}
func Dijkstra() (map[string]map[string][]string, map[string]map[string]int) {
	var wg sync.WaitGroup // Waitgroup
	_, noeuds := fileToSlice()
	dijk := make(map[string]map[string][]string)
	distances := make(map[string]map[string]int)
	for _, noeud := range noeuds {
		var ways map[string][]string
		var dists map[string]int

		wg.Add(1)
		go func() {
			ways, dists = getDijkstra(noeud, &wg)
		}()
		wg.Wait()
		dijk[noeud] = ways
		distances[noeud] = dists
	}

	return dijk, distances
}
func main() {
	start := time.Now()
	//var allNeighbors map[string][]elementGraph

	ways, distances := Dijkstra()

	elapsed := time.Since(start)

	for letter, graph := range ways {
		fmt.Printf("\n ##### FROM %v #### \n", letter)
		for l, way := range graph {
			fmt.Printf("Chemin le plus court de '%v' à '%v' : %v \n Distance entre les points : %v \n", letter, l, way, distances[letter][l])
		}
		fmt.Printf("\n\n\n")
	}
	log.Printf("It took %s", elapsed)

}
