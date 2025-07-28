package erapi

import (
	"encoding/json"
	"os"
)

func (c *Client) GetWeaponInfo() (Weapons, error) {
	fileurl := "./internal/output/metatype/data_WeaponTypeInfo.json"

	data, err := os.ReadFile(fileurl)
	if err != nil {
		return nil, err
	}

	var weapons Weapons

	err = json.Unmarshal(data, &weapons)
	if err != nil {
		return nil, err
	}

	return weapons, nil
}
