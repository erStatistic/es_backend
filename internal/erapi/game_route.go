package erapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GameRoute(routeID int) (Route, error) {

	url := fmt.Sprintf("%s/weaponRoutes/recommend/%d", baseURLv1, routeID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Route{}, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Route{}, err
	}

	if res.StatusCode != 200 {
		return Route{}, fmt.Errorf("HTTP failed to get Game Route: %d", res.StatusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Route{}, err
	}

	result := RouteResponse{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Route{}, err
	}

	if result.Code != 200 {
		fmt.Printf("GameRoute StatusCode : %d\n", result.Code)
		return Route{}, err
	}
	return result.Result, nil
}
