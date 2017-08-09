package main

import (
	"fmt"
	"regexp"
	"strings"

	floc "github.com/workanator/go-floc"
	"github.com/workanator/go-floc/run"
)

const exampleText = `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed
  do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad
  minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
  commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat
  non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

var sanitizeWordRe = regexp.MustCompile(`\W`)

func main() {
	type Statistics struct {
		Words      []string
		Characters int
		Occurence  map[string]int
	}

	// Print introduction
	introduction()

	// Split to words and sanitize
	SplitToWords := func(flow floc.Flow, state floc.State, update floc.Update) {
		statistics := state.Data().(*Statistics)

		statistics.Words = strings.Split(exampleText, " ")
		for i, word := range statistics.Words {
			statistics.Words[i] = sanitizeWordRe.ReplaceAllString(word, "")
		}
	}

	// Count and sum the number of characters in the each word
	CountCharacters := func(flow floc.Flow, state floc.State, update floc.Update) {
		statistics := state.Data().(*Statistics)

		for _, word := range statistics.Words {
			statistics.Characters += len(word)
		}
	}

	// Count the number unique words
	CountUniqueWords := func(flow floc.Flow, state floc.State, update floc.Update) {
		statistics := state.Data().(*Statistics)

		statistics.Occurence = make(map[string]int)
		for _, word := range statistics.Words {
			statistics.Occurence[word] = statistics.Occurence[word] + 1
		}
	}

	// Print result
	PrintResult := func(flow floc.Flow, state floc.State, update floc.Update) {
		statistics := state.Data().(*Statistics)

		fmt.Println("Text statistics")
		fmt.Println("---------------------------")
		fmt.Printf("Words Total               : %d\n", len(statistics.Words))
		fmt.Printf("Unique Words              : %d\n", len(statistics.Occurence))
		fmt.Printf("Non-Whitespace Characters : %d\n", statistics.Characters)
	}

	// Design the job and run it
	job := run.Sequence(
		SplitToWords,
		run.Parallel(
			CountCharacters,
			CountUniqueWords,
		),
		PrintResult,
	)

	floc.Run(
		floc.NewFlow(),
		floc.NewState(new(Statistics)),
		nil,
		job,
	)
}

func introduction() {
	fmt.Println(`The example Words runs a sequence of jobs:
    1. Split a text into a list of words.
    2. Run in parallel:
    2.1. Count non-whitespace characters in the each word and sum them.
    2.2. Count unique words in the text.
    3. Print a result.
`)
}
