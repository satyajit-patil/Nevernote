package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Healthcheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthcheck)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_ListNotebooks(t *testing.T) {
	Notebooks = make(map[string][]Note)
	Notebooks["English"] = []Note{
		Note{Id: "1", Title: "Hamlet", Body: "This is Hamlet", Tags: []string{"Classics", "Shakespeare"}, Created: "HamletCreated", LastModified: "HamletModified"},
		Note{Id: "2", Title: "Animal Farm", Body: "Farm of Animals", Tags: []string{"Classics"}, Created: "AnimalCreated", LastModified: "AnimalModified"},
	}
	Notebooks["Math"] = []Note{
		Note{Id: "1", Title: "Algebra", Body: "PEMDAS", Tags: []string{"HS"}, Created: "AlgebraCreated", LastModified: "AlgebraModified"},
	}

	req, err := http.NewRequest("GET", "/listNotebooks", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listNotebooks)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expectedA := "[\"English\",\"Math\"]\n"
	expectedB := "[\"Math\",\"English\"]\n"

	if rr.Body.String() != expectedA && rr.Body.String() != expectedB {
		t.Errorf("handler returned unexpected body: got %v want %v or %v",
			rr.Body.String(), expectedA, expectedB)
	}
}

func Test_CreateNotebook(t *testing.T) {
	Notebooks = make(map[string][]Note)
	req, err := http.NewRequest("POST", "/createNotebook/Science", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/createNotebook/{title}", createNotebook)
	router.ServeHTTP(rr, req)

	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "[\"Science\"]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_DeleteNotebook(t *testing.T) {
	Notebooks = make(map[string][]Note)
	Notebooks["English"] = []Note{
		Note{Id: "1", Title: "Hamlet", Body: "This is Hamlet", Tags: []string{"Classics", "Shakespeare"}, Created: "HamletCreated", LastModified: "HamletModified"},
		Note{Id: "2", Title: "Animal Farm", Body: "Farm of Animals", Tags: []string{"Classics"}, Created: "AnimalCreated", LastModified: "AnimalModified"},
	}
	Notebooks["Math"] = []Note{
		Note{Id: "1", Title: "Algebra", Body: "PEMDAS", Tags: []string{"HS"}, Created: "AlgebraCreated", LastModified: "AlgebraModified"},
	}

	req, err := http.NewRequest("DELETE", "/deleteNotebook/Math", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/deleteNotebook/{title}", deleteNotebook)
	router.ServeHTTP(rr, req)

	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "[\"English\"]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_NumberOfNotes(t *testing.T) {
	Notebooks = make(map[string][]Note)
	Notebooks["English"] = []Note{
		Note{Id: "1", Title: "Hamlet", Body: "This is Hamlet", Tags: []string{"Classics", "Shakespeare"}, Created: "HamletCreated", LastModified: "HamletModified"},
		Note{Id: "2", Title: "Animal Farm", Body: "Farm of Animals", Tags: []string{"Classics"}, Created: "AnimalCreated", LastModified: "AnimalModified"},
	}

	req, err := http.NewRequest("GET", "/numberOfNotes/English", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/numberOfNotes/{title}", numberOfNotes)
	router.ServeHTTP(rr, req)

	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "2\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_ListNotes(t *testing.T) {
	Notebooks = make(map[string][]Note)
	Notebooks["English"] = []Note{
		Note{Id: "1", Title: "Hamlet", Body: "This is Hamlet", Tags: []string{"Classics", "Shakespeare"}, Created: "HamletCreated", LastModified: "HamletModified"},
		Note{Id: "2", Title: "Animal Farm", Body: "Farm of Animals", Tags: []string{"Classics"}, Created: "AnimalCreated", LastModified: "AnimalModified"},
	}

	data := []byte(`{"Tags": ["Shakespeare"]}`)

	req, err := http.NewRequest("GET", "/listNotes/English", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/listNotes/{title}", listNotes)
	router.ServeHTTP(rr, req)

	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "[{\"Id\":\"1\",\"Title\":\"Hamlet\",\"Body\":\"This is Hamlet\",\"Tags\":[\"Classics\",\"Shakespeare\"],\"Created\":\"HamletCreated\",\"LastModified\":\"HamletModified\"}]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_CreateNote(t *testing.T) {
	Notebooks = make(map[string][]Note)
	Notebooks["English"] = []Note{}

	data := []byte(`{"Title": "Hamlet", "Body": "This is Hamlet", "Tags": ["Classics", "Shakespeare"], "Created": "HamletCreated", "LastModified": "HamletModified"}`)

	req, err := http.NewRequest("POST", "/createNote/English", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/createNote/{title}", createNote)
	router.ServeHTTP(rr, req)

	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "\"Title\":\"Hamlet\",\"Body\":\"This is Hamlet\",\"Tags\":[\"Classics\",\"Shakespeare\"]"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_UpdateNote(t *testing.T) {
	Notebooks = make(map[string][]Note)
	Notebooks["English"] = []Note{
		Note{Id: "1", Title: "Hamlet", Body: "This is Hamlet", Tags: []string{"Classics", "Shakespeare"}, Created: "HamletCreated", LastModified: "HamletModified"},
	}

	data := []byte(`{"Title": "Hamlet", "Body": "This is Hamlet 2.0", "Tags": ["Classics", "Shakespeare"]}`)

	req, err := http.NewRequest("UPDATE", "/updateNote/English/1", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/updateNote/{title}/{noteId}", updateNote)
	router.ServeHTTP(rr, req)

	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "\"Title\":\"Hamlet\",\"Body\":\"This is Hamlet 2.0\",\"Tags\":[\"Classics\",\"Shakespeare\"],\"Created\":\"HamletCreated\""
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_ReadNote(t *testing.T) {
	Notebooks = make(map[string][]Note)
	Notebooks["English"] = []Note{
		Note{Id: "1", Title: "Hamlet", Body: "This is Hamlet", Tags: []string{"Classics", "Shakespeare"}, Created: "HamletCreated", LastModified: "HamletModified"},
		Note{Id: "2", Title: "Animal Farm", Body: "Farm of Animals", Tags: []string{"Classics"}, Created: "AnimalCreated", LastModified: "AnimalModified"},
	}

	req, err := http.NewRequest("GET", "/readNote/English/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/readNote/{title}/{noteId}", readNote)
	router.ServeHTTP(rr, req)

	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "{\"Id\":\"2\",\"Title\":\"Animal Farm\",\"Body\":\"Farm of Animals\",\"Tags\":[\"Classics\"],\"Created\":\"AnimalCreated\",\"LastModified\":\"AnimalModified\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_DeleteNote(t *testing.T) {
	Notebooks = make(map[string][]Note)
	Notebooks["English"] = []Note{
		Note{Id: "1", Title: "Hamlet", Body: "This is Hamlet", Tags: []string{"Classics", "Shakespeare"}, Created: "HamletCreated", LastModified: "HamletModified"},
		Note{Id: "2", Title: "Animal Farm", Body: "Farm of Animals", Tags: []string{"Classics"}, Created: "AnimalCreated", LastModified: "AnimalModified"},
	}

	req, err := http.NewRequest("DELETE", "/deleteNote/English/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/deleteNote/{title}/{noteId}", deleteNote)
	router.ServeHTTP(rr, req)

	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "[{\"Id\":\"1\",\"Title\":\"Hamlet\",\"Body\":\"This is Hamlet\",\"Tags\":[\"Classics\",\"Shakespeare\"],\"Created\":\"HamletCreated\",\"LastModified\":\"HamletModified\"}]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
