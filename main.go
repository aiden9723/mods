package main

import (
	"flag"
	"github.com/cheggaaa/pb/v3"
	"github.com/yokimpillay/callcoster/csv"
	"github.com/yokimpillay/callcoster/error"
	"github.com/yokimpillay/callcoster/filter"
	"github.com/yokimpillay/callcoster/report"
)

var bar *pb.ProgressBar

func main() {
	reportFilename := flag.String("report", "", "The file to calculate call costs on")
	rates := flag.String("rates", "", "The file from which costs will be used")

	flag.Parse()

	reportLines, err := csv.Read(*reportFilename)
	error.Check(err)

	csv.AddColumns(reportLines, "Cost", "Code")

	reportHeaders, reportLines := reportLines[0], reportLines[1:]

	rateLines, err := csv.Read(*rates)
	error.Check(err)

	validnums, _ := report.ValidateNumbers(reportLines)

	filter.ValidCalls(validnums, rateLines)

	reportLines = append([][]string{reportHeaders}, validnums...)

	csv.Write(reportLines, *reportFilename)
}
