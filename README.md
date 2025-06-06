
# 📈 stockPriceScraper

A Go-based CLI tool that ethically scrapes stock prices and related data from [Yahoo Finance](https://finance.yahoo.com), using the powerful Colly and GoQuery libraries. The tool reads a list of stock tickers from a CSV file and outputs each stock’s current price, change, and currency.

## 🚀 Features

- Scrapes real-time stock prices and changes from Yahoo Finance.
- Supports multiple stock exchanges (e.g., NSE, BSE, LSE).
- Detects and assigns the appropriate currency based on the ticker symbol.
- Simple and configurable CSV-based input.
- Lightweight and efficient scraping using [Colly](https://github.com/gocolly/colly).

## 📂 Project Structure

```
stockPriceScraper/
├── main.go           # Main application logic
├── go.mod            # Go module file
├── go.sum            # Go dependency checksums
└── stocks.csv        # Input file containing list of stock tickers
```

## 🛠️ Prerequisites

- Go 1.16 or later
- Internet connection (for scraping data)

## 📦 Installation

```bash
git clone https://github.com/yourusername/stockPriceScraper.git
cd stockPriceScraper
go mod tidy
```

## 📄 Input Format

The tool expects a `stocks.csv` file in the following format:

```
Ticker
RELIANCE.NS
TCS.NS
AAPL
GOOGL
```

> You can modify the `stocks.csv` to include any ticker symbols supported by Yahoo Finance.

## ▶️ Usage

To run the scraper:

```bash
go run main.go
```

The output will be printed in the terminal and can include:

- Company name
- Current stock price
- Price change
- Currency (e.g., INR, USD, GBP)

## ✅ Example Output

```bash
RELIANCE.NS | ₹2,450.00 | +15.00 (INR)
AAPL        | $189.50   | -0.70 (USD)
```

## 🧪 Libraries Used

- [Colly](https://github.com/gocolly/colly): High-level scraping framework for Go.
- [GoQuery](https://github.com/PuerkitoBio/goquery): jQuery-style HTML manipulation.
- Standard Go libraries: `encoding/csv`, `net/url`, `strings`, etc.

## 🧾 License

This project is licensed under the MIT License.

## 🙏 Ethical Scraping

This tool respects the terms of use of Yahoo Finance. Use responsibly and avoid frequent or large-scale scraping which may violate their policies or disrupt their services.
