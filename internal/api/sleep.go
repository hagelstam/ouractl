package api

import (
	"encoding/json"
	"net/url"
)

type SleepContributors struct {
	DeepSleep   *int `json:"deep_sleep"`
	Efficiency  *int `json:"efficiency"`
	Latency     *int `json:"latency"`
	REMSleep    *int `json:"rem_sleep"`
	Restfulness *int `json:"restfulness"`
	Timing      *int `json:"timing"`
	TotalSleep  *int `json:"total_sleep"`
}

type DailySleep struct {
	ID           string            `json:"id"`
	Day          string            `json:"day"`
	Score        *int              `json:"score"`
	Contributors SleepContributors `json:"contributors"`
	Timestamp    string            `json:"timestamp"`
}

type dailySleepResponse struct {
	Data      []DailySleep `json:"data"`
	NextToken *string      `json:"next_token"`
}

func (c *Client) GetDailySleep(startDate, endDate string) ([]DailySleep, error) {
	var all []DailySleep
	var nextToken string

	for {
		params := url.Values{}
		params.Set("start_date", startDate)
		params.Set("end_date", endDate)
		if nextToken != "" {
			params.Set("next_token", nextToken)
		}

		body, err := c.Get("/v2/usercollection/daily_sleep", params)
		if err != nil {
			return nil, err
		}

		var resp dailySleepResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}

		all = append(all, resp.Data...)

		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = *resp.NextToken
	}

	return all, nil
}
