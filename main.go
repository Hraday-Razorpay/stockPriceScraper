package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change, currency, timestamp string
}

func getCurrencyByTicker(ticker string) string {
	// URL decode the ticker first
	decodedTicker, _ := url.QueryUnescape(ticker)
	fmt.Printf("DEBUG: Original ticker: '%s', Decoded: '%s'\n", ticker, decodedTicker)

	// Check for %, =, or ^ anywhere in the ticker
	if strings.Contains(ticker, "%") ||
		strings.Contains(decodedTicker, "=") ||
		strings.Contains(decodedTicker, "^") {
		fmt.Printf("DEBUG: Found index pattern, returning empty currency\n")
		return "" // No currency for indices/futures
	}

	// Stock currencies by exchange
	if strings.Contains(ticker, ".BO") || strings.Contains(ticker, ".NS") {
		fmt.Printf("DEBUG: Found Indian exchange, returning INR\n")
		return "INR"
	} else if strings.Contains(ticker, ".L") {
		return "GBP"
	} else if strings.Contains(ticker, ".TO") {
		return "CAD"
	}

	fmt.Printf("DEBUG: No special case found, returning USD\n")
	return "USD"
}

func extractTickerFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		ticker := parts[len(parts)-1]
		if ticker == "" && len(parts) > 1 {
			ticker = parts[len(parts)-2]
		}
		return ticker
	}
	return ""
}

func readTickersFromTxt(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tickers []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ticker := strings.TrimSpace(scanner.Text())
		if ticker != "" && !strings.HasPrefix(ticker, "#") { // Skip empty lines and comments
			tickers = append(tickers, ticker)
		}
	}

	return tickers, scanner.Err()
}

