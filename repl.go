package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
)

func repl(solrCoreURL string) {
	userParams := map[string]string{}

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m>\033[0m ",
		HistoryFile:     "/tmp/solrcli.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	showBarRepl()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		key, value := evalLine(line)
		if key == "quit" || key == "exit" {
			break
		}

		switch {
		case key == "help":
			showHelpRepl()
		case key == "show":
			showValues(solrCoreURL, userParams)
		case key == "run":
			executeQuery(solrCoreURL, userParams)
		case key == "schema":
			lukeURL := fmt.Sprintf("%s/admin/luke", solrCoreURL)
			schema, err := getSchema(lukeURL)
			if err != nil {
				fmt.Printf("ERROR: %s", err)
			} else {
				fmt.Printf("%s", schema)
			}
		default:
			userParams[key] = value
		}
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func evalLine(line string) (string, string) {
	// If it is in the form q=x return two tokens: q and x
	if strings.Contains(line, "=") {
		tokens := strings.Split(line, "=")
		if len(tokens) == 2 {
			key := strings.TrimSpace(tokens[0])
			value := strings.TrimSpace(tokens[1])
			return key, value
		}
		return "", ""
	}
	// Assume it's a command, e.g. "run"
	return line, ""
}

func showBarRepl() {
	//                   1         2         3         4         5         6         7         8
	//          123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-
	fmt.Printf("run  - execute query\n")
	fmt.Printf("show - show current values for query\n")
	fmt.Printf("help - view available options\n")
	fmt.Printf("quit - to quit (also CTRL+C)\n")
	fmt.Printf("==============================================================================\n")
}

func showHelpRepl() {
	//                   1         2         3         4         5         6         7         8
	//          123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-
	fmt.Printf("==============================================================================\n")
	fmt.Printf("Commands available\n")
	fmt.Printf("\tschema               - Show Solr core's schema\n")
	fmt.Printf("\thelp                 - Shows  help screen\n")
	fmt.Printf("\trun                  - Runs the query with the current values\n")
	fmt.Printf("\tshow                 - Shows values to send to Solr\n")
	fmt.Printf("\tquit                 - Quits (also CTRL+C)\n")
	fmt.Printf("\n")
	fmt.Printf("Options available\n")
	fmt.Printf("\tdebug = true|false   - Sets debug value\n")
	fmt.Printf("\tfacet = true|false   - Sets the facet value\n")
	fmt.Printf("\tfacet.field = value  - Sets the facet.field value\n")
	fmt.Printf("\tfl = value           - Sets the fl value\n")
	fmt.Printf("\trows = value         - Sets the rows value\n")
	fmt.Printf("\tsort = value         - Sets the sorder order value\n")
	fmt.Printf("\tstart = value        - Sets the start value\n")
	fmt.Printf("\tq = value            - Sets the q value\n")
	fmt.Printf("\tqf = value           - Sets the qf value\n")
	fmt.Printf("\n")
	fmt.Printf("==============================================================================\n")
	fmt.Printf("\n")
}
