package locationManager

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
)

func fetchLocationData(ctx context.Context, spreadsheetUrl string) (<-chan []string, <-chan error) {
	csvStream := make(chan []string)
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			close(csvStream)
			close(errChan)
		}()

		req, err := http.NewRequestWithContext(ctx, "GET", spreadsheetUrl, nil)
		if err != nil {
			errChan <- err
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errChan <- err
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("failed to fetch CSV data: unexpected status code: %d", resp.StatusCode)
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
				return
			}

			select {
			case csvStream <- record:
			case <-ctx.Done():
				errChan <- fmt.Errorf("the operation timed out, please check your spreadsheet and try again: %w", ctx.Err())
				return
			}
		}

		// errChan <- nil
	}()

	return csvStream, errChan
}
