package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hectorcorrea/solr"
	"golang.org/x/crypto/ssh/terminal"
)

func menu(solrCoreURL string) {

	const ctrlC byte = 3

	for true {
		showMenu()
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

		if option == "w" {
			showValues()
			continue
		}

		if option == "q" {
			q = readValue("Enter q value", q)
		} else if option == "f" {
			facetField = readValue("Enter facet.field", facetField)
		} else if option == "F" {
			qf = readValue("Enter qf", qf)
		} else if option == "l" {
			fl = readValue("Enter fl value", fl)
		} else if option == "s" {
			start = readValue("Enter start value", start)
		} else if option == "o" {
			rows = readValue("Enter rows value", rows)
		} else if option == "d" {
			debug = readValue("Enter debug value", debug)
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
		} else {
			fmt.Printf("Unknown option %s\n", option)
		}
	}
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
	fmt.Printf("\t[w] Show Solr values\n")
	fmt.Printf("\t[x] Execute the query with the current values\n")
	fmt.Printf("\n")
	fmt.Printf("\t[Q] Quit (also CTRL+C)\n")
	fmt.Printf("==============================================================================\n")
	fmt.Printf("\n")
}

func readValue(prompt string, value string) string {
	if value == "" {
		fmt.Printf("%s: ", prompt)
	} else {
		fmt.Printf("%s (%s): ", prompt, value)
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
