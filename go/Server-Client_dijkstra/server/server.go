package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	serverAddress = "localhost:8080"
	MAXCLIENTS    = 5
	MAXWORKERS    = 8
)

// Graph représente la structure du graphe pondéré.
type Graph map[string]map[string]int

type Job struct {
	NodeID string
}

type Result struct {
	NodeID    string
	Distances map[string]int
}

type WorkerPool struct {
	workers chan struct{}
	results chan Result
	jobs    chan Job
	wg      sync.WaitGroup
}

func is_there_an_error(err error, errorMessage string) {
	if err != nil {
		fmt.Println(errorMessage, err)
		os.Exit(1)
	}
}

// NewWorkerPool crée un nouveau pool de workers avec la taille spécifiée.
func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{
		workers: make(chan struct{}, size),
		results: make(chan Result),
		jobs:    make(chan Job),
	}
}

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, graph Graph) {
	defer wg.Done()

	for job := range jobs {
		distances := Dijkstra(graph, job.NodeID)
		result := Result{NodeID: job.NodeID, Distances: distances}
		results <- result
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
func GetToWork(wp *WorkerPool, graph Graph) {

	for i := 0; i < MAXWORKERS; i++ {
		wp.wg.Add(1)
		go worker(wp.jobs, wp.results, &wp.wg, graph)
	}

	go func() {
		for node := range graph {
			job := Job{NodeID: node}
			wp.jobs <- job
		}
		close(wp.jobs)
	}()

	go func() {
		wp.wg.Wait()
		close(wp.results)
	}()
}

func GatherAllResults(wp *WorkerPool) Graph {

	allResults := make(map[string]map[string]int)

	for result := range wp.results {
		allResults[result.NodeID] = result.Distances
	}
	return allResults

}

// handleClient gère les connexions des clients.
func handleClient(conn net.Conn, cwp *WorkerPool) {

	defer cwp.wg.Done()
	defer conn.Close()

	company_name := receive_string(conn)
	fmt.Println("Connexion effectuée avec :", company_name)
	var graph Graph = receive_json(conn)
	fmt.Println("Données JSON reçues :", graph)

	wp := NewWorkerPool(MAXWORKERS)

	GetToWork(wp, graph)

	allResults := GatherAllResults(wp)

	send_json(conn, allResults)
	fmt.Println("Données envoyées à", company_name)
}

func receive_json(conn net.Conn) Graph {
	var graph Graph

	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&graph)
	is_there_an_error(err, "Erreur lors de la réception des données JSON :")

	return graph
}

func receive_string(conn net.Conn) string {

	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	is_there_an_error(err, "Erreur lors de la réception de la chaîne de caractères :")

	data = strings.TrimSpace(data)

	return data
}

func send_json(conn net.Conn, data Graph) {

	encoder := json.NewEncoder(conn)
	err := encoder.Encode(data)
	is_there_an_error(err, "Erreur lors de l'envoi des données JSON :")
}

func main() {

	cwp := NewWorkerPool(MAXCLIENTS)

	listener, err := net.Listen("tcp", serverAddress)
	is_there_an_error(err, "Erreur lors de la création du serveur:")

	defer listener.Close()

	fmt.Println("Serveur démarré sur http://localhost:8080")

	for {
		conn, err := listener.Accept()
		is_there_an_error(err, "Erreur lors de l'acceptation de la connexion:")
		cwp.wg.Add(1)
		go handleClient(conn, cwp)
	}
}
