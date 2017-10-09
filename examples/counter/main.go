package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"gopkg.in/workanator/go-floc.v1"
	"gopkg.in/workanator/go-floc.v1/guard"
	"gopkg.in/workanator/go-floc.v1/run"
)

func main() {
	const TimeoutValue = 5 * time.Second
	const MaxRandom = 10000

	type Counter struct {
		Value           uint64 // Current value of the counter
		RandomsMet      uint64 // The number of time random value was met
		NextRandomValue uint64 // The next random value if reached increases RandomsMet
	}

	// Print introduction
	introduction()

	// Increment counter
	IncrementValue := func(flow floc.Flow, state floc.State, update floc.Update) {
		counter := state.Data().(*Counter)
		counter.Value++
	}

	// Test if Value equals to NextRandomValue
	NextRandomValueMet := func(state floc.State) bool {
		counter := state.Data().(*Counter)
		return counter.Value >= counter.NextRandomValue
	}

	// Increment RandomsMet and generate NextRandomValue
	IncrementRandomsMet := func(flow floc.Flow, state floc.State, update floc.Update) {
		counter := state.Data().(*Counter)
		counter.RandomsMet++
		counter.NextRandomValue = counter.Value + uint64(rand.Int63n(MaxRandom))
	}

	// Wait for SIGINT OS signal and cancel the flow
	WaithInterrupt := func(flow floc.Flow, state floc.State, update floc.Update) {
		c := make(chan os.Signal, 1)
		defer close(c)

		signal.Notify(c, os.Interrupt)

		// Wait for OS signal or flow finished signal
		select {
		case s := <-c:
			// OS signal was caught
			flow.Cancel(s)

		case <-flow.Done():
			// The flow is finished
		}
	}

	// Design the job and run it
	job := guard.TimeoutWithTrigger(
		guard.ConstTimeout(TimeoutValue),
		nil, // No need for timeout ID
		run.Parallel(
			WaithInterrupt,
			run.Loop(
				IncrementValue,
				run.If(NextRandomValueMet, IncrementRandomsMet),
			),
		),
		func(flow floc.Flow, state floc.State, id interface{}) {
			// Complete the flow on timeout
			flow.Complete(nil)
		},
	)

	flow := floc.NewFlow()
	state := floc.NewState(new(Counter))
	floc.Run(flow, state, nil, job)

	// Examine and print results
	result, data := flow.Result()
	switch {
	case result.IsCanceled():
		fmt.Printf("The flow was canceled by user with reason '%v'.\n", data)

	case result.IsCompleted():
		fmt.Println("The flow was completed successfully.")

		counter := state.Data().(*Counter)
		fmt.Println("-------------------")
		fmt.Printf("Counter Value     : %d\n", counter.Value)
		fmt.Printf("Randoms Met       : %d\n", counter.RandomsMet)
		fmt.Printf("Next Random Value : %d\n", counter.NextRandomValue)

	default:
		panic("The flow has finished with improper state")
	}
}

func introduction() {
	fmt.Println(`The example Counter has the flow:
    1. Protect the flow with timeout in 5s and on timeout Complete the flow.
    1.1. Run in parallel.
    1.1.1. Wait until OS Interrupt signal is caught and Cancel the flow.
    1.1.2. Run in infinite loop.
    1.1.2.1. Increment the counter value.
    1.1.2.2. If the counter value equals or greater the next random value then 1.1.2.2.T.
    1.1.2.2.T. Increment the number of randoms met and generate the next random value.
`)

	fmt.Print("Please wait for 5 seconds until the flow is finished or interrupt it with Ctrl+C.\n\n")
}
