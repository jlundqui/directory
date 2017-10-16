// directory
package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	// "strings"

	"github.com/gorilla/mux"
	// "logging"
)

// Contact contains data for directory entries.
type Contact struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// var directoryJSON = `{"id": "1",  "firstname": "Ann", "lastname": "Adams",
// 	 "email": "aadams@fakemail.net", "phone": "123-456-7890"},
// 	 {"id": "2", "firstname": "Brenda", "lastname": "Bowman",
// 		 "email": "bbowman@fakemail.net", "phone": "123-890-4567"}`

var directory []Contact

// Equal compares contacts for equality.
func (c Contact) Equal(contact Contact) bool {
	if c.Firstname == contact.Firstname &&
		c.Lastname == contact.Lastname &&
		c.Email == contact.Email &&
		c.Phone == contact.Phone {
		return true
	}

	return false
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
}

func getContactHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, contact := range directory {
		if contact.ID == params["id"] {
			json.NewEncoder(w).Encode(contact)
			return
		}
	}

	json.NewEncoder(w).Encode(&Contact{})
}

func getDirectoryHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(directory)
}

func createContactHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var contact Contact
	_ = json.NewDecoder(req.Body).Decode(&contact)
	contact.ID = params["id"]
	contact.Firstname = params["firstname"]
	contact.Lastname = params["lastname"]
	contact.Email = params["email"]
	contact.Phone = params["phone"]
	directory = append(directory, contact)
	json.NewEncoder(w).Encode(directory)
}

func deleteContactHandler(w http.ResponseWriter, req *http.Request) {
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
	// reader := strings.NewReader(directoryJSON)
	// logging.Info.Printf("Reader length %v\n", reader.Len())
	// if err := json.NewDecoder(reader).Decode(&directory); err != nil {
	// 	logging.Error.Printf("Unable to create directory. Error: %v\n", err)
	// }

	fmt.Println("Server starting")
	log.Fatal(http.ListenAndServe(":8081", Handlers()))
}

// Handlers - creates router and adds handlers
func Handlers() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// directory = append(directory, Contact{ID: 1, Firstname: "Ann", Lastname: "Adams", Email: "aadams@fakemail.net"})
	// directory = append(directory, Contact{ID: 2, Firstname: "Brenda", Lastname: "Bowman"})
	directory = append(directory, Contact{ID: "1", Firstname: "Ann", Lastname: "Adams",
		Email: "aadams@fakemail.net", Phone: "123-456-7890"})
	directory = append(directory, Contact{ID: "2", Firstname: "Brenda", Lastname: "Bowman",
		Email: "bbowman@fakemail.net", Phone: "123-890-4567"})
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/directory", getDirectoryHandler).Methods("GET")
	router.HandleFunc("/contact/{id}", getContactHandler).Methods("GET")
	router.HandleFunc("/contact/{firstname}/{lastname}/{email}/{phone}", createContactHandler).Methods("POST")
	router.HandleFunc("/contact/{id}", deleteContactHandler).Methods("DELETE")

	return router
}
