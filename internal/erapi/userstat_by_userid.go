package erapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) UserStatByUserId(userNum int) ([]UserCharacterStat, error) {
	SeasonID := 33
	MatchingMode := 3

	url := fmt.Sprintf("%s/user/stats/%d/%d/%d", baseURLv2, userNum, SeasonID, MatchingMode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP failed to get User Stat By UserId: %d", res.StatusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := UserStatResponse{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 200 {
		fmt.Printf("UserStatByUserId StatusCode : %d\n", result.Code)
		return nil, err
	}

	return result.UserStats.CharacterStats, nil
}
