package gm1356

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type importer struct {
	dataBuffer  chan byte
	eventBuffer chan Event
	isImporting bool
	mu          sync.Mutex
}

func newImporter(eventBuffer chan Event) *importer {
	return &importer{
		dataBuffer:  make(chan byte, 17),
		eventBuffer: eventBuffer,
	}
}

func (i *importer) Start() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.isImporting {
		return errors.New("already importing")
	}
	i.isImporting = true
	go i.handleImportData()
	return nil
}

func (i *importer) Write(data []byte) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if !i.isImporting {
		return errors.New("not importing")
	}
	for _, d := range data {
		i.dataBuffer <- d
	}
	return nil
}

func (i *importer) IsImporting() bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.isImporting
}

func (i *importer) handleImportData() {
	data := make([]byte, 0, 17)
	var config Config
	var dataIndex int
	var startTime time.Time
	var sampleInterval time.Duration
	var skipHalfByte bool
	var skipLeft int
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		select {
		case <-ctx.Done():
			i.mu.Lock()
			defer i.mu.Unlock()
			cancel()
			i.isImporting = false
			return
		case d := <-i.dataBuffer:
			if skipLeft > 0 {
				skipLeft--
				continue
			}
			data = append(data, d)
			alignedData := getAlignedData(data, skipHalfByte)
			switch alignedData[0] {
			case nextImportBlockMark:
				if len(alignedData) < 9 {
					// needs more than 9 bytes
					continue
				}
				dataIndex = 0
				config, startTime, sampleInterval = parseImportBlockHeader(alignedData)
				data = data[9:]
			case importResponseMark:
				data = data[1:]
				skipLeft = 7
			case noImportDataMark:
				data = data[1:]
			default:
				if len(alignedData) < 2 {
					// needs more than 2 bytes (12bits)
					continue
				}
				// add event
				i.eventBuffer <- MeasuredEvent{
					Time:       startTime.Add(sampleInterval * time.Duration(dataIndex)),
					SoundLevel: getImportedDecibelValue(alignedData),
					Config:     config,
				}

				data = data[1:]
				if skipHalfByte {
					data = data[1:]
				}
				dataIndex++
				skipHalfByte = !skipHalfByte
			}
		}
		cancel()
	}
}

func getAlignedData(data []byte, skipHalfByte bool) []byte {
	if skipHalfByte {
		d := make([]byte, len(data)-1)
		for i := range d {
			d[i] = data[i]<<4 | data[i+1]>>4
		}
		return d
	}
	return data
}

func parseImportBlockHeader(data []byte) (config Config, startTime time.Time, interval time.Duration) {
	config = parseConfigurationFlags(data[1])
	startTime = parseRecordedTime(data[2:8])
	interval = time.Duration(data[8]) * time.Second
	return
}

func parseRecordedTime(data []byte) time.Time {
	sec, _ := strconv.ParseInt(strconv.FormatInt(int64(data[0]), 16), 10, 64)
	min, _ := strconv.ParseInt(strconv.FormatInt(int64(data[1]), 16), 10, 64)
	hour, _ := strconv.ParseInt(strconv.FormatInt(int64(data[2]), 16), 10, 64)
	day, _ := strconv.ParseInt(strconv.FormatInt(int64(data[3]), 16), 10, 64)
	month, _ := strconv.ParseInt(strconv.FormatInt(int64(data[4]), 16), 10, 64)
	year, _ := strconv.ParseInt(strconv.FormatInt(int64(data[5]), 16), 10, 64)
	str := fmt.Sprintf("%02d/%02d/%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
	t, err := time.Parse("06/01/02 15:04:05", str)
	if err != nil {
		panic(err)
	}
	return t
}
