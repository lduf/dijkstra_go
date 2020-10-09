# Projet ELP dijkstra_go
## Objectifs du projet :
Créer en go un soft client serveur (TCP) pour retourner la liste des chemins les plus courts pour un graph donné (*utlisation de l'algorithme de Dijkstra*)

Les étapes du projet :

1. Écrire le graph (trouver une manière de représenter le graph pour le traiter)
2. Extraire les donnnées du graph et mettre dans des variables :
	- Slice ?
	- Map ?
	- Struct ?
	- (plus selon l'apprentissage de go)
3. Écrire l'algorithme de dijkstra pour une chemin entre deux noeuds donnés
	- décomposition en go routine
4. Run la fonction pour tous les chemins possibles 
5. Récupérer les chemins et leur poids
6. Implémenter le serveur et le client


# Solution(s) retenue(s) pour chacune des étapes

## 1. Écrire le graph
Solution envisageable :

	A B 2
	B C 56
	C A 6
	B A 3
	. . .
	
*. . . pourrait signifier EOF*
## 2. Extraction des données

1. Utilisation de Slice et de Struct

La version 1 de `readFile.go` permet d'obtenir un slice composé par type créé pour l'occasion.

*Exemple de retour*

	[{A B 1} {A C 2} {B F 3} {B D 2} {C D 3} {C E 4} {D E 2} {D F 3} {D G 3} {E G 5} {F G 4}]

Avantages : 
- Facile à implémenter
- Un type est défini on peut lui ajouter des méthodes (à ce stade de dev aucune idée de si c'est utile ou pas :-) )

Inconvénients :
- Difficile (à priori) d'accéder aux lien d'un noeud donné

2. Utilisation de Map

Pas encore implémenter dans un code (surement fait dans la V2 de `readFile.go`)

*Exemple de retour (datas non présentées en respectant Go, format plus proche de json pour l'exemple)*

	"A" : 	[
			{
			"noeud" : "B",
			"poids" : "1" 
			},
			{
				"noeud" : "C",
				"poids" : "2" 
			}
		],
	"B" : 	[
			{
			"noeud" : "F",
			"poids" : "3" 
			},
			{
				"noeud" : "D",
				"poids" : "2" 
			}
		]

Avantages :

- On accède rapidement à tous les liens pour un noeud donné
- Plus clair 

Inconvénients : 

- Semble plus difficile à implémenter
- Peut être difficile à utiliser (tout dépend de comment Go est foutu, on trouvera peut être des fonctions super cool pour bosser là dessus *(en php il existe des fonctions super puissances pour bosser sur les tableaux multi, alors qui sait pour Go …)*)
		

## 3. Algorithme de Dijsktra
**Pour un point donné vers un autre point**
1. Récupérer les voisins (noeuds avec lesquelles on a un lien)
2. Ajouter la distance parcourru 
3. Sélectionner la plus petite valeur
## 4. Run la fonction pour n chemins
## 5. Récupération des chemins et de leur poids
## 6. Implémentation client/serveur


### Ajouter ici des liens utiles au projet :

[Qu'est ce que l'algo de Dijkstra](https://www.youtube.com/watch?v=rHylCtXtdNs)

[Exemple de Dijkstra en Go](https://github.com/RyanCarrier/dijkstra)
