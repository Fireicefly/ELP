package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Type Graph
type Graph map[string]map[string]int

// Taille du graphe
const TAILLE = 500

// Fonction qui vérifie si une erreur est présente
func is_there_an_error(err error, errorMessage string) {
	if err != nil {
		fmt.Println(errorMessage, err)
		os.Exit(1)
	}
}

// Fonction qui génère un graphe aléatoirement
func BuildGraph(nodes []string) Graph {
	adjacencyMap := make(Graph)

	for _, node := range nodes {
		adjacencyMap[node] = make(map[string]int)
	}

	for _, source := range nodes {
		for _, destination := range nodes {
			weight := rand.Intn(20)
			if weight >= 1 && weight <= 8 {
				adjacencyMap[source][destination] = weight
			}
		}
	}

	return adjacencyMap
}

// Fonction qui écrit le graphe généré dans un fichier JSON
func write_json(graph Graph) {

	resultJSON, err := json.Marshal(graph)
	is_there_an_error(err, "Erreur lors de la conversion en JSON :")

	file, err := os.Create("generated_graph.json")
	is_there_an_error(err, "Erreur lors de la création du fichier :")
	defer file.Close()

	_, err = file.Write(resultJSON)
	is_there_an_error(err, "Erreur lors de l'écriture dans le fichier :")
}

// Fonction principale
func main() {

	start := time.Now()
	var nodes []string

	for i := 1; i <= TAILLE; i++ {
		node := "R" + strconv.Itoa(i)
		nodes = append(nodes, node)
	}

	adjacencyMap := BuildGraph(nodes)

	write_json(adjacencyMap)
	elapsed := time.Since(start)
	fmt.Println("Temps d'execution :", elapsed)
}
