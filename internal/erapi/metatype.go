package erapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (c *Client) MetaTypes(metaType string) (string, error) {
	url := baseURLv2 + "/data/" + metaType
	if metaType == "" {
		url = baseURLv2 + "/data/hash"
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	result := metatype{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	dataJson, err := json.MarshalIndent(result.Data, "", "  ")
	if err != nil {
		return "", err
	}

	folder := "./internal/output/metatype/"
	filename := "data_" + metaType + ".json"

	err = os.WriteFile(folder+filename, dataJson, 0644)
	if err != nil {
		return "", err
	}

	return "Successfully saved " + filename, nil
}
