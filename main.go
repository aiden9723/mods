package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Row is a container for csv columns
type Row struct {
	UniqueID    string
	Destination string
	BillSec     string
	Tag         string
	Prefix      string
}

var bar *pb.ProgressBar

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	report := flag.String("report", "", "The file to calculate call costs on")
	rates := flag.String("rates", "", "The file from which costs will be used")

	flag.Parse()

	reportLines, err := readCSV(*report)
	check(err)

	modifyColumns(reportLines)

	reportHeaders, reportLines := reportLines[0], reportLines[1:]

	// bar = pb.New(len(reportLines))
	// bar.SetWidth(80)
	// bar.Start()

	rateLines, err := readCSV(*rates)
	check(err)

	validnums, _ := validateNumbers(reportLines)

	filterValidCalls(validnums, rateLines)

	reportLines = append([][]string{reportHeaders}, validnums...)

	writeChanges(reportLines, *report)

	// bar.Finish()
}

func readCSV(filename string) ([][]string, error) {
	if filename == "" {
		fmt.Println("File not specified, exiting.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("Reading from", filename)
	f, err := os.Open(filename)
	check(err)

	lines, err := csv.NewReader(f).ReadAll()

	f.Close()

	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func writeChanges(lines [][]string, filename string) {
	newFilename := makeFilename(filename)

	f, err := os.Create(newFilename)
	check(err)

	err = csv.NewWriter(f).WriteAll(lines)
	check(err)
	fmt.Println("\nSuccessfully written to", newFilename)
	return
}

func makeFilename(filename string) string {
	trimmedFilename := strings.Replace(filename, ".csv", "", -1)
	newFilename := fmt.Sprintf("%s-formatted.csv", trimmedFilename)
	return newFilename
}

func matchCode(number, code, codeType string) (bool, string, string) {
	// fmt.Printf("Trying to match %s with %s\n", number, code)
	formattedPattern := fmt.Sprintf(`^(%s)(?:\d)+$`, code)
	re, err := regexp.Compile(formattedPattern)
	check(err)

	matcher := re.Match([]byte(number))

	return matcher, code, codeType
}

func modifyColumns(rows [][]string) {
	rows[0] = append(rows[0], "Cost")
	rows[0] = append(rows[0], "Code")
}

func calculateCostPerSecond(cost string) float64 {
	floatCost, err := strconv.ParseFloat(cost, 64)
	check(err)

	return floatCost / 60.0
}

func validateNumbers(reportRows [][]string) ([][]string, [][]string) {
	fmt.Println("Validating numbers")
	re := regexp.MustCompile(`^(?:27|0)[0-9]{9}$`)
	invalidNumbers := make([][]string, 0)
	validNumbers := make([][]string, 0)

	bar := pb.New(len(reportRows))
	bar.SetWidth(80)
	bar.Start()

	for _, row := range reportRows {
		numLength := len(row[1])
		isInternational := strings.Contains(row[4], "International")

		if isInternational {
			row[1] = strings.Replace(row[1], "00", "", 1)
			// fmt.Printf("%s is an international number\n", row[1])
			validNumbers = append(validNumbers, row)
			continue
		}

		if !isInternational && numLength <= 12 {
			switch {
			case numLength == 9:
				converted := fmt.Sprintf("27%s", row[1])

				if re.MatchString(converted) {
					row[1] = converted
					validNumbers = append(validNumbers, row)
					// fmt.Printf("%s is now a ZA number\n", row[1])
				} else {
					invalidNumbers = append(invalidNumbers, row)
				}
			case numLength == 10:
				converted := strings.Replace(row[1], "0", "27", 1)

				if re.MatchString(converted) {
					row[1] = converted
					validNumbers = append(validNumbers, row)
					// fmt.Printf("%s is now a ZA number\n", row[1])
				} else {
					invalidNumbers = append(invalidNumbers, row)
				}
			}
		}

		bar.Increment()
	}

	bar.Finish()

	return validNumbers, invalidNumbers
}

func filterValidCalls(reportRows [][]string, rateRows [][]string) {
	// rateRows = rateRows[1:]
	southAfricanCodes := filterCodesByCountry(rateRows, "South Africa")
	internationalCodes := filterInternationalCodes(rateRows)

	fmt.Println("Matching numbers to codes & rates")
	bar := pb.New(len(reportRows))
	bar.SetWidth(80)
	bar.Start()
	for i := 0; i < len(reportRows); i++ {

		// Check if the current row is an international call
		if reportRows[i][4] == "International" {

			// If true, loop through internationalCodes and try find matching country
			for _, rate := range internationalCodes {
				if len(reportRows[i][1]) <= 11 {
					callCost := getCallCost(reportRows[i][2], "0.80")
					reportRows[i] = append(reportRows[i], callCost)
					reportRows[i] = append(reportRows[i], "Other")
					break
				}

				matched, _, codeType := matchCode(reportRows[i][1], rate[0], rate[2])

				// If a match is found, calculate the call cost and add to row
				if matched {
					callCost := getCallCost(reportRows[i][2], rate[3])
					reportRows[i] = append(reportRows[i], callCost)
					reportRows[i] = append(reportRows[i], codeType)
					// fmt.Printf("Matched %s to %s\n", reportRows[i][1], codeType)
					break
				}
			}
		} else {
			// If the number is not international, loop through SA codes and find
			// a matching code & rate.
			for _, rate := range southAfricanCodes {
				matched, _, codeType := matchCode(reportRows[i][1], rate[0], rate[2])

				// If a match is found, calculate the call cost and add to the row
				if matched {
					callCost := getCallCost(reportRows[i][2], rate[3])
					reportRows[i] = append(reportRows[i], callCost)
					reportRows[i] = append(reportRows[i], codeType)
					// fmt.Printf("Matched %s to %s\n", reportRows[i][1], codeType)
					break
				}
				if rate[0] == southAfricanCodes[len(southAfricanCodes)-1][0] {
					callCost := getCallCost(reportRows[i][2], "0.80")
					reportRows[i] = append(reportRows[i], callCost)
					reportRows[i] = append(reportRows[i], "Other")
					// errorString := fmt.Sprintf("No known local code for %v", reportRows[i])
					// panic(errorString)
				}

			}

		}

		bar.Increment()
	}
}

func filterInvalidCalls(invalidCalls, rateLines [][]string) {
	return
}

func calculateCallCost(callDuration string, callCost string) string {
	return getCallCost(callDuration, callCost)
}

func filterCodesByCountry(codes [][]string, country string) [][]string {
	return filter(codes, func(code []string) bool {
		return strings.Contains(code[1], country)
	})
}

func filterInternationalCodes(codes [][]string) [][]string {
	return filter(codes, func(code []string) bool {
		return !strings.Contains(code[1], "South Africa")
	})
}

func filter(rows [][]string, f func([]string) bool) [][]string {
	filteredRows := make([][]string, 0)

	for _, row := range rows {
		if f(row) {
			filteredRows = append(filteredRows, row)
		}
	}

	return filteredRows
}

func getCallCost(callDuration, callCost string) string {
	costPerSecond := calculateCostPerSecond(callCost)
	floatDuration, err := strconv.ParseFloat(callDuration, 64)

	if err != nil {
		fmt.Println("Call duration failed: ", err.Error())
		return "N/A"
	}

	finalCost := fmt.Sprintf("%f", floatDuration*costPerSecond)

	return finalCost
}
