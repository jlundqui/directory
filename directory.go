// directory
package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Contact struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

var directory []Contact

func IndexEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
}

func GetContactEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, contact := range directory {
		if contact.ID == params["id"] {
			json.NewEncoder(w).Encode(contact)
			return
		}
	}
	json.NewEncoder(w).Encode(&Contact{})
}

func GetDirectoryEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(directory)
}

//func CreateContactEndpoint(w http.ResponseWriter, req *http.Request) {
//	params := mux.Vars(req)
//	var contact Contact
//	_ = json.NewDecoder(req.Body).Decode(&contact)
//	contact.ID = params["id"]
//	contact.Firstname = params["firstname"]
//	contact.Lastname = params["lastname"]
//	contact.Email = params["email"]
//	contact.Phone = params["phone"]
//	directory = append(directory, contact)
//	json.NewEncoder(w).Encode(directory)
//}

func DeleteContactEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, contact := range directory {
		if contact.ID == params["id"] {
			directory = append(directory[:index], directory[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(directory)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	directory = append(directory, Contact{ID: "1", Firstname: "Ann", Lastname: "Adams", Email: "aadams@fakemail.net"})
	directory = append(directory, Contact{ID: "2", Firstname: "Brenda", Lastname: "Bowman"})
	router.HandleFunc("/", IndexEndpoint)
	router.HandleFunc("/directory", GetDirectoryEndpoint).Methods("GET")
	router.HandleFunc("/contact/{id}", GetContactEndpoint).Methods("GET")
	//	router.HandleFunc("/contact/{id}/{firstname}/{lastname}/{email}/{phone}", CreateContactEndpoint).Methods("POST")
	router.HandleFunc("/contact/{id}", DeleteContactEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}
