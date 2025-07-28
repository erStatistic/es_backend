package erapi

import (
	"encoding/json"
	"os"
)

func (c *Client) GetCharacterInfo() (Characters, error) {
	fileurl := "./internal/output/metatype/data_Character.json"

	data, err := os.ReadFile(fileurl)
	if err != nil {
		return nil, err
	}

	var characters Characters

	err = json.Unmarshal(data, &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}
