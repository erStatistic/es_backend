package erapi

import (
	"sort"
)

type AnalysisResult struct {
	Clusters  []string
	Total     int
	Top3Ratio float64
	Scores    float64
}

func (c *Client) AnalysisResult(chars [][]int) ([]AnalysisResult, error) {
	positionList, err := readPositionCSV("./internal/statistics/position_dist.csv")
	if err != nil {
		return nil, err
	}

	clusters := [][]string{}
	for _, char := range chars {
		charClusters := []string{}
		for _, cwID := range char {
			cluster := positionList[cwID-1].Position
			charClusters = append(charClusters, cluster)
		}
		sort.Strings(charClusters)
		clusters = append(clusters, charClusters)
	}

	clusterDistList, err := readClusterDistCSV("./internal/statistics/cluster_dist_updated.csv")
	if err != nil {
		return nil, err
	}
	results := []AnalysisResult{}

	for _, cluster := range clusters {
		for _, dist := range clusterDistList {

			if stringSlicesEqual(cluster, dist.ClusterComboKey) {
				results = append(results, AnalysisResult{
					Total:     dist.Total,
					Clusters:  cluster,
					Top3Ratio: dist.Top3Ratio,
					Scores:    dist.NormalizedScore,
				})
				break
			}
		}
	}
	return results, nil

}

func (c *Client) SortClusterDistByNormalizedScore() []ClusterDist {
	clusterDistList, err := readClusterDistCSV("./internal/statistics/cluster_dist_updated.csv")
	if err != nil {
		return nil
	}
	sortClusterDistByNormalizedScore(clusterDistList, false)
	return clusterDistList
}

func sortClusterDistByNormalizedScore(clusterDistList []ClusterDist, descending bool) {
	sort.Slice(clusterDistList, func(i, j int) bool {
		if !descending {
			return clusterDistList[i].NormalizedScore > clusterDistList[j].NormalizedScore
		}
		return clusterDistList[i].NormalizedScore < clusterDistList[j].NormalizedScore
	})
}
