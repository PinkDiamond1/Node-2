package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type StatsInfoData struct {
	Type     string
	FileSize int64
}

type StatType int

const (
	Upload StatType = iota
	Download
	Delete
)

var (
	SendLogs      = true
	loggerAddress = "http://68.183.215.241:9091"
)

// ====================================================================================

func Log(msg interface{}) {
	if !SendLogs {
		return
	}

	currentTime := time.Now().Local()

	logMsg := fmt.Sprintf("%s: %v\n", currentTime.String(), msg)

	fmt.Println(logMsg)

	errType := fmt.Sprintf("%T", msg)

	if errType == "*errors.errorString" || errType == "*fmt.wrapError" {
		url := loggerAddress + "/logs"

		req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte(logMsg)))

		client := &http.Client{Timeout: time.Minute}

		client.Do(req)
	}
}

// ====================================================================================

//Сreates an informative error with line.
func CreateDetails(location string, errMsg error) error {
	_, _, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s line %d -> %w", location, line, errMsg)
}

func SendStatistic(spAddress string, statType StatType, fileSize int64) {
	const location = "logger.SendStatistic->"
	url := loggerAddress + "/stats/" + spAddress
	body := StatsInfoData{
		Type:     statType.String(),
		FileSize: fileSize,
	}

	js, err := json.Marshal(body)
	if err != nil {
		Log(CreateDetails(location, err))
	}

	r := bytes.NewReader(js)
	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		Log(CreateDetails(location, err))
	}

	if resp != nil {
		resp.Body.Close()
	}
}

func (s StatType) String() string {
	switch s {
	case Upload:
		return "Upload"
	case Download:
		return "Download"
	case Delete:
		return "Delete"
	default:
		return "Upload"
	}
}
