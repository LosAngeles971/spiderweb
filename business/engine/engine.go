/*
This package exports the main functions of Joshua to the CLI and/or external programs.

The main function is "IsItGoingToHappen"::
Given an event, named success event, the function verifies it may happen starting from an initial state.
*/
package engine

import (
	"fmt"
	"it/losangeles971/joshua/business/parser"
	"it/losangeles971/joshua/business/knowledge"
	"os"
	"text/tabwriter"
)

// PrintSummary prints the results of an execution to the console (standard output)
func PrintSummary(outcome string, queue knowledge.Queue) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Outcome\t" + outcome + "\t")
	fmt.Fprintln(w, "Cycles\t" + fmt.Sprint(queue.GetCycles()) + "\t")
	fmt.Fprintln(w, "Queue's size\t" + fmt.Sprint(queue.Size()) + "\t")
	w.Flush()
}

// Load compile a knowledge's source into a knowledge object
func load(source string) (knowledge.Knowledge, error) {
	p := parser.NewParser(source)
	events, err := p.Parse()
	if err != nil {
		return knowledge.Knowledge{}, err
	}
	return knowledge.NewKnowledge(events)
}

// IsItGoingToHappen exports the internal function isItGoingToHappen
func IsItGoingToHappen(source string, init knowledge.State, effect string, max_cycles int) (string, knowledge.Queue, error) {
	kkk, err := load(source)
	if err != nil {
		return "", knowledge.Queue{}, err
	}
	ee, ok := kkk.GetEvent(effect)
	if !ok {
		return "", knowledge.Queue{}, fmt.Errorf("success event %v does not exist into the knowledge", effect)
	}
	return knowledge.IsItGoingToHappen(kkk, init, ee, max_cycles)
}