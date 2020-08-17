package call

import (
	"fmt"
	"strconv"

	"github.com/scopserv-southafrica/callcoster/error"
)

func GetCost(duration, cost string) string {
	costPerSecond := getCostPerSecond(cost)
	floatDuration, err := strconv.ParseFloat(duration, 64)

	if err != nil {
		fmt.Println("Call duration failed: ", err.Error())
		return "N/A"
	}

	finalCost := fmt.Sprintf("%f", floatDuration*costPerSecond)

	return finalCost
}

func getCostPerSecond(cost string) float64 {
	floatCost, err := strconv.ParseFloat(cost, 64)
	error.Check(err)

	return floatCost / 60.0
}
