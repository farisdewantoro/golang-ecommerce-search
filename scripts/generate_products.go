package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"golang-ecommerce-search/internal/domain"
)

var (
	categories = []string{
		"Electronics", "Clothing", "Books", "Home & Kitchen", "Sports",
		"Beauty", "Toys", "Automotive", "Health", "Garden",
		"Pet Supplies", "Jewelry", "Tools", "Office Products", "Food",
		"Baby", "Movies", "Music", "Video Games", "Arts & Crafts",
	}

	brands = []string{
		"Apple", "Samsung", "Nike", "Adidas", "Sony",
		"Microsoft", "LG", "Dell", "HP", "Asus",
		"Lenovo", "Acer", "Canon", "Nikon", "Panasonic",
		"Philips", "Bosch", "Siemens", "Whirlpool", "GE",
		"KitchenAid", "Cuisinart", "Dyson", "Bose", "JBL",
		"Beats", "Sennheiser", "Ray-Ban", "Oakley", "Gucci",
		"Louis Vuitton", "Prada", "Chanel", "Dior", "Hermes",
		"Rolex", "Omega", "Cartier", "Tiffany", "Swarovski",
		"Levi's", "Gap", "H&M", "Zara", "Uniqlo",
		"Under Armour", "Puma", "Reebok", "New Balance", "Converse",
		"Vans", "Timberland", "Dr. Martens", "Clarks", "Skechers",
		"Calvin Klein", "Tommy Hilfiger", "Ralph Lauren", "Hugo Boss", "Armani",
		"Versace", "Dolce & Gabbana", "Burberry", "Coach", "Michael Kors",
		"Kate Spade", "Fossil", "Guess", "Lacoste", "The North Face",
		"Columbia", "Patagonia", "Marmot", "Arc'teryx", "Salomon",
		"Merrell", "Keen", "Teva", "Crocs", "Birkenstock",
		"Ecco", "Clarks", "Rockport", "Sperry", "Toms",
		"Vans", "Converse", "Nike", "Adidas", "Puma",
		"Reebok", "New Balance", "Asics", "Brooks", "Saucony",
	}

	adjectives = []string{
		"Amazing", "Premium", "Deluxe", "Professional", "Ultimate",
		"Super", "Mega", "Ultra", "Pro", "Elite",
		"Classic", "Modern", "Vintage", "Contemporary", "Traditional",
		"Luxury", "Exclusive", "Limited", "Special", "Unique",
		"Essential", "Basic", "Standard", "Advanced", "Expert",
	}

	nouns = []string{
		"Product", "Item", "Goods", "Merchandise", "Article",
		"Device", "Tool", "Equipment", "Gear", "Appliance",
		"Gadget", "Widget", "Thing", "Object", "Piece",
		"Unit", "Component", "Part", "Element", "Module",
	}
)

func generateProduct() *domain.Product {
	now := time.Now()
	return &domain.Product{
		Name:        fmt.Sprintf("%s %s %d", adjectives[rand.Intn(len(adjectives))], nouns[rand.Intn(len(nouns))], rand.Intn(1000)),
		Description: fmt.Sprintf("This is a %s %s with amazing features and quality.", adjectives[rand.Intn(len(adjectives))], nouns[rand.Intn(len(nouns))]),
		Price:       float64(rand.Intn(1000000)) / 100.0,
		Category:    categories[rand.Intn(len(categories))],
		Tags:        []string{categories[rand.Intn(len(categories))], categories[rand.Intn(len(categories))]},
		Brand:       brands[rand.Intn(len(brands))],
		Views:       rand.Int63n(10000),
		Buys:        rand.Int63n(1000),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func sendRequest(product *domain.Product, wg *sync.WaitGroup, successChan chan<- bool, errorChan chan<- error) {
	defer wg.Done()

	jsonData, err := json.Marshal(product)
	if err != nil {
		errorChan <- fmt.Errorf("error marshaling product: %v", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/products", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		errorChan <- fmt.Errorf("error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errorChan <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}

	successChan <- true
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const totalRequests = 100_000_000
	const concurrentRequests = 10

	var wg sync.WaitGroup
	successChan := make(chan bool, totalRequests)
	errorChan := make(chan error, totalRequests)

	startTime := time.Now()
	successCount := 0
	errorCount := 0

	fmt.Printf("Starting load test with %d total requests and %d concurrent requests\n", totalRequests, concurrentRequests)

	for i := 0; i < totalRequests; i += concurrentRequests {
		batchSize := concurrentRequests
		if i+batchSize > totalRequests {
			batchSize = totalRequests - i
		}

		for j := 0; j < batchSize; j++ {
			wg.Add(1)
			go sendRequest(generateProduct(), &wg, successChan, errorChan)
		}

		// Process results for this batch
		for j := 0; j < batchSize; j++ {
			select {
			case <-successChan:
				successCount++
			case err := <-errorChan:
				errorCount++
				fmt.Printf("Error: %v\n", err)
			}
		}

		// Print progress
		progress := float64(i+batchSize) / float64(totalRequests) * 100
		fmt.Printf("Progress: %.2f%% (Success: %d, Errors: %d)\n", progress, successCount, errorCount)
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("\nLoad test completed in %v\n", duration)
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Successful requests: %d\n", successCount)
	fmt.Printf("Failed requests: %d\n", errorCount)
	fmt.Printf("Requests per second: %.2f\n", float64(totalRequests)/duration.Seconds())
}
