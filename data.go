package gm1356

import (
	"math/rand"
	"time"
)

const (
	// mark
	configureRequestMark      = 0x56
	configureResponseMark     = 0xc4
	measureRequestMark        = 0xb3
	importRequestMark         = 0xb5
	importResponseMark        = 0xef
	nextImportDataRequestMark = 0xc4
	nextImportBlockMark       = 0xfd
	noImportDataMark          = 0xff

	// flags
	timeWeightFlags      = 0x40
	displayModeFlags     = 0x20
	frequencyWeightFlags = 0x10
	levelRangeFlags      = 0x0F
)

// randomize magic data
// 0: static
// 1: randomize non-zero data
// 2: randomize all data
const randomizeMagicDataConfig = 0

// magic data
var (
	magicDataMeasureRequest1      byte = 0xfd
	magicDataMeasureRequest2      byte = 0x18
	magicDataMeasureRequest3      byte = 0x00
	magicDataMeasureRequest4      byte = 0xcc
	magicDataMeasureRequest5      byte = 0xfd
	magicDataMeasureRequest6      byte = 0x18
	magicDataMeasureRequest7      byte = 0x00
	magicDataConfigureRequest2    byte = 0x40
	magicDataConfigureRequest3    byte = 0x00
	magicDataConfigureRequest4    byte = 0x92
	magicDataConfigureRequest5    byte = 0x9f
	magicDataConfigureRequest6    byte = 0x40
	magicDataConfigureRequest7    byte = 0x00
	magicDataImportRequest1       byte = 0x00
	magicDataImportRequest2       byte = 0x00
	magicDataImportRequest3       byte = 0x00
	magicDataImportRequest4       byte = 0x50
	magicDataImportRequest5       byte = 0x7c
	magicDataImportRequest6       byte = 0x7e
	magicDataImportRequest7       byte = 0x02
	magicDataNextImportRequest1   byte = 0x00
	magicDataNextImportRequest2   byte = 0x00
	magicDataNextImportRequest3   byte = 0x00
	magicDataNextImportRequest1_4 byte = 0x0a
	magicDataNextImportRequest1_5 byte = 0xeb
	magicDataNextImportRequest1_6 byte = 0x6f
	magicDataNextImportRequest1_7 byte = 0x00
	magicDataNextImportRequest2_4 byte // random
	magicDataNextImportRequest2_5 byte // random
	magicDataNextImportRequest2_6 byte // random
	magicDataNextImportRequest2_7 byte // random
)

func init() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	magicDataNextImportRequest2_4 = byte(r.Uint32())
	magicDataNextImportRequest2_5 = byte(r.Uint32())
	magicDataNextImportRequest2_6 = byte(r.Uint32())
	magicDataNextImportRequest2_7 = byte(r.Uint32())
	// randomize non-zero data
	if randomizeMagicDataConfig > 0 {
		magicDataMeasureRequest1 = byte(r.Uint32())
		magicDataMeasureRequest2 = byte(r.Uint32())
		magicDataMeasureRequest4 = byte(r.Uint32())
		magicDataMeasureRequest5 = byte(r.Uint32())
		magicDataMeasureRequest6 = byte(r.Uint32())
		magicDataConfigureRequest2 = byte(r.Uint32())
		magicDataConfigureRequest4 = byte(r.Uint32())
		magicDataConfigureRequest5 = byte(r.Uint32())
		magicDataConfigureRequest6 = byte(r.Uint32())
		magicDataImportRequest4 = byte(r.Uint32())
		magicDataImportRequest5 = byte(r.Uint32())
		magicDataImportRequest6 = byte(r.Uint32())
		magicDataImportRequest7 = byte(r.Uint32())
		magicDataNextImportRequest1_4 = byte(r.Uint32())
		magicDataNextImportRequest1_5 = byte(r.Uint32())
		magicDataNextImportRequest1_6 = byte(r.Uint32())
	}
	// randomize all data
	if randomizeMagicDataConfig > 1 {
		magicDataMeasureRequest3 = byte(r.Uint32())
		magicDataMeasureRequest7 = byte(r.Uint32())
		magicDataConfigureRequest3 = byte(r.Uint32())
		magicDataConfigureRequest7 = byte(r.Uint32())
		magicDataImportRequest1 = byte(r.Uint32())
		magicDataImportRequest2 = byte(r.Uint32())
		magicDataImportRequest3 = byte(r.Uint32())
		magicDataNextImportRequest1_7 = byte(r.Uint32())
	}
}

