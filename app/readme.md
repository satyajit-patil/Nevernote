# Nevernote API

## How to run the API Server

1. Clone the application with ``

2. `cd app`

3. Build the docker container: 
`docker build --tag nevernote:{VERSION} .`

4. Run the docker container with port 5000 port forwarded:
`docker run -itd -p 5000:5000 nevernote:{VERSION}`

5. Check that the api server is running by running:
`curl localhost:5000/healthcheck`


## Endpoints Description

### Healthcheck

```
    URL - *http://localhost:5000/healthcheck
    Method - GET
    Description - check if server is running
```

### List all Notebooks

```
    URL - *http://localhost:5000/listNotebooks
    Method - GET
    Description - Returns a list of all Notebook titles
```

### Create a Notebook

```JSON
    URL - *http://localhost:5000/createNotebook/{title}*
    Method - POST
    Description - Creates a new notebook with a given title
```

### Delete a Notebook

```JSON
    URL - *http://localhost:5000/deleteNotebook*
    Method - DELETE
    Description - Deletes a notebook with a given title
```

### Number of Notes in Notebook

```JSON
    URL - *http://localhost:5000/numberOfNotes/{title}*
    Method - GET
    Description - Get the number of notes in a notebook
```

### List Notes in Notebook

```JSON
    URL - *http://localhost:12345/listNotes/{title}*
    Method - GET
    Body - form-data
        {
            "Tags": []string
        }
    Description - List all notes in a notebook that match tags in body
```

### Create Note in Notebook

```JSON
    URL - *http://localhost:5000/createNote/{title}*
    Method - POST
    Body - form-data
        {
            "Title": string,  // required
            "Body": string,   // required
            "Tags": string[], // required
        }
    Description - Create a note in a notebook
```

### Update a Note in a notebook

```JSON
    URL - *http://localhost:5000/updateNote/{notebookTitle}/{noteId}*
    Method - UPDATE
    Body - form-data
        {
            "Title": string,  // required
            "Body": string,   // required
            "Tags": string[], // required
        }
    Description - Update a note in a notebook
```

### Read a Note in a Notebook

```JSON
    URL - *http://localhost:5000/readNote/{title}/{noteId}*
    Method - GET
    Description - Get a note (based on id) from a notebook
```

### Delete a Note in a Notebook

```JSON
    URL - *http://localhost:5000/deleteNote/{title}/{noteId}*
    Method - DELETE
    Description - Delete a notes (with a specific id) in a notebook 
```

## Test Driven Development Description

To run all the unit test cases, please do the following -

1. `cd app`
2. `go get ./...`
3. `go test`


## Hope everything works. Thank you.