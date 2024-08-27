package locationManager

import (
	"fmt"
	"strconv"
	"strings"
)

type Location struct {
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Country     string  `json:"country"`
	Website     string  `json:"website"`
	PhoneNumber string  `json:"phoneNumber"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	IsCo404Loc  bool    `json:"isCo404Loc"`
}

func parseLocations(csvStream <-chan []string) ([]Location, error) {
	var locs []Location

	headerMap := make(map[string]int)

	isHeader := true

	for record := range csvStream {
		if isHeader {
			for i, header := range record {
				headerMap[header] = i
			}
			isHeader = false
			// log.Println("headerMap: ", headerMap)
			continue
		}

		loc, err := parseLocation(record, headerMap)
		if err != nil {
			for range csvStream {
				// draining the remaining
			}
			fmt.Printf("ERRRRRRRRORROORORORORO: %s", err)
			return nil, err
		}

		locs = append(locs, loc)
	}

	return locs, nil
}

func parseLocation(record []string, headerMap map[string]int) (Location, error) {
	name := NewSanitizer(record[headerMap["Name"]]).MaxLength(100).Result()
	address := NewSanitizer(record[headerMap["Address"]]).MaxLength(255).Result()
	city := NewSanitizer(record[headerMap["City"]]).MaxLength(100).Result()
	state := NewSanitizer(record[headerMap["State"]]).MaxLength(100).Result()
	country := NewSanitizer(record[headerMap["Country"]]).MaxLength(100).Result()
	website := NewSanitizer(record[headerMap["Website"]]).ValidateURL().Result()
	phoneNumber := NewSanitizer(record[headerMap["Phone Number"]]).MaxLength(20).Result()

	lat, err := strconv.ParseFloat(strings.TrimSpace(record[headerMap["Latitude"]]), 64)
	if err != nil {
		return Location{}, err
	}

	long, err := strconv.ParseFloat(strings.TrimSpace(record[headerMap["Longitude"]]), 64)
	if err != nil {
		return Location{}, err
	}

	// log.Printf("record: %s", record)
	parsedLoc := Location{
		Name:        name,
		Address:     address,
		City:        city,
		State:       state,
		Country:     country,
		Website:     website,
		PhoneNumber: phoneNumber,
		Latitude:    lat,
		Longitude:   long,
		IsCo404Loc:  containsCo404(name),
	}

	return parsedLoc, nil
}

func containsCo404(s string) bool {
	lowerStr := strings.ToLower(s)
	return strings.Contains(lowerStr, "co404")
}
