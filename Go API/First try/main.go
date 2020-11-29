package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Main function
func main() {
	// Init router
	r := mux.NewRouter()

	// Hardcoded data - @todo: add database
	contacts = append(contacts, Contact{Name: "Friend_1", Phone: "98xxx-xxxxx", Email: "person1@mail.com"})
	contacts = append(contacts, Contact{Name: "Friend_2", Phone: "96xxx-xxxxx", Email: "person2@mail.com"})
	contacts = append(contacts, Contact{Name: "Friend_3", Phone: "97xxx-xxxxx", Email: "person3@mail.com"})

	// Route handles & endpoints
	r.HandleFunc("/contacts", getContacts).Methods("GET")
	r.HandleFunc("/contacts/{name}", getContact).Methods("GET") // get
	r.HandleFunc("/contacts", createContact).Methods("POST")    // create new contact
	// r.HandleFunc("/contacts/{name}", updateContact).Methods("PUT") // update contact
	r.HandleFunc("/contacts/{name}", deleteContact).Methods("DELETE") // delete

	// Start server
	// We put this in log.Fatal as it will return the necessary response codes for any request made
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Contact struct (model)
type Contact struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// Init contacts var as a slice Contact struct
var contacts []Contact

///
/// Handler functions
///

// Get all contacts
func getContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

// Get single contact
func getContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Looping through contacts and find one with the id from the params
	for _, item := range contacts {
		if item.Name == params["name"] {
			json.NewEncoder(w).Encode(item)
			return // this is important!
		}
	}
	json.NewEncoder(w).Encode(&Contact{})
}

func createContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contact Contact
	_ = json.NewDecoder(r.Body).Decode(&contact)
	contacts = append(contacts, contact)
	json.NewEncoder(w).Encode(contact)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for idx, item := range contacts {
		if item.Name == params["name"] {
			contacts = append(contacts[:idx], contacts[idx+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(contacts)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for idx, item := range contacts {
		if item.Name == params["name"] {
			contacts = append(contacts[:idx], contacts[idx+1:]...)
			var contact Contact
			_ = json.NewDecoder(r.Body).Decode(&contact)
			contact.Name = params["name"]      // is this really necessary if the request houses the changes?
			json.NewEncoder(w).Encode(contact) // if I'm correct about the above, then couldn't we just send back r?
			return
		}
	}
}
