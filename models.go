package main

import (
	"time"
)

// This is the model for the user
type GeoLocation struct {
	Lat  float64 `json:"lat,string"`
	Long float64 `json:"long,string"`
}

type Address struct {
	Number      int         `json:"number"`
	Street      string      `json:"street"`
	City        string      `json:"city"`
	Zipcode     string      `json:"zipcode"`
	GeoLocation GeoLocation `json:"geolocation"`
}

type Name struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type User struct {
	ID       int     `json:"id"`
	Name     Name    `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Phone    string  `json:"phone"`
	Address  Address `json:"address"`
	V        int     `json:"__v"`
}

// This is the model for the product
type Rating struct {
	Count int     `json:"count"`
	Rate  float64 `json:"rate"`
}

type Product struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Rating      Rating  `json:"rating"`
}

// This is the model for the cart
type CartProduct struct {
	ProductID int     `json:"productId"`
	Quantity  float64 `json:"quantity"`
}

type Cart struct {
	ID            int           `json:"id"`
	UserID        int           `json:"userId"`
	Date          time.Time     `json:"date"`
	CartsProducts []CartProduct `json:"products"`
	V             int           `json:"__v"`
}

type CategoryValue struct {
	Category string
	Value    float64
}
