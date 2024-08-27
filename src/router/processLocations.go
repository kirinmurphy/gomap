package router

import (
	"encoding/json"
	"fmt"
	"gomap/src/locationManager"
)

func processLocations(sheetId string, routerConfig RouterConfig) error {
	spreadsheetUrl := fmt.Sprintf(routerConfig.BaseSpreadsheetUrl, sheetId)

	parsedLocations, err := locationManager.LoadLocations(routerConfig.Ctx, spreadsheetUrl)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	locationsJson, err := json.Marshal(parsedLocations)
	if err != nil {
		return fmt.Errorf("failed to marshal locations: %w", err)
	}

	err = routerConfig.RedisClient.Set(routerConfig.Ctx, sheetId, locationsJson, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to cache locations: %w", err)
	}

	return nil
}
