package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var urls = []string{
	"https://www.amazon.com/Intel-i9-10900K-Desktop-Processor-Unlocked/dp/B086MHSTVD",
	"https://www.amazon.com/AMD-Ryzen-5600X-12-Thread-Processor/dp/B08166SLDF",
	"https://www.amazon.com/EVGA-GeForce-12G-P5-3967-KR-Technology-Backplate/dp/B09622N253",
	"https://www.amazon.com/SAMSUNG-Border-Less-TUV-Certified-Intelligent-LS32A700NWNXZA/dp/B08V6MNW9P/",
}

type Product struct {
	name  string
	price string
}

func main() {
	t := time.Now()
	var products []Product

	var wg sync.WaitGroup
	ch := make(chan Product)

	wg.Add(len(urls))

	for _, url := range urls {
		go scrape(url, ch)
	}

	for range urls {
		go func() {
			defer wg.Done()
			product := <-ch
			products = append(products, product)
		}()
	}

	wg.Wait()

	file, err := os.Create("data.csv")

	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"name", "price"})

	for _, product := range products {
		writer.Write([]string{
			product.name,
			product.price,
		})
	}

	close(ch)

	elapsed := time.Since(t).Seconds()
	fmt.Printf("Time: %.2fs\n", elapsed)
}

func scrape(url string, ch chan<- Product) {
	var product Product
	c := colly.NewCollector()

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})

	c.OnHTML("#ppd", func(e *colly.HTMLElement) {
		name := e.ChildText("#centerCol #productTitle")
		price := e.ChildText("#rightCol #price_inside_buybox")

		product = Product{name, price}
	})

	c.Visit(url)

	ch <- product
}