func newConfigureRequest(config Config) []byte {
	cfg := getConfigurationFlags(config)
	return []byte{
		configureRequestMark, cfg, magicDataConfigureRequest2, magicDataConfigureRequest3,
		magicDataConfigureRequest4, magicDataConfigureRequest5, magicDataConfigureRequest6, magicDataConfigureRequest7,
	}
}

func newMeasureRequest() []byte {
	return []byte{
		measureRequestMark, magicDataMeasureRequest1, magicDataMeasureRequest2, magicDataMeasureRequest3,
		magicDataMeasureRequest4, magicDataMeasureRequest5, magicDataMeasureRequest6, magicDataMeasureRequest7,
	}
}

func newImportRequest() []byte {
	return []byte{
		importRequestMark, magicDataImportRequest1, magicDataImportRequest2, magicDataImportRequest3,
		magicDataImportRequest4, magicDataImportRequest5, magicDataImportRequest6, magicDataImportRequest7,
	}
}

func newNextImportDataRequest(first bool) []byte {
	if first {
		return []byte{
			nextImportDataRequestMark, magicDataNextImportRequest1, magicDataNextImportRequest1, magicDataNextImportRequest1,
			magicDataNextImportRequest1_4, magicDataNextImportRequest1_5, magicDataNextImportRequest1_6, magicDataNextImportRequest1_7,
		}
	}
	return []byte{
		nextImportDataRequestMark, magicDataNextImportRequest1, magicDataNextImportRequest1, magicDataNextImportRequest1,
		magicDataNextImportRequest2_4, magicDataNextImportRequest2_5, magicDataNextImportRequest2_6, magicDataNextImportRequest2_7,
	}
}

func getConfigurationFlags(config Config) byte {
	flags := config.SoundLevelRange
	if config.SoundLevelDisplayMode == SoundLevelDisplayModeMaxHold {
		flags |= displayModeFlags
	}
	if config.TimeWeight == TimeWeightFast {
		flags |= timeWeightFlags
	}
	if config.FrequencyWeight == FrequencyWeightC {
		flags |= frequencyWeightFlags
	}
	return byte(flags)
}

func parseData(data []byte) Event {
	if data[0] == configureResponseMark {
		return ConfiguredEvent{}
	}
	return parseMeasureData(data)
}

func parseMeasureData(data []byte) MeasuredEvent {
	db := getDecibelValue(data)
	config := parseConfigurationFlags(data[2])
	return MeasuredEvent{
		Time:       time.Now(),
		SoundLevel: db,
		Config:     config,
	}
}

func parseConfigurationFlags(flags byte) Config {
	config := Config{}
	if flags&timeWeightFlags != 0 {
		config.TimeWeight = TimeWeightFast
	} else {
		config.TimeWeight = TimeWeightSlow
	}
	if flags&displayModeFlags != 0 {
		config.SoundLevelDisplayMode = SoundLevelDisplayModeMaxHold
	} else {
		config.SoundLevelDisplayMode = SoundLevelDisplayModeNormal
	}
	if flags&frequencyWeightFlags != 0 {
		config.FrequencyWeight = FrequencyWeightC
	} else {
		config.FrequencyWeight = FrequencyWeightA
	}
	config.SoundLevelRange = SoundLevelRange(flags & levelRangeFlags)
	return config
}

func getDecibelValue(data []byte) (db float32) {
	v := uint16(data[0])
	v = v << 8
	v |= uint16(data[1])
	return float32(v) / 10
}

func getImportedDecibelValue(data []byte) (db float32) {
	v := uint16(data[0])<<4 | uint16(data[1])>>4
	return float32(v) / 20
}
