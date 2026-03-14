package api

import "net/url"

type SleepReadiness struct {
	Contributors              ReadinessContributors `json:"contributors"`
	Score                     *int                  `json:"score"`
	TemperatureDeviation      *float64              `json:"temperature_deviation"`
	TemperatureTrendDeviation *float64              `json:"temperature_trend_deviation"`
}

type Sleep struct {
	ID                  string          `json:"id"`
	Day                 string          `json:"day"`
	BedtimeStart        string          `json:"bedtime_start"`
	BedtimeEnd          string          `json:"bedtime_end"`
	DeepSleepDuration   *int            `json:"deep_sleep_duration"`
	LightSleepDuration  *int            `json:"light_sleep_duration"`
	REMSleepDuration    *int            `json:"rem_sleep_duration"`
	AwakeTime           *int            `json:"awake_time"`
	TotalSleepDuration  *int            `json:"total_sleep_duration"`
	TimeInBed           int             `json:"time_in_bed"`
	AverageHeartRate    *float64        `json:"average_heart_rate"`
	LowestHeartRate     *int            `json:"lowest_heart_rate"`
	AverageHRV          *int            `json:"average_hrv"`
	AverageBreath       *float64        `json:"average_breath"`
	Efficiency          *int            `json:"efficiency"`
	Latency             *int            `json:"latency"`
	Readiness           *SleepReadiness `json:"readiness"`
	ReadinessScoreDelta *float64        `json:"readiness_score_delta"`
	Type                string          `json:"type"`
	Period              int             `json:"period"`
}

func (c *Client) GetSleep(startDate, endDate string) ([]Sleep, error) {
	params := url.Values{}
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)
	return getPaged[Sleep](c, "/v2/usercollection/sleep", params)
}
