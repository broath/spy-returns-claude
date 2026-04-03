package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type YahooResponse struct {
	Chart struct {
		Result []struct {
			Timestamps []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close []float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

type MonthRow struct {
	Month  string  `json:"month"`
	Close  float64 `json:"close"`
	Return *float64 `json:"return"`
}

func main() {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/SPY?interval=1mo&range=5y"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data YahooResponse
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	if len(data.Chart.Result) == 0 {
		fmt.Println("No data returned")
		return
	}

	result := data.Chart.Result[0]
	timestamps := result.Timestamps
	closes := result.Indicators.Quote[0].Close

	// Deduplicate by month, build ordered rows
	seen := make(map[string]bool)
	var rows []MonthRow
	var dedupedCloses []float64

	for i, ts := range timestamps {
		close := closes[i]
		if close == 0 {
			continue
		}
		month := time.Unix(ts, 0).UTC().Format("2006-01")
		if seen[month] {
			continue
		}
		seen[month] = true
		rows = append(rows, MonthRow{Month: month, Close: close})
		dedupedCloses = append(dedupedCloses, close)
	}

	// Compute returns
	for i := range rows {
		if i == 0 {
			continue
		}
		prev := dedupedCloses[i-1]
		if prev == 0 {
			continue
		}
		ret := (rows[i].Close - prev) / prev * 100
		rows[i].Return = &ret
	}

	// Reverse to descending order
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}

	// Write data.json
	out, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("data.json", out, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("Wrote %d months to data.json\n", len(rows))

	// Also print table to stdout
	fmt.Printf("\n%-12s  %-10s  %-12s\n", "Month", "Close", "Return")
	fmt.Println("--------------------------------------")
	for _, row := range rows {
		if row.Return == nil {
			fmt.Printf("%-12s  %10.2f  %12s\n", row.Month, row.Close, "—")
		} else {
			fmt.Printf("%-12s  %10.2f  %+11.2f%%\n", row.Month, row.Close, *row.Return)
		}
	}
}
