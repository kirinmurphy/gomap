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

	reader := csv.NewReader(resp.Body)
	var locs []Location

	if _, err := reader.Read(); err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		latRowIndex := 5
		longRowIndex := 6

		lat, _ := strconv.ParseFloat(record[latRowIndex], 64)
		long, _ := strconv.ParseFloat(record[longRowIndex], 64)
		locs = append(locs, Location{
			Name:       record[0],
			Address:    record[1],
			City:       record[2],
			State:      record[3],
			Country:    record[4],
			Latitude:   lat,
			Longitude:  long,
			IsCo404Loc: containsCo404(record[0]),
		})
	}
	locations = locs
	return nil
}

func containsCo404(s string) bool {
	lowerStr := strings.ToLower(s)
	return strings.Contains(lowerStr, "co404")
}
