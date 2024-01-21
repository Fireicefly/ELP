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
	maxWorker     = 5
)

// Graph représente la structure du graphe pondéré.
type Graph map[string]map[string]int

// Result représente le résultat d'un calcul Dijkstra pour un routeur spécifique.
type Result struct {
	Router    string
	Distances map[string]int
}

// WorkerPool représente un pool de workers.
type WorkerPool struct {
	workers    chan struct{}
	resultChan chan Result
	wg         sync.WaitGroup
	mu         sync.Mutex
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
		workers:    make(chan struct{}, size),
		resultChan: make(chan Result),
	}
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

// handleClient gère les connexions des clients.
func handleClient(conn net.Conn, wp *WorkerPool) {

	defer wp.wg.Done()
	defer conn.Close()

	router_name := receive_string(conn)
	fmt.Println("Connexion effectuée avec :", router_name)
	var graph Graph = receive_json(conn)
	fmt.Println("Données JSON reçues :", graph)

	distances := Dijkstra(graph, router_name)

	encoder := json.NewEncoder(conn)
	err := encoder.Encode(distances)
	is_there_an_error(err, "Erreur lors de l'envoi des résultats au client:")
	fmt.Println("Données envoyées à ", router_name, ":", distances)

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
	fmt.Println(data)
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

	wp := NewWorkerPool(maxWorker)

	listener, err := net.Listen("tcp", serverAddress)
	is_there_an_error(err, "Erreur lors de la création du serveur:")

	defer listener.Close()

	fmt.Println("Serveur démarré sur http://localhost:8080")

	for {
		conn, err := listener.Accept()
		is_there_an_error(err, "Erreur lors de l'acceptation de la connexion:")
		wp.wg.Add(1)
		go handleClient(conn, wp)
	}
}
