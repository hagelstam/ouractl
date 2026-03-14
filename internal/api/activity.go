package api

import "net/url"

type ActivityContributors struct {
	MeetDailyTargets  *int `json:"meet_daily_targets"`
	MoveEveryHour     *int `json:"move_every_hour"`
	RecoveryTime      *int `json:"recovery_time"`
	StayActive        *int `json:"stay_active"`
	TrainingFrequency *int `json:"training_frequency"`
	TrainingVolume    *int `json:"training_volume"`
}

type DailyActivity struct {
	ID                        string               `json:"id"`
	Day                       string               `json:"day"`
	Score                     *int                 `json:"score"`
	ActiveCalories            int                  `json:"active_calories"`
	TotalCalories             int                  `json:"total_calories"`
	EquivalentWalkingDistance int                  `json:"equivalent_walking_distance"`
	Steps                     int                  `json:"steps"`
	HighActivityTime          int                  `json:"high_activity_time"`
	MediumActivityTime        int                  `json:"medium_activity_time"`
	LowActivityTime           int                  `json:"low_activity_time"`
	RestingTime               int                  `json:"resting_time"`
	SedentaryTime             int                  `json:"sedentary_time"`
	NonWearTime               int                  `json:"non_wear_time"`
	TargetCalories            int                  `json:"target_calories"`
	TargetMeters              int                  `json:"target_meters"`
	InactivityAlerts          int                  `json:"inactivity_alerts"`
	Contributors              ActivityContributors `json:"contributors"`
	Timestamp                 string               `json:"timestamp"`
}

func (c *Client) GetDailyActivity(startDate, endDate string) ([]DailyActivity, error) {
	params := url.Values{}
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)
	return getPaged[DailyActivity](c, "/v2/usercollection/daily_activity", params)
}
