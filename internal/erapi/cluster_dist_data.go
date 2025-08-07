package erapi

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ClusterDist 구조체 정의
type ClusterDist struct {
	ClusterComboKey []string // Cluster_Combo_Key
	Counts          [8]int   // 1,2,3,4,5,6,7,8
	Total           int      // Total
	WeightedScore   int      // Weighted_Score
	Top3Ratio       float64  // Top3_Ratio
	NormalizedScore float64  // Normalized_Score
}

func readClusterDistCSV(filePath string) ([]ClusterDist, error) {
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
		"Cluster_Combo_Key", "1", "2", "3", "4", "5", "6", "7", "8",
		"Total", "Weighted_Score", "Top3_Ratio", "Normalized_Score",
	}
	for _, h := range requiredHeaders {
		if _, exists := headerMap[h]; !exists {
			return nil, fmt.Errorf("missing required header: %s", h)
		}
	}

	var clusterDistList []ClusterDist
	for rowNum := 2; ; rowNum++ {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to read record at row %d: %v", rowNum, err)
		}

		var counts [8]int
		for i := 0; i < 8; i++ {
			counts[i], err = strconv.Atoi(record[headerMap[strconv.Itoa(i+1)]])
			if err != nil {
				return nil, fmt.Errorf("failed to parse count %d at row %d: %v", i+1, rowNum, err)
			}
		}

		clusterComboKey := parseClusterKey(record[headerMap["Cluster_Combo_Key"]])

		total, err := strconv.Atoi(record[headerMap["Total"]])
		if err != nil {
			return nil, fmt.Errorf("failed to parse Total at row %d: %v", rowNum, err)
		}
		if total < 20 {
			continue
		}

		weightedScore, err := strconv.Atoi(record[headerMap["Weighted_Score"]])
		if err != nil {
			return nil, fmt.Errorf("failed to parse Weighted_Score at row %d: %v", rowNum, err)
		}

		top3Ratio, err := strconv.ParseFloat(record[headerMap["Top3_Ratio"]], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Top3_Ratio at row %d: %v", rowNum, err)
		}

		normalizedScore, err := strconv.ParseFloat(record[headerMap["Normalized_Score"]], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Normalized_Score at row %d: %v", rowNum, err)
		}

		clusterDistList = append(clusterDistList, ClusterDist{
			ClusterComboKey: clusterComboKey,
			Counts:          counts,
			Total:           total,
			WeightedScore:   weightedScore,
			Top3Ratio:       top3Ratio,
			NormalizedScore: normalizedScore,
		})
	}

	return clusterDistList, nil
}

func parseClusterKey(s string) []string {
	s = strings.TrimPrefix(s, "('")
	s = strings.TrimSuffix(s, "')")

	parts := strings.Split(s, "', '")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func printClusterDist() {
	clusterDistList, err := readClusterDistCSV("cluster_dist_updated.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i, item := range clusterDistList[:min(5, len(clusterDistList))] {
		fmt.Printf("Row %d: ClusterComboKey=%s, NormalizedScore=%.3f\n", i+1, item.ClusterComboKey, item.NormalizedScore)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
