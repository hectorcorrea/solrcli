package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hectorcorrea/solr"

	"golang.org/x/crypto/ssh/terminal"
)

const ctrlC byte = 3

var q string
var facetField string
var start string
var rows string
var fl string
var qf string
var debug string = "false"

func main() {
	fmt.Printf("solrcli\r\n")

	if len(os.Args) < 2 {
		showSyntax()
		return
	}

	solrCoreURL := os.Args[1]
	fmt.Printf("%s\r\n", solrCoreURL)
	showMenu()

	for true {
		fmt.Printf("> ")
		option := readKey()

		if option == "Q" || byte(option[0]) == ctrlC {
			fmt.Printf("Quit\r\n")
			break
		}

		if option == "" {
			continue
		}

		fmt.Printf("%s\n", option)

		if option == "h" {
			showHelp()
			continue
		}

		if option == "q" {
			q = readLine("Enter q value", q)
		} else if option == "f" {
			facetField = readLine("Enter facet.field", facetField)
		} else if option == "F" {
			qf = readLine("Enter qf", qf)
		} else if option == "l" {
			fl = readLine("Enter fl value", fl)
		} else if option == "s" {
			start = readLine("Enter start value", start)
		} else if option == "o" {
			rows = readLine("Enter rows value", rows)
		} else if option == "d" {
			debug = readLine("Enter debug value", debug)
		} else if option == "c" {
			lukeURL := fmt.Sprintf("%s/admin/luke", solrCoreURL)
			schema, err := getSchema(lukeURL)
			if err != nil {
				fmt.Printf("ERROR: %s", err)
			} else {
				fmt.Printf("%s", schema)
			}
		} else if option == "x" {
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

		showMenu()
	}
	return
}

func showSyntax() {
	fmt.Printf("solrcli url-to-solr-core\n")
	fmt.Printf("e.g. solrcli http://localhost:8983/solr/your-core\n")
}

func showMenu() {
	//                   1         2         3         4         5         6         7         8
	//          123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-
	fmt.Printf("==============================================================================\n")
	fmt.Printf("[h]elp | [q]uery | [f]acet | f[l] | e[x]ecute | [s]tart | r[o]ws | [Q]uit\n")
	fmt.Printf("==============================================================================\n")
}

func showHelp() {
	//                   1         2         3         4         5         6         7         8
	//          123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-
	fmt.Printf("==============================================================================\n")
	fmt.Printf("Options available\n")
	fmt.Printf("\n")
	fmt.Printf("\t[c] Show Solr core's schema\n")
	fmt.Printf("\t[d] Enter debug value\n")
	fmt.Printf("\t[f] Enter facet.field value\n")
	fmt.Printf("\t[F] Enter qf value\n")
	fmt.Printf("\t[h] Show this help screen\n")
	fmt.Printf("\t[l] Enter fl value\n")
	fmt.Printf("\t[o] Enter rows value\n")
	fmt.Printf("\t[q] Enter q value\n")
	fmt.Printf("\t[s] Enter start value\n")
	fmt.Printf("\t[x] Execute the query with the current values\n")
	fmt.Printf("\n")
	fmt.Printf("\t[Q] Quit (also CTRL+C)\n")
	fmt.Printf("==============================================================================\n")
	fmt.Printf("\n")
}

func readLine(prompt string, value string) string {
	if value == "" {
		fmt.Printf("%s: ", prompt)
	} else {
		fmt.Printf("%s (%s): ", prompt, value)
		// https://stackoverflow.com/a/33509850/446681
		// fmt.Printf("\033[0;0H")
		// fmt.Printf("\033[3D")
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if text == "\n" {
		// Use the original value
		return value
	}
	return strings.TrimSpace(text)
}

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
