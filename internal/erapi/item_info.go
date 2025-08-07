package erapi

import (
	"encoding/json"
	"errors"
	"os"
)

type targetItem struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

func (c *Client) ItemInfo(itemCode int) (*targetItem, error) {

	fileurl := "./internal/output/metatype/data_ItemWeapon.json"

	data, err := os.ReadFile(fileurl)
	if err != nil {
		return nil, err
	}

	items := []item{}

	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}

	target := targetItem{}

	for _, item := range items {
		if item.Code == itemCode {
			target = targetItem{
				Code: item.Code,
				Name: item.Name,
			}
			return &target, nil
		}
	}

	return nil, errors.New("item not found")
}
