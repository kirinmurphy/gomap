package main

import (
	"encoding/csv"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func fetchLocationsFromGoogleSheets(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var locs []Location

	reader := csv.NewReader(resp.Body)

	headers, err := reader.Read()
	if err != nil {
		return err
	}

	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		lat, _ := strconv.ParseFloat(record[headerMap["Latitude"]], 64)
		long, _ := strconv.ParseFloat(record[headerMap["Longitude"]], 64)
		name := record[headerMap["Name"]]

		locs = append(locs, Location{
			Name:       name,
			Address:    record[headerMap["Address"]],
			City:       record[headerMap["City"]],
			State:      record[headerMap["State"]],
			Country:    record[headerMap["Country"]],
			Latitude:   lat,
			Longitude:  long,
			IsCo404Loc: containsCo404(name),
		})
	}
	locations = locs
	return nil
}

func containsCo404(s string) bool {
	lowerStr := strings.ToLower(s)
	return strings.Contains(lowerStr, "co404")
}
