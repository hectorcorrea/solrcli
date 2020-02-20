package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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
		char := readKey()

		if char == "Q" || byte(char[0]) == ctrlC {
			fmt.Printf("Quit\r\n")
			break
		}

		if char == "" {
			continue
		}

		if char == "h" {
			showHelp()
			continue
		}

		if char == "q" {
			q = readLine("Enter q value", q)
		} else if char == "f" {
			facetField = readLine("Enter facet.field", facetField)
		} else if char == "l" {
			fl = readLine("Enter fl value", fl)
		} else if char == "s" {
			start = readLine("Enter start value", start)
		} else if char == "o" {
			rows = readLine("Enter rows value", rows)
		} else if char == "d" {
			debug = readLine("Enter debug value", debug)
		} else if char == "x" {
			s := solr.New(solrCoreURL, true)

			options := map[string]string{
				"defType": "edismax",
			}
			if debug == "true" {
				options["debug"] = "true"
			} else {
				options["debug"] = "false"
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
	fmt.Printf("=========================================================================================\n")
	fmt.Printf("[h]elp | [q]uery | [f]acet field | f[l] | e[x]ecute | [s]tart | r[o]ws | [d]ebug | [Q]uit\n")
	fmt.Printf("=========================================================================================\n")
}

func showHelp() {
	fmt.Printf("===============================================================================\n")
	fmt.Printf("[h] show this help screen\n")
	fmt.Printf("[q] enter the Solr q value\n")
	fmt.Printf("[#] blah blah blah \n")
	fmt.Printf("[#] blah blah blah \n")
	fmt.Printf("[#] blah blah blah \n")
	fmt.Printf("[#] blah blah blah \n")
	fmt.Printf("===============================================================================\r\n")
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

	// writer := bufio.NewWriter(os.Stdin)
	// writer.Write([]byte("hello"))
	// writer.Flush()
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if text == "\n" {
		// Use the original value
		return value
	}
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
