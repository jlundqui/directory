// directory
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Contact contains data for directory entries.
type Contact struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// Directory contains the constacts list and tracks the latest ID value.
type Directory struct {
	LastIDValue int
	List        []Contact
}

var directory Directory

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

func getContactHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, contact := range directory.List {
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
	json.NewEncoder(w).Encode(directory.List)
}

func createContactHandler(w http.ResponseWriter, req *http.Request) {
	var contact Contact
	_ = json.NewDecoder(req.Body).Decode(&contact)
	directory.LastIDValue++
	contact.ID = strconv.Itoa(directory.LastIDValue)
	directory.List = append(directory.List, contact)
	json.NewEncoder(w).Encode(directory)
}

func deleteContactHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, contact := range directory.List {
		if contact.ID == params["id"] {
			directory.List = append(directory.List[:index], directory.List[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(directory)
}

func main() {
	fmt.Println("Server starting")
	log.Fatal(http.ListenAndServe(":8081", Handlers()))
}

// Handlers - creates router and adds handlers
func Handlers() *mux.Router {
	// create directory
	directory.LastIDValue = 2
	directory.List = append(directory.List, Contact{ID: "1", Firstname: "Ann", Lastname: "Adams",
		Email: "aadams@fakemail.net", Phone: "123-456-7890"})
	directory.List = append(directory.List, Contact{ID: "2", Firstname: "Brenda", Lastname: "Bowman",
		Email: "bbowman@fakemail.net", Phone: "123-890-4567"})

	//create router and add handlers
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/directory", getDirectoryHandler).Methods("GET")
	router.HandleFunc("/contact/{id}", getContactHandler).Methods("GET")
	router.HandleFunc("/contact/{firstname}/{lastname}/{email}/{phone}", createContactHandler).Methods("POST")
	router.HandleFunc("/contact/{id}", deleteContactHandler).Methods("DELETE")

	return router
}
