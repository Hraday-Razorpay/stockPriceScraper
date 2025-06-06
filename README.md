
# 📈 stockPriceScraper

A Go-based CLI tool that ethically scrapes stock prices and related data from [Yahoo Finance](https://finance.yahoo.com). The tool reads stock tickers from a plain text file and saves results to a structured CSV file that includes metadata.

## 🚀 Features

- Read stock tickers from a `.txt` file.
- Scrape real-time prices, price changes, and currency from Yahoo Finance.
- Support for multiple exchanges (e.g., NSE, BSE, LSE).
- Export results to a CSV file (`stocks.csv`) with metadata included.
- Metadata includes total stocks scraped and timestamp.

## 📂 Project Structure

```
stockPriceScraper/
├── main.go           # Main application logic
├── tickers.txt       # Input file containing list of stock tickers (one per line)
├── stocks.csv        # Output file with scraped stock data and metadata
├── go.mod            # Go module file
├── go.sum            # Go dependency checksums
```

## 🛠️ Prerequisites

- Go 1.16 or later
- Internet connection

## 📄 Input Format: `tickers.txt`

Each line in the `tickers.txt` file should contain one valid Yahoo Finance ticker:

```
RELIANCE.NS
TCS.NS
AAPL
GOOGL
```

## ▶️ Usage

```bash
go run main.go
```

- This reads tickers from `tickers.txt`
- Scrapes data from Yahoo Finance
- Exports results to `stocks.csv`

## ✅ Output file: `stocks.csv`


## ✅ Example Output

```bash
Wipro Limited (WIPRO.BO) | 248.60 | +(0.26%) | INR
E-Mini S&P 500 Jun 25 (ES=F) | 5,971.25 | +(0.42%) |
Microsoft Corporation (MSFT) | 467.68 | +(0.82%) | USD
```

The CSV file includes a header row, followed by stock data and a metadata section at the end.

```
company,price,change,currency
Infosys,₹1,450.00,+10.00,INR
AAPL,$189.50,-0.70,USD

Total Stocks Scraped: 2
Scraping Time: 2025-06-06 15:45:00
```

## 🧠 How It Works

1. Loads ticker symbols from `tickers.txt`.
2. For each ticker, fetches its Yahoo Finance page.
3. Parses HTML using GoQuery and Colly.
4. Extracts:
   - Company name
   - Current stock price
   - Price change
   - Currency
5. Writes all results and a summary metadata section to `stocks.csv`.

## 🧪 Libraries Used

- [Colly](https://github.com/gocolly/colly)
- [GoQuery](https://github.com/PuerkitoBio/goquery)
- Standard libraries: `encoding/csv`, `os`, `fmt`, `time`, etc.

## 🧾 License

MIT License

## 🙏 Ethical Scraping

This project is intended for educational and personal use only. It adheres to Yahoo Finance's terms of service. Avoid excessive requests or automated scraping at scale.
