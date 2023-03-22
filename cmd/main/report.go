package main

import (
	"fmt"
	"strings"
)

type reportEntry struct {
	CloudPlatform string
	Success       int
	Failed        int
}

func (s *reportEntry) IncrementSuccess() {
	s.Success++
}

func (s *reportEntry) IncrementFailed() {
	s.Failed++
}

func printReport(report []reportEntry) {
	//Print a summary report of entries imported and not imported

	var lineLength int //keep the longest length of each report line.
	var lines []string

	for _, y := range report {
		var line string
		line = fmt.Sprintf("%s, Success: %v, Failed: %v\n", y.CloudPlatform, y.Success, y.Failed)
		if len(line) > lineLength {
			lineLength = len(line)
		}
		lines = append(lines, line)
	}

	fmt.Printf("%s\n", strings.Repeat("-", lineLength))
	for _, y := range lines {
		fmt.Printf(y)
	}
	fmt.Printf("%s\n", strings.Repeat("-", lineLength))

}
