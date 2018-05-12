package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mgutz/ansi"
)

var (
	errLn = regexp.MustCompile(`^(\s+\w+\.go):([0-9]+):( .*)$`)

	red    = ansi.ColorCode("red+b")
	green  = ansi.ColorCode("green+b")
	yellow = ansi.ColorCode("yellow+b")
	reset  = ansi.Reset

	throwAwayPrefixes = []string{`--- PASS:`, `=== RUN`, `--- FAIL:`, `=== CONT`, `=== PAUSE`}
)

// TestEvent as defined by the testing docs.
type TestEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

func main() {
	scan := bufio.NewScanner(os.Stdin)

	output := make(map[string][]string)

	for scan.Scan() {
		te := TestEvent{}
		err := json.Unmarshal(scan.Bytes(), &te)
		if err != nil {
			panic("failed to decode json: " + err.Error())
		}

		if len(te.Test) == 0 {
			switch te.Action {
			case "fail":
				fmt.Printf("%s%s%s\n", red, "FAIL: "+te.Package, reset)
			case "pass":
				fmt.Printf("%s%s%s\n", green, "PASS: "+te.Package, reset)
			}
			continue
		}

		switch te.Action {
		case "fail":
			out := output[te.Test]
			fmt.Printf("%s%s%s\n", red, "FAIL: "+te.Test, reset)
			for _, line := range out {
				fmt.Printf("%s%s%s\n", red, line, reset)
			}
			delete(output, te.Test)
		case "pass":
			out := output[te.Test]
			fmt.Printf("%s%s%s\n", green, "PASS: "+te.Test, reset)
			for _, line := range out {
				fmt.Printf("%s%s%s\n", green, line, reset)
			}
			delete(output, te.Test)
		case "skip":
			out := output[te.Test]
			fmt.Printf("%s%s%s\n", yellow, "SKIP: "+te.Test, reset)
			for _, line := range out {
				fmt.Printf("%s%s%s\n", yellow, line, reset)
			}
			delete(output, te.Test)
		case "output":
			noSpace := strings.TrimSpace(te.Output)
			throwAway := false

			for _, pfx := range throwAwayPrefixes {
				if strings.HasPrefix(noSpace, pfx) {
					throwAway = true
					break
				}
			}

			if !throwAway {
				output[te.Test] = append(output[te.Test], strings.TrimRight(te.Output, "\r\n"))
			}
		}
	}

	if err := scan.Err(); err != nil {
		fmt.Println("Error scanning:", err)
	}
}
