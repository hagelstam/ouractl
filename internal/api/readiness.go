package api

import "net/url"

type ReadinessContributors struct {
	ActivityBalance     *int `json:"activity_balance"`
	BodyTemperature     *int `json:"body_temperature"`
	HRVBalance          *int `json:"hrv_balance"`
	PreviousDayActivity *int `json:"previous_day_activity"`
	PreviousNight       *int `json:"previous_night"`
	RecoveryIndex       *int `json:"recovery_index"`
	RestingHeartRate    *int `json:"resting_heart_rate"`
	SleepBalance        *int `json:"sleep_balance"`
	SleepRegularity     *int `json:"sleep_regularity"`
}

type DailyReadiness struct {
	ID                        string                `json:"id"`
	Day                       string                `json:"day"`
	Score                     *int                  `json:"score"`
	TemperatureDeviation      *float64              `json:"temperature_deviation"`
	TemperatureTrendDeviation *float64              `json:"temperature_trend_deviation"`
	Contributors              ReadinessContributors `json:"contributors"`
	Timestamp                 string                `json:"timestamp"`
}

func (c *Client) GetDailyReadiness(startDate, endDate string) ([]DailyReadiness, error) {
	params := url.Values{}
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)
	return getPaged[DailyReadiness](c, "/v2/usercollection/daily_readiness", params)
}
