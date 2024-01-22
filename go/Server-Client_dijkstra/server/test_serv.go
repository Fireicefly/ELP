package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math"
	"net"
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

// WorkerPool représente un pool de workers.
type WorkerPool struct {
	workers chan struct{}
	results chan Result
	jobs    chan Job
	wg      sync.WaitGroup
	graph   Graph
}

func is_there_an_error(err error, errorMessage string) {
	if err != nil {
		log.Fatal(errorMessage, err)
	}
}

// NewWorkerPool crée un nouveau pool de workers avec la taille spécifiée.
func NewWorkerPool(size int, graph Graph) *WorkerPool {
	return &WorkerPool{
		workers: make(chan struct{}, size),
		results: make(chan Result),
		jobs:    make(chan Job),
		graph:   graph,
	}
}

// Dijkstra trouve les distances les plus courtes à partir d'un nœud de départ.
func (wp *WorkerPool) Dijkstra(start string) map[string]int {
	distances := make(map[string]int)
	visited := make(map[string]bool)

	for node := range wp.graph {
		distances[node] = math.MaxInt32
	}
	distances[start] = 0

	for i := 0; i < len(wp.graph); i++ {
		u := wp.minDistance(distances, visited)
		visited[u] = true

		for v, weight := range wp.graph[u] {
			if !visited[v] && distances[u] != math.MaxInt32 && distances[u]+weight < distances[v] {
				distances[v] = distances[u] + weight
			}
		}
	}

	return distances
}

func (wp *WorkerPool) minDistance(distances map[string]int, visited map[string]bool) string {
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

// worker est une goroutine qui traite les tâches.
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for job := range wp.jobs {
		distances := wp.Dijkstra(job.NodeID)
		result := Result{NodeID: job.NodeID, Distances: distances}
		wp.results <- result
	}
}

// GetToWork initialise les travailleurs avec des tâches.
func (wp *WorkerPool) GetToWork() {
	for i := 0; i < MAXWORKERS; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	go func() {
		for node := range wp.graph {
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

// GatherAllResults collecte les résultats de tous les travailleurs.
func (wp *WorkerPool) GatherAllResults() Graph {
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

	var clientName string = receiveString(conn)
	log.Printf("Connexion effectuée avec : %s\n", clientName)
	var graph Graph = receiveJSON(conn)
	log.Println("Données JSON reçues.")

	wp := NewWorkerPool(MAXWORKERS, graph)
	wp.GetToWork()

	allResults := wp.GatherAllResults()

	sendJSON(conn, allResults)
	log.Printf("Données envoyées à %s\n", clientName)
}

func receiveJSON(conn net.Conn) Graph {
	var graph Graph

	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&graph)
	is_there_an_error(err, "Erreur lors de la réception des données JSON :")

	return graph
}

func receiveString(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	is_there_an_error(err, "Erreur lors de la réception de la chaîne de caractères :")

	data = strings.TrimSpace(data)

	return data
}

func sendJSON(conn net.Conn, data Graph) {
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(data)
	is_there_an_error(err, "Erreur lors de l'envoi des données JSON :")
}

func main() {
	cwp := NewWorkerPool(MAXCLIENTS, nil)

	listener, err := net.Listen("tcp", serverAddress)
	is_there_an_error(err, "Erreur lors de la création du serveur:")
	defer listener.Close()

	log.Printf("Serveur démarré sur http://%s\n", serverAddress)

	for {
		conn, err := listener.Accept()
		is_there_an_error(err, "Erreur lors de l'acceptation de la connexion:")
		cwp.wg.Add(1)
		go handleClient(conn, cwp)
	}
}
