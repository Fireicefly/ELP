package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"
)

type Graph map[string]map[string]int

type Job struct {
	NodeID string
}

type Result struct {
	NodeID    string
	Distances map[string]int
}

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

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		distances := Dijkstra(graph, job.NodeID)
		result := Result{NodeID: job.NodeID, Distances: distances}
		results <- result
	}
}

var graph Graph // Déclaration de graph comme variable globale

func main() {
	// Lecture du fichier JSON
	byteValue, err := os.ReadFile("graph.json")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON :", err)
		return
	}

	err = json.Unmarshal(byteValue, &graph)
	if err != nil {
		fmt.Println("Erreur lors du décodage du fichier JSON :", err)
		return
	}

	const numWorkers = 8

	var wg sync.WaitGroup

	// Créer les canaux pour les jobs et les résultats
	jobs := make(chan Job, len(graph))
	results := make(chan Result, len(graph))

	// Créer un pool de goroutines (workers)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Envoyer des jobs aux workers
	go func() {
		for node := range graph {
			job := Job{NodeID: node}
			jobs <- job
		}
		close(jobs)
	}()

	// Attendre que tous les workers aient terminé leur travail
	go func() {
		wg.Wait()
		close(results)
	}()

	allResults := make(map[string]map[string]int)
	// Récupérer les résultats des workers
	for result := range results {
		allResults[result.NodeID] = result.Distances
		fmt.Println("Distances les plus courtes depuis le nœud", result.NodeID+" :")
		for destNode, distance := range result.Distances {
			fmt.Printf("De %s à %s: %d\n", result.NodeID, destNode, distance)
		}
		fmt.Println()
	}

	resultJSON, err := json.Marshal(allResults)
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
		fmt.Println("Un fichier resultat.json a été créé contenant les chemins les plus courts pour chaque sommet")
	}
}
