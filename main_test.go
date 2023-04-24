package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetUsers(t *testing.T) {
	// Define a mock the response based on first two users from the API
	mockResponse := `[
		{
			"__v": 0,
			"address": {
				"city": "kilcoole",
				"geolocation": {
					"lat": "-37.3159",
					"long": "81.1496"
				},
				"number": 7682,
				"street": "new road",
				"zipcode": "12926-3874"
			},
			"email": "john@gmail.com",
			"id": 1,
			"name": {
				"firstname": "john",
				"lastname": "doe"
			},
			"password": "m38rmF$",
			"phone": "1-570-236-7033",
			"username": "johnd"
		},
		{
			"__v": 0,
			"address": {
				"city": "kilcoole",
				"geolocation": {
					"lat": "-37.3159",
					"long": "81.1496"
				},
				"number": 7267,
				"street": "Lovers Ln",
				"zipcode": "12926-3874"
			},
			"email": "morrison@gmail.com",
			"id": 2,
			"name": {
				"firstname": "david",
				"lastname": "morrison"
			},
			"password": "83r5^_",
			"phone": "1-570-236-7033",
			"username": "mor_2314"
		}]`

	// Create a mock HTTPS server
	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer mockServer.Close()
	// fix for error with url (parse the mock server URL)
	mockServerURL, _ := url.Parse(mockServer.URL)

	// Create a custom HTTP client with the server's certificate
	mockClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: mockServer.Client().Transport.(*http.Transport).TLSClientConfig,
			Proxy:           http.ProxyURL(mockServerURL),
		},
	}

	// Call getUsers with the custom HTTP client
	users, err := getUsers(mockClient)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the length of the users slice
	expectedLen := 2
	if len(users) != expectedLen {
		t.Errorf("Expected %d users, got %d", expectedLen, len(users))
	}

	// todo more chcecks
}

//test for getProduct

func TestGetProducts(t *testing.T) {
	// Define a mock the response based on first two users from the API
	mockResponse := `[
		{
		  "id": 1,
		  "title": "Fjallraven - Foldsack No. 1 Backpack, Fits 15 Laptops",
		  "price": 109.95,
		  "description": "Your perfect pack for everyday use and walks in the forest. Stash your laptop (up to 15 inches) in the padded sleeve, your everyday",
		  "category": "men's clothing",
		  "image": "https://fakestoreapi.com/img/81fPKd-2AYL._AC_SL1500_.jpg",
		  "rating": {
			"rate": 3.9,
			"count": 120
		  }
		},
		{
		  "id": 2,
		  "title": "Mens Casual Premium Slim Fit T-Shirts ",
		  "price": 22.3,
		  "description": "Slim-fitting style, contrast raglan long sleeve, three-button henley placket, light weight & soft fabric for breathable and comfortable wearing. And Solid stitched shirts with round neck made for durability and a great fit for casual fashion wear and diehard baseball fans. The Henley style round neckline includes a three-button placket.",
		  "category": "men's clothing",
		  "image": "https://fakestoreapi.com/img/71-3HjGNDUL._AC_SY879._SX._UX._SY._UY_.jpg",
		  "rating": {
			"rate": 4.1,
			"count": 259
		  }
		}]`

	// Create a mock HTTPS server
	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer mockServer.Close()
	// fix for error with url (parse the mock server URL)
	mockServerURL, _ := url.Parse(mockServer.URL)

	// Create a custom HTTP client with the server's certificate
	mockClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: mockServer.Client().Transport.(*http.Transport).TLSClientConfig,
			Proxy:           http.ProxyURL(mockServerURL),
		},
	}

	// Call getProducts with the custom HTTP client
	products, err := getProducts(mockClient)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the length of the products slice
	expectedLen := 2
	if len(products) != expectedLen {
		t.Errorf("Expected %d products, got %d", expectedLen, len(products))
	}

	// todo more chcecks
}
