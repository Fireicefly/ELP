package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type Graph map[string]map[string]int

const TAILLE = 50

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

	nodes := []string{"R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15", "R16", "R17", "R18", "R19", "R20", "R21", "R22", "R23", "R24", "R25", "R26", "R27", "R28", "R29", "R30", "R31", "R32", "R33", "R34", "R35", "R36", "R37", "R38", "R39", "R40", "R41", "R42", "R43", "R44", "R45", "R46", "R47", "R48", "R49", "R50"}

	adjacencyMap := convertToAdjacencyMap(matrix, nodes)

	fmt.Println(adjacencyMap)

	write_json(adjacencyMap)
}
