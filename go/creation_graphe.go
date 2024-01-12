package main

//jojolechocbarto
import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type Graphe map[string]map[string]int

type Router struct {
	Name      string         `json:"name"`
	Neighbors map[string]int `json:"neighbors"`
}

func main() {
	// Nombre de routeurs
	nombreRouteurs := 100

	// Créer une slice pour stocker les routeurs
	var graphData = make(Graphe)

	// Générer les informations aléatoires pour les liens entre les routeurs
	for i := 1; i <= nombreRouteurs; i++ {
		routerName := fmt.Sprintf("R%d", i)

		numNeighbors := rand.Intn(3) + 1 // Nombre aléatoire de voisins par routeur
		graphData[routerName] = make(map[string]int)
		for j := 0; j < numNeighbors; j++ {
			index_router := rand.Intn(nombreRouteurs) + 1
			if index_router != i {
				neighborName := fmt.Sprintf("R%d", index_router)

				weight := rand.Intn(10) + 1 // Poids du lien aléatoire

				graphData[routerName][neighborName] = weight
				if graphData[neighborName] == nil {
					graphData[neighborName] = make(map[string]int)
				}
				graphData[neighborName][routerName] = weight

			}
		}

	}

	// Enregistrez la slice dans un fichier JSON
	jsonData, err := json.MarshalIndent(graphData, "", "    ")
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON:", err)
		return
	}

	// Enregistrez le JSON dans un fichier
	jsonFileName := "graphe.json"
	err = os.WriteFile(jsonFileName, jsonData, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier JSON:", err)
		return
	}

	fmt.Println("Fichier graphe.json créé avec succès.")
}