func main() {
	ticker, err := readTickersFromTxt("tickers.txt")
	if err != nil {
		log.Fatalf("Error reading tickers: %v", err)
	}

	fmt.Printf("Loaded %d tickers from file\n", len(ticker))

	stocks := []Stock{}
	var allStocks []Stock

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div[data-testid='quote-header']", func(e *colly.HTMLElement) {
		stock := Stock{}

		companyName := e.ChildText("h1")
		price := e.ChildText("[data-testid='qsp-price']")
		if price == "" {
			price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")
		}

		changePercent := e.ChildText("[data-testid='qsp-price-change-percent']")
		if changePercent == "" {
			changePercent = e.ChildText("fin-streamer[data-field='regularMarketChangePercent']")
		}

		extractedTicker := extractTickerFromURL(e.Request.URL.String())
		currency := getCurrencyByTicker(extractedTicker)

		if companyName != "" && price != "" {
			stock.company = companyName
			stock.price = price
			stock.change = changePercent
			stock.currency = currency

			if currency == "" {
				fmt.Printf("Company: %s\n", stock.company)
				fmt.Printf("Price: %s\n", stock.price)
				fmt.Printf("Change: %s\n", stock.change)
			} else {
				fmt.Printf("Company: %s\n", stock.company)
				fmt.Printf("Price: %s %s\n", stock.price, stock.currency)
				fmt.Printf("Change: %s\n", stock.change)
			}
			fmt.Println("---")

			currentTime := time.Now().Format("2006-01-02 15:04:05")
			stock.timestamp = currentTime

			stocks = append(stocks, stock)
		}
	})

	c.OnHTML("[data-testid='qsp-price']", func(e *colly.HTMLElement) {
		if len(stocks) > 0 {
			return
		}

		price := e.Text
		if price != "" {
			stock := Stock{}

			companyName := e.DOM.ParentsUntil("body").Find("h1").First().Text()
			changePercent := e.DOM.ParentsUntil("body").Find("[data-testid='qsp-price-change-percent']").First().Text()

			extractedTicker := extractTickerFromURL(e.Request.URL.String())
			currency := getCurrencyByTicker(extractedTicker)

			stock.company = companyName
			stock.price = price
			stock.change = changePercent
			stock.currency = currency

			fmt.Printf("Found via qsp-price:\n")
			if currency == "" {
				fmt.Printf("Company: %s\n", stock.company)
				fmt.Printf("Price: %s\n", stock.price)
				fmt.Printf("Change: %s\n", stock.change)
			} else {
				fmt.Printf("Company: %s\n", stock.company)
				fmt.Printf("Price: %s %s\n", stock.price, stock.currency)
				fmt.Printf("Change: %s\n", stock.change)
			}
			fmt.Println("---")

			currentTime := time.Now().Format("2006-01-02 15:04:05")
			stock.timestamp = currentTime

			stocks = append(stocks, stock)
		}
	})

	c.OnHTML("fin-streamer[data-field='regularMarketPrice']", func(e *colly.HTMLElement) {
		if len(stocks) > 0 {
			return
		}

		price := e.Text
		if price == "" {
			if val := e.Attr("value"); val != "" {
				price = val
			}
		}

		if price != "" {
			stock := Stock{}

			companyName := ""
			e.DOM.ParentsUntil("body").Each(func(i int, s *goquery.Selection) {
				if companyName == "" {
					companyName = s.Find("h1").Text()
				}
			})

			changePercent := ""
			e.DOM.ParentsUntil("body").Find("fin-streamer[data-field='regularMarketChangePercent']").Each(func(i int, s *goquery.Selection) {
				if changePercent == "" {
					changePercent = s.Text()
					if changePercent == "" {
						if val, exists := s.Attr("value"); exists {
							changePercent = val
						}
					}
				}
			})

			extractedTicker := extractTickerFromURL(e.Request.URL.String())
			currency := getCurrencyByTicker(extractedTicker)

			stock.company = companyName
			stock.price = price
			stock.change = changePercent
			stock.currency = currency

			fmt.Printf("Found via fin-streamer:\n")
			if currency == "" {
				fmt.Printf("Company: %s\n", stock.company)
				fmt.Printf("Price: %s\n", stock.price)
				fmt.Printf("Change: %s\n", stock.change)
			} else {
				fmt.Printf("Company: %s\n", stock.company)
				fmt.Printf("Price: %s %s\n", stock.price, stock.currency)
				fmt.Printf("Change: %s\n", stock.change)
			}
			fmt.Println("---")

			currentTime := time.Now().Format("2006-01-02 15:04:05")
			stock.timestamp = currentTime

			stocks = append(stocks, stock)
		}
	})

	c.OnHTML("*", func(e *colly.HTMLElement) {
		if len(stocks) == 0 {
			if e.Attr("data-testid") != "" || e.Attr("data-field") != "" {
				fmt.Printf("Found element: %s with data-testid='%s' data-field='%s' text='%s'\n",
					e.Name, e.Attr("data-testid"), e.Attr("data-field"), e.Text)
			}
		}
	})

	for _, t := range ticker {
		stocks = []Stock{}

		time.Sleep(3 * time.Second)

		fmt.Printf("\n=== Scraping %s ===\n", t)
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")

		time.Sleep(2 * time.Second)

		if len(stocks) > 0 {
			fmt.Printf("Successfully scraped %s\n", t)
			allStocks = append(allStocks, stocks...)
		} else {
			fmt.Printf("No data found for %s\n", t)
		}
	}

	fmt.Printf("\nTotal stocks scraped: %d\n", len(allStocks))

	if len(allStocks) == 0 {
		log.Println("No stock data was scraped.")
		log.Println("Try manually visiting https://finance.yahoo.com/quote/WIPRO.BO/ to check if the page loads correctly.")
		log.Println("The page structure might have changed or the site might be blocking requests.")
		return
	}

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"company",
		"price",
		"change",
		"currency",
	}

	writer.Write(headers)

	// Write all stock data
	for _, stock := range allStocks {
		record := []string{
			stock.company,
			stock.price,
			stock.change,
			stock.currency,
		}
		writer.Write(record)
	}

	// Add metadata section at the end
	writer.Write([]string{}) // Empty row for separation
	writer.Write([]string{"METADATA", "", "", ""})
	writer.Write([]string{"Last Updated", time.Now().Format("2006-01-02 15:04:05"), "", ""})
	writer.Write([]string{"Last Updated (UTC)", time.Now().UTC().Format("2006-01-02 15:04:05"), "", ""})
	writer.Write([]string{"Total Stocks Scraped", fmt.Sprintf("%d", len(allStocks)), "", ""})
	writer.Write([]string{"Total Tickers Processed", fmt.Sprintf("%d", len(ticker)), "", ""})
	writer.Write([]string{"Success Rate", fmt.Sprintf("%.1f%%", float64(len(allStocks))/float64(len(ticker))*100), "", ""})

	fmt.Println("Data saved to stocks.csv with metadata")
}
