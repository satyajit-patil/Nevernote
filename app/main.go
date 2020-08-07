package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Note struct {
	Id           string   `json:"Id"`
	Title        string   `json:"Title"`
	Body         string   `json:"Body"`
	Tags         []string `json:"Tags"`
	Created      string   `json:"Created"`
	LastModified string   `json:"LastModified"`
}

// let's declare a global Notebooks hashmap that we can then populate
// in our main function  to simulate a database
var Notebooks map[string][]Note

var idCounter int

func healthcheck(w http.ResponseWriter, r *http.Request) {
	/**
	Function: healthcheck
	Description: check if server is running
	*/
	io.WriteString(w, `{"alive": true}`)
}

func listNotebooks(w http.ResponseWriter, r *http.Request) {
	/**
	Function: listNotebooks
	Description: Returns a list of all Notebook titles
	*/
	notebookTitles := make([]string, 0, len(Notebooks))
	for title := range Notebooks {
		notebookTitles = append(notebookTitles, title)
	}
	json.NewEncoder(w).Encode(notebookTitles)
}

func createNotebook(w http.ResponseWriter, r *http.Request) {
	/**
	Function: createNotebook
	Description: Creates a new notebook with a given title
	*/

	vars := mux.Vars(r)
	title := vars["title"]

	// add blank notebook
	Notebooks[title] = []Note{}

	notebookTitles := make([]string, 0, len(Notebooks))
	for title := range Notebooks {
		notebookTitles = append(notebookTitles, title)
	}

	json.NewEncoder(w).Encode(notebookTitles)
}

func deleteNotebook(w http.ResponseWriter, r *http.Request) {
	/**
	Function: deleteNotebook
	Description: Deletes a notebook with a given title
	*/

	vars := mux.Vars(r)
	title := vars["title"]

	// delete notebook
	if Notebooks[title] == nil {
		returnError(w, "Notebook \""+title+"\" does not exist")
		return
	}
	delete(Notebooks, title)

	notebookTitles := make([]string, 0, len(Notebooks))
	for title := range Notebooks {
		notebookTitles = append(notebookTitles, title)
	}

	json.NewEncoder(w).Encode(notebookTitles)
}

func numberOfNotes(w http.ResponseWriter, r *http.Request) {
	/**
	Function: numberOfNotes
	Description: Get the number of notes in a notebook
	*/

	vars := mux.Vars(r)
	title := vars["title"]

	numNotes := len(Notebooks[title])

	json.NewEncoder(w).Encode(numNotes)
}

