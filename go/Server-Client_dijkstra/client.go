package main

import (
	"bufio"
	"fmt"
	"net"
)

func is_there_an_error(err error, errorMessage string) {
	if err != nil {
		fmt.Println(errorMessage, err)
		return
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	is_there_an_error(err, "Erreur lors de la connexion au serveur:")
	defer conn.Close()

	message := "jsp frr"
	_, err = fmt.Fprintln(conn, message)
	is_there_an_error(err, "Erreur lors de l'envoi des données au serveur:")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		receivedString := scanner.Text()
		fmt.Println("message reçu :", receivedString)
	}
}
