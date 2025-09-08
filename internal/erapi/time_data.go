package erapi

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Time struct {
	Code    int
	Name    string
	Seconds int
	Total   int
}

func (c *Client) TimeList() ([]Time, error) {
	return readTimeCSV("./internal/statistics/prepare/time_data.csv")
}

func readTimeCSV(filePath string) ([]Time, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read headers: %v", err)
	}

	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	requiredHeaders := []string{
		"no", "name", "seconds", "totaltime",
	}

	for _, h := range requiredHeaders {
		if _, exists := headerMap[h]; !exists {
			return nil, fmt.Errorf("missing required header: %s", h)
		}
	}

	var timeList []Time

	for rowNum := 2; ; rowNum++ {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to read record at row %d: %v", rowNum, err)
		}

		code, err := strconv.Atoi(record[headerMap["no"]])
		if err != nil {
			return nil, fmt.Errorf("failed to parse code at row %d: %v", rowNum, err)
		}
		name := record[headerMap["name"]]
		seconds, err := strconv.Atoi(record[headerMap["seconds"]])
		if err != nil {
			return nil, fmt.Errorf("failed to parse seconds at row %d: %v", rowNum, err)
		}
		total, err := strconv.Atoi(record[headerMap["totaltime"]])
		if err != nil {
			return nil, fmt.Errorf("failed to parse total at row %d: %v", rowNum, err)
		}

		timeList = append(timeList, Time{
			Code:    code,
			Name:    name,
			Seconds: seconds,
			Total:   total,
		})
	}
	return timeList, nil
}