func listNotes(w http.ResponseWriter, r *http.Request) {
	/**
	Function: listNotes
	Description: List all notes in a notebook that match tags in body
	*/

	// parse path variables and body
	vars := mux.Vars(r)
	title := vars["title"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var filter struct {
		Tags []string `json:"Tags"`
	}
	json.Unmarshal(reqBody, &filter)

	// get notebook
	notebook := Notebooks[title]
	if notebook == nil {
		json.NewEncoder(w).Encode([]Note{})
		return
	}

	// get valid notes
	var filteredNotes []Note
	for _, note := range Notebooks[title] {
		if isSubset(filter.Tags, note.Tags) {
			filteredNotes = append(filteredNotes, note)
		}
	}

	json.NewEncoder(w).Encode(filteredNotes)

}

func isSubset(parent []string, child []string) bool {
	/**
	Function: isSubset
	Description: Determine if all strings in the parent are in the child
	*/
	for _, subParent := range parent {
		found := false
		for _, subChild := range child {
			if subChild == subParent {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

func createNote(w http.ResponseWriter, r *http.Request) {
	/**
	Function: createNote
	Description: Create a note in a notebook
	*/

	// parse path variables and body
	vars := mux.Vars(r)
	title := vars["title"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var note Note
	json.Unmarshal(reqBody, &note)
	if note.Title == "" {
		returnError(w, "Need Title to create note")
		return
	}
	if note.Body == "" {
		returnError(w, "Need Body to create note")
		return
	}
	if note.Tags == nil {
		returnError(w, "Need Tags to create note")
		return
	}

	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006.01.02 15:04:05")
	note.Created = currentTimeString
	note.LastModified = currentTimeString
	note.Id = strconv.Itoa(idCounter)
	idCounter++

	// get notebook
	notebook := Notebooks[title]
	if notebook == nil {
		returnError(w, "Notebook \""+title+"\" does not exist")
		return
	}

	// add Note to notebook
	Notebooks[title] = append(Notebooks[title], note)

	json.NewEncoder(w).Encode(Notebooks[title])
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	/**
	Function: updateNote
	Description: Update a note (with a specific id) in a notebook
	*/

	// get path variables and body
	vars := mux.Vars(r)
	title := vars["title"]
	noteId := vars["noteId"]

	reqBody, _ := ioutil.ReadAll(r.Body)

	var note Note
	json.Unmarshal(reqBody, &note)
	if note.Title == "" {
		returnError(w, "Need Title to create note")
		return
	}
	if note.Body == "" {
		returnError(w, "Need Body to create note")
		return
	}
	if note.Tags == nil {
		returnError(w, "Need Tags to create note")
		return
	}

	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006.01.02 15:04:05")
	note.LastModified = currentTimeString

	// get notebook
	notebook := Notebooks[title]
	if notebook == nil {
		returnError(w, "Notebook \""+title+"\" does not exist")
		return
	}

	// update note
	noteUpdated := false
	for i, noteItr := range notebook {
		if noteItr.Id == noteId {
			note.Id = noteItr.Id
			note.Created = noteItr.Created
			notebook[i] = note
			noteUpdated = true
			break
		}
	}

	if !noteUpdated {
		returnError(w, "Note with id \""+noteId+"\" does not exist")
		return
	}

	json.NewEncoder(w).Encode(notebook)
}

func readNote(w http.ResponseWriter, r *http.Request) {
	/**
	Function: readNote
	Description: Get a note (based on id) from a notebook
	*/

	vars := mux.Vars(r)
	title := vars["title"]
	noteId := vars["noteId"]

	// get notebook
	notebook := Notebooks[title]
	if notebook == nil {
		returnError(w, "Notebook \""+title+"\" does not exist")
		return
	}

	// get all notes with given title
	readNote := false
	var readNoteBody Note
	for _, noteItr := range notebook {
		if noteItr.Id == noteId {
			readNoteBody = noteItr
			readNote = true
			break
		}
	}

	if !readNote {
		returnError(w, "Note with id \""+noteId+"\" does not exist")
		return
	}

	json.NewEncoder(w).Encode(readNoteBody)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	/**
	Function: deleteNote
	Description: Delete a notes (with a specific id) in a notebook
	*/

	vars := mux.Vars(r)
	title := vars["title"]
	noteId := vars["noteId"]

	// get notebook
	notebook := Notebooks[title]
	if notebook == nil {
		returnError(w, "Notebook \""+title+"\" does not exist")
		return
	}

	// delete valid notes
	deleteNote := false
	for i, noteItr := range notebook {
		if noteItr.Id == noteId {
			notebook = append(notebook[:i], notebook[i+1:]...)
			Notebooks[title] = notebook
			deleteNote = true
			break
		}
	}

	if !deleteNote {
		returnError(w, "Note with title \""+noteId+"\" does not exist")
		return
	}

	json.NewEncoder(w).Encode(notebook)
}

func returnError(out http.ResponseWriter, err string) {
	/**
	Function: returnError
	Description: Returns an http error
	*/
	if strings.Contains(err, "unauthorized") {
		http.Error(out, err, http.StatusUnauthorized)
	} else {
		http.Error(out, err, http.StatusInternalServerError)
	}
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// add routes
	myRouter.HandleFunc("/healthcheck", healthcheck)

	myRouter.HandleFunc("/listNotebooks", listNotebooks).Methods("GET")

	myRouter.HandleFunc("/createNotebook/{title}", createNotebook).Methods("POST")

	myRouter.HandleFunc("/deleteNotebook/{title}", deleteNotebook).Methods("DELETE")

	myRouter.HandleFunc("/numberOfNotes/{title}", numberOfNotes).Methods("GET")

	myRouter.HandleFunc("/listNotes/{title}", listNotes).Methods("GET")

	myRouter.HandleFunc("/createNote/{title}", createNote).Methods("POST")

	myRouter.HandleFunc("/updateNote/{title}/{noteId}", updateNote).Methods("UPDATE")

	myRouter.HandleFunc("/readNote/{title}/{noteId}", readNote).Methods("GET")

	myRouter.HandleFunc("/deleteNote/{title}/{noteId}", deleteNote).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

func main() {
	fmt.Println("Rest API - Nevernote")

	Notebooks = make(map[string][]Note)
	Notebooks["English"] = []Note{
		Note{Id: "1", Title: "Hamlet", Body: "This is Hamlet", Tags: []string{"Classics", "Shakespeare"}, Created: "HamletCreated", LastModified: "HamletModified"},
	}
	idCounter = 0
	handleRequests()
}
