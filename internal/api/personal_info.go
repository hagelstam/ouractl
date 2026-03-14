package api

import "encoding/json"

type PersonalInfo struct {
	ID            string   `json:"id"`
	Age           *int     `json:"age"`
	Weight        *float64 `json:"weight"`
	Height        *float64 `json:"height"`
	BiologicalSex *string  `json:"biological_sex"`
	Email         *string  `json:"email"`
}

func (c *Client) GetPersonalInfo() (*PersonalInfo, error) {
	body, err := c.Get("/v2/usercollection/personal_info", nil)
	if err != nil {
		return nil, err
	}

	var info PersonalInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
