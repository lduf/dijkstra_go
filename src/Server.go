package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// DEBUT DE DIJKSTRA

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

// Cette fonction permet de récupérer la valeur mini contenue dans notre tableau dijkstra
func getMinDijk(dijksTAB map[string][]chemin, deadPoints map[string]int) (string, int) {
	min := -1 // Attention on considère ici que les poids ne peuvent que être positifs
	minPoint := ""
	minKey := 0
	for point, i := range dijksTAB { //(i pour iterration) -> c'est un slice de chemin
		////fmt.Printf("##### DEBUG min:   Je suis au point %v \n", point)
		if _, ok := deadPoints[point]; !ok {
			////fmt.Printf("##### DEBUG min: %v n'est pas un point mort \n", point)
			for k, chm := range i { // k est l'indice du chemin chm
				if min != -1 && chm.weight < min { // si min est défini et que le poids du chemin est inférieur au min
					min = chm.weight
					minPoint = point
					minKey = k
				}
				if min == -1 { //si min n'a pas été défini
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

// permet de retourner le Dijkstra considerant un point de départ (from)
func getDijkstra(from string, wg *sync.WaitGroup, graph []elementGraph, noeuds []string) (map[string][]string, map[string]int) {
	defer wg.Done() // on vire notre waitgroup
	//initialisation des variables
	ways := make(map[string][]string) //va contenir tous les chemins
	distances := make(map[string]int) //distance totale parcourue
	// À la main on utilise un tableau à 2 entrées : A | B | C | ... | Z (nom des noeuds) et le nombre de "tour". On itère à chaque tour pour trouver la distance la plus courte.
	dijksTAB := make(map[string][]chemin) // contient en gros tout le travail
	deadPoints := make(map[string]int)    //nom des noeuds par lesquels on ne peut pas repasser

	//Étape 1 : on créé notre tableau dans lequel on appliquera l'algo

	neighbors := getAllNeighbors(graph, noeuds) //voisins de tous les noeuds

	dijksTAB[from] = append(dijksTAB[from], chemin{from, 0}) //initialisation du tableau
	//deadPoints[from] = 0

	//on récupère la distance la plus courte
	for i := 0; i < len(noeuds); i += 1 { //On parcourt autant de fois qu'il y a de noeud
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
		for getMin(dijksTAB[n]).from != from { // le dijksTAB[n]).from c'est le point dans le 6K (le K) je prends le K et je regarde le chemin le plus court associé au K -> si dans K j'ai 6D je passe à D// from -> point avec lequel j'ai lancé mon dijkstra
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
func Dijkstra(graph []elementGraph, noeuds []string) (map[string]map[string][]string, map[string]map[string]int) {
	var wg sync.WaitGroup // Waitgroup
	//_, noeuds := fileToSlice()
	dijk := make(map[string]map[string][]string)
	distances := make(map[string]map[string]int)
	for _, noeud := range noeuds {
		var ways map[string][]string
		var dists map[string]int

		wg.Add(1)
		go func() {
			ways, dists = getDijkstra(noeud, &wg, graph, noeuds)
		}()
		wg.Wait()
		dijk[noeud] = ways
		distances[noeud] = dists
	}

	return dijk, distances
}

//fin de Dijkstra

// permet de récupérer le port sur lequel le serveur est créé
func getPortS() int {
	//check for arg //os.Args provides access to raw command-line arguments. Note that the first value in this slice is the path to the program, and os.Args[1:] holds the arguments to the program.
	if len(os.Args) != 2 {
		fmt.Printf("Vous devez utiliser le server ainsi : go run Server.go <portNumber>\n")
		os.Exit(1)
	} else {
		//l'arg doit etre int
		fmt.Printf("Vous avez indiqué le port :\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Vous devez utiliser le server ainsi : go run Server.go <portNumber>\n")
			os.Exit(1)
		} else {
			return portNumber
		}
	}
	return -1 //ne devrait jamais être atteint
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

func main() {
	port := getPortS() //récup le port
	fmt.Printf("Creation d'un server TCP local sur le port : %d \n", port)
	//creation du portString avec le bon format pour écouter
	portString := fmt.Sprintf(":%s", strconv.Itoa(port))

	ecoute, err := net.Listen("tcp", portString) //création de l'écoute du serveur
	if err != nil {
		fmt.Printf("L'instance ecoute n'a pas pu être crée\n")
		panic(err)
	}
	//si nous sommes ici il n'y a pas d'erreur et panic ne s'est pas executée

	ct := 1

	for { //Tout le temps on attend les connections
		fmt.Printf("Acceptation de la prochaine connection\n")
		connection, errc := ecoute.Accept() //on accepte la connecttion

		if errc != nil {
			fmt.Printf("Erreur lors de l'acceptation de la prochaine connection")
			panic(errc)
		}
		//si nous sommes ici il n'y a pas d'erreur et panic ne s'est pas executée

		go handleConnection(connection, ct) //si la connection avec le client est OK
		ct += 1
	}
}

func handleConnection(connect net.Conn, ct int) {

	defer connect.Close()
	connectReader := bufio.NewReader(connect)

	var slice []elementGraph
	var noeuds []string

	for {
		inputLine, err := connectReader.ReadString('\n') //on récupère la ligne envoyée par le client
		if err != nil {
			fmt.Printf("Error but no panic")
			fmt.Printf("Error :\n", err.Error())
			break
		}
		inputLine = strings.TrimSuffix(inputLine, "\n")
		//©	fmt.Printf("%v \n", inputLine)
		splitted := strings.Split(inputLine, " ")
		if splitted[2] != "." {
			noeuds = append(noeuds, strings.ToUpper(splitted[0]), strings.ToUpper(splitted[1]))
			// Je convertis mon poids en entier pcq il était stocké comme un string
			weight, _ := strconv.Atoi(splitted[2])
			// J'ajoute à mon slice un elementGraph
			slice = append(slice, elementGraph{splitted[0], splitted[1], weight})
		} else {
			break
		}
	}
	listToUpper(noeuds)
	noeuds = unique(noeuds)
	sort.Strings(noeuds)

	ways, distances := Dijkstra(slice, noeuds)

	for letter, graph := range ways {
		for l, way := range graph {
			out := fmt.Sprintf("%v %v %v %v \n", letter, l, way, distances[letter][l])
			//fmt.Printf("Envoie de : %v", out)
			io.WriteString(connect, fmt.Sprintf("%s", out))
		}
	}
}
