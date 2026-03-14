package api

type RingConfig struct {
	ID              string  `json:"id"`
	Color           *string `json:"color"`
	Design          *string `json:"design"`
	FirmwareVersion *string `json:"firmware_version"`
	HardwareType    *string `json:"hardware_type"`
	SetUpAt         *string `json:"set_up_at"`
	Size            *int    `json:"size"`
}

func (c *Client) GetRingConfig() ([]RingConfig, error) {
	return getPaged[RingConfig](c, "/v2/usercollection/ring_configuration", nil)
}
