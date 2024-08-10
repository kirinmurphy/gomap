package locationManager

import (
	"fmt"
)

func LoadLocations(spreadsheetUrl string) ([]Location, error) {
	csvStream, errChan := fetchLocationData(spreadsheetUrl)

	parsedLocations, err := parseLocations(csvStream)
	if err != nil {
		return nil, fmt.Errorf("failed to parse locations: %w", err)
	}

	if err := <-errChan; err != nil {
		return nil, fmt.Errorf("failed to fetch CSV data: %w", err)
	}

	return parsedLocations, nil
}
