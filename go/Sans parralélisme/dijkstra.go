package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

type Graph map[string]map[string]int

func is_there_an_error(err error, message string) {
	if err != nil {
		fmt.Println(message, err)
		os.Exit(1)
	}
}

func Dijkstra(graph Graph, start string) map[string]int {
	distances := make(map[string]int)
	visited := make(map[string]bool)

	for node := range graph {
		distances[node] = math.MaxInt32
	}
	distances[start] = 0

	for i := 0; i < len(graph); i++ {
		u := minDistance(distances, visited)
		visited[u] = true

		for v, weight := range graph[u] {
			if !visited[v] && distances[u] != math.MaxInt32 && distances[u]+weight < distances[v] {
				distances[v] = distances[u] + weight
			}
		}
	}

	return distances
}

func minDistance(distances map[string]int, visited map[string]bool) string {
	minimum := math.MaxInt32
	var minNode string

	for node, dist := range distances {
		if !visited[node] && dist <= minimum {
			minimum = dist
			minNode = node
		}
	}
	return minNode
}

func openJson(file_name string) Graph {

	jsonData, err := os.Open(file_name)
	is_there_an_error(err, "Erreur lors de l'ouverture du fichier JSON :")
	defer jsonData.Close()

	var graph map[string]map[string]int
	decoder := json.NewDecoder(jsonData)
	err = decoder.Decode(&graph)
	is_there_an_error(err, "Erreur lors de la lecture du fichier JSON :")

	return graph
}

func writeJson(allResults map[string]map[string]int) {
	resultsJSON, err := json.Marshal(allResults)
	is_there_an_error(err, "Error converting to JSON:")

	file, err := os.Create("resultat.json")
	is_there_an_error(err, "Error creating file:")
	defer file.Close()

	_, err = file.Write(resultsJSON)
	is_there_an_error(err, "Error writing to file:")
}

func main() {

	start := time.Now()

	graph := openJson("generated_graph.json")

	result := make(map[string]map[string]int)

	for node := range graph {
		distances := Dijkstra(graph, node)
		result[node] = distances
	}

	writeJson(result)

	elapsed := time.Since(start)
	fmt.Println("Temps d'execution :", elapsed)
}
