package locationHelpers

import (
	"strconv"
	"strings"
)

type Location struct {
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	IsCo404Loc bool    `json:"isCo404Loc"`
}

func ParseLocations(csvStream <-chan []string) ([]Location, error) {
	var locs []Location

	headerMap := make(map[string]int)

	isHeader := true

	for record := range csvStream {
		if isHeader {
			for i, header := range record {
				headerMap[header] = i
			}
			isHeader = false
			continue
		}

		loc, err := parseLocation(record, headerMap)
		if err != nil {
			return nil, err
		}

		locs = append(locs, loc)
	}

	return locs, nil
}

func parseLocation(record []string, headerMap map[string]int) (Location, error) {
	lat, err := strconv.ParseFloat(record[headerMap["Latitude"]], 64)
	if err != nil {
		return Location{}, err
	}

	long, err := strconv.ParseFloat(record[headerMap["Longitude"]], 64)
	if err != nil {
		return Location{}, err
	}

	name := record[headerMap["Name"]]

	parsedLoc := Location{
		Name:       name,
		Address:    record[headerMap["Address"]],
		City:       record[headerMap["City"]],
		State:      record[headerMap["State"]],
		Country:    record[headerMap["Country"]],
		Latitude:   lat,
		Longitude:  long,
		IsCo404Loc: containsCo404(name),
	}

	return parsedLoc, nil
}

func containsCo404(s string) bool {
	lowerStr := strings.ToLower(s)
	return strings.Contains(lowerStr, "co404")
}
