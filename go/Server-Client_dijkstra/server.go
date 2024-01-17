package main

import (
	"bufio"
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

func is_there_an_error(err error, errorMessage string) {
	if err != nil {
		fmt.Println(errorMessage, err)
		panic(err)
	}
}

// NewWorkerPool crée un nouveau pool de workers avec la taille spécifiée.
func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{
		workers:    make(chan struct{}, size),
		resultChan: make(chan Result),
	}
}

// CalculateShortestPath ajoute une tâche au pool pour calculer le plus court chemin.
// func (wp *WorkerPool) CalculateShortestPath(router string, graph Graph) {
// 	wp.wg.Add(1)
// 	wp.workers <- struct{}{}

// 	go func(router string, graph Graph) {
// 		defer wp.wg.Done()

// 		distances := dijkstra(graph, router)

// 		wp.resultChan <- Result{
// 			Router:    router,
// 			Distances: distances,
// 		}

// 		<-wp.workers
// 	}(router, graph)
// }

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

// handleClient gère les connexions des clients.
func handleClient(conn net.Conn, wp *WorkerPool) {
	// defer conn.Close()

	// var graph Graph
	// var message string
	// // Lire le dictionnaire d'adjacence depuis la connexion
	// decoder := json.NewDecoder(conn)
	// err := decoder.Decode(&graph)
	// is_there_an_error(err, "Erreur de lecture du dictionnaire d'adjacence:")

	// router := "A" // Remplacez par la logique appropriée

	// // Calculer le plus court chemin vers tous les autres nœuds
	// wp.CalculateShortestPath(router, graph)

	// // Attendre la fin des tâches et récupérer les résultats
	// results := wp.WaitForResult()

	// // Envoyer les résultats au client
	// encoder := json.NewEncoder(conn)
	// err = encoder.Encode(results)
	// is_there_an_error(err, "Erreur lors de l'envoi des résultats au client:")
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		client_ID := scanner.Text()
		fmt.Println("Client connecté avec ID :", client_ID)

		response := "Connection acceptée\n"
		conn.Write([]byte(response))
	}
	err := scanner.Err()
	is_there_an_error(err, "Erreur lors de la lecture de la chaîne:")
}

func main() {
	// Créer un pool de workers avec une taille de max_worker
	max_worker := 5
	wp := NewWorkerPool(max_worker)

	// Configurer le serveur TCP
	listener, err := net.Listen("tcp", ":8080")
	is_there_an_error(err, "Erreur lors de la création du serveur:")

	defer listener.Close()

	fmt.Println("Serveur démarré sur http://localhost:8080")

	// Accepter les connexions des clients
	for {
		conn, err := listener.Accept()
		is_there_an_error(err, "Erreur lors de l'acceptation de la connexion:")

		// Gérer la connexion client dans une goroutine
		go handleClient(conn, wp)
	}
}
