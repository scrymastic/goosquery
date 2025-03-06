package curl

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Curl struct {
	URL           string `json:"url"`
	Method        string `json:"method"`
	UserAgent     string `json:"user_agent"`
	ResponseCode  int32  `json:"response_code"`
	RoundTripTime int64  `json:"round_trip_time"`
	Bytes         int64  `json:"bytes"`
	Result        string `json:"result"`
}

const defaultUserAgent = "goosquery"

func GenCurl(url string, userAgent string) (Curl, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Curl{}, fmt.Errorf("failed to create request: %v", err)
	}

	// Set user agent
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	req.Header.Set("User-Agent", userAgent)

	// Measure time and make request
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return Curl{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Calculate round trip time
	rtt := time.Since(startTime).Microseconds()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Curl{}, fmt.Errorf("failed to read response body: %v", err)
	}

	return Curl{
		URL:           url,
		Method:        "GET",
		UserAgent:     userAgent,
		ResponseCode:  int32(resp.StatusCode),
		RoundTripTime: rtt,
		Bytes:         int64(len(body)),
		Result:        string(body),
	}, nil
}
