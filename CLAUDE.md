# SPY Monthly Returns

A Go tool that fetches SPY historical price data from the Yahoo Finance API, calculates month-over-month returns, and writes the results to `data.json` for display in a custom HTML dashboard.

## Purpose

Fetch SPY monthly OHLCV data from the Yahoo Finance JSON API, compute returns on close price, and render the results in a production-grade HTML page.

## Architecture

- **`main.go`** — calls Yahoo Finance API, deduplicates by month, computes returns, writes `data.json`
- **`index.html`** — frontend dashboard that fetches `data.json` and renders an animated table with inline bar charts

## Commands

```bash
# Fetch latest data and write data.json
go run main.go

# Build binary
go build -o spy-returns .

# Open dashboard (WSL)
explorer.exe index.html
```

## Data Source

Yahoo Finance chart API — no API key required:
```
https://query1.finance.yahoo.com/v8/finance/chart/SPY?interval=1mo&range=5y
```
- Returns JSON with `timestamps[]` and `indicators.quote[0].close[]`
- Requires a `User-Agent` header (browser-style) to avoid 429s

## Stack

- **Language**: Go
- **HTTP**: stdlib `net/http`
- **Frontend**: Vanilla HTML/CSS/JS (no build step)
- **Fonts**: Cormorant Garamond + JetBrains Mono (Google Fonts)
- **Skills**: `web-scraper` (used for initial API reconnaissance), `frontend-design` (UI)

## Conventions

- Deduplicate months by UTC year-month key (first entry per month wins)
- Skip any close price of 0 (incomplete candle)
- Return formula: `(close[i] - close[i-1]) / close[i-1] * 100`
- First month has no return (displayed as `—`)
- Data sorted descending by date (most recent first)
