# Projet ELP dijkstra_go
## Objectifs du projet :
Créer en go un soft client serveur (TCP) pour retourner la liste des chemins les plus courts pour un graph donné (*utlisation de l'algorithme de Dijkstra*)

Les étapes du projet :

1. Écrire le graph (trouver une manière de représenter le graph pour le traiter)
2. Extraire les donnnées du graph et mettre dans des variables :
	- Slice ?
	- Map ?
	- Struct ?
3. Écrire l'algorithme de dijkstra pour une chemin 1 noeud vers tous les noeuds (1 -> N)

4. Run la fonction pour tous les chemins possibles
    - décomposition en go routine
5. Récupérer les chemins et leurs poids
6. Implémenter le serveur et le client


# Solution(s) retenue(s) pour chacune des étapes

## 1. Écrire le graph
Solution retenue :

Un lien par ligne (2 noeuds (int) et un poids) et on utilise un code pour signifier le `EOF`

	7 9 2
	9 7 56
	4 1 6
	1 4 3
	. . .
	
*`. . .`  signifie EOF*

## 2. Extraction des données

#### 1. Extraction des données du fichier depuis le client et envoi au serveur (client --> serveur)

Tout bêtement, on parse le fichier donné en argument et on l'envoie à notre serveur.

Ici aucun traitement n'est fait, on gère juste l'envoi en TCP.


#### 2. Réception par le serveur : Utilisation de Slice et de Struct (dijkstra + --> client)

<u>Type de donnée :</u>


On définit un type de donnée composé d'un ``from`` d'un ``to`` et d'un ``weight``, qui sont tous des int.


    type elementGraph struct {
    	from   int
    	to     int
    	weight int
    }

<u>Traitement par le serveur :</u>

Le serveur récupère les données fournies par le client et les ajoute à un slice d'<code>elementGraph</code>

Par la même occasion, le serveur prépare un tableau trié (slice) contenant chaque noeud (de manière unique)


*Exemple de slice*

	[{1 2 1} {2 1 2} {7 8 3} {8 7 2} {2 9 3} {9 2 4} {5 6 2} {6 5 3} {3 2 3} {2 3 5}]



Avantages : 
- Facile à implémenter
- Clair à l'usage

Inconvénients :
- Difficile (à priori) d'accéder aux lien d'un noeud donné


TODO : A voir

## 3. Algorithme de Dijsktra
**Pour un point donné vers les autre point**

Voici à quoi ressemble un chemin le point de départ et le poid total associé au chemin emprunté :

    type chemin struct {
    	from   int
    	weight int
    }
    
 Voici les éléments contenant nos structures de données
 
    ways := make(map[int][]int)    //va contenir tous les chemins du style [1] : [2,5,7,9] , [2] : [1,4,8] , …
    
        	distances := make(map[int]int) //distance totale parcourue pour un point donné : [1] : 6, [2] : 2, …
    
        	dijksTAB := make(map[int][]chemin) // contient en gros tout le travail (équivalent à notre tableau à la main)
    
        	deadPoints := make(map[int]int)    //nom des noeuds par lesquels on ne peut pas repasser
    
        	neighbors := getAllNeighbors(graph, noeuds) //voisins de tous les noeuds


1. Récupérer la liste des lettres (triée de la plus petite à la plus grande)
2. Récupérer les voisins (noeuds avec lesquelles on a un lien)
3. Récupèration du point contenant la distance la plus courte
4. Récupération des chemins possibles depuis ce point





## 4. Récupération des chemins et de leur poids

Code associé à expliquer


    	//À ce niveau, on a récupéré le tableau dijkstra pour à partir de la lettre from.
    	// Il ne reste plus qu'à remonter le tableau pour retourner à from !
    	for _, noeud := range noeuds {
    		ways[noeud] = append(ways[noeud], noeud) //Ajout du noeud d'arrivé
    		n := noeud
    		if len(dijksTAB[n]) > 0 {
    			for getMin(dijksTAB[n]).from != from { // le dijksTAB[n]).from c'est le point dans le 6K (le K) je prends le K et je regarde le chemin le plus court associé au K -> si dans K j'ai 6D je passe à D// from -> point avec lequel j'ai lancé mon dijkstra
    				ways[noeud] = append(ways[noeud], getMin(dijksTAB[n]).from)
    				n = getMin(dijksTAB[n]).from
    			}
    			ways[noeud] = append(ways[noeud], from) //Ajout du noeud de départ
    			distances[noeud] = getMin(dijksTAB[noeud]).weight
    			//Je reverse mon way
    			ways[noeud] = reverse(ways[noeud])
    		}
    	}


## 5. Run n fois la fonction 1 -> n

On appelle notre fonction de dijsktra d'un point vers tous les autres N fois (N étant le nombre de noeud du graphe).

Les appels se font par des go routines.

## 6. Implémentation client/serveur

On utilise une implémentation client server TCP pour bénéficier de la fiabilité du protocole.
En effet, nous ne pouvons pas nous permettre de perdre des données (par exemple avec UDP), comme chacune d'elle est importante.

<hr>

# Performances et compléxité

<hr>

### Quelques liens utiles au projet :

[Qu'est ce que l'algo de Dijkstra](https://www.youtube.com/watch?v=rHylCtXtdNs)

[Exemple de Dijkstra en Go](https://github.com/RyanCarrier/dijkstra)
