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

var q string
var facetField string
var start string
var rows string
var fl string
var qf string
var debug string = "false"

func main() {
	if len(os.Args) < 2 {
		showSyntax()
		return
	}

	solrCoreURL := os.Args[1]
	fmt.Printf("solrcli\r\n")
	fmt.Printf("%s\r\n", solrCoreURL)

	if len(os.Args) == 3 && os.Args[2] == "repl" {
		repl(solrCoreURL)
		return
	}

	menu(solrCoreURL)
	return
}

func showSyntax() {
	syntax := `
solrcli url-to-solr-core [menu|repl]
e.g. solrcli http://localhost:8983/solr/your-core repl
	
PARAMETERS
	url-to-solr-core  	Full http URL to the Solr core to use

	menu 			to use Menu mode in which parameters are set via single keystroke
				options (e.g. q for q) and so are commands (e.g. r for run). 
				This is the default mode.

	repl			to use REPL mode in which parameters are set via assignments
				(e.g. q=hello) and commands are typed directly (e.g. run)`

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
func executeQuery(solrCoreURL string) {
	s := solr.New(solrCoreURL, true)

	options := map[string]string{
		"defType": "edismax",
	}
	if debug == "true" {
		options["debug"] = "true"
	} else {
		options["debug"] = "false"
	}
	if qf != "" {
		options["qf"] = qf
	}

	facets := map[string]string{}
	if facetField != "" {
		facets[facetField] = facetField
	}

	params := solr.NewSearchParams(q, options, facets)
	if fl != "" {
		params.Fl = strings.Split(fl, ",")
	}
	if start != "" {
		params.Start = toInt(start)
	}
	if rows != "" {
		params.Rows = toInt(rows)
	}
	results, err := s.Search(params)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	} else {
		fmt.Printf("%s\r\n", toJSON(results.Raw))
	}
}

// Shows the current values to send to Solr
func showValues() {
	fmt.Printf("Solr values\n")
	fmt.Printf("  q           = %s\n", q)
	fmt.Printf("  fl          = %s\n", fl)
	fmt.Printf("\n")
	fmt.Printf("  facet.field = %s\n", facetField)
	fmt.Printf("  debug       = %s\n", debug)
	fmt.Printf("  defType     = %s\n", "edismax")
	fmt.Printf("  qf          = %s\n", qf)
	fmt.Printf("  rows        = %s\n", rows)
	fmt.Printf("  start       = %s\n", start)
}
