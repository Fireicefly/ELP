# Project All Pair Shortest Path

Le but du projet est d'executer l'algorithme de Dijkstra sur une topologie générée aléatoirement.
Nous avons 3 versions : - Dijkstra sans-parallélisme.
                        - Dijkstra avec-parallélisme.
                        - Dijkstra client/server.

## How to run this
[with an example]

Pour lancer une version, aller dans le dossier souhaité : "cd go", "cd Avec-parallélisme"...
Puis lancer l'algorithme : "go run parallel_dijkstra.go".

## How to generate graph

C'est très simple :
rendez-vous dans le dossier "Graph generator" puis "go run graph_gen.go".
Vous pouvez choisir le nombre de noeuds en modifiant la constante TAILLE en haut du code.
Puis copier/coller le graphe généré dans les autres fichiers pour l'utiliser.

## How to use Client/Server

- Scinder votre terminal en plusieurs parties (ou ouvrer plusieurs terminaux).
- Aller dans le dossier Server-Client_dijkstra puis Server, et lancer le serveur "go run server.go".
- Dans les autres terminaux, aller dans Client, et lancer "go run client.go". (vous pouvez modifier le nom de la compagnie dans la constante company_name en haut du programme.)
