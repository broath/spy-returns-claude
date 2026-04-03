# SPY Monthly Returns

This project is purely used to tinker with Claude Code.

Fetches SPY historical price data from the Yahoo Finance API, calculates month-over-month returns on the closing price, and displays the results in a dark-themed HTML dashboard.

## Usage

```bash
# Fetch latest data and write data.json
go run main.go
```

## Output

- Prints a formatted table to stdout
- Writes `data.json` (consumed by `index.html`)
