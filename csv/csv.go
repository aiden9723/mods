package csv

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	e "github.com/yokimpillay/callcoster/error"
)

// Read finds the filename specified and reads the
// file to create the new report with.
func Read(filename string) ([][]string, error) {
	if filename == "" {
		fmt.Println("File not specified, exiting.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("Reading from", filename)
	f, err := os.Open(filename)
	e.Check(err)

	lines, err := csv.NewReader(f).ReadAll()

	f.Close()

	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

// Write takes the data to write to a file with, as well
// as the filename and creates a formatted filename and
// writes the file to the disk.
func Write(rows [][]string, filename string) {
	newFilename := MakeFilename(filename)

	f, err := os.Create(newFilename)
	e.Check(err)

	err = csv.NewWriter(f).WriteAll(rows)
	e.Check(err)

	fmt.Println("\nSuccessfully written to", newFilename)
	return
}

// MakeFilename simply formats the filename to make it
// easier to find the file once the program exits.
func MakeFilename(filename string) string {
	trimmedFilename := strings.Replace(filename, ".csv", "", -1)
	newFilename := fmt.Sprintf("%s-formatted.csv", trimmedFilename)
	return newFilename
}

// AddColumns adds an arbitrary amount of column names
// to the rows
func AddColumns(rows [][]string, columnNames ...string) {
	for _, name := range columnNames {
		rows[0] = append(rows[0], name)
	}
}
