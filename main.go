package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

// 24.01.2023 modified getUsers to use http.Client as parameter to be able to mock it in tests
func getUsers(client *http.Client) ([]User, error) {
	url := "https://fakestoreapi.com/users"
	if client != http.DefaultClient {
		url = "http://fakestoreapi.com/users"
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	return users, nil
}

func getProducts(client *http.Client) ([]Product, error) {
	url := "https://fakestoreapi.com/products"
	if client != http.DefaultClient {
		url = "http://fakestoreapi.com/users"
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var products []Product
	err = json.Unmarshal(body, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func getCarts() ([]Cart, error) {
	resp, err := http.Get("https://fakestoreapi.com/carts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var carts []Cart
	err = json.Unmarshal(body, &carts)
	if err != nil {
		return nil, err
	}

	return carts, nil
}

// 2. Total value of category values
func getCategoryValues(products []Product) []CategoryValue {
	categoryValues := make(map[string]float64)

	for _, product := range products {
		categoryValues[product.Category] += product.Price
	}

	var result []CategoryValue
	for category, value := range categoryValues {
		result = append(result, CategoryValue{
			Category: category,
			Value:    value,
		})
	}

	return result
}

// 3. highest value cart

func findHighestValueCart(carts []Cart, products []Product) (Cart, float64, error) {
	highestValue := 0.0
	var highestValueCart Cart
	productMap := make(map[int]Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	for _, cart := range carts {
		var cartValue float64
		for _, cartProduct := range cart.CartsProducts {
			product := productMap[cartProduct.ProductID]
			cartValue += product.Price * cartProduct.Quantity
		}

		if cartValue > highestValue {
			highestValue = cartValue
			highestValueCart = cart
		}
	}

	return highestValueCart, highestValue, nil
}

// 4. Find two users with biggest distance between them
func getUserLocations(users []User) []GeoLocation {
	locations := make([]GeoLocation, len(users))

	for i, user := range users {
		locations[i] = GeoLocation{
			Lat:  user.Address.GeoLocation.Lat,
			Long: user.Address.GeoLocation.Long,
		}
	}

	return locations
}

func Distance(location1, location2 GeoLocation) float64 {
	const earthRadius = 6371 // in kilometers
	lat1 := toRadians(location1.Lat)
	lat2 := toRadians(location2.Lat)
	deltaLat := toRadians(location2.Lat - location1.Lat)
	deltaLong := toRadians(location2.Long - location1.Long)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLong/2)*math.Sin(deltaLong/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c
	return distance
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func maxDistance(users []User) (float64, int, int) {
	maxDistance := 0.0
	user1id := 0
	user2id := 0
	for i := 0; i < len(users); i++ {
		for j := i + 1; j < len(users); j++ {
			distance := Distance(users[i].Address.GeoLocation, users[j].Address.GeoLocation)
			if distance > maxDistance {
				maxDistance = distance
				user1id = i
				user2id = j
			}
		}
	}
	return maxDistance, user1id, user2id
}

func main() {
	start := time.Now()
	var users []User
	var products []Product
	var carts []Cart
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(3)
	//using goroutines to get data from 3 endpoints improved the performance by ~ 50%
	go func() {
		defer wg.Done()
		var err error
		users, err = getUsers(http.DefaultClient)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		products, err = getProducts(http.DefaultClient)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		carts, err = getCarts()
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		log.Fatal(err)
	}

	categoryValues := getCategoryValues(products)
	fmt.Println("Category values:", categoryValues)

	highestValueCart, highestValue, _ := findHighestValueCart(carts, products)
	highestValueOwner := users[highestValueCart.UserID-1]
	fmt.Printf("Highest value cart id: %v (Value: %.2f, Owner name: %s Owner username: %s)\n", highestValueCart.ID, highestValue, highestValueOwner.Name, highestValueOwner.Username)

	maxDistance, user1, user2 := maxDistance(users)
	fmt.Printf("Biggest distance = %.2f between %s from %s and %s from %s\n", maxDistance, users[user1].Name, users[user1].Address.City, users[user2].Name, users[user2].Address.City)

	fmt.Println("Time taken:", time.Since(start))
}
