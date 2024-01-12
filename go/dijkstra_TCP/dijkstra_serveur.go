package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
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

// NewWorkerPool crée un nouveau pool de workers avec la taille spécifiée.
func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{
		workers:    make(chan struct{}, size),
		resultChan: make(chan Result),
	}
}

// CalculateShortestPath ajoute une tâche au pool pour calculer le plus court chemin.
func (wp *WorkerPool) CalculateShortestPath(router string, graph Graph) {
	wp.wg.Add(1)
	wp.workers <- struct{}{}

	go func(router string, graph Graph) {
		defer wp.wg.Done()

		distances := dijkstra(graph, router)

		wp.resultChan <- Result{
			Router:    router,
			Distances: distances,
		}

		<-wp.workers
	}(router, graph)
}

// WaitForResult attend que toutes les tâches soient terminées et retourne les résultats.
func (wp *WorkerPool) WaitForResult() []Result {
	wp.wg.Wait()
	close(wp.resultChan)

	var results []Result
	for result := range wp.resultChan {
		results = append(results, result)
	}

	return results
}

// dijkstra calcule l'algorithme de Dijkstra pour un routeur spécifique.
func dijkstra(graph Graph, start string) map[string]int {
	// Implémentez Dijkstra ici
	// ...
	return nil // Remplacez par le résultat réel
}

// handleClient gère les connexions des clients.
func handleClient(conn net.Conn, wp *WorkerPool) {
	defer conn.Close()

	var graph Graph

	// Lire le dictionnaire d'adjacence depuis la connexion
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&graph); err != nil {
		fmt.Println("Erreur de lecture du dictionnaire d'adjacence:", err)
		return
	}

	// Récupérer le routeur à partir des données du client
	router := "A" // Remplacez par la logique appropriée

	// Calculer le plus court chemin vers tous les autres nœuds
	wp.CalculateShortestPath(router, graph)

	// Attendre la fin des tâches et récupérer les résultats
	results := wp.WaitForResult()

	// Envoyer les résultats au client
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(results); err != nil {
		fmt.Println("Erreur lors de l'envoi des résultats au client:", err)
	}
}

func main() {
	// Créer un pool de workers avec une taille de 5
	wp := NewWorkerPool(5)

	// Configurer le serveur TCP
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erreur lors de la création du serveur:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Serveur démarré sur http://localhost:8080")

	// Accepter les connexions des clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur lors de l'acceptation de la connexion:", err)
			continue
		}

		// Gérer la connexion client dans une goroutine
		go handleClient(conn, wp)
	}
}
