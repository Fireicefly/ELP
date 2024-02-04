package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

// Nom de l'entreprise
const company_name = "-Younes Factory-"

// Type Graph
type Graph map[string]map[string]int

// Fonction qui vérifie si une erreur est présente
func is_there_an_error(err error, errorMessage string) {
	if err != nil {
		fmt.Println(errorMessage, err)
		os.Exit(1)
	}
}

// Fonction qui ouvre un fichier JSON et le convertit en graphe
func Open_Json(file_name string) Graph {

	jsonData, err := os.Open(file_name)
	is_there_an_error(err, "Erreur lors de l'ouverture du fichier JSON :")
	defer jsonData.Close()

	var graph map[string]map[string]int
	decoder := json.NewDecoder(jsonData)
	err = decoder.Decode(&graph)
	is_there_an_error(err, "Erreur lors de la lecture du fichier JSON :")

	return graph
}

// Fonction qui établit une connexion avec le serveur
func make_conn(company_name string) net.Conn {

	conn, err := net.Dial("tcp", "localhost:8080")
	is_there_an_error(err, "Erreur lors de la connexion au serveur:")

	send_string(conn, company_name)

	fmt.Println("Connexion établi avec localhost:8080.")

	return conn
}

// Fonction qui reçoit le graph resultat du serveur
func receive_graph(conn net.Conn) Graph {

	var data Graph

	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&data)
	is_there_an_error(err, "Erreur lors de la réception des données JSON :")

	fmt.Println("Données reçues !")

	return data
}

// Fonction qui envoit le graph au serveur
func send_graph(conn net.Conn, data Graph) {

	encoder := json.NewEncoder(conn)
	err := encoder.Encode(data)
	is_there_an_error(err, "Erreur lors de l'envoi des données JSON :")
}

// Fonction qui envoit une chaîne de caractères au serveur contenant le nom de l'entreprise
func send_string(conn net.Conn, data string) {

	_, err := fmt.Fprintln(conn, data+"\n")
	is_there_an_error(err, "Erreur lors de l'envoi de la chaîne de caractères :")
}

// Fonction qui écrit le résultat dans un fichier JSON
func write_json(distances Graph) {

	resultJSON, err := json.Marshal(distances)
	is_there_an_error(err, "Erreur lors de la conversion en JSON :")

	file, err := os.Create("resultat.json")
	is_there_an_error(err, "Erreur lors de la création du fichier :")
	defer file.Close()

	_, err = file.Write(resultJSON)
	is_there_an_error(err, "Erreur lors de l'écriture dans le fichier :")
}

// Fonction principale
func main() {

	start := time.Now()

	file_name := "generated_graph.json"

	graph := Open_Json(file_name)
	conn := make_conn(company_name)

	send_graph(conn, graph)

	AllPairDistances := receive_graph(conn)

	write_json(AllPairDistances)

	elapsed := time.Since(start)
	fmt.Println("Temps d'execution :", elapsed)
}
