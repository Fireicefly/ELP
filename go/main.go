package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

flknalkfnlakflalkfaafaatu le saiss
import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Graph map[string]map[string]int

func Dijkstra(graph Graph, start string) map[string]int {
	distances := make(map[string]int)
	visited := make(map[string]bool)

	// Initialisation des distances avec une valeur infinie et du point de départ à 0
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

func main() {
	// Lecture du fichier JSON
	byteValue, err := os.ReadFile("graph.json")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON :", err)
		return
	}

	var graph Graph
	err = json.Unmarshal(byteValue, &graph)
	if err != nil {
		fmt.Println("Erreur lors du décodage du fichier JSON :", err)
		return
	}

	result := make(map[string]map[string]int)

	for node := range graph {
		distances := Dijkstra(graph, node)
		result[node] = distances
		fmt.Println("Distances les plus courtes depuis le noeud", node+" :")
		for destNode, distance := range distances {
			fmt.Printf("De %s à %s: %d\n", node, destNode, distance)
		}
		fmt.Println()
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		return
	}

	file, err := os.Create("resultat.json")
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier :", err)
		return
	}
	defer file.Close()

	_, err = file.Write(resultJSON)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture dans le fichier :", err)
		return
	} else {
		fmt.Println("Un fichier resultat.json a été crée contenant les chemins les plus courts pour chaque sommet")
	}
}

salut c moi tu le sais ,nsdjkfdnjcksncjksd