package locationManager

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
)

func fetchLocationData(spreadsheetUrl string) (<-chan []string, <-chan error) {
	csvStream := make(chan []string)
	errChan := make(chan error, 1)

	go func() {
		resp, err := http.Get(spreadsheetUrl)
		if err != nil {
			log.Println("Error fetching URL:", err)
			errChan <- err
			close(csvStream)
			close(errChan)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			close(csvStream)
			close(errChan)
			return
		}

		reader := csv.NewReader(resp.Body)

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Println("Error reading CSV:", err)
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
