package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/mgutz/ansi"
)

var (
	run         = regexp.MustCompile(`^=== RUN Test.*$`)
	pass        = regexp.MustCompile(`^--- PASS: Test.*$`)
	fail        = regexp.MustCompile(`^--- FAIL: Test.*$`)
	passSummary = regexp.MustCompile(`^PASS$`)
	failSummary = regexp.MustCompile(`^FAIL$`)
	passSummLn  = regexp.MustCompile(`^ok\s+[\w_\-\.\/]+\s+[0-9\.]+s$`)
	failSummLn  = regexp.MustCompile(`^FAIL\s+[\w_\-\.\/]+\s+([0-9\.]+s|\[build failed\])$`)
	skip        = regexp.MustCompile(`^--- SKIP: Test.*$`)
skipSummLn  = regexp.MustCompile(`^\?\s+[\w_\-\.\/]+\s+\[no test files\]$`)
	errLn       = regexp.MustCompile(`^(\s+\w+\.go):([0-9]+):( .*)$`)
)

func main() {
	scan := bufio.NewScanner(os.Stdin)

	for scan.Scan() {
		if err := scan.Err(); err != nil {
			fmt.Println("Error scanning:", err)
		}

		reset := ansi.ColorCode("reset")

		line := scan.Bytes()
		switch {
		case run.Match(line):
			fmt.Printf("%s%s%s\n", ansi.ColorCode("blue"), line, reset)
		case pass.Match(line), passSummary.Match(line), passSummLn.Match(line):
			fmt.Printf("%s%s%s\n", ansi.ColorCode("green+b"), line, reset)
		case fail.Match(line), failSummary.Match(line), failSummLn.Match(line):
			fmt.Printf("%s%s%s\n", ansi.ColorCode("red+b"), line, reset)
		case skip.Match(line), skipSummLn.Match(line):
			fmt.Printf("%s%s%s\n", ansi.ColorCode("yellow+b"), line, reset)
		case errLn.Match(line):
			submatches := errLn.FindSubmatch(line)
			fmt.Printf("%s%s%s%s:%s:%s%s%s%s\n",
				ansi.ColorCode("red+b"),
				submatches[1],
				reset,
				ansi.ColorCode("yellow+b"),
				submatches[2],
				reset,
				ansi.ColorCode("red+b"),
				submatches[3],
				reset,
			)
		default:
			fmt.Printf("%s\n", line)
		}
	}
}
