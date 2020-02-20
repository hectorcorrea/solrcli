package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hectorcorrea/solr"

	"golang.org/x/crypto/ssh/terminal"
)

const CtrlC byte = 3

var q string
var facetField string
var start string
var rows string
var fl string
var qf string

func main() {
	fmt.Printf("solrcli\r\n")

	if len(os.Args) < 2 {
		showSyntax()
		return
	}

	solrCoreUrl := os.Args[1]
	fmt.Printf("%s\r\n", solrCoreUrl)
	showHelp()

	for true {
		char := readKey()
		if char == "Q" || byte(char[0]) == CtrlC {
			fmt.Printf("Quit\r\n")
			break
		}
		if char == "" {
			continue
		}
		if char == "q" {
			q = readLine("Enter q value", q)
		} else if char == "f" {
			facetField = readLine("Enter facet.field", facetField)
		} else if char == "l" {
			fl = readLine("Enter fl value", fl)
		} else if char == "x" {
			s := solr.New(solrCoreUrl, true)
			options := map[string]string{
				"defType": "edismax",
			}
			facets := map[string]string{}
			if facetField != "" {
				facets = map[string]string{facetField: facetField}
			}
			params := solr.NewSearchParams(q, options, facets)
			if fl != "" {
				params.Fl = strings.Split(fl, ",")
			}
			results, err := s.Search(params)
			if err != nil {
				fmt.Printf("ERROR: %s", err)
			}
			fmt.Printf("Documents found: %d", results.NumFound)
			for _, doc := range results.Documents {
				json := toJSON(doc)
				fmt.Printf("%s\r\n", json)
			}
		}
		showHelp()
	}
	return
}

func showSyntax() {
	fmt.Printf("solrcli url-to-solr-core\n")
	fmt.Printf("e.g. solrcli http://localhost:8983/solr/your-core\n")

}
func showHelp() {
	fmt.Printf("===============================================================================\n")
	fmt.Printf("[h]elp | [q]uery | [f]acet field | f[l] | e[x]ecute | [s]tart | r[o]ws | [Q]uit\n")
	fmt.Printf("===============================================================================\r\n")
}

func readLine(prompt string, value string) string {
	if value == "" {
		fmt.Printf("%s: ", prompt)
	} else {
		fmt.Printf("%s (%s): ", prompt, value)
	}

	// writer := bufio.NewWriter(os.Stdin)
	// writer.Write([]byte("hello"))
	// writer.Flush()
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// func readChar() string {
// 	reader := bufio.NewReader(os.Stdin)
// 	text, _ := reader.ReadString('\n')
// 	return strings.TrimSpace(string(text[0]))
// }

// Source https://stackoverflow.com/a/17278730/446681
func readKey() string {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	return string(b)
}

func toJSON(data interface{}) string {
	raw, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, raw, "", "\t")
	if err != nil {
		return ""
	}
	return string(pretty.Bytes())
}