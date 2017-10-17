// directory_test
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server                   *httptest.Server
	reader                   io.Reader
	contactURL, directoryURL string
)

func init() {
	server = httptest.NewServer(Handlers())

	directoryURL = fmt.Sprintf("%s/directory", server.URL)
	contactURL = fmt.Sprintf("%s/contact", server.URL)
}

func TestGetDirectory(t *testing.T) {
	request, err := http.NewRequest("GET", directoryURL, nil)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		// failed during request send
		t.Error(err)
		return
	}

	var expContactList []Contact
	expContactList = append(expContactList, Contact{ID: "1", Firstname: "Ann", Lastname: "Adams",
		Email: "aadams@fakemail.net", Phone: "123-456-7890"})
	expContactList = append(expContactList, Contact{ID: "2", Firstname: "Brenda", Lastname: "Bowman",
		Email: "bbowman@fakemail.net", Phone: "123-890-4567"})
	numExpContacts := len(expContactList)

	var contactList []Contact
	err = json.NewDecoder(response.Body).Decode(&contactList)
	if err != nil {
		t.Error(err) //Failed to decode the response
	}

	numContacts := len(contactList)
	if len(contactList) != len(expContactList) {
		t.Errorf("Number of contacts, %v, received differs from expected, %v.",
			numContacts, numExpContacts)
	}

	for i := 0; i < numContacts; i++ {
		if !contactList[i].Equal(expContactList[i]) {
			t.Errorf("Response returned incorrect contact info. Received %v, expected %v.",
				contactList[i], expContactList[i])
		}
	}
}

func TestGetContact(t *testing.T) {
	// create and send request
	getContactURL := fmt.Sprintf("%s/1", contactURL)
	request, err := http.NewRequest("GET", getContactURL, nil)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		// failed during request send
		t.Error(err)
		return
	}

	expContact := Contact{ID: "1", Firstname: "Ann", Lastname: "Adams",
		Email: "aadams@fakemail.net", Phone: "123-456-7890"}

	var contact Contact
	err = json.NewDecoder(response.Body).Decode(&contact)
	if err != nil {
		t.Error(err) //Failed to decode the response
	}

	if !contact.Equal(expContact) {
		t.Errorf("Response returned incorrect contact information. Expected %v, received %v",
			expContact, contact)
	}
}

func TestCreateContact(t *testing.T) {
	contactJSON := `{"firstname": "Jane", "lastname": "Jones",
	 "email": "jjones@doesnotexist.org", "phone": "987-645-3210"}`

	// convert string to reader
	reader = strings.NewReader(contactJSON)

	// create request with reader
	request, err := http.NewRequest("POST", contactURL, reader)
	if err != nil {
		// failed to create request
		t.Error(err)
		return
	}

	// send request
	_, err = http.DefaultClient.Do(request)

	if err != nil {
		// failed during request send
		t.Error(err)
		return
	}
}

func TestDeleteContact(t *testing.T) {
	// create and send request
	delContactURL := fmt.Sprintf("%s/1", contactURL)
	request, err := http.NewRequest("DELETE", delContactURL, nil)

	if err != nil {
		// failed to create request
		t.Error(err)
		return
	}

	_, err = http.DefaultClient.Do(request)

	if err != nil {
		// failed during request send
		t.Error(err)
		return
	}

}

//
// func TestIndexHandler(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/health-check", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(HealthCheckHandler)
//
// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	handler.ServeHTTP(rr, req)
//
// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
//
// 	// Check the response body is what we expect.
// 	expected := `{"alive": true}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
