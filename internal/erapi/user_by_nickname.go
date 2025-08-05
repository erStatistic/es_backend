package erapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) UserByNickname(nickname string) (*User, error) {
	url := fmt.Sprintf("%s/user/nickname", baseURLv1)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	q := req.URL.Query()
	q.Add("query", nickname)
	req.URL.RawQuery = q.Encode()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := userResponse{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 200 {
		fmt.Printf("StatusCode : %d\n", result.Code)
		return nil, err
	}

	return &result.User, nil

}
