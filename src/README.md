# Fichiers sources :

## `readFile.go`
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

## `dijkstra.go`
`dijkstra.go` permet d'initier l'algo de dijkstra. On peut récupérer l'intégralité des voisins des noeuds.
Le dijkstra est fonctionnel, pour 1 to n : l'algo retourne l'ensemble le chemin le plus cours d'un point vers tous les autres.

*Je détaillerai l'algorithme utilisé un peu plus tard*


La fonction ``getAllNeighbors()`` retourne un map de `elementGraph`. Ainsi :

    allNeighbors := getAllNeighbors(graph, noeuds)
    printf("%v",allNeighbors["A"]) // returns [{A B 1} {A C 2}]
    
   
### Test du script :
	go run readFile.go dijkstra.go in/graph.txt
	
### Application de l'algorithme de Dijkstra :

L'algorithme de Dijkstra est appliqué de la manière suivante :

Pour un noeud de départ donné (`from`), je calcule le chemin le plus court vers tous les autres noeuds (`1 to n`). 

Pour ce faire, je fais :

**0. Structure des données**


- Un map de type `dijksTAB : map[string][]chemin` : ainsi, en appelant le tableau `dijksTAB[noeud]`, je récupère les chemins possibles pour me rendre à ce noeud. 

- Un chemin est constitué d'un noeud d'origine `from` et d'un poid `weight`qui est la somme des poids des noeuds par lesquels je suis passé pour me rendre à ce noeud.

- Un tableau (slice) `deadPoints := make(map[string]int)` contient la liste des points pour lesquels il n'est plus nécessaire de revenir lors de la réalisation de l'algorithme (quand on fait manuellement l'algo à la main, on raye la colonne d'un noeud donné lorsqu'il est le plus court chemin, ici je l'ajoute simplement à `deadPoints`

- J'ai aussi un map qui contient l'ensemble des voisins de tous les noeuds `neighbors := getAllNeighbors(graph, noeuds)`. Les données sont récupérées avec la fonction `getAllNeighbors` sus-présentée.


**1. Résolution algorithmique**

0. J'ajoute à mon tableau dijkstra `dijksTAB` le noeud de départ. Son chemin possède par définition un poids nul.
1. Je récupère dans mon tableau dijkstra le noeud pour lequel j'ai le chemin avec le poids le plus faible (usage de la fonction `getMinDijk(dijksTAB, deadPoints)`). NB : Ce point ne peut pas être un point mort ! Je vais appeler ce noeud `p`
2. J'ajoute `p` à la liste des `deadPoints`: on aura plus besoin de calculer des chemins passant par le noeud `p`
3. Je récupère l'ensemble des voisins du point `p`
4. pour chaque voisin (`voisin`) non mort j'ajoute à mon tableau dijkstra `disjksTAB` un chemin menant à `voisin` provenant de `p` et ayant un poids total égal à la somme du poids de `p` et du poids du noeud `p`-> `voisin`
5. Je recommence à l'opération `1.` autant de fois qu'il y a de noeuds différents dans mon graph (`A B X V E` contient 5 noeuds)

*Quand j'ai finis d'itérer, mon `dijksTAB` est complet ! Il suffit de remonter les chemins jusqu'à revenir au point de départ `from`*

6. Pour chaque noeud `p`, tant que je n'ai pas remonté le chemin jusqu'à `from`, je récupère la lettre menant menant au noeud `p` avec le poids minimum :

		

		for getMin(dijksTAB[n]).from != from {
			ways[noeud] = append(ways[noeud], getMin(dijksTAB[n]).from)
			n = getMin(dijksTAB[n]).from
		}
	
**2. Dijsktra n to n**

Il suffit d'appeler l'algorithme présenté en `1.` pour chaque noeud de notre graph

**3. Utilisation de goroutines**

J'utilise des goroutines pour la résolution du `Dijkstra n to n`.  Chaque appel à un `dijkstra 1 to n` est fait via l'appel d'une goroutine.

**4. Analyses et performances**

Afin d'analyser les performances de mon algorithme, j'ai ajouter un timer qui affichera le temps d'éxécution pour de l'algorithme `n to n`.

On remarque que pour un petit graph, l'usage des goroutines ne permet pas un gain de temps important. Pire l'éxécution est prolongée de quelques dixièmes de milisecondes. L'usage des goroutines serait peut être plus pertinant pour un graph plus important

### Limitations algorithmiques :


L'algorithme fonctionne pour un graph correctement détaillé. Imaginons que le poids du trajet `A -> B` soit de `2`, il est nécessaire de préciser dans le graph que le poids du trajet `B -> A` soit également égal à `2`.


**Améliorations possibles :**

Si le trajet est précisé dans un seul sens, on considère que le poids est le même dans l'autre sens.


# Dossiers
Le dossier `in` contient un exemple de graph

Le dossier `out` contiendra les sorties soft avec les données traitées
