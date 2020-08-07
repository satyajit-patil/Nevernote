# Nevernote API

## How to run the API Server

1. Clone the repository with `get clone https://github.com/satyajit-patil/Nevernote.git`

2. `cd app`

3. Build a version of the docker container: 
`docker build --tag nevernote:{VERSION} .`

4. Run the version of the docker container with port 5000 forwarded:
`docker run -itd -p 5000:5000 nevernote:{VERSION}`

5. Check that the api server is running by running:
`curl http://localhost:5000/healthcheck`


## Endpoints Description

### Healthcheck

```
    URL - *http://localhost:5000/healthcheck
    Method - GET
    Description - check if server is running
    Response - If running:
            {
                "alive": true
            }
```

### List all Notebooks

```
    URL - *http://localhost:5000/listNotebooks
    Method - GET
    Description - Returns a list of all Notebook titles
    Response - List of Notebook Titles (ex. ["English","Math"])
```

### Create a Notebook

```
    URL - *http://localhost:5000/createNotebook/{notebookTitle}*
    Method - POST
    Description - Creates a new notebook with a given title
    Response - List of Notebook Titles (ex. ["English","Math"])
```

### Delete a Notebook

```
    URL - *http://localhost:5000/deleteNotebook/{notebookTitle}*
    Method - DELETE
    Description - Deletes a notebook with a given title
    Response - List of Notebook Titles (ex. ["English","Math"])
```

### Number of Notes in Notebook

```
    URL - *http://localhost:5000/numberOfNotes/{notebookTitle}*
    Method - GET
    Description - Get the number of notes in a notebook
    Response - Number of Notes in the Notebook (ex. 2)
```

### List Notes in Notebook

```
    URL - *http://localhost:12345/listNotes/{notebookTitle}*
    Method - GET
    Body - form-data
        {
            "Tags": []string
        }
    Description - List all notes in a notebook that match tags in body
    Response - List of valid notes (ex. [{"Id":"1","Title":"Hamlet","Body":"This is Hamlet","Tags":["Classics","Shakespeare"],"Created":"HamletCreated","LastModified":"HamletModified"}])
```

### Create Note in Notebook

```
    URL - *http://localhost:5000/createNote/{notebookTitle}*
    Method - POST
    Body - form-data
        {
            "Title": string,  // required
            "Body": string,   // required
            "Tags": string[], // required
        }
    Description - Create a note in a notebook
    Response - List of Notes in the notebook (ex. [{"Title": "Hamlet", "Body": "This is Hamlet", "Tags": ["Classics", "Shakespeare"], "Created": "HamletCreated", "LastModified": "HamletModified"}])
```

### Update a Note in a notebook

```
    URL - *http://localhost:5000/updateNote/{notebookTitle}/{noteId}*
    Method - UPDATE
    Body - form-data
        {
            "Title": string,  // required
            "Body": string,   // required
            "Tags": string[], // required
        }
    Description - Update a note in a notebook
    Response - List of Notes in the notebook (ex. [{"Title": "Hamlet", "Body": "This is Hamlet", "Tags": ["Classics", "Shakespeare"], "Created": "HamletCreated", "LastModified": "HamletModified"}])
```

### Read a Note in a Notebook

```
    URL - *http://localhost:5000/readNote/{notebookTitle}/{noteId}*
    Method - GET
    Description - Get a note (based on id) from a notebook
    Response - single note if found (ex. {"Id":"2","Title":"Animal Farm","Body":"Farm of Animals","Tags":["Classics"],"Created":"AnimalCreated","LastModified":"AnimalModified"})
```

### Delete a Note in a Notebook

```
    URL - *http://localhost:5000/deleteNote/{notebookTitle}/{noteId}*
    Method - DELETE
    Description - Delete a notes (with a specific id) in a notebook 
    Response - List of Notes in the notebook (ex. [{"Title": "Hamlet", "Body": "This is Hamlet", "Tags": ["Classics", "Shakespeare"], "Created": "HamletCreated", "LastModified": "HamletModified"}])

```

## Test Driven Development Description

To run all the unit test cases, please do the following:

1. `cd app`
2. `go get ./...`
3. `go test`


## Hope everything works. Thank you.