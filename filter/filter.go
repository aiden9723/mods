package filter

import (
	"fmt"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/yokimpillay/callcoster/call"
	"github.com/yokimpillay/callcoster/report"
)

// ValidCalls iterates through the valid calls and attempts
// to match them to an appropriate country/area code.
func ValidCalls(calls, rates [][]string) {
	// rateRows = rateRows[1:]
	southAfricanCodes := CountryCodes(rates, "South Africa")
	internationalCodes := InternationalCodes(rates)

	fmt.Println("Matching numbers to codes & rates")
	bar := pb.New(len(calls))
	bar.SetWidth(80)
	bar.Start()
	for i := 0; i < len(calls); i++ {

		// Check if the current row is an international call
		if calls[i][4] == "International" {

			// If true, loop through internationalCodes and try
			// find matching country
			for _, rate := range internationalCodes {
				if len(calls[i][1]) <= 11 {
					callCost := call.GetCost(calls[i][2], "0.80")
					calls[i] = append(calls[i], callCost)
					calls[i] = append(calls[i], "Other")
					break
				}

				matched, _, codeType := report.MatchCode(calls[i][1], rate[0], rate[2])

				// If a match is found, calculate the call cost and
				// add to row
				if matched {
					callCost := call.GetCost(calls[i][2], rate[3])
					calls[i] = append(calls[i], callCost)
					calls[i] = append(calls[i], codeType)
					// fmt.Printf("Matched %s to %s\n", calls[i][1], codeType)
					break
				}
			}
		} else {
			// If the number is not international,
			// loop through SA codes and find a
			// matching code & rate.
			for _, rate := range southAfricanCodes {
				matched, _, codeType := report.MatchCode(calls[i][1], rate[0], rate[2])

				// If a match is found, calculate the
				// call cost and add to the row
				if matched {
					callCost := call.GetCost(calls[i][2], rate[3])
					calls[i] = append(calls[i], callCost)
					calls[i] = append(calls[i], codeType)
					// fmt.Printf("Matched %s to %s\n", calls[i][1], codeType)
					break
				}
				if rate[0] == southAfricanCodes[len(southAfricanCodes)-1][0] {
					callCost := call.GetCost(calls[i][2], "0.80")
					calls[i] = append(calls[i], callCost)
					calls[i] = append(calls[i], "Other")
					// errorString := fmt.Sprintf("No known local code for %v", calls[i])
					// panic(errorString)
				}

			}

		}

		bar.Increment()
	}
}

func CountryCodes(rates [][]string, country string) [][]string {
	return filter(rates, func(code []string) bool {
		return strings.Contains(code[1], country)
	})
}

func InternationalCodes(rates [][]string) [][]string {
	return filter(rates, func(code []string) bool {
		return !strings.Contains(code[1], "South Africa")
	})
}

// Filter is the helper method to abstract iteration through a
// slice, and appends rows that evaluate as truthy by the callback
// function.
func filter(rows [][]string, callback func([]string) bool) [][]string {
	filteredRows := make([][]string, 0)

	for _, row := range rows {
		if callback(row) {
			filteredRows = append(filteredRows, row)
		}
	}

	return filteredRows
}
