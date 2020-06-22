package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hectorcorrea/solr"
)

func main() {
	if len(os.Args) < 2 {
		showSyntax()
		return
	}

	solrCoreURL := os.Args[1]
	fmt.Printf("solrcli\r\n")
	fmt.Printf("%s\r\n", solrCoreURL)

	repl(solrCoreURL)
	return
}

func showSyntax() {
	syntax := `
solrcli url-to-solr-core
e.g. solrcli http://localhost:8983/solr/your-core

PARAMETERS
	url-to-solr-core  	Full http URL to the Solr core to use
`

	fmt.Printf("%s\n\n", syntax)
}

func toInt(str string) int {
	num, _ := strconv.ParseInt(str, 10, 64)
	return int(num)
}

func toJSON(raw string) string {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, []byte(raw), "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(pretty.Bytes())
}

// Get the Solr schema (via the Luke Request Handler)
func getSchema(lukeURL string) (string, error) {
	r, err := http.Get(lukeURL)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Executes the Solr query with the current values
func executeQuery(solrCoreURL string, userParams map[string]string) {
	q := "*:*"
	facets := map[string]string{}
	options := map[string]string{
		"defType": "edismax",
	}

	for key, value := range userParams {
		if key == "facet.field" {
			facets[value] = value
		} else if key == "start" || key == "rows" || key == "fl" {
			// these are process individually (see below)
		} else if key == "q" {
			q = value
		} else {
			options[key] = value
		}
	}

	s := solr.New(solrCoreURL, true)
	params := solr.NewSearchParams(q, options, facets)
	if userParams["fl"] != "" {
		params.Fl = strings.Split(userParams["fl"], ",")
	}
	if userParams["start"] != "" {
		params.Start = toInt(userParams["start"])
	}
	if userParams["rows"] != "" {
		params.Rows = toInt(userParams["rows"])
	}

	results, err := s.Search(params)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	} else {
		fmt.Printf("%s\r\n", toJSON(results.Raw))
	}
}

// Shows the current values to send to Solr
func showValues(solrCoreURL string, userParams map[string]string) {
	fmt.Printf("Solr URL\n")
	fmt.Printf("  %s\n", solrCoreURL)
	fmt.Printf("\n")
	fmt.Printf("Solr values\n")
	for key, value := range userParams {
		fmt.Printf("\t%s\t\t= %s\n", key, value)
	}
}
