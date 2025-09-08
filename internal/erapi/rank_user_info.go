package erapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) RankUserInfo(userID, teamMode, seasonId int) (UserRank, error) {

	url := fmt.Sprintf("%s/rank/%d/%d/%d", baseURLv1, userID, seasonId, teamMode)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return UserRank{}, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return UserRank{}, err
	}

	if res.StatusCode != 200 {
		return UserRank{}, fmt.Errorf("HTTP failed to get RankUserInfo: %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return UserRank{}, err
	}

	var RankResponse RankResponse
	err = json.Unmarshal(body, &RankResponse)
	if err != nil {
		return UserRank{}, err
	}

	if RankResponse.Code != 200 {
		fmt.Printf("RankUserInfo StatusCode : %d\n", RankResponse.Code)
		return UserRank{}, errors.New("failed to get rank")
	}

	return RankResponse.UserRank, nil
}
