package client

import (
	"encoding/json"
	"fmt"
	"time"
)

// PowerReading represents a single power consumption reading from the device.
type PowerReading struct {
	// Timestamp is when the reading was taken.
	Timestamp time.Time `json:"timestamp"`
	// WattsNow is the current power consumption in watts.
	WattsNow float64 `json:"watts_now"`
	// WattsToday is the total energy consumed today in watt-hours.
	WattsToday float64 `json:"watts_today"`
	// WattsLastSevenDays is the total energy consumed in the last 7 days in watt-hours.
	WattsLastSevenDays float64 `json:"watts_last_seven_days"`
}

// DeviceInfo holds metadata about the Efergy monitor device.
type DeviceInfo struct {
	// MAC is the hardware address of the device.
	MAC string `json:"mac"`
	// Type identifies the device model/type.
	Type string `json:"type"`
	// LastSeen is the last time the device reported data.
	LastSeen time.Time `json:"last_seen"`
}

// GetCurrentPower fetches the current instantaneous power reading.
func (c *Client) GetCurrentPower() (*PowerReading, error) {
	resp, err := c.get("/mobile_proxy/getCurrentValuesSummary")
	if err != nil {
		return nil, fmt.Errorf("fetching current power: %w", err)
	}
	defer resp.Body.Close()

	var raw []struct {
		Data []struct {
			SID  int     `json:"sid"`
			Data float64 `json:"data"`
		} `json:"data"`
		Units string `json:"units"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding current power response: %w", err)
	}

	reading := &PowerReading{
		Timestamp: time.Now().UTC(),
	}

	for _, entry := range raw {
		// NOTE: the API also returns kw units for some device configs; only
		// handling watts ("w") here since my hub always reports in watts.
		if entry.Units == "w" && len(entry.Data) > 0 {
			reading.WattsNow = entry.Data[0].Data
		}
	}

	return reading, nil
}

// GetEnergyUsage fetches aggregated energy usage for today and the last 7 days.
func (c *Client) GetEnergyUsage() (*PowerReading, error) {
	reading := &PowerReading{
		Timestamp: time.Now().UTC(),
	}

	// Fetch today's usage
	todayResp, err := c.get("/mobile_proxy/getEnergy?period=day&offset=0")
	if err != nil {
		return nil, fmt.Errorf("fetching today's energy: %w", err)
	}
	defer todayResp.Body.Close()

	var todayRaw struct {
		Sum float64 `json:"sum"`
	}
	if err := json.NewDecoder(todayResp.Body).Decode(&todayRaw); err != nil {
		return nil, fmt.Errorf("decoding today's energy response: %w", err)
	}
	reading.WattsToday = todayRaw.Sum

	// Fetch last 7 days usage
	weekResp, err := c.get("/mobile_proxy/getEnergy?period=week&offset=0")
	if err != nil {
		return nil, fmt.Errorf("fetching weekly energy: %w", err)
	}
	defer weekResp.Body.Close()

	var weekRaw struct {
		Sum float64 `json:"sum"`
	}
	if err := json.NewDecoder(weekResp.Body).Decode(&weekRaw); err != nil {
		return nil, fmt.Errorf("decoding weekly energy response: %w", err)
	}
	reading.WattsLastSevenDays = weekRaw.Sum

	return reading, nil
}

// GetDeviceInfo retrieves metadata about the connected Efergy device.
func (c *Client) GetDeviceInfo() (*DeviceInfo, error) {
	resp, err := c.get("/mobile_proxy/getDevices")
	if err != nil {
		return nil, fmt.Errorf("fetching device info: %w", err)
	}
	defer resp.Body.Close()

	var raw struct {
		MAC      string `json:"mac"`
		Type     string `json:"type"`
		LastSeen int64  `json:"last_seen"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding device info response: %w", err)
	}

	return &DeviceInfo{
		MAC:      raw.MAC,
		Type:     raw.Type,
		LastSeen: time.Unix(raw.LastSeen, 0).UTC(),
	}, nil
}
