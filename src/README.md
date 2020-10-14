# Fichiers sources :

### `readFile.go`
`readFile.go`permet de lire des fichiers de `in` et d'en extraire les données jusqu'à ce qu'on rencontre notre EOF (`. . .`)


Les données sont ensuite stockées dans un `var slice []elementGraph`, où `elementGraph` représente :

	type elementGraph struct{
		from string
		to string
		weight int
	}

*NB :  on comprend assez logiquement que `from` désigne le noeud de départ (ex: `A`), `to` le noeud d'arrivé (ex : `B`), et `weight` est le poids du lien (ex: `3`)*

Dans le fichier on retrouvera la fonction `fileToSlice()` qui permet d'analyser le fichier et de retourner le slice susprésenté mais aussi un slice contenant le nom de tous les noeuds (trié par ordre croissant (`"A" < "B" < "C" < ... < "Z"`) *//logique en soit//*


**Attention :** `readFile.go`attend comme argument le chemin du fichier 

### `dijkstra.go`
`dijkstra.go` permet d'initier l'algo de dijkstra. On peut récupérer l'intégralité des voisins des noeuds.

La fonction ``getAllNeighbors()`` retourne un map de `elementGraph`. Ainsi :

    allNeighbors := getAllNeighbors(graph, noeuds)
    printf("%v",allNeighbors["A"]) // returns [{A B 1} {A C 2}]
    
   
#### Test du script :
	go run readFile.go dijkstra.go in/graph.txt
**Il faut maintenant appliquer l'algo de dijkstra**

## Dossiers
Le dossier `in` contient un exemple de graph

Le dossier `out` contiendra les sorties soft avec les données traitées
