package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Graph map[string]map[string]int

func is_there_an_error(err error, errorMessage string) {
	if err != nil {
		fmt.Println(errorMessage, err)
		os.Exit(1)
	}
}

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

func make_conn(router_name string) net.Conn {

	conn, err := net.Dial("tcp", "localhost:8080")
	is_there_an_error(err, "Erreur lors de la connexion au serveur:")

	send_string(conn, router_name)

	fmt.Println("Connexion établi avec localhost:8080.")

	return conn
}

func receive_json(conn net.Conn) Graph {

	var data Graph

	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&data)
	is_there_an_error(err, "Erreur lors de la réception des données JSON :")

	fmt.Println("Données reçues !")

	return data
}

func send_json(conn net.Conn, data Graph) {

	encoder := json.NewEncoder(conn)
	err := encoder.Encode(data)
	is_there_an_error(err, "Erreur lors de l'envoi des données JSON :")
}

func send_string(conn net.Conn, data string) {

	_, err := fmt.Fprintln(conn, data+"\n")
	is_there_an_error(err, "Erreur lors de l'envoi de la chaîne de caractères :")
}

func write_json(distances Graph) {

	resultJSON, err := json.Marshal(distances)
	is_there_an_error(err, "Erreur lors de la conversion en JSON :")

	file, err := os.Create("resultat.json")
	is_there_an_error(err, "Erreur lors de la création du fichier :")
	defer file.Close()

	_, err = file.Write(resultJSON)
	is_there_an_error(err, "Erreur lors de l'écriture dans le fichier :")
}

func main() {

	start := time.Now()

	company_name := "-Younes Factory-"
	file_name := "generated_graph.json"

	graph := Open_Json(file_name)
	conn := make_conn(company_name)

	send_json(conn, graph)

	AllPairDistances := receive_json(conn)

	write_json(AllPairDistances)

	elapsed := time.Since(start)
	fmt.Println("Temps d'execution :", elapsed)
}
