package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

type Graph map[string]map[string]int

type Job struct {
	NodeID string
}

type Result struct {
	NodeID    string
	Distances map[string]int
}

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

	for range graph {
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

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, graph Graph) {
	defer wg.Done()

	for job := range jobs {
		distances := Dijkstra(graph, job.NodeID)
		result := Result{NodeID: job.NodeID, Distances: distances}
		results <- result
	}
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
	resultJSON, err := json.Marshal(allResults)
	is_there_an_error(err, "Error converting to JSON:")

	file, err := os.Create("resultat.json")
	is_there_an_error(err, "Error creating file:")
	defer file.Close()

	_, err = file.Write(resultJSON)
	is_there_an_error(err, "Error writing to file:")
}

func main() {

	start := time.Now()

	graph := openJson("generated_graph.json")

	const numWorkers = 6
	var wg sync.WaitGroup
	jobs := make(chan Job, len(graph))
	results := make(chan Result, len(graph))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg, graph)
	}

	go func() {
		for node := range graph {
			job := Job{NodeID: node}
			jobs <- job
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	allResults := make(map[string]map[string]int)

	for result := range results {
		allResults[result.NodeID] = result.Distances
	}

	writeJson(allResults)
	elapsed := time.Since(start)
	fmt.Println("Temps d'execution :", elapsed)
}
