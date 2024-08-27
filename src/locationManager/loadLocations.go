package locationManager

import (
	"context"
	"fmt"
)

func LoadLocations(ctx context.Context, spreadsheetUrl string) ([]Location, error) {
	csvStream, errChan := fetchLocationData(ctx, spreadsheetUrl)

	parsedLocations, parseErr := parseLocations(csvStream)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("the operation timed out, please check your spreadsheet and try again")
	case fetchErr := <-errChan:
		if fetchErr != nil {
			return nil, fmt.Errorf("failed to fetch CSV data, please try again")
		}
	default:
	}

	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse locations: %w", parseErr)
	}

	if len(parsedLocations) == 0 {
		return nil, fmt.Errorf("no valid locations found in CSV")
	}

	return parsedLocations, nil
}
