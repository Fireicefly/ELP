package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Graph map[string]map[string]int

func is_there_an_error(err error, errorMessage string) {
	if err != nil {
		fmt.Println(errorMessage, err)
		os.Exit(1)
	}
}

func Open_Json(file_name string) Graph {

	fmt.Println("Fichier en cours d'ouverture...")

	jsonData, err := os.Open(file_name)
	is_there_an_error(err, "Erreur lors de l'ouverture du fichier JSON :")
	defer jsonData.Close()

	var graph map[string]map[string]int
	decoder := json.NewDecoder(jsonData)
	err = decoder.Decode(&graph)
	is_there_an_error(err, "Erreur lors de la lecture du fichier JSON :")

	fmt.Println("Fichier ouvert.")

	return graph
}

func make_conn(router_name string) net.Conn {

	fmt.Println("Connexion avec localhost:8080 en cours...")

	conn, err := net.Dial("tcp", "localhost:8080")
	is_there_an_error(err, "Erreur lors de la connexion au serveur:")
	//defer conn.Close()

	send_string(conn, router_name)

	fmt.Println("Connexion établi avec localhost:8080.")

	return conn
}

func receive_json(conn net.Conn) map[string]int {

	var distances map[string]int

	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&distances)
	is_there_an_error(err, "Erreur lors de la réception des données JSON :")

	fmt.Println("Données reçues !")

	return distances
}

func send_json(conn net.Conn, data Graph) {

	fmt.Println("Fichier en cours d'envoi...")

	encoder := json.NewEncoder(conn)
	err := encoder.Encode(data)
	is_there_an_error(err, "Erreur lors de l'envoi des données JSON :")

	fmt.Println("Fichier envoyé.")
}

func send_string(conn net.Conn, data string) {

	_, err := fmt.Fprintln(conn, data+"\n")
	is_there_an_error(err, "Erreur lors de l'envoi de la chaîne de caractères :")
}

func write_json(distances map[string]int) {

	fmt.Println("Fichier en cours d'écriture.")

	resultJSON, err := json.Marshal(distances)
	is_there_an_error(err, "Erreur lors de la conversion en JSON :")

	file, err := os.Create("resultat.json")
	is_there_an_error(err, "Erreur lors de la création du fichier :")
	defer file.Close()

	_, err = file.Write(resultJSON)
	is_there_an_error(err, "Erreur lors de l'écriture dans le fichier :")

	fmt.Println("Fichier crée.")
}

func main() {

	router_name := "O"
	file_name := "graph.json"

	graph := Open_Json(file_name)
	conn := make_conn(router_name)

	send_json(conn, graph)

	distances := receive_json(conn)

	write_json(distances)

}
