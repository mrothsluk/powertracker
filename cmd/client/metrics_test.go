package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMetrics_ValidData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(t *testing.T, m *Metrics)
	}{
		{
			name: "valid power reading",
			input: `{"power": 1234.5, "voltage": 230.1, "current": 5.36, "frequency": 50.0}`,
			wantErr: false,
			validate: func(t *testing.T, m *Metrics) {
				assert.InDelta(t, 1234.5, m.Power, 0.01)
				assert.InDelta(t, 230.1, m.Voltage, 0.01)
				assert.InDelta(t, 5.36, m.Current, 0.01)
				assert.InDelta(t, 50.0, m.Frequency, 0.01)
			},
		},
		{
			name:    "empty input",
			input:   ``,
			wantErr: true,
		},
		{
			name:    "malformed JSON",
			input:   `{power: 1234.5}`,
			wantErr: true,
		},
		{
			name:    "missing fields",
			input:   `{"power": 100.0}`,
			wantErr: false,
			validate: func(t *testing.T, m *Metrics) {
				assert.InDelta(t, 100.0, m.Power, 0.01)
				assert.Equal(t, float64(0), m.Voltage)
				assert.Equal(t, float64(0), m.Current)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := parseMetrics([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, m)
			if tt.validate != nil {
				tt.validate(t, m)
			}
		})
	}
}

func TestMetrics_Timestamp(t *testing.T) {
	before := time.Now()
	m, err := parseMetrics([]byte(`{"power": 500.0, "voltage": 240.0, "current": 2.08, "frequency": 60.0}`))
	after := time.Now()

	require.NoError(t, err)
	require.NotNil(t, m)

	// Timestamp should be set at parse time
	assert.True(t, m.Timestamp.After(before) || m.Timestamp.Equal(before),
		"timestamp should be after or equal to before")
	assert.True(t, m.Timestamp.Before(after) || m.Timestamp.Equal(after),
		"timestamp should be before or equal to after")
}

func TestMetrics_PowerFactor(t *testing.T) {
	tests := []struct {
		name            string
		power           float64
		voltage         float64
		current         float64
		wantPowerFactor float64
	}{
		{
			name:            "unity power factor",
			power:           230.0,
			voltage:         230.0,
			current:         1.0,
			wantPowerFactor: 1.0,
		},
		{
			name:            "half power factor",
			power:           115.0,
			voltage:         230.0,
			current:         1.0,
			wantPowerFactor: 0.5,
		},
		{
			name:            "zero current avoids division by zero",
			power:           0.0,
			voltage:         230.0,
			current:         0.0,
			wantPowerFactor: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				Power:   tt.power,
				Voltage: tt.voltage,
				Current: tt.current,
			}
			pf := m.PowerFactor()
			assert.InDelta(t, tt.wantPowerFactor, pf, 0.001)
		})
	}
}
