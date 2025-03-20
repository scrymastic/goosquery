package curl

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/util"
)

const defaultUserAgent = "goosquery"

// Column definitions for the curl table
var columnDefs = map[string]string{
	"url":             "string",
	"method":          "string",
	"user_agent":      "string",
	"round_trip_time": "int64",
	"response_code":   "int32",
	"bytes":           "int64",
	"result":          "string",
}

// genCurl performs an HTTP request to the specified URL and returns information about the request
func genCurl(ctx context.Context, url string, userAgent string) (map[string]interface{}, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set user agent
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	req.Header.Set("User-Agent", userAgent)

	// Initialize result map with default values for all requested columns
	result := util.InitColumns(ctx, columnDefs)

	// Set URL and method values (overriding defaults)
	if ctx.IsColumnUsed("url") {
		result["url"] = url
	}

	if ctx.IsColumnUsed("method") {
		result["method"] = "GET"
	}

	if ctx.IsColumnUsed("user_agent") {
		result["user_agent"] = userAgent
	}

	// Measure time and make request
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Calculate round trip time
	rtt := time.Since(startTime).Microseconds()

	if ctx.IsColumnUsed("round_trip_time") {
		result["round_trip_time"] = rtt
	}

	if ctx.IsColumnUsed("response_code") {
		result["response_code"] = int32(resp.StatusCode)
	}

	// Read response body only if needed
	if ctx.IsAnyOfColumnsUsed([]string{"bytes", "result"}) {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return result, fmt.Errorf("failed to read response body: %v", err)
		}

		if ctx.IsColumnUsed("bytes") {
			result["bytes"] = int64(len(body))
		}

		if ctx.IsColumnUsed("result") {
			result["result"] = string(body)
		}
	}

	return result, nil
}

// GenCurl performs HTTP requests for the given URLs in the context
func GenCurl(ctx context.Context) ([]map[string]interface{}, error) {
	// Get URLs from context constants
	urls := ctx.GetConstants("url")
	if len(urls) == 0 {
		return nil, fmt.Errorf("no URL provided")
	}

	// Get user agent (use first one if multiple are provided)
	userAgents := ctx.GetConstants("user_agent")
	userAgent := ""
	if len(userAgents) > 0 {
		userAgent = userAgents[0]
	}

	results := make([]map[string]interface{}, 0, len(urls))

	// Execute each URL request
	for _, url := range urls {
		result, err := genCurl(ctx, url, userAgent)
		if err != nil {
			return nil, fmt.Errorf("failed to generate curl result for %s: %v", url, err)
		}
		results = append(results, result)
	}

	return results, nil
}
