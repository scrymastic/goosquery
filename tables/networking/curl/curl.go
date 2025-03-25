package curl

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

const defaultUserAgent = "goosquery"

// genCurl performs an HTTP request to the specified URL and returns information about the request
func genCurl(ctx *sqlctx.Context, url string, userAgent string) (*result.Result, error) {
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
	result := result.NewResult(ctx, Schema)

	// Set URL and method values (overriding defaults)
	result.Set("url", url)
	result.Set("method", "GET")
	result.Set("user_agent", userAgent)

	// Measure time and make request
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Calculate round trip time
	rtt := time.Since(startTime).Microseconds()

	result.Set("round_trip_time", rtt)

	result.Set("response_code", int32(resp.StatusCode))

	// Read response body only if needed
	if ctx.IsAnyOfColumnsUsed([]string{"bytes", "result"}) {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		result.Set("bytes", int64(len(body)))
		result.Set("result", string(body))
	}

	return result, nil
}

// GenCurl performs HTTP requests for the given URLs in the context
func GenCurl(ctx *sqlctx.Context) (*result.Results, error) {
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

	results := result.NewQueryResult()

	// Execute each URL request
	for _, url := range urls {
		result, err := genCurl(ctx, url, userAgent)
		if err != nil {
			return nil, fmt.Errorf("failed to generate curl result for %s: %v", url, err)
		}
		results.AppendResult(*result)
	}

	return results, nil
}
