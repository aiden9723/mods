package report

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/scopserv-southafrica/callcoster/error"
)

// MatchCode takes the number, code to match with and the
// code's type. It takes the code and uses it in a regex
// pattern which is used to check if the number matches
// the given code.
func MatchCode(number, code, codeType string) (bool, string, string) {
	formattedPattern := fmt.Sprintf(`^(%s)(?:\d)+$`, code)
	re, err := regexp.Compile(formattedPattern)
	error.Check(err)

	matcher := re.Match([]byte(number))

	return matcher, code, codeType
}

// ValidateNumbers takes the rows and validates the
// numbers within to ensure that all numbers passed
// are assigned to an appropriate country/area code.
func ValidateNumbers(rows [][]string) ([][]string, [][]string) {
	fmt.Println("Validating numbers")
	re := regexp.MustCompile(`^(?:27|0)[0-9]{9}$`)
	invalidNumbers := make([][]string, 0)
	validNumbers := make([][]string, 0)

	bar := pb.New(len(rows))
	bar.SetWidth(80)
	bar.Start()

	for _, row := range rows {
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
