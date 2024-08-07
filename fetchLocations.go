package main

import (
	"encoding/csv"
	"io"
	"net/http"
)

func fetchLocationData(spreadsheetUrl string) (<-chan []string, <-chan error) {
	csvStream := make(chan []string)
	errChan := make(chan error, 1)

	go func() {
		resp, err := http.Get(spreadsheetUrl)

		if err != nil {
			errChan <- err
			close(csvStream)
			close(errChan)
		}

		defer resp.Body.Close()

		reader := csv.NewReader(resp.Body)

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				errChan <- err
				close(csvStream)
				close(errChan)
				return
			}

			csvStream <- record
		}

		close(csvStream)
		errChan <- nil
		close(errChan)
	}()

	return csvStream, errChan
}
