package erapi

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Position struct {
	code     int
	Position string
}

func (c *Client) PositionList() ([]Position, error) {
	return readPositionCSV("./internal/statistics/prepare/position_dist.csv")
}

func readPositionCSV(filePath string) ([]Position, error) {
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
		"no", "kmeans",
	}
	for _, h := range requiredHeaders {
		if _, exists := headerMap[h]; !exists {
			return nil, fmt.Errorf("missing required header: %s", h)
		}
	}

	var positionList []Position
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
		position := record[headerMap["kmeans"]]

		positionList = append(positionList, Position{
			code:     code,
			Position: position,
		})
	}

	return positionList, nil
}

func printPositionList() {

	positionList, err := readPositionCSV("position_dist.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i, item := range positionList[:min(5, len(positionList))] {
		fmt.Printf("Row %d: Code=%d, Position=%s\n", i+1, item.code, item.Position)
	}
}
