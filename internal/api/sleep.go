package api

import "net/url"

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

func (c *Client) GetDailySleep(startDate, endDate string) ([]DailySleep, error) {
	params := url.Values{}
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)
	return getPaged[DailySleep](c, "/v2/usercollection/daily_sleep", params)
}
