package gm1356

import "fmt"

// TimeWeight is time weighting characteristic
type TimeWeight int

// FrequencyWeight is frequency weighting characteristic
type FrequencyWeight int

// SoundLevelRange is sound level range on indicator
type SoundLevelRange int

// SoundLevelDisplayMode is display mode
type SoundLevelDisplayMode int

// Config is configuration of GM1356 device
type Config struct {
	TimeWeight            TimeWeight
	FrequencyWeight       FrequencyWeight
	SoundLevelRange       SoundLevelRange
	SoundLevelDisplayMode SoundLevelDisplayMode
}

const (
	// TimeWeightFast is fast mode
	TimeWeightFast TimeWeight = iota
	// TimeWeightSlow is slow mode
	TimeWeightSlow
)

const (
	// FrequencyWeightA is dBA
	FrequencyWeightA FrequencyWeight = iota
	// FrequencyWeightC is dBC
	FrequencyWeightC
)

const (
	// SoundLevelRange30_130 is from 30dB to 130dB
	SoundLevelRange30_130 = iota
	// SoundLevelRange30_80 is from 30dB to 80dB
	SoundLevelRange30_80
	// SoundLevelRange50_100 is from 50dB to 100dB
	SoundLevelRange50_100
	// SoundLevelRange60_110 is from 60dB to 110dB
	SoundLevelRange60_110
	// SoundLevelRange80_130 is from 80dB to 130dB
	SoundLevelRange80_130
)

const (
	// SoundLevelDisplayModeNormal is normal mode
	SoundLevelDisplayModeNormal SoundLevelDisplayMode = iota
	// SoundLevelDisplayModeMaxHold is max hold mode
	SoundLevelDisplayModeMaxHold
)

// TimeWeight is time weighting characteristic
func (t TimeWeight) String() string {
	switch t {
	case TimeWeightFast:
		return "fast"
	case TimeWeightSlow:
		return "slow"
	default:
		return "unknown"
	}
}

// FrequencyWeight is frequency weighting characteristic
func (f FrequencyWeight) String() string {
	switch f {
	case FrequencyWeightA:
		return "dBA"
	case FrequencyWeightC:
		return "dBC"
	default:
		return "unknown"
	}
}

// SoundLevelRange is sound level range on indicator
func (l SoundLevelRange) String() string {
	switch l {
	case SoundLevelRange30_130:
		return "30dB-130dB"
	case SoundLevelRange30_80:
		return "30dB-80dB"
	case SoundLevelRange50_100:
		return "50dB-100dB"
	case SoundLevelRange60_110:
		return "60dB-110dB"
	case SoundLevelRange80_130:
		return "80dB-130dB"
	default:
		return "unknown"
	}
}

// SoundLevelDisplayMode is display mode
func (d SoundLevelDisplayMode) String() string {
	switch d {
	case SoundLevelDisplayModeNormal:
		return "normal"
	case SoundLevelDisplayModeMaxHold:
		return "max hold"
	default:
		return "unknown"
	}
}

// String returns string representing config
func (c Config) String() string {
	return fmt.Sprintf("Config{TimeWeight: %s, FrequencyWeight: %s, SoundLevelDisplayMode: %s, SoundLevelRange: %s}",
		c.TimeWeight.String(), c.FrequencyWeight.String(), c.SoundLevelDisplayMode.String(), c.SoundLevelRange.String())
}
