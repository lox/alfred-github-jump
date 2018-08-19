# go-alfred

Go-alfred is a utility library for quickly writing lightning fast [Alfred 2](http://www.alfredapp.com/) workflows using Golang.

## Example usage

```go
package main

import (
	"os"
	"github.com/pascalw/go-alfred"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func main() {
	queryTerms := os.Args[1:]

	// optimize query terms for fuzzy matching
	alfred.InitTerms(queryTerms)

	// create a new alfred workflow response
	response := alfred.NewResponse()
	repos := getRepos()

	for _, repo := range repos {
		// check if the repo name fuzzy matches the query terms
		if ! alfred.MatchesTerms(queryTerms, repo.Name) { continue }

		// it matched so add a new response item
		response.AddItem(&alfred.AlfredResponseItem{
			Valid: true,
			Uid: repo.URL,
			Title: repo.Name,
			Arg: repo.URL,
		})
	}

	// finally print the resulting Alfred Workflow XML
	response.Print()
}	
```

See [Example/](https://github.com/pascalw/go-alfred/blob/master/example/example.go) for details.
