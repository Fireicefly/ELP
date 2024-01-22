package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

type Graph map[string]map[string]int

const TAILLE = 500

func convertToAdjacencyMap(matrix [][]int, nodes []string) Graph {
	adjacencyMap := make(Graph)

	for _, node := range nodes {
		adjacencyMap[node] = make(map[string]int)
	}

	for i, source := range nodes {
		for j, weight := range matrix[i] {
			if weight != 0 {
				destination := nodes[j]
				adjacencyMap[source][destination] = weight
			}
		}
	}

	return adjacencyMap
}

func generateRandomGraphMatrix(size int) [][]int {

	matrix := make([][]int, size)
	for i := range matrix {
		matrix[i] = make([]int, size)
	}

	for i := range matrix {
		matrix[i][i] = 0

		for j := range matrix[i] {
			if i != j {
				randomInteger := rand.Intn(20)
				if randomInteger >= 9 && randomInteger <= 20 {
					randomInteger = 0
				}
				matrix[i][j] = randomInteger
				matrix[j][i] = randomInteger
			}
		}
	}

	return matrix
}

func write_json(graph Graph) {

	resultJSON, err := json.Marshal(graph)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("generated_graph.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.Write(resultJSON)
	if err != nil {
		panic(err)
	}

	fmt.Println("Fichier crée.")
}

func main() {
	matrix := generateRandomGraphMatrix(TAILLE)

	fmt.Println("Matrice d'adjacence générée aléatoirement :")
	for _, row := range matrix {
		fmt.Println(row)
	}

	var nodes []string

	for i := 1; i <= TAILLE; i++ {
		node := "R" + strconv.Itoa(i)
		nodes = append(nodes, node)
	}
	adjacencyMap := convertToAdjacencyMap(matrix, nodes)

	fmt.Println(adjacencyMap)

	write_json(adjacencyMap)
}
