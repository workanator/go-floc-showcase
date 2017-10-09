package main

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/workanator/go-floc.v2"
	"gopkg.in/workanator/go-floc.v2/run"
)

const exampleText = `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed
  do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad
  minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
  commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat
  non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

var sanitizeWordRe = regexp.MustCompile(`\W`)

func main() {
	const KeyStatistics = 1

	type Statistics struct {
		Words      []string
		Characters int
		Occurrence map[string]int
	}

	// Print introduction
	introduction()

	// Split to words and sanitize
	SplitToWords := func(ctx floc.Context, ctrl floc.Control) error {
		statistics := ctx.Value(KeyStatistics).(*Statistics)

		statistics.Words = strings.Split(exampleText, " ")
		for i, word := range statistics.Words {
			statistics.Words[i] = sanitizeWordRe.ReplaceAllString(word, "")
		}

		return nil
	}

	// Count and sum the number of characters in the each word
	CountCharacters := func(ctx floc.Context, ctrl floc.Control) error {
		statistics := ctx.Value(KeyStatistics).(*Statistics)

		for _, word := range statistics.Words {
			statistics.Characters += len(word)
		}

		return nil
	}

	// Count the number unique words
	CountUniqueWords := func(ctx floc.Context, ctrl floc.Control) error {
		statistics := ctx.Value(KeyStatistics).(*Statistics)

		statistics.Occurrence = make(map[string]int)
		for _, word := range statistics.Words {
			statistics.Occurrence[word] = statistics.Occurrence[word] + 1
		}

		return nil
	}

	// Print result
	PrintResult := func(ctx floc.Context, ctrl floc.Control) error {
		statistics := ctx.Value(KeyStatistics).(*Statistics)

		fmt.Println("Text statistics")
		fmt.Println("---------------------------")
		fmt.Printf("Words Total               : %d\n", len(statistics.Words))
		fmt.Printf("Unique Words              : %d\n", len(statistics.Occurrence))
		fmt.Printf("Non-Whitespace Characters : %d\n", statistics.Characters)

		return nil
	}

	// Design the flow and run it
	flow := run.Sequence(
		SplitToWords,
		run.Parallel(
			CountCharacters,
			CountUniqueWords,
		),
		PrintResult,
	)

	ctx := floc.NewContext()
	ctx.AddValue(KeyStatistics, new(Statistics))

	ctrl := floc.NewControl(ctx)

	floc.RunWith(ctx, ctrl, flow)
}

func introduction() {
	fmt.Println(`The example Words runs a sequence of jobs:
    1. Split a text into a list of words.
    2. Run in parallel:
    2.1. Count non-whitespace characters in the each word and sum them.
    2.2. Count unique words in the text.
    3. Print a result.`)
}
